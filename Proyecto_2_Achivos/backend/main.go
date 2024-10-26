package main

import (
	functions_test "backend/Functions"
	"backend/config"
	"backend/lexer"
	"encoding/json"
	"fmt"
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

type LoginRequest struct {
	IDParticion string `json:"idParticion"`
	Usuario     string `json:"usuario"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	Success bool `json:"success"`
}

const cuenta = false

// Handler para obtener particiones por ruta de disco
func particionesHandler(w http.ResponseWriter, r *http.Request) {
	rutaDisco := r.URL.Query().Get("rutaDisco")
	var particionesFiltradas []functions_test.Particion

	// Filtrar particiones por ruta de disco
	for _, particion := range functions_test.Particiones {
		if particion.RutaDisco == rutaDisco {
			particionesFiltradas = append(particionesFiltradas, particion)
		}
	}

	// Convertir a JSON y enviar respuesta
	jsonData, err := json.Marshal(particionesFiltradas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func discosHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := functions_test.ObtenerDiscosJSON()
	if err != nil {
		http.Error(w, "Error al obtener discos", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(jsonResponse)) // Enviar la respuesta JSON
}

// New handler for logging out
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	functions_test.LOGOUT()
	// Here you can add logic to clear session data or token if applicable
	functions_test.VERIFICARLOGIN = false // Reset login status

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
func submitHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	config.ErrorMessage = ""
	config.GeneralMessage = ""
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

	if config.ErrorMessage == "" {
		output = config.GeneralMessage
	} else {
		output = config.ErrorMessage
	}

	responseBody := ResponseBody{Output: output}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	functions_test.VERIFICARLOGIN = false
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Aquí puedes añadir la lógica de autenticación real
	// Por simplicidad, asumiendo autenticación siempre exitosa
	functions_test.LOGIN(loginReq.Usuario, loginReq.Password, loginReq.IDParticion)

	response := LoginResponse{Success: functions_test.VERIFICARLOGIN}
	fmt.Println(response.Success)
	fmt.Println(response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	functions_test.CargarDatos()
	functions_test.ListarDiscos()
	functions_test.ListarParticiones()
	mux := http.NewServeMux()
	mux.HandleFunc("/particiones", particionesHandler)
	mux.HandleFunc("/submit", submitHandler)
	mux.HandleFunc("/discos", discosHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)
	http.ListenAndServe(":8080", handler)
}
