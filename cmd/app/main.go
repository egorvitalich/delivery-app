package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	v1 "DostavkaProject/internal/api/v1"
	"DostavkaProject/internal/api/v1/handlers"
	"DostavkaProject/internal/configs"
	"DostavkaProject/internal/repository"
	"DostavkaProject/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := configs.Load("internal/configs/.env")
	if err != nil {
		log.Fatal("configs:", err)
	}

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatal("db open:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("db ping:", err)
	}
	log.Println("Подключено к БД")

	userRepo := repository.NewUserRepo(db)
	jwtSvc := service.NewJWTService(cfg.JWTSecret, cfg.JWTTTL)
	authSvc := service.NewAuthService(userRepo, jwtSvc)

	authH := handlers.NewAuthHandler(authSvc)
	versionH := handlers.NewVersionHandler(db)

	router := v1.NewRouter(v1.Deps{
		Auth:    authH,
		Version: versionH,
		JWT:     jwtSvc,
	})

	// 4. HTTP-сервер
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("http://localhost%s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
