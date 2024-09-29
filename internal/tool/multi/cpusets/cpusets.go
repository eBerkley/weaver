package cpusets

import (
	"fmt"
	"time"

	"github.com/ServiceWeaver/weaver/internal/must"
	"github.com/ServiceWeaver/weaver/runtime"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/manager"
	"github.com/opencontainers/runc/libcontainer/configs"
	"k8s.io/utils/cpuset"
)

const (
	USE_CGROUPS                    = true
	shutdownAttempts               = 5
	shutdownTimeout  time.Duration = 500 * time.Millisecond
)

var (
	unused     cpuset.CPUSet
	headCgroup cgroups.Manager
	subCgroups = make(map[string]cgroups.Manager)
)

func init() {
	if USE_CGROUPS {
		runtime.OnExitSignal(Cleanup)
	}
}

// Cleanup removes the cgroups created by this package.
//
// Because all processes within cgroups need to be terminated before the cgroup can be deleted, this function will attempt to shutdown multiple times on a timeout.
func Cleanup() {
	if !USE_CGROUPS {
		return
	}

	var err error
	if headCgroup == nil {
		return
	}

	// fmt.Println("trying to delete...")
	for i := 0; i < shutdownAttempts; i++ {
		err = headCgroup.Destroy()
		if err == nil {
			return
		}
		time.Sleep(shutdownTimeout)
	}
	fmt.Printf("could not delete cgroups because of error %v", err)
}

// InitCPUs initializes the cpusets package, and must be called before any other function.
//
// all_cpus is the list of cpus that processes will be allowed to pull from, and should be a Linux CPU list formatted string.
func InitCPUs(all_cpus string) {
	if !USE_CGROUPS {
		return
	}
	// assign initial set of CPU cores to be used
	unused = must.Must(cpuset.Parse(all_cpus))

	// create base cgroup with all_cpus as the cpuset
	headCgroup = must.Must(manager.New(&configs.Cgroup{
		Name:      "head",
		Resources: &configs.Resources{CpusetCpus: all_cpus},
	}))

	// create vfs.
	if err := headCgroup.Apply(-1); err != nil {
		panic(err)
	}
	// set cpuset resource config
	if err := headCgroup.Set(nil); err != nil {
		panic(err)
	}
}

func RemainingCpus() string {
	return unused.String()
}

// NewCgroup creates a new cgroup that restricts the CPUs processes in the cgroup are allowed to use.
// No other feature of cgroups are currently used.
//
// cpu_req should be a Linux CPU list formatted string.
func NewCgroup(name string, cpu_req string) error {
	if !USE_CGROUPS {
		return nil
	}

	cpus, err := requestCPUs(cpu_req)
	if err != nil {
		return err
	}

	return createCgroup(name, cpus)
}

// AddPidToCgroup restrics a process to a cgroup created with NewCgroup.
// name should be a valid cgroup that was passed to NewCgroup prior without error.
func AddPidToCgroup(name string, pid int) error {
	if !USE_CGROUPS {
		return nil
	}

	c, ok := subCgroups[name]
	if !ok || !c.Exists() {
		return fmt.Errorf("cgroup name %s not found", name)
	}

	return c.Apply(pid)
}

func requestCPUs(s string) (cpuset.CPUSet, error) {
	if !USE_CGROUPS {
		return cpuset.CPUSet{}, nil
	}
	req := must.Must(cpuset.Parse(s))
	if !req.IsSubsetOf(unused) {
		return cpuset.CPUSet{}, fmt.Errorf("CPUs %v already used", req.Difference(unused).List())
	}
	unused = unused.Difference(req)
	return req, nil
}

// Assumes that RequestCPUs has been called prior.
func createCgroup(name string, cpus cpuset.CPUSet) error {
	if !USE_CGROUPS {
		return nil
	}

	if _, ok := subCgroups[name]; ok {
		return fmt.Errorf("cgroup name %s already used", name)
	}

	c, err := manager.New(&configs.Cgroup{
		Name:   name,
		Parent: "head",
	})
	if err != nil {
		return err
	}

	if err = c.Apply(-1); err != nil {
		return err
	}
	if err = c.Set(&configs.Resources{CpusetCpus: cpus.String()}); err != nil {
		return err
	}
	subCgroups[name] = c
	return nil
}
