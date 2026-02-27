package main

import (
	"lively-backend/internal/handler"
	"lively-backend/internal/repository"
	"lively-backend/internal/service"
	"lively-backend/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := repository.GetDSN()
	db := repository.InitDB(dsn)
	hub := websocket.NewHub()
	go hub.Run()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	roomRepo := repository.NewRoomRepository(db)
	deezerService := service.NewDeezerService()
	roomService := service.NewRoomService(roomRepo, hub, deezerService)
	roomHandler := handler.NewRoomHandler(roomService)
	musicHandler := handler.NewMusicHandler(deezerService)

	wsHandler := handler.NewWebSocketHandler(hub, roomRepo)

	router := gin.Default()

	router.Use(handler.CORSMiddleware())

	api := router.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.POST("/rooms/sync", roomHandler.SyncTrack)
		api.GET("/search", musicHandler.Search)
		api.GET("/artist/:id", musicHandler.GetArtist)
		api.GET("/radios", musicHandler.GetRadios)
		api.GET("/track/:id", musicHandler.GetTrack)
		api.GET("/time", musicHandler.GetServerTime)

	}

	router.GET("/ws/:roomID/:userID", wsHandler.HandleWS)

	router.Run(":8080")
}
