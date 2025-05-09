package main

import (
	"1337b04rd/config"
	"1337b04rd/internal/adapters/handler"
	"1337b04rd/internal/adapters/middleware"
	"1337b04rd/internal/adapters/repo/postgresql"
	"1337b04rd/internal/service"
	imageuploader "1337b04rd/internal/service/image_uploader"
	"1337b04rd/pkg/logger"
	"1337b04rd/pkg/utils"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := flag.String("port", "", "Port number")
	help := flag.Bool("help", false, "Show this screen.")

	// Override default usage text
	flag.Usage = utils.PrintUsage

	flag.Parse()

	// Show help and exit
	if *help {
		utils.PrintUsage()
		os.Exit(0)
	}

	// Confid and logger
	cfg := config.LoadConfig()
	MyLogger := logger.GetLoggerObject(cfg.LogFilePath)

	// Init PostgreSQL
	db := utils.InitPostgres()

	// Repositories
	sessionRepo := postgresql.NewPostgresSessionRepo(db, MyLogger)
	postRepo := postgresql.NewPostgresPostRepo(db, MyLogger)
	commentRepo := postgresql.NewPostgresCommentRepo(db, MyLogger)
	uploader := imageuploader.NewLocalUploader(cfg.UploadDir, MyLogger)

	// Services
	sessionService := service.NewSessionServiceImpl(sessionRepo, postRepo, commentRepo, MyLogger)
	postService := service.NewPostServiceImpl(postRepo, commentRepo, db, uploader, MyLogger)
	commentService := service.NewCommentServiceImpl(postRepo, commentRepo, uploader, MyLogger)

	// Handlers
	h := handler.NewHandler(postService, commentService, sessionService, MyLogger)

	// Middleware
	sessionMiddleware := middleware.SessionMiddleware(sessionService)

	// Match requests to corresponding handlers
	mux := http.NewServeMux()

	// Static assets and templates, serves them to client
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Converts h.Catalog(w, r) --> http.Handler
	mux.Handle("/", http.HandlerFunc(h.Catalog))
	mux.Handle("/archive", http.HandlerFunc(h.Archive)) // GET /archive
	mux.Handle("/posts/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.Post(w, r)
		} else if r.Method == http.MethodPost {
			h.SubmitComment(w, r)
		} else {
			utils.LogWarn(MyLogger, "MuxRouter", "invalid method for /posts/", "method", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.Handle("/create", http.HandlerFunc(h.CreatePostForm))  // GET /create
	mux.Handle("/submit-post", http.HandlerFunc(h.SubmitPost)) // POST /posts
	mux.Handle("/error", http.HandlerFunc(h.ErrorPage))        // GET /error

	// If flag is not from CLI, then use environment
	finalPort := *port
	if finalPort == "" {
		finalPort = os.Getenv("PORT")
		if finalPort == "" {
			finalPort = "8080"
		}
	}

	// Apply middlewares: CORS → Session → mux
	handler := middleware.CORSMiddleware()(mux)
	handler = sessionMiddleware(handler)

	server := &http.Server{
		Addr:         ":" + finalPort,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Timers for archivation logic and cleaning expired sessions
	go func() {
		ticker := time.NewTicker(1 * time.Minute) // check every minute
		defer ticker.Stop()

		for range ticker.C {
			ctx := context.Background()
			posts, err := postService.GetAllPosts(ctx, false)
			if err == nil {
				for _, p := range posts {
					_ = postService.ArchivePost(ctx, p.PostID)
				}
			}
			_ = sessionService.DeleteExpiredSessions(ctx)
		}
	}()

	log.Printf("Server running on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
