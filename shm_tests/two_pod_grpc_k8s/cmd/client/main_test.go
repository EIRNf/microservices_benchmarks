package main

import (
	"testing"
	"time"

	"github.com/fullstorydev/grpchan/shmgrpc"
	"github.com/fullstorydev/grpchan/test_hello_service"
)

// func main() {

// 	testing.Init()
// 	b := testing.B{}

// 	testing.Benchmark(BenchmarkGrpcOverSharedMemory(*b))

// }

func BenchmarkGrpcOverSharedMemory(b *testing.B) {

	cc := shmgrpc.NewChannel("localhost", "http://127.0.0.1:8080/hello")

	time.Sleep(10 * time.Second)

	test_hello_service.RunChannelBenchmarkCases(b, cc, false)

}
