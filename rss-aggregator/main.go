package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/pcaokhai/rssaggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main()  {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	go startScraping(dbQueries, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options {
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: false,
		MaxAge: 30,
	}))

	v1Router := chi.NewRouter()

	// users
	v1Router.Post("/users", apiCfg.handlerUsersCreate)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	// feeds
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedCreate))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	// feed follows
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet))
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowCreate))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowDelete))

	// posts
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	// testing server
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	router.Mount("/v1", v1Router)

	log.Printf("Server staring on port %v", port)
	server := &http.Server {
		Handler: router,
		Addr: ":" + port,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}