package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"stakeholders_service/handl"
	"stakeholders_service/proto/auth"
	"syscall"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	//connStr := "user=postgres dbname=explorer-v1 password=super sslmode=disable"
	connStr := "user=postgres dbname=explorer password=super sslmode=disable port=5432 host=database"
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
func Conn() *nats.Conn {
	conn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

type Message struct {
	Id     int    `json:"id"`
	Body   string `json:"body"`
	UserId int    `json:"userId"`
}

func main() {

	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}

	listener, err := net.Listen("tcp", "stakeholder_service:8001")
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

	authHandler := handl.AuthHandler{DatabaseConnection: database, Key: "explorer_secret_key"}
	auth.RegisterStakeholderServiceServer(grpcServer, authHandler)

	conn := Conn()
	_, errSub := conn.Subscribe("subStakeholders", func(message *nats.Msg) {
		var messageRec Message
		err := json.Unmarshal(message.Data, &messageRec)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			return
		}
		fmt.Printf("RECEIVED MESSAGE: %s\n", messageRec.Body)
		authHandler.ChangeRating(messageRec.Id, messageRec.UserId)
	})
	if errSub != nil {
		log.Fatal(err)
	}
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	grpcServer.Stop()
}
