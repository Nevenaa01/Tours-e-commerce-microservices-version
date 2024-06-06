package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
	"tours_service/handler"
	"tours_service/proto/tours"
	"tours_service/repository"
	"tours_service/service"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func initMongoDb() *mongo.Client {

	dburi := "mongodb://mongo:27017" //ili localhost ili mongo

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		fmt.Print(err)
	}

	return client
}

func Conn() *nats.Conn {
	conn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// func manageRouter(client *mongo.Client) http.Server {

// 	// FacilityRepository := &repository.FacilityRepository{FacilityClient: client}
// 	// FacilityService := &service.FacilityService{FacilityRepository: FacilityRepository}
// 	// FacilityHandler := &handler.FacilityHandler{FacilityService: FacilityService}

// 	// KeypointRepository := &repository.KeypointRepository{KeypointClient: client}
// 	// KeypointService := &service.KeypointService{KeypointRepository: KeypointRepository}
// 	// KeypointHandler := &handler.KeypointHandler{KeypointService: KeypointService}

// 	// tourRatingRepository := &repository.TourRatingRepository{TourRatingClient: client}
// 	// tourRatingService := &service.TourRatingService{TourRatingRepository: tourRatingRepository}
// 	// tourRatingHandler := &handler.TourRatingHandler{TourRatingService: tourRatingService}

// 	// tourProblemRepository := &repository.TourProblemRepository{TourProblemClient: client}
// 	// tourProblemService := &service.TourProblemService{TourProblemRepository: tourProblemRepository,
// 	// 	TourService: tourService}
// 	// tourProblemHandler := &handler.TourProblemHandler{TourProblemService: tourProblemService}

// 	// router := mux.NewRouter().StrictSlash(true)

// 	// postFacilityRouter := router.Methods(http.MethodPost).Subrouter()
// 	// postFacilityRouter.HandleFunc("/facilities", FacilityHandler.CreateFacility)
// 	// postFacilityRouter.Use(FacilityHandler.MiddlewareFacilityDeserialization)

// 	// deleteRouter := router.Methods(http.MethodDelete).Subrouter()
// 	// deleteRouter.HandleFunc("/facilities/{id}", FacilityHandler.DeleteFacility)

// 	// postKeypointRouter := router.Methods(http.MethodPost).Subrouter()
// 	// postKeypointRouter.HandleFunc("/keypoints", KeypointHandler.CreateKeypoint)
// 	// postKeypointRouter.Use(KeypointHandler.MiddlewareKeypointDeserialization)

// 	// postTourRouter := router.Methods(http.MethodPost).Subrouter()
// 	// postTourRouter.HandleFunc("/createTour", tourHandler.CreateTour)
// 	// postTourRouter.Use(tourHandler.MiddlewareTourDeserialization)

// 	// putTourRouter := router.Methods(http.MethodPut).Subrouter()
// 	// putTourRouter.HandleFunc("/tours", tourHandler.UpdateTour)
// 	// putTourRouter.Use(tourHandler.MiddlewareTourDeserialization)

// 	// getTourRouter := router.Methods(http.MethodGet).Subrouter()
// 	// getTourRouter.HandleFunc("/tours/{id}", tourHandler.GetTourById)
// 	// getTourRouter.HandleFunc("/tours/author/{id}", tourHandler.GetToursByAuthorId)

// 	// postTourRatingRouter := router.Methods(http.MethodPost).Subrouter()
// 	// postTourRatingRouter.HandleFunc("/createTourRating", tourRatingHandler.CreateTourRating)
// 	// postTourRatingRouter.Use(tourRatingHandler.MiddlewareTourRatingDeserialization)

// 	// getTourProblemRouter := router.Methods(http.MethodGet).Subrouter()
// 	// getTourProblemRouter.HandleFunc("/getByAuthorId/{authorId}", tourProblemHandler.GetByAuthorId)

// 	// cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

// 	// server := http.Server{
// 	// 	Addr:         ":8080",
// 	// 	Handler:      cors(router),
// 	// 	IdleTimeout:  120 * time.Second,
// 	// 	ReadTimeout:  1 * time.Second,
// 	// 	WriteTimeout: 1 * time.Second,
// 	// }
// 	// return server
// }

const serviceName = "tours_service"

var tp *trace.TracerProvider

func initJaegerTracer() (*trace.TracerProvider, error) {
	log.Printf("Initializing tracing to jaeger at %s\n", "http://jaeger:14268/api/traces")
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://jaeger:14268/api/traces")))
	if err != nil {
		return nil, err
	}
	fmt.Println("error nije null")
	return trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	), nil
}

type Message struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// OpenTelemetry
	var err1 error
	tp, err1 = initJaegerTracer()
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Println("prosao err u mainu")
	defer func() {
		if err1 := tp.Shutdown(context.Background()); err1 != nil {
			log.Printf("Error shutting down tracer provider: %v", err1)
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	fmt.Println("prosao otel")

	client := initMongoDb()

	err := client.Connect(timeoutContext)
	if err != nil {
		fmt.Print(err)
	}

	logger := log.New(os.Stdout, "[logger-main] ", log.LstdFlags)

	listener, err := net.Listen("tcp", "tour_service:8004")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	tourRepository := &repository.TourRepository{TourClient: client}
	keypointRepository := &repository.KeypointRepository{KeypointClient: client}

	tourService := &service.TourService{TourRepository: tourRepository, KeypointRepository: keypointRepository}
	keypointService := &service.KeypointService{KeypointRepository: keypointRepository}

	tourHandler := &handler.TourHandler{TourService: tourService}
	keypointHandler := &handler.KeypointHandler{KeypointService: keypointService}

	tours.RegisterTourServiceServer(grpcServer, tourHandler)
	tours.RegisterKeypointServiceServer(grpcServer, keypointHandler)

	conn := Conn()
	_, errSub := conn.Subscribe("subTours", func(message *nats.Msg) {
		var messageRec Message
		err := json.Unmarshal(message.Data, &messageRec)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			return
		}
		fmt.Printf("RECEIVED MESSAGE: %s\n", messageRec.Body)
		tour, errGet := tourRepository.GetById(messageRec.Id)
		if errGet != nil {
			fmt.Println(errGet)
			return
		}
		if messageRec.Body == "Failed" {
			tour.State = 1
		} else {
			tour.State = 2
		}
		tourRepository.Update(tour)
	})
	if errSub != nil {
		log.Fatal(err)
	}
	//server := manageRouter(client)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	// if server.Shutdown(timeoutContext) != nil {
	// 	logger.Fatal("Cannot gracefully shutdown...")
	// }
	// logger.Println("Server stopped")

	print("ok")
}
