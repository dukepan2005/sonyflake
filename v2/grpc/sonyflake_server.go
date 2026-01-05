package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/sony/sonyflake/v2"
	"github.com/sony/sonyflake/v2/awsutil"
)

var sf *sonyflake.Sonyflake

var opts struct {
	Port         int  `short:"p" long:"port" description:"the port grpc server listen" required:"true"`
	AwsMachineID bool `long:"aws" description:"if use awsutil.AmazonEC2MachineID to canculate the unique ID of the Sonyflake instance"`
}

// https://github.com/grpc/grpc-go/tree/master/examples/features/keepalive
// https://blog.ivansli.com/2022/02/05/grpc-keepalive/
// https://pandaychen.github.io/2020/09/01/GRPC-CLIENT-CONN-LASTING/

// MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
// PermitWithoutStream: true,            // Allow pings even when there are no active streams
var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

// MaxConnectionIdle:     15 * time.Second,  // If a client is idle for 15 seconds, send a GOAWAY
// MaxConnectionAge:      300 * time.Second, // If any connection is alive for more than 300 seconds, send a GOAWAY
// MaxConnectionAgeGrace: 5 * time.Second,   // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
// Time:                  5 * time.Second,   // Ping the client if it is idle for 30 seconds to ensure the connection is still active
// Timeout:               1 * time.Second,   // Wait 10 second for the ping ack before assuming the connection is dead
var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second,  // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      300 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,   // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,   // Ping the client if it is idle for 30 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,   // Wait 10 second for the ping ack before assuming the connection is dead
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

	sf, err = sonyflake.New(st)
	if err != nil {
		panic(err)
	}
	if sf == nil {
		panic("sonyflake not created")
	}
}

// set channel buffer size to 4096 or higher to reduce the impact of GC pauses, if you run on machines with more CPU cores
// set it according to your QPS and GC performance
// currently we pre-generate ids into cache in a separate goroutine, and set cache size to 1024
var idCache = make(chan int64, 1024)

func main() {
	// config log
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// pre-generate ids into cache
	go func() {
		for {
			id, err := sf.NextID()
			if err != nil {
				time.Sleep(time.Millisecond) // wait a monent if error occurs
				continue
			}
			idCache <- int64(id) // wait here if blocked, until there is space in the channel, it don't affect grpc goroutines serving
		}
	}()

	// 1. new grpc server
	rpcServer := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))

	// 2. new service
	grpcSonyService := &sonyflake.GrpcSonyflakeService{IDCache: idCache}

	// 3. register service into grpc server
	sonyflake.RegisterSonyflakeServiceServer(rpcServer, grpcSonyService)

	// https://blog.51cto.com/u_14592069/5711900
	// no defrer close, because when the process exits, its file descriptors will be automatically cleaned up
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
	// https://www.cnblogs.com/sky-heaven/p/10176422.html
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Server running")
	s := <-quit

	log.Printf("Got signal %v, attempting graceful shutdown\n", s)

	rpcServer.GracefulStop()
	log.Println("Server exit successfully")
}
