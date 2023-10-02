module geo

go 1.20

require github.com/fullstorydev/grpchan v1.1.1

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.1
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/hailocab/go-geoindex v0.0.0-20160127134810-64631bfe9711
	github.com/rs/zerolog v1.30.0
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20230913181813-007df8e322eb // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230803162519-f966b187b2e5 // indirect
	google.golang.org/grpc v1.58.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

replace github.com/fullstorydev/grpchan v1.1.1 => /home/estebanramos/app_testing/grpchan/
replace github.com/hailocab/go-geoindex  => /home/estebanramos/projects/microservices_benchmarks/hotelReservation/vendor/github.com/hailocab/go-geoindex
