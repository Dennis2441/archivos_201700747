package lexer

import (
	"backend/commands"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	// COMANDOS
	MKDISK string = "mkdisk"
	FDISK  string = "fdisk"
	// PARAMETROS
	SIZE string = "-size"
	FIT  string = "-fit"
	UNIT string = "-unit"
	PATH string = "-path"
	TYPE string = "-type"
	NAME string = "-name"
)

func ParseLine(str string) string { // SOLO RECIBE UNA LINEA
	// Si es espacio realizar el split
	re := regexp.MustCompile(`[ |\t]+`)
	parts := re.Split(str, -1)
	var respuesta string

	if len(parts) > 0 {
		command := parts[0]
		fmt.Println(command)
		// Comando mkdisk
		if strings.ToLower(command) == MKDISK {
			var size int
			var fit string
			var unit string
			var path string
			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, SIZE) {
						size, _ = strconv.Atoi(params[1])
					}
					if strings.Contains(param, FIT) {
						fit = strings.TrimSpace(params[1])
					}
					if strings.Contains(param, UNIT) {
						unit = strings.TrimSpace(params[1])
					}
					if strings.Contains(param, PATH) {
						path = strings.Trim(params[1], "\"")
					}
				}
			}
			if unit == "" {
				unit = "m"
			}
			if fit == "" {
				fit = "f"
			}
			commands.CmdMkdisk(size, fit, unit, path)
			commands.CreateDisk(size, unit, fit, path)
		} else if strings.ToLower(command) == FDISK {
			fmt.Println("entro fdisklexer")
			var size int
			var fit string
			var unit string
			var path string
			var name string
			var tipo string
			for _, part := range parts {
				params := strings.Split(part, "=")

				fmt.Println(params)
				if len(params) > 0 {
					fmt.Println("dentro param fdkislexer")
					fmt.Println(params)
					param := strings.ToLower(params[0])

					if strings.Contains(param, SIZE) {
						size, _ = strconv.Atoi(params[1])
						fmt.Println(size)
					}
					if strings.Contains(param, FIT) {
						fit = strings.TrimSpace(params[1])

						valorFit := fit
						if valorFit == "" {
							fit = "f"
						} else if strings.ToLower(valorFit) != "bf" && strings.ToLower(valorFit) != "ff" && strings.ToLower(valorFit) != "wf" {
							fmt.Println("Error: Fit no reconocido")
							respuesta += "Error: Fit no reconocido\n"
							return respuesta
						} else {
							if strings.ToLower(valorFit) == "bf" {
								fit = "b"
							} else if strings.ToLower(valorFit) == "ff" {
								fit = "f"
							} else if strings.ToLower(valorFit) == "wf" {
								fit = "w"
							}

						}
					}

					if strings.Contains(param, UNIT) {
						unit = strings.TrimSpace(params[1])

						if unit == "" {
							unit = "k"
						} else if strings.ToLower(unit) != "k" && strings.ToLower(unit) != "m" && strings.ToLower(unit) != "b" {
							fmt.Println("Error: unit no reconocido")
							respuesta += "Error: unit no reconocido\n"
							return respuesta

						} else {

							if strings.ToLower(unit) == "k" {
								unit = "K"
							} else if strings.ToLower(unit) == "b" {
								unit = "B"

							} else if strings.ToLower(unit) == "m" {
								unit = "M"

							}
						}
					}

					fmt.Println("ver")
					if strings.Contains(param, TYPE) {
						tipo = strings.TrimSpace(params[1])

						if tipo == "" {
							tipo = "P"
						} else if strings.ToLower(tipo) != "p" && strings.ToLower(tipo) != "e" && strings.ToLower(tipo) != "l" {
							fmt.Println("Error: tipo no reconocido")
							respuesta += "Error: tipo no reconocido\n"
							return respuesta

						} else {

							if strings.ToLower(tipo) == "p" {
								tipo = "p"

							} else if strings.ToLower(tipo) == "e" {
								tipo = "e"

							} else if strings.ToLower(tipo) == "l" {
								tipo = "l"
							}
						}
					}
					fmt.Println("name" + name)
					if strings.Contains(param, NAME) {
						name = strings.Trim(params[1], "\"")
						fmt.Println("nombre:" + name)
						if name == "" {
							fmt.Println("Error: nombre vacio poner nombre ")
							respuesta += "Error: nombre vacio poner nombre\n"
							return respuesta
						}
					}

					if strings.Contains(param, PATH) {
						path = strings.Trim(params[1], "\"")
						fmt.Println(path)
						fmt.Println("es el path del lexico:")
						fmt.Println(path)
						if path == "" {
							fmt.Println("Error: agregar path ")
							respuesta += "Error: agregar path\n"
							return respuesta
						}
					}
				}

			}
			//commands.CmdFdiskk(size, fit, unit, path, tipo, name)
			if unit == "" {
				unit = "k"

			}
			if fit == "" {
				fit = "f"
			}
			if tipo == "" {
				tipo = "P"
			}
			commands.Fdisk(size, unit, fit, path, name, tipo)
		}
	}
	return respuesta

}
