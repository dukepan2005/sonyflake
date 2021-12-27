package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dukepan2005/sonyflake"
	"github.com/dukepan2005/sonyflake/awsutil"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
)

var sf *sonyflake.Sonyflake

var opts struct {
	Port int `short:"p" long:"port" description:"the port grpc server listen" required:"true"`
	// Debug        bool `long:"debug" description:"debug mode to print debug information"`
	AwsMachineID bool `long:"aws" description:"if use awsutil.AmazonEC2MachineID to canculate the unique ID of the Sonyflake instance"`
}

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	var st sonyflake.Settings
	if opts.AwsMachineID {
		st.MachineID = awsutil.AmazonEC2MachineID
	}

	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func main() {
	// 1. new grpc server
	rpcServer := grpc.NewServer()

	// 2. new service
	grpcSonyService := &sonyflake.GrpcSonyflakeService{
		Sf: sf,
	}

	// 3. register service into grpc server
	sonyflake.RegisterSonyflakeServiceServer(rpcServer, grpcSonyService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.Port))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Server listen on: %d", opts.Port)
	}

	go func() {
		// service connections
		if err := rpcServer.Serve(listener); err != nil {
			log.Panic(err.Error())
		}
	}()

	// 4. granceful stop
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	log.Println("Server running")
	s := <-quit

	log.Printf("Got signal %v, attempting graceful shutdown\n", s)

	rpcServer.GracefulStop()
	log.Println("Server exit successfully")
}
