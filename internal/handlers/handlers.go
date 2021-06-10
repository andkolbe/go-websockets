package handlers

import (
	"log"
	"net/http"

	"github.com/andkolbe/go-websockets/internal/helpers"
)

func Login(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "login.jet.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "register.jet.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func Chat(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "chat.jet.html", nil)
	if err != nil {
		log.Println(err)
	}
}
