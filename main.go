package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vipulvpatil/airetreat-go/internal/config"
	"github.com/vipulvpatil/airetreat-go/internal/health"
	"github.com/vipulvpatil/airetreat-go/internal/server"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	"github.com/vipulvpatil/airetreat-go/internal/tls"
	pb "github.com/vipulvpatil/airetreat-go/protos"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting Service")

	cfg, errs := config.NewConfigFromEnvVars()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		log.Fatal("Unable to load config. Required Env vars are missing")
	}

	db, err := storage.InitDb(cfg)
	if err != nil {
		log.Fatalf("Unable to initialize database: %v", err)
	}

	dbStorage, err := storage.NewDbStorage(
		storage.StorageOptions{
			Db: db,
		},
	)
	if err != nil {
		log.Fatalf("Unable to initialize storage: %v", err)
	}

	serverDeps := server.ServerDependencies{
		Storage: dbStorage,
	}

	s, err := server.NewServer(serverDeps)
	if err != nil {
		log.Fatalf("Unable to initialize new server: %v", err)
	}
	grpcServer := setupGrpcServer(s, cfg)

	hs := &health.AiRetreatGoHealthService{}
	grpcHealthServer := setupGrpcHealthServer(hs, cfg)

	var wg sync.WaitGroup
	startGrpcServerAsync("ai retreat go", &wg, grpcServer, "9000")
	startGrpcServerAsync("health", &wg, grpcHealthServer, "9090")

	osTermSig := make(chan os.Signal, 1)
	signal.Notify(osTermSig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Everything started correctly")

	<-osTermSig

	grpcHealthServer.GracefulStop()
	grpcServer.GracefulStop()
	wg.Wait()
	fmt.Println("Stopping Service")
}

func setupGrpcServer(s *server.AiRetreatGoService, cfg *config.Config) *grpc.Server {
	serverOpts := make([]grpc.ServerOption, 0)
	tlsServerOpts := tlsGrpcServerOptions(cfg)
	if tlsServerOpts != nil {
		serverOpts = append(serverOpts, tlsServerOpts)
	}
	serverOpts = append(serverOpts, grpc.UnaryInterceptor(s.RequestingUserInterceptor))
	grpcServer := grpc.NewServer(serverOpts...)
	pb.RegisterAiRetreatGoServer(grpcServer, s)
	return grpcServer
}

func setupGrpcHealthServer(hs *health.AiRetreatGoHealthService, cfg *config.Config) *grpc.Server {
	serverOpts := make([]grpc.ServerOption, 0)
	tlsServerOpts := tlsGrpcServerOptions(cfg)
	if tlsServerOpts != nil {
		serverOpts = append(serverOpts, tlsServerOpts)
	}
	grpcServer := grpc.NewServer(serverOpts...)
	pb.RegisterAiRetreatGoHealthServer(grpcServer, hs)
	return grpcServer
}

func startGrpcServerAsync(name string, wg *sync.WaitGroup, grpcServer *grpc.Server, port string) {
	wg.Add(1)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		fmt.Printf("Starting GRPC Server: %s\n", name)
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("GrpcServer %s failed to start: %v", name, err)
		}
		fmt.Printf("Stopping GRPC Server: %s\n", name)
		wg.Done()
	}()
}

func tlsGrpcServerOptions(cfg *config.Config) grpc.ServerOption {
	if cfg.EnableTls {
		tlsCredentials, err := tls.LoadTLSCredentials(cfg)
		if err != nil {
			log.Fatal("cannot load TLS credentials: ", err)
		}
		return grpc.Creds(tlsCredentials)
	}
	return nil
}
