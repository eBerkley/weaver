// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package weaver

import (
	"context"
	"crypto/tls"
	"math/rand"
	"sync"

	"github.com/eberkley/weaver/internal/cond"
	"github.com/eberkley/weaver/internal/net/call"
	"github.com/eberkley/weaver/runtime/protos"
	"github.com/google/uuid"
)

// routingBalancer balances requests according to a routing assignment.
type routingBalancer struct {
	balancer  call.Balancer // balancer to use for non-routed calls
	tlsConfig *tls.Config   // tls config to use; may be nil.

	stateful      bool // If we do routing using statefulSlice or assignment/index.
	statefulSlice []string
	mu            sync.RWMutex
	assignment    *protos.Assignment
	index         index

	// Map from address to connection. We currently allow just one
	// connection per address.
	// Guarded by mu.
	conns map[string]call.ReplicaConnection
}

// newRoutingBalancer returns a new routingBalancer.
func newRoutingBalancer(tlsConfig *tls.Config) *routingBalancer {
	return &routingBalancer{
		balancer:  call.RoundRobin(),
		tlsConfig: tlsConfig,
		conns:     map[string]call.ReplicaConnection{},
	}
}

// Add adds c to the set of connections we are balancing across.
func (rb *routingBalancer) Add(c call.ReplicaConnection) {
	rb.balancer.Add(c)

	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.conns[c.Address()] = c
}

// Remove removes c from the set of connections we are balancing across.
func (rb *routingBalancer) Remove(c call.ReplicaConnection) {
	rb.balancer.Remove(c)

	rb.mu.Lock()
	defer rb.mu.Unlock()
	delete(rb.conns, c.Address())
}

// makes rb use stateful routing.
func (rb *routingBalancer) updateStateful(replicas []string) {
	rb.mu.Lock()
	rb.statefulSlice = replicas
	rb.stateful = true
	rb.mu.Unlock()
}

// update updates the balancer with the provided assignment
// makes rb not use stateful routing
func (rb *routingBalancer) update(assignment *protos.Assignment) {
	if assignment == nil {
		return
	}

	index := newIndex(assignment)
	rb.mu.Lock()
	rb.stateful = false
	defer rb.mu.Unlock()
	rb.assignment = assignment
	rb.index = index
}

// Used for routed method calls that may be local.
// If the shardKey points to the same address as dialAddr, returns true.
func (rb *routingBalancer) IsLocal(shardKey uint64, dialAddr string) bool {
	rb.mu.RLock()
	if rb.stateful {
		ret := rb.statefulSlice[shardKey] == dialAddr
		rb.mu.RUnlock()
		return ret
	}

	if shardKey == 0 {
		rb.mu.RUnlock()
		return true
	}

	assignment := rb.assignment
	index := rb.index
	rb.mu.RUnlock()

	if assignment == nil {
		// There is no assignment. This is possible if we haven't received an
		// assignment from the assigner yet.
		// Probably will end up in a failure.
		return true
	}

	slice, ok := index.find(shardKey)
	if !ok {
		// Should be impossible.
		return true
	}

	rb.mu.RLock()
	defer rb.mu.RUnlock()

	for _, addr := range slice.replicas {
		if addr == dialAddr {
			return true
		}
	}
	return false
}

// Pick implements the call.Balancer interface.
func (rb *routingBalancer) Pick(opts call.CallOptions) (call.ReplicaConnection, bool) {
	rb.mu.RLock()
	if rb.stateful {
		defer rb.mu.RUnlock()
		if len(rb.statefulSlice) <= int(opts.ShardKey) {
			return nil, false
		}
		if c, ok := rb.conns[rb.statefulSlice[opts.ShardKey]]; ok {
			return c, true
		}
		return nil, false
	}

	if opts.ShardKey == 0 {
		// If the method we're calling is not sharded (which is guaranteed to
		// be true for nonsharded components), then the shard key is 0.
		rb.mu.RUnlock()
		return rb.balancer.Pick(opts)
	}

	// Grab the current assignment. It's possible that the current assignment
	// changes between when we release the lock and when we pick an endpoint,
	// but using a slightly stale assignment is okay.
	assignment := rb.assignment
	index := rb.index
	rb.mu.RUnlock()

	if assignment == nil {
		// There is no assignment. This is possible if we haven't received an
		// assignment from the assigner yet.
		return rb.balancer.Pick(opts)
	}

	slice, ok := index.find(opts.ShardKey)
	if !ok {
		// TODO(mwhittaker): Shouldn't this be impossible. Understand better
		// when this happens.
		return rb.balancer.Pick(opts)
	}

	// Search for an available ReplicConnection starting at a random offset.
	// TODO(sanjay):Precompute the set of available ReplicaConnections per slice.
	offset := rand.Intn(len(slice.replicas))
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	for i, n := 0, len(slice.replicas); i < n; i++ {
		offset++
		if offset == n {
			offset = 0
		}
		if c, ok := rb.conns[slice.replicas[offset]]; ok {
			return c, true
		}
	}
	return nil, false
}

// routingResolver is a dummy resolver that returns whatever endpoints are
// passed to the update method.
type routingResolver struct {
	m         sync.Mutex      // guards all of the following fields
	changed   cond.Cond       // fires when endpoints changes
	version   *call.Version   // the current version of endpoints
	endpoints []call.Endpoint // the endpoints returned by Resolve
}

// newRoutingResolver returns a new routingResolver.
func newRoutingResolver() *routingResolver {
	r := &routingResolver{
		version: &call.Version{Opaque: call.Missing.Opaque},
	}
	r.changed.L = &r.m
	return r
}

// IsConstant implements the call.Resolver interface.
func (rr *routingResolver) IsConstant() bool { return false }

// update updates the resolver with the provided endpoints.
func (rr *routingResolver) update(endpoints []call.Endpoint) {
	rr.m.Lock()
	defer rr.m.Unlock()
	rr.version = &call.Version{Opaque: uuid.New().String()}
	rr.endpoints = endpoints
	rr.changed.Broadcast()
}

// Resolve implements the call.Resolver interface.
func (rr *routingResolver) Resolve(ctx context.Context, version *call.Version) ([]call.Endpoint, *call.Version, error) {
	rr.m.Lock()
	defer rr.m.Unlock()

	if version == nil {
		return rr.endpoints, rr.version, nil
	}

	for *version == *rr.version {
		if err := rr.changed.Wait(ctx); err != nil {
			return nil, nil, err
		}
	}
	return rr.endpoints, rr.version, nil
}
