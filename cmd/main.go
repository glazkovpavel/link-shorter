package main

import (
	"fmt"
	"go/link-shorter/configs"
	"go/link-shorter/internal/auth"
	"go/link-shorter/internal/link"
	"go/link-shorter/pkg/db"
	"go/link-shorter/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{Config: conf})

	//Repositories
	linkRepository := link.NewLinkRepository(db)
	//Handler
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
