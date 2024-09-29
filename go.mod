module github.com/ServiceWeaver/weaver

go 1.22.0

toolchain go1.22.3

require (
	github.com/BurntSushi/toml v1.4.0
	github.com/DataDog/hyperloglog v0.0.0-20240529073211-ec1b7f6ebabb
	github.com/alecthomas/chroma/v2 v2.2.0
	github.com/fsnotify/fsnotify v1.7.0
	github.com/go-sql-driver/mysql v1.7.1
	github.com/goburrow/cache v0.1.4
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/cel-go v0.21.0
	github.com/google/go-cmp v0.6.0
	github.com/google/pprof v0.0.0-20240925223930-fa3061bff0bc
	github.com/google/uuid v1.6.0
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/lightstep/varopt v1.4.0
	github.com/opencontainers/runc v1.1.12
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c
	github.com/testcontainers/testcontainers-go v0.24.1
	github.com/testcontainers/testcontainers-go/modules/mysql v0.24.1
	github.com/yuin/goldmark v1.4.15
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20220924101305-151362477c87
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.55.0
	go.opentelemetry.io/otel v1.30.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.30.0
	go.opentelemetry.io/otel/sdk v1.30.0
	go.opentelemetry.io/otel/trace v1.30.0
	golang.org/x/crypto v0.27.0
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0
	golang.org/x/image v0.10.0
	golang.org/x/sync v0.8.0
	golang.org/x/term v0.24.0
	golang.org/x/text v0.18.0
	golang.org/x/tools v0.25.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240924160255-9d4c2d233b61
	google.golang.org/protobuf v1.34.2
	gorm.io/driver/postgres v1.5.3
	gorm.io/gorm v1.25.5
	k8s.io/utils v0.0.0-20240902221715-702e33fdd3c3
	modernc.org/sqlite v1.33.1
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/DataDog/mmh3 v0.0.0-20210722141835-012dc69a9e49 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/Microsoft/hcsshim v0.11.4 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cilium/ebpf v0.9.1 // indirect
	github.com/containerd/containerd v1.7.11 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/cpuguy83/dockercfg v0.3.1 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v24.0.9+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dustin/randbo v0.0.0-20140428231429-7f1b564ca724 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.4 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.16.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc4 // indirect
	github.com/opencontainers/runtime-spec v1.1.0-rc.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/shirou/gopsutil/v3 v3.23.7 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240924160255-9d4c2d233b61 // indirect
	google.golang.org/grpc v1.67.0 // indirect
	modernc.org/gc/v3 v3.0.0-20240801135723-a856999a2e4a // indirect
	modernc.org/libc v1.61.0 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.8.0 // indirect
	modernc.org/strutil v1.2.0 // indirect
	modernc.org/token v1.1.0 // indirect
)
