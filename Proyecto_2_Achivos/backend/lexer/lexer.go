package lexer

import (
	functions_test "backend/Functions"
	"backend/config"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	// COMANDOS
	MKDISK  string = "mkdisk"
	FDISK   string = "fdisk"
	RDISK   string = "rmdisk"
	MOUNT   string = "mount"
	MKFS    string = "mkfs"
	LOGIN   string = "login"
	LOGOUT  string = "logout"
	MKGRP   string = "mkgrp"
	RMGRP   string = "rmgrp"
	MKUSR   string = "mkusr"
	RMUSR   string = "rmusr"
	CHGRP   string = "chgrp"
	REP     string = "rep"
	MKFILE  string = "mkfile"
	CAT     string = "cat"
	UNMOUNT string = "unmount"
	MKDIR   string = "mkdir"

	// PARAMETROS
	SIZE   string = "-size"
	FIT    string = "-fit"
	UNIT   string = "-unit"
	PATH   string = "-path"
	TYPE   string = "-type"
	NAME   string = "-name"
	ID     string = "-id"
	USER   string = "-user"
	PASS   string = "-pass"
	GRP    string = "-grp"
	RUTA   string = "-path_file_ls "
	R      string = "-r"
	COUNT1 string = "-cont"
	DELETE string = "-delete"
	ADD    string = "-add"
)

