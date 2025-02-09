package main

import (
	"fmt"
	"go/link-shorter/configs"
	"go/link-shorter/internal/auth"
	"go/link-shorter/internal/link"
	"go/link-shorter/internal/stat"
	"go/link-shorter/internal/user"
	"go/link-shorter/pkg/db"
	"go/link-shorter/pkg/event"
	"go/link-shorter/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	//Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)
	// Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})
	//Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		Config:         conf,
		EventBus:       eventBus,
	})

	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}

	go statService.AddClick()
	fmt.Println("Listening on port 8080")
	server.ListenAndServe()
}
