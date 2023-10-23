FROM golang:1.21

RUN git config --global http.sslverify false
COPY  . /go/src/two_pod_grpc_k8s
WORKDIR /go/src/two_pod_grpc_k8s
RUN go mod tidy

RUN go build cmd/server/main.go
RUN go test -c cmd/client/main_test.go 