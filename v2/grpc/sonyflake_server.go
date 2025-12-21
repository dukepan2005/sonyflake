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
	Port int `short:"p" long:"port" description:"the port grpc server listen" required:"true"`
	// Debug        bool `long:"debug" description:"debug mode to print debug information"`
	AwsMachineID bool `long:"aws" description:"if use awsutil.AmazonEC2MachineID to canculate the unique ID of the Sonyflake instance"`
}

// https://github.com/grpc/grpc-go/tree/master/examples/features/keepalive
// https://blog.ivansli.com/2022/02/05/grpc-keepalive/
// https://pandaychen.github.io/2020/09/01/GRPC-CLIENT-CONN-LASTING/

// MinTime：如果客户端两次 ping 的间隔小于 5s，则关闭连接
// PermitWithoutStream： 即使没有 active stream, 也允许 ping
var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

// MaxConnectionIdle：如果一个 client 空闲超过 15s, 发送一个 GOAWAY, 为了防止同一时间发送大量 GOAWAY, 会在 15s 时间间隔上下浮动 15*10%, 即 15+1.5 或者 15-1.5
// MaxConnectionAge：如果任意连接存活时间超过 30s, 发送一个 GOAWAY
// MaxConnectionAgeGrace：在强制关闭连接之间, 允许有 5s 的时间完成 pending 的 rpc 请求
// Time： 如果一个 client 空闲超过 5s, 则发送一个 ping 请求
// Timeout： 如果 ping 请求 1s 内未收到回复, 则认为该连接已断开
var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 30 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,  // Wait 10 second for the ping ack before assuming the connection is dead
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

func main() {
	// config log
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// 1. new grpc server
	rpcServer := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))

	// 2. new service
	grpcSonyService := &sonyflake.GrpcSonyflakeService{Sf: sf}

	// 3. register service into grpc server
	sonyflake.RegisterSonyflakeServiceServer(rpcServer, grpcSonyService)

	// https://blog.51cto.com/u_14592069/5711900
	// 不用defer listener.close(), 因为进程退出后，它的文件描述符会被自动清理
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
