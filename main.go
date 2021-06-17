package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {

	args := os.Args

	var port string
	if len(args) > 1 {
		port = args[1]
	}

	if len(port) < 1 {
		port = "8081"
	}

	controller()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func controller() {
	http.HandleFunc("/one", stepOne)
	http.HandleFunc("/authorise", renderAuth)
	http.HandleFunc("/access_token", accessToken)
	http.HandleFunc("/", redirect)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	stepOne(w, r)
}

func renderAuth(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")
	authorise(w, q, state)
}

func stepOne(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("layout.html", "stepone.html")
	err := parsedTemplate.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}

func accessToken(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().Unix()
	data := map[string]interface{}{
		"timestamp":    timestamp,
		"access_token": token(),
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}

func authorise(w http.ResponseWriter, q, state string) {
	parsedTemplate, _ := template.ParseFiles("layout.html", "authorise.html")
	err := parsedTemplate.ExecuteTemplate(w, "layout", struct {
		State,
		RedirectURL string
	}{
		State:       state,
		RedirectURL: q,
	})
	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}

func render404(w http.ResponseWriter, r *http.Request) {
	sendError("not found", w, r)
}

func sendError(message string, w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("layout.html", "error.html")
	err := parsedTemplate.ExecuteTemplate(w, "layout", struct {
		ErrorMessage string
	}{
		ErrorMessage: message,
	})
	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

func token() string {
	bytes := randBytes(32)
	return string(bytes)
}
