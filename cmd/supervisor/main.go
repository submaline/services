package main

import (
	"context"
	"database/sql"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/bwmarrin/snowflake"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/submaline/services/database"
	"github.com/submaline/services/gen/protocol/supervisor/v1/supervisorv1connect"
	"github.com/submaline/services/interceptor"
	"github.com/submaline/services/server"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// firebaseの準備
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed to setup firebase app: %v", err)
	}

	// firebase authの準備
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("failed to setup firebase auth: %v", err)
	}

	// databaseの準備
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MARIADB_USER"),
		os.Getenv("MARIADB_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("MARIADB_DATABASE")))
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	dbClient := &database.DBClient{DB: db}

	// id generatorの準備
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("failed to setup id gererator: %v", err)
	}

	// rabbitmqの準備
	rabbitConn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("failed to setup rabbitmq: %v", err)
	}
	defer rabbitConn.Close()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to setup zap: %v", err)
	}
	defer logger.Sync()

	// サービスの準備
	supervisorServer := &server.SupervisorServer{
		DB:     dbClient,
		Auth:   authClient,
		Id:     node,
		Rb:     rabbitConn,
		Logger: logger,
	}

	// ハンドラの準備
	mux := http.NewServeMux()
	interceptors := connect.WithInterceptors(
		interceptor.NewAuthInterceptor(authClient))
	mux.Handle(supervisorv1connect.NewSupervisorServiceHandler(
		supervisorServer,
		interceptors,
	))

	addr := fmt.Sprintf("0.0.0.0:%v", os.Getenv("SUPERVISOR_SERVICE_PORT"))

	// 起動
	if err := http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatal(err)
	}
}