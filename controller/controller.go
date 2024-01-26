// controller/controller.go
package controller

import (
	"encoding/json"
	"fmt"
	"guessgame/game"
	"net/http"
	"text/template"
)

type UIController interface {
	UpdateDisplay(word string)
	RegisterLetterInputHandler(handler func(rune))
	DisplayWinMessage()
}

type Controller struct {
	gameState *game.GameState
	ui        UIController // Add this line
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageData struct {
	Word []int
}

func NewController(gameState *game.GameState) *Controller {
	return &Controller{gameState: gameState}
}

func (c *Controller) HandleGuess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {

		data := PageData{
			Word: c.gameState.GetTheWord(),
		}

		tmpl, err := template.ParseFiles("./static/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	var guess struct {
		Letter string `json:"letter"`
	}

	if err := json.NewDecoder(r.Body).Decode(&guess); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	correct, err := c.gameState.GuessLetter(rune(guess.Letter[0]))

	if c.gameState.AllGuessed() {
		sendJSONResponse(w, http.StatusOK, true, "Processed successfully", "Yay! You have guessed the Name")

	} else {
		sendJSONResponse(w, http.StatusOK, true, "Processed successfully", correct)
	}

	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}
}

func sendJSONResponse(w http.ResponseWriter, status int, success bool, message string, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}

func (c *Controller) OnLetterGuessed(letter rune) {
	correct, err := c.gameState.GuessLetter(letter)
	if err != nil {
		// Handle error (e.g., show in UI)
		return
	}

	fmt.Println(correct)

	if c.gameState.AllGuessed() {
		// ui.DisplayWinMessage()
		c.ui.DisplayWinMessage()
		// Handle win (e.g., show win message in UI)
	} else {
		c.ui.UpdateDisplay(c.gameState.CurrentWordState())
		// Update display with the correct guess
	}
}

func (c *Controller) Start(ui UIController) {
	ui.RegisterLetterInputHandler(c.OnLetterGuessed)
	ui.UpdateDisplay(c.gameState.CurrentWordState())
}
