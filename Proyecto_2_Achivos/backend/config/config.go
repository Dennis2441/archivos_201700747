package config

// Variable global exportada para almacenar mensajes de error
var ErrorMessage string

// Variable global exportada para almacenar mensajes generales
var GeneralMessage string

// Función para establecer un mensaje de error
func SetErrorMessage(message string) {
	ErrorMessage = ErrorMessage + message + " \n"
}

// Función para establecer un mensaje general
func SetGeneralMessage(message string) {
	GeneralMessage = GeneralMessage + message + " \n"
}
