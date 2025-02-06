package main

import (
	"fmt"
	"go/link-shorter/configs"
	"go/link-shorter/internal/auth"
	"go/link-shorter/internal/link"
	"go/link-shorter/internal/user"
	"go/link-shorter/pkg/db"
	"go/link-shorter/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	//Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	// Services
	authService := auth.NewAuthService(userRepository)
	//Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
	})

	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}

	fmt.Println("Listening on port 8080")
	server.ListenAndServe()
}
