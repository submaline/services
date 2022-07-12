package main

import (
	"context"
	"database/sql"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/bufbuild/connect-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/submaline/services/db"
	"github.com/submaline/services/gen/supervisor/v1/supervisorv1connect"
	"github.com/submaline/services/gen/talk/v1/talkv1connect"
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
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

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
	_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MARIADB_USER"),
		os.Getenv("MARIADB_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("MARIADB_DATABASE")))
	if err != nil {
		log.Fatalf("failed to setup _db: %v", err)
	}
	defer _db.Close()
	_db.SetConnMaxLifetime(time.Minute * 3)
	_db.SetMaxOpenConns(10)
	_db.SetMaxIdleConns(10)
	dbClient := &db.DBClient{DB: _db}

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
	talkServer := &server.TalkServer{
		DB:       dbClient,
		Auth:     authClient,
		Logger:   logger,
		SvClient: &supervisorClient,
	}

	// ハンドラの準備
	mux := http.NewServeMux()
	interceptors := connect.WithInterceptors(
		interceptor.NewAuthInterceptor(authClient, interceptor.AuthPolicy{}))
	mux.Handle(talkv1connect.NewTalkServiceHandler(
		talkServer,
		interceptors,
	))

	addr := fmt.Sprintf("0.0.0.0:%v", os.Getenv("TALK_SERVICE_PORT"))

	// 起動
	if err := http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatal(err)
	}
}
