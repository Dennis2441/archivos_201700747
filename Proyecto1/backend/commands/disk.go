package commands

import (
	"backend/structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func CreateDisk(size int, unit string, fit string, pathv string) string {
	var respuesta string
	fmt.Println("entro creardisk")
	filename := path.Base(pathv)
	dirpath := strings.TrimSuffix(pathv, filename)
	fmt.Println("path" + dirpath)
	fmt.Println("nombre" + filename)

	//size of the disk
	if strings.ToLower(unit) == "k" {
		size = size * 1024
	} else if strings.ToLower(unit) == "m" {
		size = size * 1024 * 1024
	} else {
		fmt.Println("Error: Unit no reconocido")
		respuesta += "Error: Unit no reconocido\n"
		return respuesta
	}

	//create filedirectory

	err := os.MkdirAll(dirpath, 0664)
	if err != nil {
		fmt.Println("Error: No se pudo crear el directorio")
		respuesta += "Error: No se pudo crear el directorio\n"
		return respuesta
	}

	//create file
	archivo, err := os.Create(pathv)
	if err != nil {
		fmt.Println("Error: No se pudo crear el archivo")
		respuesta += "Error: No se pudo crear el archivo\n"
		return respuesta
	}
	defer archivo.Close()

	//write file
	randomNum := rand.Intn(99) + 1
	var disk structs.MBR

	disk.Mbr_tamano = int32(size)
	disk.Mbr_disk_signature = int32(randomNum)
	if fit == "" {
		fit = "F"
	} else if fit == "WF" {
		fit = "W"
	} else if fit == "BF" {
		fit = "B"
	} else if fit == "FF" {
		fit = "F"
	}
	fitAux := []byte(fit)
	disk.Dsk_fit = [1]byte{fitAux[0]}
	fechaA := time.Now()
	fecha := fechaA.Format("2006-01-02 15:04:05")
	copy(disk.Mbr_fecha_creacion[:], fecha)
	disk.Mbr_partition_1.Part_status = [1]byte{'0'}
	disk.Mbr_partition_2.Part_status = [1]byte{'0'}
	disk.Mbr_partition_3.Part_status = [1]byte{'0'}
	disk.Mbr_partition_4.Part_status = [1]byte{'0'}

	disk.Mbr_partition_1.Part_fit = [1]byte{'0'}
	disk.Mbr_partition_2.Part_fit = [1]byte{'0'}
	disk.Mbr_partition_3.Part_fit = [1]byte{'0'}
	disk.Mbr_partition_4.Part_fit = [1]byte{'0'}

	disk.Mbr_partition_1.Part_type = [1]byte{'0'}
	disk.Mbr_partition_2.Part_type = [1]byte{'0'}
	disk.Mbr_partition_3.Part_type = [1]byte{'0'}
	disk.Mbr_partition_4.Part_type = [1]byte{'0'}

	disk.Mbr_partition_1.Part_start = 0
	disk.Mbr_partition_2.Part_start = 0
	disk.Mbr_partition_3.Part_start = 0
	disk.Mbr_partition_4.Part_start = 0

	disk.Mbr_partition_1.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	disk.Mbr_partition_2.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	disk.Mbr_partition_3.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	disk.Mbr_partition_4.Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}

	buffer := new(bytes.Buffer)
	for i := 0; i < 1024; i++ {
		buffer.WriteByte(0)
	}
	var totalBytes int = 0
	for totalBytes < size {
		c, err := archivo.Write(buffer.Bytes())
		if err != nil {
			fmt.Println("Error: No se pudo escribir en el archivo")
			respuesta += "Error: No se pudo escribir en el archivo\n"
			return respuesta
		}
		totalBytes += c
	}
	fmt.Println("Archivo llenado")

	//write mbr in file
	archivo.Seek(0, 0)
	err = binary.Write(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error: No se pudo escribir en el archivo")
		respuesta += "Error: No se pudo escribir en el archivo\n"
		return respuesta
	}
	fmt.Println("Disco " + filename + " creado correctamente")
	respuesta += "Disco " + filename + " creado correctamente\n"

	return respuesta
}

