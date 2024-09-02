package main

import (
	"backend/lexer"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

type RequestBody struct {
	Input string `json:"input"`
}

type ResponseBody struct {
	Output string `json:"output"`
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	output := ""
	output = "Salida procesada para: " + reqBody.Input

	inputLines := strings.Split(reqBody.Input, "\n")

	for _, line := range inputLines {
		lexer.ParseLine(line)
	}

	responseBody := ResponseBody{Output: output}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/submit", submitHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)
	http.ListenAndServe(":8080", handler)
}
