package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/vipulvpatil/airetreat-go/internal/config"
	"github.com/vipulvpatil/airetreat-go/internal/health"
	"github.com/vipulvpatil/airetreat-go/internal/server"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	"github.com/vipulvpatil/airetreat-go/internal/tls"
	"github.com/vipulvpatil/airetreat-go/internal/workers"
	pb "github.com/vipulvpatil/airetreat-go/protos"
	"google.golang.org/grpc"
)

const WORKER_NAMESPACE = "airetreat_go"

func main() {
	fmt.Println("Starting Service")
	rand.Seed(time.Now().UTC().UnixNano())

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

	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(cfg.RedisUrl)
		},
	}
	jobStarter := workers.NewJobStarter(WORKER_NAMESPACE, redisPool)

	serverDeps := server.ServerDependencies{
		Storage: dbStorage,
		Config:  cfg,
	}

	s, err := server.NewServer(serverDeps)
	if err != nil {
		log.Fatalf("Unable to initialize new server: %v", err)
	}
	grpcServer := setupGrpcServer(s, cfg)

	hs := &health.AiRetreatGoHealthService{}
	grpcHealthServer := setupGrpcHealthServer(hs, cfg)

	workerPooldeps := workers.PoolDependencies{
		RedisPool: redisPool,
		Namespace: WORKER_NAMESPACE,
		Storage:   dbStorage,
	}
	workerPool := workers.NewPool(workerPooldeps)
	workerPool.Start()

	var wg sync.WaitGroup
	startGrpcServerAsync("ai retreat go", &wg, grpcServer, "9000")
	startGrpcServerAsync("health", &wg, grpcHealthServer, "9090")

	gameHandlerLoopCtx, cancelGameHandlerLoop := context.WithCancel(context.Background())
	gameHandlerLoop(gameHandlerLoopCtx, &wg, jobStarter)

	osTermSig := make(chan os.Signal, 1)
	signal.Notify(osTermSig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Everything started correctly")

	<-osTermSig

	cancelGameHandlerLoop()
	grpcHealthServer.GracefulStop()
	grpcServer.GracefulStop()
	workerPool.Stop()
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
		fmt.Println("using TLS")
		return grpc.Creds(tlsCredentials)
	}
	return nil
}

func gameHandlerLoop(ctx context.Context, wg *sync.WaitGroup, jobStarter workers.JobStarter) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			job, err := jobStarter.EnqueueUnique(workers.START_GAME_ONCE_PLAYERS_HAVE_JOINED, work.Q{"id": "1"})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(job)
			}
		case <-ctx.Done():
			return
		}
	}
}