func Fdisk(size int, unit string, fit string, pathValor string, name string, typePart string) string {

	var respuesta string

	//Abir el archivo
	archivo, err := os.OpenFile(pathValor, os.O_RDWR, 0664)
	if err != nil {
		fmt.Print("path:")
		fmt.Println(pathValor)
		fmt.Println("Error: No se pudo abrir el archivo")
		respuesta += "Error: No se pudo abrir el archivo\n"
		return respuesta
	}
	defer archivo.Close()
	// Leer el MBR
	var disk structs.MBR
	archivo.Seek(0, 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error: No se pudo leer el archivo")
		respuesta += "Error: No se pudo leer el archivo\n"
		return respuesta
	}

	Desplazamiento := 1 + binary.Size(structs.MBR{})
	ParticionExtendida := structs.NewPartition()
	indiceParticion := 0
	nombreRepetido := false
	verificarEspacio := false
	fmt.Println("entro a mbr")
	if disk.Mbr_partition_1.Part_size != 0 {
		if disk.Mbr_partition_1.Part_type == [1]byte{'e'} {
			ParticionExtendida = disk.Mbr_partition_1
		}
		if strings.Contains(string(disk.Mbr_partition_1.Part_name[:]), name) {
			nombreRepetido = true
		}
		Desplazamiento += int(disk.Mbr_partition_1.Part_size) + 1
	} else {
		indiceParticion = 1
		verificarEspacio = true
	}
	if disk.Mbr_partition_2.Part_size != 0 {
		if disk.Mbr_partition_2.Part_type == [1]byte{'e'} {
			ParticionExtendida = disk.Mbr_partition_2
		}
		//Pasar el arreglo de bytes a string
		if strings.Contains(string(disk.Mbr_partition_2.Part_name[:]), name) {
			nombreRepetido = true
		}
		Desplazamiento += int(disk.Mbr_partition_2.Part_size) + 1
	} else if !verificarEspacio {
		indiceParticion = 2
		verificarEspacio = true
	}
	if disk.Mbr_partition_3.Part_size != 0 {
		if disk.Mbr_partition_3.Part_type == [1]byte{'e'} {
			ParticionExtendida = disk.Mbr_partition_3
		}
		//Pasar el arreglo de bytes a string
		if strings.Contains(string(disk.Mbr_partition_3.Part_name[:]), name) {
			nombreRepetido = true
		}
		Desplazamiento += int(disk.Mbr_partition_3.Part_size) + 1
	} else if !verificarEspacio {
		indiceParticion = 3
		verificarEspacio = true
	}
	if disk.Mbr_partition_4.Part_size != 0 {
		if disk.Mbr_partition_4.Part_type == [1]byte{'e'} {
			ParticionExtendida = disk.Mbr_partition_4
		}
		//Pasar el arreglo de bytes a string
		if strings.Contains(string(disk.Mbr_partition_4.Part_name[:]), name) {
			nombreRepetido = true
		}
		Desplazamiento += int(disk.Mbr_partition_4.Part_size) + 1
	} else if !verificarEspacio {
		indiceParticion = 4
		verificarEspacio = true
	}
	// Si el indice es 0, no hay espacio para crear la particion y el tipo es diferente a lógica
	if indiceParticion == 0 && typePart != "l" {
		fmt.Println("Error: No hay espacio para crear la particion primaria o extendida")
		respuesta = "Error: No hay espacio para crear la particion primaria o extendida"
		return respuesta
	}
	// Si el nombre de la particion ya existe
	if nombreRepetido {
		fmt.Println("Error: El nombre de la particion ya existe")
		respuesta = "Error: El nombre de la particion ya existe"
		return respuesta
	}
	// Si el tipo es extendida y ya existe una extendida
	if typePart == "e" && ParticionExtendida.Part_type == [1]byte{'e'} {
		fmt.Println("Error: Ya existe una particion extendida")
		respuesta = "Error: Ya existe una particion extendida"
		return respuesta
	}
	// Si es diferente a lógica
	///ver
	if strings.ToLower(typePart) != "l" {
		particionNueva := structs.NewPartition()
		particionNueva.Part_type = [1]byte{typePart[0]}
		particionNueva.Part_fit = [1]byte{fit[0]}
		particionNueva.Part_start = int32(Desplazamiento)

		if unit == "k" {
			size = size * 1024
		} else if unit == "m" {
			size = size * 1024 * 1024
		}

		particionNueva.Part_size = int32(size)
		fmt.Println("Size: ", size)
		copy(particionNueva.Part_name[:], name)
		//Verificar si hay espacio para crear la particion
		if int32(Desplazamiento)+particionNueva.Part_size+1 > disk.Mbr_tamano {
			fmt.Println("Error: No hay espacio para crear la particion")
			respuesta = "Error: No hay espacio para crear la particion"
			return respuesta
		}
		if indiceParticion == 1 {
			disk.Mbr_partition_1 = particionNueva
		} else if indiceParticion == 2 {
			disk.Mbr_partition_2 = particionNueva
		} else if indiceParticion == 3 {
			disk.Mbr_partition_3 = particionNueva
		} else if indiceParticion == 4 {
			disk.Mbr_partition_4 = particionNueva
		}
		archivo.Seek(0, 0)
		binary.Write(archivo, binary.LittleEndian, &disk)
		archivo.Close()
		if strings.ToLower(typePart) == "e" {
			fmt.Println("Se creo la particion extendida " + name)
			respuesta = "Se creo la particion extendida " + name
			return respuesta
		} else {
			fmt.Println("Se creo la particion primaria " + name)
			respuesta = "Se creo la particion " + name
			return respuesta
		}
	} else {
		//Verificar si existe una particion extendida
		if ParticionExtendida.Part_type != [1]byte{'e'} {
			fmt.Println("Error: No existe una particion extendida")
			respuesta = "Error: No existe una particion extendida"
			return respuesta
		}
		ebr := structs.NewEBR()
		Desplazamiento = int(ParticionExtendida.Part_start)
		//Do while
		for {
			archivo.Seek(int64(Desplazamiento), 0)
			binary.Read(archivo, binary.LittleEndian, &ebr)
			if ebr.Part_size != 0 {
				//Comprobar si el nombre de la particion ya existe
				if strings.Contains(string(ebr.Part_name[:]), name) {
					fmt.Println("Error: El nombre de la particion ya existe")
					respuesta = "Error: El nombre de la particion ya existe"
					return respuesta
				}
				Desplazamiento += int(ebr.Part_size) + 1 + binary.Size(structs.EBR{})
			}
			if ebr.Part_size == 0 {
				break
			}
		}
		//Crear la particion logica

		if unit == "k" {
			size = size * 1024
		} else if unit == "m" {
			size = size * 1024 * 1024
		}
		//Verificar si hay espacio para crear la particion
		if int32(Desplazamiento)+int32(size) > ParticionExtendida.Part_start+ParticionExtendida.Part_size {
			fmt.Println("Error: No hay espacio para crear la particion")
			respuesta = "Error: No hay espacio para crear la particion"
			return respuesta
		}
		ebrNueva := structs.NewEBR()
		ebrNueva.Part_fit = [1]byte{fit[0]}
		ebrNueva.Part_start = int32(Desplazamiento) + 1 + int32(binary.Size(structs.EBR{}))
		ebrNueva.Part_size = int32(size)
		ebrNueva.Part_next = int32(Desplazamiento) + 1 + int32(binary.Size(structs.EBR{})) + ebrNueva.Part_size
		copy(ebrNueva.Part_name[:], name)
		archivo.Seek(int64(Desplazamiento), 0)
		binary.Write(archivo, binary.LittleEndian, &ebrNueva)
		archivo.Close()
		fmt.Println("Se creo la particion logica " + name)
		respuesta = "Se creo la particion logica " + name
		return respuesta
	}
}
func Rdisk(pathvalor string) error {
	// Intenta eliminar el archivo
	err := os.Remove(pathvalor)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func MountPartition(pathValor string, name string) string {
	var respuesta string
	//Abir el archivo
	archivo, err := os.OpenFile(pathValor, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo")
		respuesta += "Error: No se pudo abrir el archivo\n"
		return respuesta
	}
	defer archivo.Close()
	//Leer el MBR
	var disk structs.MBR
	archivo.Seek(0, 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error: No se pudo leer el archivo")
		respuesta += "Error: No se pudo leer el archivo\n"
		return respuesta
	}
	//Verificar MBR sea valido
	if disk.Mbr_tamano == 0 {
		fmt.Println("Error: Disco no valido")
		respuesta += "Error: Disco no valido\n"
		return respuesta
	}
	//Verificar si la particion existe
	indiceParticion := 0

	if strings.Contains(string(disk.Mbr_partition_1.Part_name[:]), name) {
		indiceParticion = 1
	} else if strings.Contains(string(disk.Mbr_partition_2.Part_name[:]), name) {
		indiceParticion = 2
	} else if strings.Contains(string(disk.Mbr_partition_3.Part_name[:]), name) {
		indiceParticion = 3
	} else if strings.Contains(string(disk.Mbr_partition_4.Part_name[:]), name) {
		indiceParticion = 4
	}

	if indiceParticion == 0 {
		fmt.Println("Error: La particion no existe")
		respuesta += "Error: La particion no existe\n"
		return respuesta
	}

	for i := 0; i < len(particionesMontadas); i++ {
		if particionesMontadas[i].Path == pathValor && particionesMontadas[i].Name == name {
			fmt.Println("Error: La particion ya esta montada")
			respuesta += "Error: La particion ya esta montada\n"
			return respuesta
		}
	}

	if indiceParticion == 1 {
		particion := disk.Mbr_partition_1

		//Montar la particion
		var particionMontada Mount
		particionMontada.Path = pathValor
		particionMontada.Name = name
		particionMontada.Part_type = particion.Part_type

		/* Últimos dos dígitos del Carnet + Número de Partición + LetraEjemplo: carnet = 202401234
		Id´s = 341A, 341B, 341C, 342A, 343A
		Al montar la partición debe validar:
		● Si es una partición del mismo disco se incrementa en uno el número
		de partición.
		● Si la partición es de otro disco se debe colocar la letra siguiente en el
		orden del abecedario (A,B,C,D, etc) y volver a iniciar en uno el número
		de partición al utilizar la siguiente letra.
		● El número de particiones se obtien */

		contador := 1
		for i := 0; i < len(particionesMontadas); i++ {
			if particionesMontadas[i].Path == pathValor {
				contador++
			}
		}

		pathExiste := false
		for i := 0; i < len(pathsParticiones); i++ {
			if pathsParticiones[i] == pathValor {
				pathExiste = true
				break
			}
		}
		if !pathExiste {
			pathsParticiones = append(pathsParticiones, pathValor)
		}

		alfabeto := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

		letra := alfabeto[len(pathsParticiones)-1]

		particionMontada.Id = "45" + strconv.Itoa(contador) + string(letra)
		particionMontada.Size = particion.Part_size
		particionMontada.Start = particion.Part_start
		particionesMontadas = append(particionesMontadas, particionMontada)
		fmt.Println("Particion montada correctamente " + name + " con id " + particionMontada.Id)
		respuesta += "Particion montada correctamente " + name + " con id " + particionMontada.Id + "\n"

		for i := 0; i < len(particionesMontadas); i++ {
			respuesta += "Id: " + particionesMontadas[i].Id + " Path: " + particionesMontadas[i].Path + " Name: " + particionesMontadas[i].Name + " Part_type: " + string(particionesMontadas[i].Part_type[:]) + " Start: " + strconv.Itoa(int(particionesMontadas[i].Start)) + " Size: " + strconv.Itoa(int(particionesMontadas[i].Size)) + "\n"
		}

		disk.Mbr_partition_1.Part_status = [1]byte{'1'}
		disk.Mbr_partition_1.Part_correlative = int32(contador)
		archivo.Seek(0, 0)
		binary.Write(archivo, binary.LittleEndian, &disk)

		return respuesta
	} else if indiceParticion == 2 {
		particion := disk.Mbr_partition_2

		//Montar la particion
		var particionMontada Mount
		particionMontada.Path = pathValor
		particionMontada.Name = name
		particionMontada.Part_type = particion.Part_type

		/* Últimos dos dígitos del Carnet + Número de Partición + LetraEjemplo: carnet = 202401234
		Id´s = 341A, 341B, 341C, 342A, 343A
		Al montar la partición debe validar:
		● Si es una partición del mismo disco se incrementa en uno el número
		de partición.
		● Si la partición es de otro disco se debe colocar la letra siguiente en el
		orden del abecedario (A,B,C,D, etc) y volver a iniciar en uno el número
		de partición al utilizar la siguiente letra.
		● El número de particiones se obtien */

		contador := 1
		for i := 0; i < len(particionesMontadas); i++ {
			if particionesMontadas[i].Path == pathValor {
				contador++
			}
		}

		pathExiste := false
		for i := 0; i < len(pathsParticiones); i++ {
			if pathsParticiones[i] == pathValor {
				pathExiste = true
				break
			}
		}
		if !pathExiste {
			pathsParticiones = append(pathsParticiones, pathValor)
		}

		alfabeto := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

		letra := alfabeto[len(pathsParticiones)-1]

		particionMontada.Id = "45" + strconv.Itoa(contador) + string(letra)
		particionMontada.Size = particion.Part_size
		particionMontada.Start = particion.Part_start
		particionesMontadas = append(particionesMontadas, particionMontada)
		fmt.Println("Particion montada correctamente " + name + " con id " + particionMontada.Id)
		respuesta += "Particion montada correctamente " + name + " con id " + particionMontada.Id + "\n"

		for i := 0; i < len(particionesMontadas); i++ {
			respuesta += "Id: " + particionesMontadas[i].Id + " Path: " + particionesMontadas[i].Path + " Name: " + particionesMontadas[i].Name + " Part_type: " + string(particionesMontadas[i].Part_type[:]) + " Start: " + strconv.Itoa(int(particionesMontadas[i].Start)) + " Size: " + strconv.Itoa(int(particionesMontadas[i].Size)) + "\n"
		}
		disk.Mbr_partition_2.Part_status = [1]byte{'1'}
		disk.Mbr_partition_2.Part_correlative = int32(contador)
		archivo.Seek(0, 0)
		binary.Write(archivo, binary.LittleEndian, &disk)

		return respuesta
	} else if indiceParticion == 3 {
		particion := disk.Mbr_partition_3

		//Montar la particion
		var particionMontada Mount
		particionMontada.Path = pathValor
		particionMontada.Name = name
		particionMontada.Part_type = particion.Part_type

		/* Últimos dos dígitos del Carnet + Número de Partición + LetraEjemplo: carnet = 202401234
		Id´s = 341A, 341B, 341C, 342A, 343A
		Al montar la partición debe validar:
		● Si es una partición del mismo disco se incrementa en uno el número
		de partición.
		● Si la partición es de otro disco se debe colocar la letra siguiente en el
		orden del abecedario (A,B,C,D, etc) y volver a iniciar en uno el número
		de partición al utilizar la siguiente letra.
		● El número de particiones se obtien */

		contador := 1
		for i := 0; i < len(particionesMontadas); i++ {
			if particionesMontadas[i].Path == pathValor {
				contador++
			}
		}

		pathExiste := false
		for i := 0; i < len(pathsParticiones); i++ {
			if pathsParticiones[i] == pathValor {
				pathExiste = true
				break
			}
		}
		if !pathExiste {
			pathsParticiones = append(pathsParticiones, pathValor)
		}

		alfabeto := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

		letra := alfabeto[len(pathsParticiones)-1]

		particionMontada.Id = "45" + strconv.Itoa(contador) + string(letra)
		particionMontada.Size = particion.Part_size
		particionMontada.Start = particion.Part_start
		particionesMontadas = append(particionesMontadas, particionMontada)
		fmt.Println("Particion montada correctamente " + name + " con id " + particionMontada.Id)
		respuesta += "Particion montada correctamente " + name + " con id " + particionMontada.Id + "\n"

		for i := 0; i < len(particionesMontadas); i++ {
			respuesta += "Id: " + particionesMontadas[i].Id + " Path: " + particionesMontadas[i].Path + " Name: " + particionesMontadas[i].Name + " Part_type: " + string(particionesMontadas[i].Part_type[:]) + " Start: " + strconv.Itoa(int(particionesMontadas[i].Start)) + " Size: " + strconv.Itoa(int(particionesMontadas[i].Size)) + "\n"
		}
		disk.Mbr_partition_3.Part_status = [1]byte{'1'}
		disk.Mbr_partition_3.Part_correlative = int32(contador)
		archivo.Seek(0, 0)
		binary.Write(archivo, binary.LittleEndian, &disk)
		return respuesta
	} else if indiceParticion == 4 {
		particion := disk.Mbr_partition_4

		//Montar la particion
		var particionMontada Mount
		particionMontada.Path = pathValor
		particionMontada.Name = name
		particionMontada.Part_type = particion.Part_type

		/* Últimos dos dígitos del Carnet + Número de Partición + LetraEjemplo: carnet = 202401234
		Id´s = 341A, 341B, 341C, 342A, 343A
		Al montar la partición debe validar:
		● Si es una partición del mismo disco se incrementa en uno el número
		de partición.
		● Si la partición es de otro disco se debe colocar la letra siguiente en el
		orden del abecedario (A,B,C,D, etc) y volver a iniciar en uno el número
		de partición al utilizar la siguiente letra.
		● El número de particiones se obtien */

		contador := 1
		for i := 0; i < len(particionesMontadas); i++ {
			if particionesMontadas[i].Path == pathValor {
				contador++
			}
		}

		pathExiste := false
		for i := 0; i < len(pathsParticiones); i++ {
			if pathsParticiones[i] == pathValor {
				pathExiste = true
				break
			}
		}
		if !pathExiste {
			pathsParticiones = append(pathsParticiones, pathValor)
		}

		alfabeto := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

		letra := alfabeto[len(pathsParticiones)-1]

		particionMontada.Id = "45" + strconv.Itoa(contador) + string(letra)
		particionMontada.Size = particion.Part_size
		particionMontada.Start = particion.Part_start
		particionesMontadas = append(particionesMontadas, particionMontada)
		fmt.Println("Particion montada correctamente " + name + " con id " + particionMontada.Id)
		respuesta += "Particion montada correctamente " + name + " con id " + particionMontada.Id + "\n"

		for i := 0; i < len(particionesMontadas); i++ {
			respuesta += "Id: " + particionesMontadas[i].Id + " Path: " + particionesMontadas[i].Path + " Name: " + particionesMontadas[i].Name + " Part_type: " + string(particionesMontadas[i].Part_type[:]) + " Start: " + strconv.Itoa(int(particionesMontadas[i].Start)) + " Size: " + strconv.Itoa(int(particionesMontadas[i].Size)) + "\n"
		}
		disk.Mbr_partition_4.Part_status = [1]byte{'1'}
		disk.Mbr_partition_4.Part_correlative = int32(contador)
		archivo.Seek(0, 0)
		binary.Write(archivo, binary.LittleEndian, &disk)
		return respuesta
	}
	return respuesta
}
