package main

import (
	"github.com/fullstorydev/grpchan/shmgrpc"
	"github.com/fullstorydev/grpchan/test_hello_service"
)

func main() {

	// svr := &grpchantesting.TestServer{}
	svc := &test_hello_service.TestServer{}
	svr := shmgrpc.NewServer("/hello")

	//Register Server and instantiate with necessary information
	test_hello_service.RegisterTestServiceServer(svr, svc)

	//Create Listener
	lis := shmgrpc.Listen("http://127.0.0.1:8080/hello")

	svr.Serve(lis)
	defer svr.Stop()

}