func ParseLine(str string) string { // SOLO RECIBE UNA LINEA
	// Si es espacio realizar el split
	functions_test.LoadMount()
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
					} // Verificar si el parámetro comienza con '-'
					if strings.HasPrefix(param, "-") {
						// Comparar con los parámetros válidos
						if strings.ToLower(param) == SIZE {
							continue // Ya se ha procesado
						} else if strings.ToLower(param) == FIT {
							continue // Ya se ha procesado
						} else if strings.ToLower(param) == UNIT {
							continue // Ya se ha procesado
						} else if strings.ToLower(param) == PATH {
							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para MKDISK \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}

				}
			}

			if size <= 0 {
				fmt.Println("Error: Size must be greater than 0")
				config.SetErrorMessage("Error: Size must be greater than 0 \n")
				return respuesta
			}
			if unit == "" {
				unit = "m"
			} else if strings.ToLower(unit) != "k" && strings.ToLower(unit) != "m" && strings.ToLower(unit) != "b" {
				fmt.Println("Error: Unit must be (k/m)")
				config.SetGeneralMessage("Error: Unit must be (b/k/m)")
				return respuesta
			}
			if fit == "" {
				fit = "f"
			} else if fit != "BF" && fit != "FF" && fit != "WF" {
				fmt.Println("Error: Fit must be (bf/ff/wf)")
				config.SetErrorMessage("Error: Fit must be (bf/ff/wf)")
				return respuesta
			}

			if fit == "BF" {
				fit = "b"
			}
			if fit == "FF" {
				fit = "f"
			}
			if fit == "WF" {
				fit = "w"
			}
			fmt.Println("size:" + strconv.Itoa(size))

			functions_test.MKDISK(path, size, fit, unit)
			fmt.Println(size, fit, unit, path)

		} else if strings.ToLower(command) == FDISK {
			fmt.Println("entro fdisklexer")
			var size int
			var fit string
			var unit string
			var path string
			var name string
			var tipo string
			var dele string
			var add int
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
					if strings.Contains(param, ADD) {
						add, _ = strconv.Atoi(params[1])
						fmt.Println(size)
					}
					if strings.Contains(param, DELETE) {
						dele = strings.Trim(params[1], "\"")
						fmt.Println("nombre:" + name)
						if name == "" {
							fmt.Println("Error: nombre vacio poner nombre ")
							respuesta += "Error: nombre vacio poner nombre \n"
							config.ErrorMessage = config.ErrorMessage + respuesta
							return respuesta
						}
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
							respuesta += "Error: unit no reconocido \n"
							config.ErrorMessage = config.ErrorMessage + respuesta
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
							respuesta += "Error: tipo no reconocido \n"
							config.ErrorMessage = config.ErrorMessage + respuesta
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
							respuesta += "Error: nombre vacio poner nombre \n"
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
							respuesta += "Error: agregar path \n"
							config.ErrorMessage = config.ErrorMessage + respuesta
							return respuesta
						}
					}

					if strings.HasPrefix(param, "-") {
						// Comparar con los parámetros válidos
						if strings.ToLower(param) == SIZE {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == FIT {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == UNIT {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == DELETE {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == PATH {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == NAME {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == ADD {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == TYPE {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para FDISK \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
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
			if fit == "BF" {
				fit = "b"
			}
			if fit == "FF" {
				fit = "f"
			}
			if fit == "WF" {
				fit = "w"
			}
			functions_test.FDISK(size, path, name, unit, tipo, fit, dele, add)
			fmt.Println(size, unit, fit, path, name, tipo)

		} else if strings.ToLower(command) == RDISK {
			var path string
			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, PATH) {
						path = strings.Trim(params[1], "\"")
					}
					// Verificación de parámetros que comienzan con '-'
					if strings.HasPrefix(param, "-") {
						if strings.ToLower(param) == PATH {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para RDISK \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}
				}
			}
			functions_test.RMDISK(path)
			fmt.Println(path)

			fmt.Println()
			//MOUNT
		} else if strings.ToLower(command) == MOUNT {
			fmt.Println()
			var path string
			var name string
			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, PATH) {
						path = strings.Trim(params[1], "\"")
						if path == "" {
							fmt.Println("Error: path vacio poner  ")
							respuesta += "Error: path vacio \n"
							return respuesta
						}
					}
					if strings.Contains(param, NAME) {
						name = strings.Trim(params[1], "\"")
						fmt.Println("nombre:" + name)
						if name == "" {
							fmt.Println("Error: nombre vacio poner nombre ")
							respuesta += "Error: nombre vacio poner nombre\n"
							return respuesta
						}
					}
					// Verificación de parámetros que comienzan con '-'
					if strings.HasPrefix(param, "-") {
						if strings.ToLower(param) == PATH {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == NAME {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para MOUNT \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}

				}
			}
			fmt.Println("entro")
			fmt.Print(path, name)
			functions_test.MOUNT(path, name)
			//MKDISJ
		} else if strings.ToLower(command) == MKFS {
			var id string
			var tipo string
			var valor string
			valor = "2fs"
			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, ID) {
						id = strings.Trim(params[1], "\"")
						if id == "" {
							fmt.Println("Error: id vacio poner  ")
							respuesta += "Error: id vacio \n"
							return respuesta
						}
					}
					if strings.Contains(param, TYPE) {
						tipo = strings.Trim(params[1], "\"")
						fmt.Println("nombre:" + tipo)
						if tipo == "" {
							fmt.Println("Error: tipo vacio poner  ")
							respuesta += "Error: tipo vacio poner \n"
							return respuesta
						}
					}
					// Verificación de parámetros que comienzan con '-'
					if strings.HasPrefix(param, "-") {
						if strings.ToLower(param) == ID {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == TYPE {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para MKFS \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}
				}
			}
			functions_test.MKFS(id, tipo, valor)
			fmt.Println(tipo, id, valor)

		} else if strings.ToLower(command) == LOGIN {

			var id string
			var user string
			var pass string
			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, USER) {
						user = strings.TrimSpace(params[1])
					}
					if strings.Contains(param, PASS) {
						pass = strings.TrimSpace(params[1])
					}
					if strings.Contains(param, ID) {
						id = strings.Trim(params[1], "\"")
					}
					// Verificación de parámetros que comienzan con '-'
					if strings.HasPrefix(param, "-") {
						if strings.ToLower(param) == USER {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == PASS {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == ID {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para LOGIN \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}
				}
			}

			functions_test.LOGIN(user, pass, id)
		} else if strings.ToLower(command) == LOGOUT {

			functions_test.LOGOUT()
		} else if strings.ToLower(command) == MKGRP {
			var name string

			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, NAME) {
						name = strings.Trim(params[1], "\"")
						fmt.Println("nombre:" + name)
						if name == "" {
							fmt.Println("Error: tipo vacio poner  ")
							respuesta += "Error: tipo vacio poner \n"
							return respuesta
						}
					}
					// Verificación de parámetros que comienzan con '-'
					if strings.HasPrefix(param, "-") {
						if strings.ToLower(param) == NAME {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para MKGRP \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}
				}
			}
			functions_test.MKGRP(name)

		} else if strings.ToLower(command) == RMGRP {
			var name string

			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, NAME) {
						name = strings.Trim(params[1], "\"")
						fmt.Println("nombre:" + name)
						if name == "" {
							fmt.Println("Error: tipo vacio poner  ")
							respuesta += "Error: tipo vacio poner \n"
							return respuesta
						}
					}
					// Verificación de parámetros que comienzan con '-'
					if strings.HasPrefix(param, "-") {
						if strings.ToLower(param) == NAME {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para RMGRP \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}
				}
			}
			functions_test.RMGRP(name)

		} else if strings.ToLower(command) == MKUSR {

			var grp string
			var user string
			var pass string
			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, USER) {
						user = strings.TrimSpace(params[1])
					}
					if strings.Contains(param, PASS) {
						pass = strings.TrimSpace(params[1])
					}
					if strings.Contains(param, GRP) {
						grp = strings.Trim(params[1], "\"")
					}
					// Verificación de parámetros que comienzan con '-'
					if strings.HasPrefix(param, "-") {
						if strings.ToLower(param) == USER {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == PASS {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == GRP {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para MKUSR \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}

				}
			}

			functions_test.MKUSR(user, pass, grp)
		} else if strings.ToLower(command) == RMUSR {
			var name string

			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, USER) {
						name = strings.Trim(params[1], "\"")
						fmt.Println("nombre:" + name)
						if name == "" {
							fmt.Println("Error: tipo vacio poner  ")
							respuesta += "Error: tipo vacio poner \n"
							return respuesta
						}
					}
					// Verificación de parámetros que comienzan con '-'
					if strings.HasPrefix(param, "-") {
						if strings.ToLower(param) == USER {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para RMUSR \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}
				}
			}
			functions_test.RMUSR(name)

		} else if strings.ToLower(command) == CHGRP {

			var grp string
			var user string

			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, USER) {
						user = strings.TrimSpace(params[1])
					}

					if strings.Contains(param, GRP) {
						grp = strings.Trim(params[1], "\"")
					}
					// Verificación de parámetros que comienzan con '-'
					if strings.HasPrefix(param, "-") {
						if strings.ToLower(param) == USER {

							continue // Ya se ha procesado
						} else if strings.ToLower(param) == GRP {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para CHGRP \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}
				}
			}

			functions_test.CHGRP(user, grp)
		} else if strings.ToLower(command) == REP {
			var name string
			var path string
			var id string
			var ruta string
			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, RUTA) {
						ruta = strings.Trim(params[1], "\"")
					}
					if strings.Contains(param, ID) {
						id = strings.Trim(params[1], "\"")
					}
					if strings.Contains(param, NAME) {
						name = strings.Trim(params[1], "\"")
						fmt.Println("nombre:" + name)
						if name == "" {
							fmt.Println("Error: tipo vacio poner  ")
							respuesta += "Error: tipo vacio poner \n"
							return respuesta
						}
					}

					if strings.Contains(param, PATH) {
						path = strings.Trim(params[1], "\"")
					}

					// Verificación de parámetros que comienzan con '-'
					if strings.HasPrefix(param, "-") {
						if strings.ToLower(param) == RUTA {

							continue // Ya se ha procesado
						}
						if strings.ToLower(param) == ID {

							continue // Ya se ha procesado
						}
						if strings.ToLower(param) == NAME {

							continue // Ya se ha procesado
						}
						if strings.ToLower(param) == PATH {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para REP \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}
				}
			}

			if name == "" || path == "" || id == "" {
				if name == "" {
					config.SetErrorMessage("Error: campo vacio nombre")

				} else if path == "" {
					config.SetErrorMessage("Error: campo vacio path")
				} else if id == "" {
					config.SetErrorMessage("Error: Ncampo vacio id")
				}
				return ""
			}

			var id_temp string
			id_temp = functions_test.FindID2(functions_test.MountedDiskList, id)
			if id == id_temp {

			} else {
				config.SetErrorMessage("Error: No existe id")
				return ""

			}

			functions_test.GenerateReports(name, path, id, ruta)

		} else if strings.ToLower(command) == MKFILE {
			var size int
			var path string
			var r bool
			r = false
			var cont string
			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, COUNT1) {
						cont = strings.Trim(params[1], "\"")
					}
					if strings.Contains(param, R) {
						r = true
					}
					if strings.Contains(param, SIZE) {
						size, _ = strconv.Atoi(params[1])
					}

					if strings.Contains(param, PATH) {
						path = strings.Trim(params[1], "\"")
					}
				}
			}
			if path == "" {
				println("Error: path no puede estar vacio")
				config.SetErrorMessage("Error: path no puede estar vacio \n")
				return ""
			}

			if cont != "" {
				// Verificar si la ruta existe
				if _, err := os.Stat(cont); err == nil {
					fmt.Println("--------------------------------------------------------------------------")
					fmt.Println("                        MKFILE: LA RUTA EXISTE                            ")
					fmt.Println("--------------------------------------------------------------------------")
					fmt.Println("La ruta existe en el sistema.")

				} else if os.IsNotExist(err) {
					fmt.Println("Error: La ruta no existe en el sistema.")
					config.SetErrorMessage("Error: La ruta no existe en el sistema.")
					return ""
				} else {
					fmt.Println("Error: No se logro verificar la ruta:", err)
					teml := "Error: No se logro verificar la ruta:" + err.Error()
					config.SetErrorMessage(teml)
					return ""
				}

			}

			if size < 0 {
				println("Error: size negativo")
				config.SetErrorMessage("Error: size negativo")

				return ""
			}
			fmt.Println(size, path, r, cont)
			functions_test.MKFILE(path, r, size, cont)

		} else if strings.ToLower(command) == CAT {

			functions_test.Cat("/home/dennis/user.txt")
		} else if strings.ToLower(command) == UNMOUNT {
			var id *string

			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, ID) {
						*id = strings.Trim(params[1], "\"")
						fmt.Println("nombre:" + *id)
						if *id == "" {
							fmt.Println("Error: tipo vacio poner  ")
							respuesta += "Error: tipo vacio poner \n"
							return respuesta
						}
					}
					// Verificación de parámetros que comienzan con '-'
					if strings.HasPrefix(param, "-") {
						if strings.ToLower(param) == ID {

							continue // Ya se ha procesado
						} else {
							respuesta += fmt.Sprintf("Error: '%s' no es un parámetro válido para UNMOUNT \n", param)
							config.ErrorMessage = config.ErrorMessage + respuesta
							return ""
						}
					}
				}
			}
			functions_test.UNMOUNT_Partition(id)

		} else if strings.ToLower(command) == MKDIR {
			var size *int
			var path *string
			var r *bool
			var cont *string
			for _, part := range parts {
				params := strings.Split(part, "=")
				if len(params) > 0 {
					param := strings.ToLower(params[0])

					if strings.Contains(param, COUNT1) {
						*cont = strings.Trim(params[1], "\"")
					}
					if strings.Contains(param, R) {
						*r = strings.Trim(params[1], "\"") == "true" // Asignar true o false
					}
					if strings.Contains(param, SIZE) {
						*size, _ = strconv.Atoi(params[1])
					}

					if strings.Contains(param, PATH) {
						*path = strings.Trim(params[1], "\"")
					}
				}
			}
			if *path == "" {
				println("Error: path no puede estar vacio")
				config.SetErrorMessage("Error: path no puede estar vacio \n")
				return ""
			}

			if *cont != "" {
				// Verificar si la ruta existe
				if _, err := os.Stat(*cont); err == nil {
					fmt.Println("--------------------------------------------------------------------------")
					fmt.Println("                        MKFILE: LA RUTA EXISTE                            ")
					fmt.Println("--------------------------------------------------------------------------")
					fmt.Println("La ruta existe en el sistema.")

				} else if os.IsNotExist(err) {
					fmt.Println("Error: La ruta no existe en el sistema.")
					config.SetErrorMessage("Error: La ruta no existe en el sistema.")
					return ""
				} else {
					fmt.Println("Error: No se logro verificar la ruta:", err)
					teml := "Error: No se logro verificar la ruta:" + err.Error()
					config.SetErrorMessage(teml)
					return ""
				}

			}

			if *size < 0 {
				println("Error: size negativo")
				config.SetErrorMessage("Error: size negativo")

				return ""
			}
			fmt.Println(size, path, r, cont)
			functions_test.MKDIR(path, r)

		}

	}
	return respuesta

}
