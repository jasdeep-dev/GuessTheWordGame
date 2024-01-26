package main

import (
	"guessgame/controller"
	"guessgame/game"
	"log"
	"net/http"
)

type PageData struct {
	Word string
}

func main() {
	gameState := game.NewGame("continent")
	gameController := controller.NewController(gameState)

	http.HandleFunc("/guess", gameController.HandleGuess)

	http.HandleFunc("/", gameController.HandleGuess)

	//	fs := http.FileServer(http.Dir("./static"))
	//	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
