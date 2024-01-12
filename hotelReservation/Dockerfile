FROM golang:1.20

# Custom cache invalidation
ARG CACHEBUST=1

RUN git config --global http.sslverify false
COPY . /go/src/github.com/harlow/go-micro-services
WORKDIR /go/src/github.com/harlow/go-micro-services
# RUN go mod init
RUN go get gopkg.in/mgo.v2
RUN go get github.com/bradfitz/gomemcache/memcache
RUN go get github.com/google/uuid
RUN go get github.com/grafana/pyroscope-go
# RUN go get -u github.com/EIRNf/grpc_lb@v0.0.2
# RUN go get github.com/mbobakov/grpc-consul-resolver
# RUN go get github.com/opentracing-contrib/go-grpc
# RUN go get github.com/EIRNf/grpc_lb/lb/consul
# RUN go get github.com/harlow/go-micro-services/dialer
RUN go mod tidy
# COPY ./vendor/github.com/olivere/grpc /go/src/github.com/olivere/grpc
# RUN go get google.golang.org/grpc/naming@v1.29.0

# RUN echo "replace google.golang.org/grpc/naming => google.golang.org/grpc/naming v1.29.0" >> go.mod
# RUN echo "replace github.com/EIRNf/grpc_lb => ./vendor/github.com/EIRNf/grpc_lb" >> go.mod
# RUN echo "replace github.com/olivere/grpc/lb/consul => github.com/EIRNf/grpc_lb/lb/consul v1.0.0"
# RUN echo "replace github.com/EIRNf/grpc_lb/lb/naming => ./vendor/github.com/EIRNf/grpc_lb/lb/naming"
RUN go mod vendor
# # use local modified github.com/hailocab/go-geoindex instead of upstream one
COPY ./vendor/github.com/hailocab/go-geoindex /go/src/github.com/harlow/go-micro-services/vendor/github.com/hailocab/go-geoindex 
COPY ./vendor/github.com/fullstorydev/grpchan /go/src/github.com/fullstorydev/grpchan
RUN go install -ldflags="-s -w" ./cmd/...
