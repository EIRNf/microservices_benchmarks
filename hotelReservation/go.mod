module hotelReservation

go 1.21

toolchain go1.22.1

replace github.com/harlow/go-micro-services => ./

replace github.com/EIRNf/notnets_grpc => ./notnets_grpc

// replace github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.5.1

// replace go.uber.org/atomic => github.com/uber-go/atomic v1.11.0

require (
	github.com/EIRNf/notnets_grpc v0.0.0-20240401041235-0157d8f3a918
	github.com/bradfitz/gomemcache v0.0.0-20230905024940-24af94b03874
	github.com/golang/protobuf v1.5.4
	github.com/google/uuid v1.6.0
	github.com/grafana/pyroscope-go v1.1.1
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/hailocab/go-geoindex v0.0.0-20160127134810-64631bfe9711
	github.com/harlow/go-micro-services v0.0.0-20231215003052-c9a9a17f2007
	github.com/hashicorp/consul/api v1.28.2
	github.com/mbobakov/grpc-consul-resolver v1.5.3
	github.com/opentracing-contrib/go-stdlib v1.0.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/rs/zerolog v1.32.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	golang.org/x/net v0.22.0
	google.golang.org/grpc v1.62.1
	google.golang.org/protobuf v1.33.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

require (
	github.com/bufbuild/protocompile v0.9.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
)

require (
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/fullstorydev/grpchan v1.1.1
	github.com/go-playground/form v3.1.4+incompatible // indirect
	github.com/grafana/pyroscope-go/godeltaprof v0.1.7 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/ianlancetaylor/cgosymbolizer v0.0.0-20240326020559-581a3f7c677f
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda // indirect
	modernc.org/libc v1.49.0 // indirect
)

replace github.com/uber-go/atomic => github.com/uber-go/atomic v1.4.0
