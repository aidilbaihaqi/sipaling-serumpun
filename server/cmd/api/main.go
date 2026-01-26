package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"serumpun-data-api/internal/cache"
	"serumpun-data-api/internal/db"
	httpx "serumpun-data-api/internal/http"
	"serumpun-data-api/internal/queries"

	"github.com/joho/godotenv"
)

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing env: %s", key)
	}
	return v
}

func main() {
	_ = godotenv.Load()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	ttlSec := 60
	if v := os.Getenv("CACHE_TTL_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			ttlSec = n
		}
	}

	databaseURL := mustEnv("DATABASE_URL")
	workspaceName := mustEnv("WORKSPACE_NAME")
	projectName := mustEnv("PROJECT_NAME")
	kabkotaKey := mustEnv("KABKOTA_KEY")

	ctx := context.Background()
	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	srv := &httpx.Server{
		DB:            pool,
		Cache:         cache.New(time.Duration(ttlSec) * time.Second),
		Queries:       queries.New("queries"),
		WorkspaceName: workspaceName,
		ProjectName:   projectName,
		KabkotaKey:    kabkotaKey,
	}

	h := httpx.NewRouter(srv)

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, h))
}
