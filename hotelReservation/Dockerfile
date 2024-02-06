FROM golang:1.21.3

# Custom cache invalidation
ARG CACHEBUST=1
ENV CGO_ENABLED=1
# ENV GRPC_GO_LOG_VERBOSITY_LEVEL=99
# ENV GRPC_GO_LOG_SEVERITY_LEVEL=info
ENV GODEBUG=http2debug=1   
ENV GODEBUG=http2debug=2


RUN git config --global http.sslverify false
COPY . /go/src/github.com/harlow/go-micro-services
WORKDIR /go/src/github.com/harlow/go-micro-services
RUN go mod tidy

# # use local modified github.com/EIRNf/notnets_grpc instead of upstream one
COPY ./vendor/github.com/EIRNf/notnets_grpc /go/src/github.com/EIRNf/notnets_grpc
# RUN go mod vendor

# # use local modified github.com/hailocab/go-geoindex instead of upstream one
COPY ./vendor/github.com/hailocab/go-geoindex /go/src/github.com/harlow/go-micro-services/vendor/github.com/hailocab/go-geoindex 
RUN go install -ldflags="-s -w" ./cmd/...
