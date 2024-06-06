package main

import (
	"api_gateway/proto/auth"
	"api_gateway/proto/blogs"
	"api_gateway/proto/followings"
	"api_gateway/proto/tours"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Address                   string
	StakeholderServiceAddress string
	BlogServiceAddress        string
	FollowingServiceAdress    string
	ToursServiceAddress       string
}

func main() {
	cfg := Config{
		Address:                   "api_gateway:8000",
		StakeholderServiceAddress: "stakeholder_service:8001",
		BlogServiceAddress:        "blog_service:8002",
		FollowingServiceAdress:    "user_management_service:8003",
		ToursServiceAddress:       "tour_service:8004",
	}

	gwmux := runtime.NewServeMux()

	// Connect to the Stakeholder Service
	conn, err := grpc.DialContext(
		context.Background(),
		cfg.StakeholderServiceAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial stakeholder server:", err)
	}
	defer conn.Close()

	stakeholderClient := auth.NewStakeholderServiceClient(conn)
	err = auth.RegisterStakeholderServiceHandlerClient(
		context.Background(),
		gwmux,
		stakeholderClient,
	)
	if err != nil {
		log.Fatalln("Failed to register stakeholder gateway:", err)
	}

	// Connect to the Blog Service
	blogConn, err := grpc.DialContext(
		context.Background(),
		cfg.BlogServiceAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial blog server:", err)
	}
	defer blogConn.Close()

	blogClient := blogs.NewBlogServiceClient(blogConn)
	err = blogs.RegisterBlogServiceHandlerClient(
		context.Background(),
		gwmux,
		blogClient,
	)
	if err != nil {
		log.Fatalln("Failed to register blog gateway:", err)
	}

	//Connect to the Followers Service
	followCon, err := grpc.DialContext(
		context.Background(),
		cfg.FollowingServiceAdress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial following service:", err)
	}
	defer followCon.Close()

	followClient := followings.NewFollowerServiceClient(followCon)
	err = followings.RegisterFollowerServiceHandlerClient(
		context.Background(),
		gwmux,
		followClient,
	)
	if err != nil {
		log.Fatalln("Failed to register following gateway:", err)
	}

	//Connect to the Tour Service
	tourConn, err := grpc.DialContext(
		context.Background(),
		cfg.ToursServiceAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial tour service:", err)
	}

	tourClient := tours.NewTourServiceClient(tourConn)
	err = tours.RegisterTourServiceHandlerClient(
		context.Background(),
		gwmux,
		tourClient,
	)
	if err != nil {
		log.Fatalln("Failed to register tourr gateway:", err)
	}

	keypointClient := tours.NewKeypointServiceClient(tourConn)
	err = tours.RegisterKeypointServiceHandlerClient(
		context.Background(),
		gwmux,
		keypointClient,
	)
	if err != nil {
		log.Fatalln("Failed to register tourr gateway:", err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                            // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, // Allow specific HTTP methods
		AllowedHeaders:   []string{"*"},                            // Allow all headers
		AllowCredentials: true,                                     // Allow sending credentials (e.g., cookies)
	})

	handler := c.Handler(gwmux)

	gwServer := &http.Server{
		Addr:    cfg.Address,
		Handler: handler,
	}

	go func() {
		if err := gwServer.ListenAndServe(); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	// Graceful shutdown

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	<-stopCh

	if err = tourGwServer.Close(); err != nil {
		log.Fatalln("error while stopping server: ", err)
	}
}
