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
	"github.com/submaline/services/gen/auth/v1/authv1connect"
	"github.com/submaline/services/gen/supervisor/v1/supervisorv1connect"
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

	supervisorClient := supervisorv1connect.NewSupervisorServiceClient(
		http.DefaultClient,
		fmt.Sprintf("http://%s:%s", "supervisor", os.Getenv("SUPERVISOR_SERVICE_PORT")),
	)

	authServer := &server.AuthServer{
		DB:       dbClient,
		Auth:     authClient,
		Logger:   logger,
		SvClient: &supervisorClient,
	}

	mux := http.NewServeMux()
	interceptors := connect.WithInterceptors(
		interceptor.NewAuthInterceptor(authClient, interceptor.AuthPolicy{"/auth.v1.AuthService/LoginWithEmail": false}),
		interceptor.NewLogInterceptor(logger),
	)
	mux.Handle(authv1connect.NewAuthServiceHandler(
		authServer,
		interceptors,
	))

	addr := fmt.Sprintf("0.0.0.0:%v", os.Getenv("AUTH_SERVICE_PORT"))

	// 起動
	if err := http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatal(err)
	}
	//if err := http.ListenAndServeTLS(
	//	addr,
	//	fmt.Sprintf("/data/ssl_certs/%s/staging/signed.crt", os.Getenv("AUTH_DOMAIN")),
	//	fmt.Sprintf("/data/ssl_certs/%s/staging/domain.key", os.Getenv("AUTH_DOMAIN")),
	//	mux,
	//); err != nil {
	//	log.Fatal(err)
	//}
}
