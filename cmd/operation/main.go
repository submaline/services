package main

import (
	"context"
	"database/sql"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/bufbuild/connect-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/submaline/services/database"
	"github.com/submaline/services/gen/protocol/operation/v1/operationv1connect"
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

	// rabbitmqの準備
	rabbitConn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("failed to setup rabbitmq: %v", err)
	}
	defer rabbitConn.Close()

	// ログ
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to setup zap: %v", err)
	}
	defer logger.Sync()

	//
	supervisorClient := supervisorv1connect.NewSupervisorServiceClient(
		http.DefaultClient,
		fmt.Sprintf("http://%s:%s", os.Getenv("SUPERVISOR_SERVICE_HOST"), os.Getenv("SUPERVISOR_SERVICE_PORT")),
	)

	// サービスの準備
	operationServer := &server.OperationServer{
		DB:       dbClient,
		Auth:     authClient,
		Rb:       rabbitConn,
		Logger:   logger,
		SvClient: &supervisorClient,
	}

	// ハンドラの準備
	mux := http.NewServeMux()
	interceptors := connect.WithInterceptors(
		interceptor.NewAuthInterceptor(authClient))
	mux.Handle(operationv1connect.NewOperationServiceHandler(
		operationServer,
		interceptors,
	))

	addr := fmt.Sprintf("0.0.0.0:%v", os.Getenv("OPERATION_SERVICE_PORT"))

	// 起動
	if err := http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatal(err)
	}
}