package functions_test

import (
	structs "backend/Structs"
	utilities_test "backend/Utilities"
	"backend/config"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Estructura para representar un disco
type Disco struct {
	Nombre string
	Ruta   string
}

// Estructura para representar una partición
type Particion struct {
	Nombre      string
	RutaDisco   string
	IDParticion *string // ID opcional
}

// Lista de discos y particiones
var Discos []Disco
var Particiones []Particion

var fileCounter int = 0
var particionesMontadasListado = "--------------------MOUNT: LISTADO DE PARTICIONES MONTADAS------------------\n"
var contadorMontadas int
var MountedDiskList []DiskMounted
var contadorPorDisco = make(map[string]int)
var contador_particion int

type DiskMounted struct {
	id    string
	PATHH string
}

// Agregar un disco

// Eliminar un disco

// Agregar un disco
func AgregarDisco(ruta string) {
	nombre := filepath.Base(ruta) // Extraer el nombre del disco de la ruta

	// Check if the disk already exists
	for _, disco := range Discos {
		if disco.Nombre == nombre && disco.Ruta == ruta {
			fmt.Printf("El disco ya existe: %s en %s\n", nombre, ruta)
			return
		}
	}

	// Add the disk if it doesn't exist
	disco := Disco{Nombre: nombre, Ruta: ruta}
	Discos = append(Discos, disco)
	fmt.Printf("Disco agregado: %s en %s\n", nombre, ruta)
}

// Eliminar un disco
func EliminarDisco(nombre string) {
	for i, disco := range Discos {
		if disco.Nombre == nombre {
			Discos = append(Discos[:i], Discos[i+1:]...) // Eliminar disco
			fmt.Printf("Disco eliminado: %s\n", nombre)
			return
		}
	}
	fmt.Printf("Disco no encontrado: %s\n", nombre)
}

// Listar discos
func ListarDiscos() {
	fmt.Println("Lista de Discos:")
	for _, disco := range Discos {
		fmt.Printf("Nombre: %s, Ruta: %s\n", disco.Nombre, disco.Ruta)
	}
}

// Guardar datos en un archivo
func GuardarDatos() {
	filename := "datos.txt"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Guardar discos únicos
	discoMap := make(map[string]bool)
	for _, disco := range Discos {
		if _, exists := discoMap[disco.Nombre]; !exists {
			_, err := writer.WriteString(fmt.Sprintf("DISCO|%s|%s\n", disco.Nombre, disco.Ruta))
			if err != nil {
				fmt.Println("Error al escribir en el archivo:", err)
				return
			}
			discoMap[disco.Nombre] = true
		}
	}

	// Guardar particiones únicas por nombre y ruta de disco
	particionMap := make(map[string]map[string]bool)
	for _, particion := range Particiones {
		id := "N/A"
		if particion.IDParticion != nil {
			id = *particion.IDParticion
		}

		// Initialize map for the disk path if it doesn't exist
		if _, exists := particionMap[particion.RutaDisco]; !exists {
			particionMap[particion.RutaDisco] = make(map[string]bool)
		}

		// Check for uniqueness based on partition name and disk path
		if _, exists := particionMap[particion.RutaDisco][particion.Nombre]; !exists {
			_, err := writer.WriteString(fmt.Sprintf("PARTICION|%s|%s|%s\n", particion.Nombre, particion.RutaDisco, id))
			if err != nil {
				fmt.Println("Error al escribir en el archivo:", err)
				return
			}
			particionMap[particion.RutaDisco][particion.Nombre] = true
		}
	}

	writer.Flush()
}

// Cargar datos desde un archivo
func CargarDatos() {
	filename := "datos.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}
		switch parts[0] {
		case "DISCO":
			//nombre := parts[1]
			ruta := parts[2]
			// Agregar disco solo si no existe
			AgregarDisco(ruta)
		case "PARTICION":
			var idParticion string
			if len(parts) == 4 {
				id := parts[3]
				idParticion = id
			}
			// Agregar partición sin verificar duplicados
			AgregarParticion2(parts[1], parts[2], idParticion)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
}

// Agregar una partición
func AgregarParticion(nombre string, rutaDisco string) {
	// Check if the partition already exists
	for _, particion := range Particiones {
		if particion.Nombre == nombre && particion.RutaDisco == rutaDisco {
			fmt.Printf("La partición ya existe: %s en %s\n", nombre, rutaDisco)
			return
		}
	}

	// Add the partition if it doesn't exist
	particion := Particion{Nombre: nombre, RutaDisco: rutaDisco}
	Particiones = append(Particiones, particion)
	fmt.Printf("Partición agregada: %s en %s\n", nombre, rutaDisco)
}
func AgregarParticion2(nombre string, rutaDisco string, id_part string) {
	particion := Particion{Nombre: nombre, RutaDisco: rutaDisco, IDParticion: &id_part}
	Particiones = append(Particiones, particion)
	fmt.Printf("Partición agregada: %s en %s\n", nombre, rutaDisco)
}

// Eliminar una partición
func EliminarParticion(nombre string) {
	for i, particion := range Particiones {
		if particion.Nombre == nombre {
			Particiones = append(Particiones[:i], Particiones[i+1:]...) // Eliminar partición
			fmt.Printf("Partición eliminada: %s\n", nombre)
			return
		}
	}
	fmt.Printf("Partición no encontrada: %s\n", nombre)
}

// Modificar ID de una partición
func ModificarIDParticion(nombre string, nuevoID string, rutaDisco string) {
	for i, particion := range Particiones {
		if particion.Nombre == nombre && particion.RutaDisco == rutaDisco {
			Particiones[i].IDParticion = &nuevoID // Modificar ID
			fmt.Printf("ID de la partición %s modificado a: %s\n", nombre, nuevoID)
			return
		}
	}
	fmt.Printf("Partición no encontrada: %s\n", nombre)
}

// Listar particiones
func ListarParticiones() {
	fmt.Println("Lista de Particiones:")
	for _, particion := range Particiones {
		id := "N/A"
		if particion.IDParticion != nil {
			id = *particion.IDParticion
		}
		fmt.Printf("Nombre: %s, Ruta Disco: %s, ID Partición: %s\n", particion.Nombre, particion.RutaDisco, id)
	}
}

// Verificar si un disco ya existe
func DiscoExiste(nombre string) bool {
	for _, disco := range Discos {
		if disco.Nombre == nombre {
			return true
		}
	}
	return false
}

// Obtener la lista de discos en formato JSON
func ObtenerDiscosJSON() (string, error) {
	jsonData, err := json.Marshal(Discos)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func ObtenerParticionesJSON() (string, error) {
	jsonData, err := json.Marshal(Particiones)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// Letra asignada a cada disco
type MountedPartition struct {
	Path       string
	Partition  int // 1 para la primera partición, 2 para la segunda, etc.
	MountOrder int // El número de orden en que fue montada
	// Identificador que incluye el ID del disco (e.g., 471A)
}

var mountedPartitionsList []MountedPartition

func saveMountedPartitions1() {
	var data []string
	for _, disk := range MountedDiskList {
		data = append(data, disk.id+","+disk.PATHH)
	}
	ioutil.WriteFile("mounted_partitions.txt", []byte(strings.Join(data, "\n")), 0644)
}
func saveMountedPartitions() {
	dataMap := make(map[string]struct{})
	var data []string

	for _, disk := range MountedDiskList {
		entry := disk.id + "," + disk.PATHH
		dataMap[entry] = struct{}{} // Use a map to eliminate duplicates
	}

	// Convert map keys to a slice
	for entry := range dataMap {
		data = append(data, entry)
	}

	ioutil.WriteFile("mounted_partitions.txt", []byte(strings.Join(data, "\n")), 0644)
}
func LoadMount() {
	data, err := ioutil.ReadFile("mounted_partitions.txt")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.Split(line, ",")
				if len(parts) == 2 {
					MountedDiskList = append(MountedDiskList, DiskMounted{id: parts[0], PATHH: parts[1]})
				}
			}
		}
	}

}
func printMountedPartitionsList() {
	if len(mountedPartitionsList) == 0 {
		fmt.Println("No hay particiones montadas.")
		return
	}

	var result []string
	for _, mounted := range mountedPartitionsList {
		result = append(result, mounted.Path) // Puedes cambiar `mounted.Path` por cualquier otro campo que desees imprimir
	}

	// Unir los elementos con comas
	output := strings.Join(result, ", ")
	fmt.Println("Particiones montadas: " + output)
}
func findMountedPartition(path string) bool {
	for _, mounted := range mountedPartitionsList {
		if mounted.Path == path {

			return true
		}
	}
	return false
}
func findMountedPartition2(path string) int {
	for _, mounted := range mountedPartitionsList {
		if mounted.Path == path {
			mounted.Partition = mounted.Partition + 1
			return mounted.Partition
		}
	}
	return 1
}
func findMountedPartition3(path string) int {
	for _, mounted := range mountedPartitionsList {
		if mounted.Path == path {

			return mounted.MountOrder
		}
	}
	return 1
}

// Función para cargar el estado de las particiones montadas
func loadMountedPartitions() {
	data, err := ioutil.ReadFile("mounted_partitions.txt")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.Split(line, ",")
				if len(parts) == 2 {
					MountedDiskList = append(MountedDiskList, DiskMounted{id: parts[0], PATHH: parts[1]})
				}
			}
		}
	}
}

func MKFS(id string, type_ string, fs string) {
	loadMountedPartitions()
	var path string
	path = FindPathByID(MountedDiskList, id)
	file, err := utilities_test.OpenFile(path)
	if err != nil {
		return
	}
	var TempMBR structs.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}
	var index int = -1
	// Print object
	structs.PrintMBR(TempMBR)
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 {
			if strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), id) {
				fmt.Println("Particion encontrada")
				if strings.Contains(string(TempMBR.Mbr_particion[i].Part_status[:]), "1") {
					fmt.Println("Particion montada")
					index = i
				} else {
					fmt.Println("Error: La particion no esta montada")
					return
				}
				break
			}
		}
	}

	if index != -1 {
		structs.PrintPartition(TempMBR.Mbr_particion[index])
	} else {
		fmt.Println("Error: No se encontro la particion")
		return
	}

	numerador := int32(TempMBR.Mbr_particion[index].Part_size - int32(binary.Size(structs.Superblock{})))
	denominador_base := int32(4 + int32(binary.Size(structs.Inode{})) + 3*int32(binary.Size(structs.Fileblock{})))
	var temp int32 = 0
	if fs == "2fs" {
		temp = 0
	} else {
		temp = int32(binary.Size(structs.Journaling{}))
	}
	denominador := denominador_base + temp
	n := int32(numerador / denominador)

	fmt.Println("N:", n)

	// var newMRB Structs.MRB
	var newSuperblock structs.Superblock
	newSuperblock.S_inodes_count = 0
	newSuperblock.S_blocks_count = 0

	newSuperblock.S_free_blocks_count = 3 * n
	newSuperblock.S_free_inodes_count = n

	// Obtener la marca de tiempo actual
	currentTime := time.Now()

	// Formatear la marca de tiempo como una cadena
	timeString := currentTime.Format("2006-01-02 15:04:05")

	// Convertir la cadena a un slice de bytes
	timeBytes := []byte(timeString)

	copy(newSuperblock.S_mtime[:], timeBytes)
	copy(newSuperblock.S_umtime[:], timeBytes)
	newSuperblock.S_mnt_count = 0

	//BlockCounter = 0
	//InodeCounter = 0

	if fs == "2fs" {
		create_ext2(n, TempMBR.Mbr_particion[index], newSuperblock, timeString, file)
	} else {
		create_ext3(n, TempMBR.Mbr_particion[index], newSuperblock, timeString, file)
	}

	// Close bin file
	defer file.Close()

}

func create_ext2(n int32, partition structs.Partition, newSuperblock structs.Superblock, date string, file *os.File) {
	fmt.Println("N:", n)
	fmt.Println("Superblock:", newSuperblock)
	fmt.Println("Date:", date)

	newSuperblock.S_filesystem_type = 2
	newSuperblock.S_bm_inode_start = partition.Part_start + int32(binary.Size(structs.Superblock{}))
	newSuperblock.S_bm_block_start = newSuperblock.S_bm_inode_start + n
	newSuperblock.S_inode_start = newSuperblock.S_bm_block_start + 3*n
	newSuperblock.S_block_start = newSuperblock.S_inode_start + n*int32(binary.Size(structs.Inode{}))
	newSuperblock.S_magic = 0xEF53
	newSuperblock.S_mnt_count = 1
	newSuperblock.S_inode_size = int32(binary.Size(structs.Inode{}))
	newSuperblock.S_block_size = int32(binary.Size(structs.Folderblock{}))

	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1
	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1

	for i := int32(0); i < n; i++ {
		err := utilities_test.WriteObject(file, byte(0), int64(newSuperblock.S_bm_inode_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	for i := int32(0); i < 3*n; i++ {
		err := utilities_test.WriteObject(file, byte(0), int64(newSuperblock.S_bm_block_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var newInode structs.Inode
	for i := int32(0); i < 15; i++ {
		newInode.I_block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		err := utilities_test.WriteObject(file, newInode, int64(newSuperblock.S_inode_start+i*int32(binary.Size(structs.Inode{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var newFileblock structs.Fileblock
	for i := int32(0); i < 3*n; i++ {
		err := utilities_test.WriteObject(file, newFileblock, int64(newSuperblock.S_block_start+i*int32(binary.Size(structs.Fileblock{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	//newSuperblock.S_inodes_count++
	var Inode0 structs.Inode //Inode 0
	Inode0.I_uid = usuario.ID
	Inode0.I_gid = 0
	Inode0.I_size = int32(binary.Size(structs.Inode{}))
	copy(Inode0.I_atime[:], date)
	copy(Inode0.I_ctime[:], date)
	copy(Inode0.I_mtime[:], date)
	Inode0.I_type = '1'
	copy(Inode0.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.I_block[i] = -1
	}

	Inode0.I_block[0] = 0
	// . | 0
	// .. | 0
	// users.txt | 1
	//

	//newSuperblock.S_blocks_count++
	var Folderblock0 structs.Folderblock //Bloque 0 -> carpetas
	copy(Folderblock0.B_content[0].B_name[:], ".")
	Folderblock0.B_content[0].B_inodo = 0
	copy(Folderblock0.B_content[1].B_name[:], "..")
	Folderblock0.B_content[1].B_inodo = 0
	copy(Folderblock0.B_content[2].B_name[:], "users.txt")
	Folderblock0.B_content[2].B_inodo = 1
	Folderblock0.B_content[3].B_inodo = -1

	newSuperblock.S_inodes_count++
	var Inode1 structs.Inode //Inode 1
	Inode1.I_uid = 1
	Inode1.I_gid = 0
	Inode1.I_size = int32(binary.Size(structs.Inode{}))
	copy(Inode1.I_atime[:], date)
	copy(Inode1.I_ctime[:], date)
	copy(Inode1.I_mtime[:], date)
	Inode1.I_type = '1'
	copy(Inode1.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode1.I_block[i] = -1
	}

	Inode1.I_block[0] = 1

	newSuperblock.S_blocks_count++
	data := "1,G,root\n1,U,root,root,123\n"
	var Fileblock1 structs.Fileblock //Bloque 1 -> archivo
	copy(Fileblock1.B_content[:], data)
	CreateUserFile(data)
	// BlockCounter++
	// InodeCounter++
	newSuperblock.S_fist_ino = int32(0)
	newSuperblock.S_first_blo = int32(1)

	// Inodo 0 -> Bloque 0 -> Inodo 1 -> Bloque 1
	// Crear la carpeta raiz /
	// Crear el archivo users.txt "1,G,root\n1,U,root,root,123\n"

	// write superblock
	err := utilities_test.WriteObject(file, newSuperblock, int64(partition.Part_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// write bitmap inodes
	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start+1))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	// write bitmap blocks
	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start+1))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Inode 0:", int64(newSuperblock.S_inode_start))
	fmt.Println("Inode 1:", int64(newSuperblock.S_inode_start+int32(binary.Size(structs.Inode{}))))

	// write inodes
	err = utilities_test.WriteObject(file, Inode0, int64(newSuperblock.S_inode_start)) //Inode 0
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities_test.WriteObject(file, Inode1, int64(newSuperblock.S_inode_start+int32(binary.Size(structs.Inode{})))) //Inode 1
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// write blocks
	err = utilities_test.WriteObject(file, Folderblock0, int64(newSuperblock.S_block_start)) //Bloque 0
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities_test.WriteObject(file, Fileblock1, int64(newSuperblock.S_block_start+int32(binary.Size(structs.Fileblock{})))) //Bloque 1
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("--------------------------------------------------------------------------")
	config.SetGeneralMessage("                         MKFS: FORMATO EXT2 APLICADO                      ")

	fmt.Println("--------------------------------------------------------------------------")

}

func create_ext3(n int32, partition structs.Partition, newSuperblock structs.Superblock, date string, file *os.File) {

	fmt.Println("N:", n)
	fmt.Println("Superblock:", newSuperblock)
	fmt.Println("Date:", date)

	newSuperblock.S_filesystem_type = 3
	newSuperblock.S_bm_inode_start = partition.Part_start + int32(binary.Size(structs.Superblock{}))
	newSuperblock.S_bm_block_start = newSuperblock.S_bm_inode_start + n
	newSuperblock.S_inode_start = newSuperblock.S_bm_block_start + 3*n
	newSuperblock.S_block_start = newSuperblock.S_inode_start + n*int32(binary.Size(structs.Inode{}))
	newSuperblock.S_magic = 0xEF53
	newSuperblock.S_mnt_count = 1
	newSuperblock.S_inode_size = int32(binary.Size(structs.Inode{}))
	newSuperblock.S_block_size = int32(binary.Size(structs.Folderblock{}))

	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1
	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1

	var err error // Declarar la variable err una sola vez

	for i := int32(0); i < n; i++ {
		err = utilities_test.WriteObject(file, byte(0), int64(newSuperblock.S_bm_inode_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	for i := int32(0); i < 3*n; i++ {
		err = utilities_test.WriteObject(file, byte(0), int64(newSuperblock.S_bm_block_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var newInode structs.Inode
	for i := int32(0); i < 15; i++ {
		newInode.I_block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		err = utilities_test.WriteObject(file, newInode, int64(newSuperblock.S_inode_start+i*int32(binary.Size(structs.Inode{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var newFileblock structs.Fileblock
	for i := int32(0); i < 3*n; i++ {
		err = utilities_test.WriteObject(file, newFileblock, int64(newSuperblock.S_block_start+i*int32(binary.Size(structs.Fileblock{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	newSuperblock.S_inodes_count++
	var Inode0 structs.Inode //Inode 0
	Inode0.I_uid = 0
	Inode0.I_gid = 0
	Inode0.I_size = 0
	copy(Inode0.I_atime[:], date)
	copy(Inode0.I_ctime[:], date)
	copy(Inode0.I_mtime[:], date)
	Inode0.I_type = '1'
	copy(Inode0.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.I_block[i] = -1
	}

	Inode0.I_block[0] = 0

	// . | 0
	// .. | 0
	// users.txt | 1
	//

	newSuperblock.S_blocks_count++
	var Folderblock0 structs.Folderblock //Bloque 0 -> carpetas
	copy(Folderblock0.B_content[0].B_name[:], ".")
	Folderblock0.B_content[0].B_inodo = 0
	copy(Folderblock0.B_content[1].B_name[:], "..")
	Folderblock0.B_content[1].B_inodo = 0
	copy(Folderblock0.B_content[2].B_name[:], "users.txt")
	Folderblock0.B_content[2].B_inodo = 1
	Folderblock0.B_content[3].B_inodo = -1

	newSuperblock.S_inodes_count++
	var Inode1 structs.Inode //Inode 1
	Inode1.I_uid = 1
	Inode1.I_gid = 0
	Inode1.I_size = int32(binary.Size(structs.Folderblock{}))
	copy(Inode1.I_atime[:], date)
	copy(Inode1.I_ctime[:], date)
	copy(Inode1.I_mtime[:], date)
	Inode1.I_type = '1'
	copy(Inode1.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode1.I_block[i] = -1
	}

	Inode1.I_block[0] = 1

	newSuperblock.S_blocks_count++
	data := "1,G,root\n1,U,root,root,123\n"
	var Fileblock1 structs.Fileblock //Bloque 1 -> archivo
	copy(Fileblock1.B_content[:], data)

	// Inodo 0 -> Bloque 0 -> Inodo 1 -> Bloque 1
	// Crear la carpeta raiz /
	// Crear el archivo users.txt "1,G,root\n1,U,root,root,123\n"

	// Write Journaling structure
	var journal structs.Journaling
	journal.Size = 50
	journal.Ultimo = -1 // Assuming this should be initialized with -1

	newSuperblock.S_inodes_count = int32(2)
	newSuperblock.S_blocks_count = int32(1)
	newSuperblock.S_fist_ino = int32(0)
	newSuperblock.S_first_blo = int32(1)

	// write superblock
	err = utilities_test.WriteObject(file, newSuperblock, int64(partition.Part_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Writing journal to disk
	err = utilities_test.WriteObject(file, journal, int64(newSuperblock.S_block_start+(3*n)*int32(binary.Size(structs.Fileblock{}))))
	if err != nil {
		fmt.Println("Error writing Journaling to disk:", err)
		return
	}

	// write bitmap inodes
	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start+1))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	// write bitmap blocks
	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start+1))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Inode 0:", int64(newSuperblock.S_inode_start))
	fmt.Println("Inode 1:", int64(newSuperblock.S_inode_start+int32(binary.Size(structs.Inode{}))))

	// write inodes
	err = utilities_test.WriteObject(file, Inode0, int64(newSuperblock.S_inode_start)) //Inode 0
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities_test.WriteObject(file, Inode1, int64(newSuperblock.S_inode_start+int32(binary.Size(structs.Inode{})))) //Inode 1
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// write blocks
	err = utilities_test.WriteObject(file, Folderblock0, int64(newSuperblock.S_block_start)) //Bloque 0
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities_test.WriteObject(file, Fileblock1, int64(newSuperblock.S_block_start+int32(binary.Size(structs.Fileblock{})))) //Bloque 1

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("                         MKFS: FORMATO EXT3 APLICADO                      ")
	config.GeneralMessage = config.GeneralMessage + "MKFS: FORMATO EXT3 APLICADO  \n"
	fmt.Println("--------------------------------------------------------------------------")
}

func MOUNT(path string, name string) {
	var iden string

	file, err := utilities_test.OpenFile(path)
	if err != nil {
		return
	}
	var TempMBR structs.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	encontrada := false

	var compareMBR structs.MBR
	copy(compareMBR.Mbr_particion[0].Part_name[:], name)
	copy(compareMBR.Mbr_particion[0].Part_status[:], "1")
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")
	alfabeto := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	// Contar montajes para asignar ID
	var currentMountOrder int
	if len(mountedPartitionsList) == 0 {
		currentMountOrder = 0
	} else {
		currentMountOrder = len(mountedPartitionsList)

	}

	for i := 0; i < 4; i++ {

		if bytes.Equal(TempMBR.Mbr_particion[i].Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {

			if bytes.Equal(TempMBR.Mbr_particion[i].Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
				println("Error: No es necesario montar la particion extendida")
				config.SetErrorMessage("Error: No es necesario montar la particion extendida \n")
				return
			}

			if bytes.Equal(TempMBR.Mbr_particion[i].Part_status[:], compareMBR.Mbr_particion[0].Part_status[:]) {
				println("Error: La particion ya esta montada")
				return
			}
			letra := alfabeto[currentMountOrder]
			var ID string
			if findMountedPartition(path) == true {
				contador_particion = findMountedPartition2(path)
				currentMountOrder = findMountedPartition3(path)
				letra = alfabeto[currentMountOrder]
				ID = "47" + strconv.Itoa(contador_particion) + string(letra)
				iden = ID

			} else {

				letra = alfabeto[currentMountOrder]
				ID = "47" + strconv.Itoa(1) + string(letra)
				iden = ID
				currentMountOrder = currentMountOrder
				mountedPartitionsList = append(mountedPartitionsList, MountedPartition{
					Path:       path,
					Partition:  1,
					MountOrder: currentMountOrder,
				})

			}

			encontrada = true
			copy(TempMBR.Mbr_particion[i].Part_status[:], "1")

			fmt.Println("id:" + ID)

			copy(TempMBR.Mbr_particion[i].Part_id[:], ID)
			particionesMontadasListado += structs.GetPartition(TempMBR.Mbr_particion[i]) + "\n"
			contadorMontadas++
			break

		}
	}

	//Validar si existe una particion extendida
	var EPartition = false
	var EPartitionStart int
	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
		}
	}

	//?EBR verificacion
	if !encontrada && EPartition {
		// Validar que si no existe una particion extendida no se puede crear una logica
		for i := 0; i < 4; i++ {
			//?EBR verificacion
			var x = 0
			for x < 1 {
				var TempEBR structs.EBR
				if err := utilities_test.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if TempEBR.Part_s != 0 {
					if bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
						if bytes.Equal(TempEBR.Part_mount[:], compareMBR.Mbr_particion[0].Part_status[:]) {
							println("Error: La particion ya esta montada")
							config.SetErrorMessage("Error: La particion ya esta montada")
							return
						}
						copy(TempEBR.Part_mount[:], "1") // Cambia a 1 (montada) es estado de la particion
						encontrada = true
						// Escribir el nuevo EBR en el archivo binario
						if err := utilities_test.WriteObject(file, TempEBR, int64(EPartitionStart)); err != nil {
							return
						}
						particionesMontadasListado += structs.GetEBR(TempEBR) + "\n"
					}
					//structs_test.PrintEBR(TempEBR)
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					x = 1
				}
			}
		}
	}
	if encontrada {

		MountedDiskList = append(MountedDiskList, DiskMounted{
			id:    iden,
			PATHH: path,
		})

		saveMountedPartitions()
		fmt.Println("--------------------------------------------------------------------------")
		config.SetGeneralMessage("                        MOUNT: PARTICION " + name + " MONTADA                       \n")
		ModificarIDParticion(name, iden, path)
		GuardarDatos()
		fmt.Println("--------------------------------------------------------------------------")
		// Overwrite the MBR
		if err := utilities_test.WriteObject(file, TempMBR, 0); err != nil {
			return
		}
		//structs_test.PrintMBR(TempMBR)

	} else {
		println("Error: no se encontro la particion")
		config.SetErrorMessage("Error: no se encontro la particion")
	}
	particionesMontadasListado += "--------------------------------------------------------------------------\n"
	println(particionesMontadasListado)
	printMountedPartitionsList()

}

func MOUNT2(path string, name string) {
	var iden string

	file, err := utilities_test.OpenFile(path)
	if err != nil {
		return
	}
	var TempMBR structs.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	encontrada := false

	var compareMBR structs.MBR
	copy(compareMBR.Mbr_particion[0].Part_name[:], name)
	copy(compareMBR.Mbr_particion[0].Part_status[:], "1")
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")
	alfabeto := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

	for i := 0; i < 4; i++ {
		if bytes.Equal(TempMBR.Mbr_particion[i].Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {

			if bytes.Equal(TempMBR.Mbr_particion[i].Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
				println("Error: No es necesario montar la particion extendida")
				config.SetErrorMessage("Error: No es necesario montar la particion extendida \n")
				return
			}

			if bytes.Equal(TempMBR.Mbr_particion[i].Part_status[:], compareMBR.Mbr_particion[0].Part_status[:]) {
				println("Error: La particion ya esta montada")
				return
			}

			letra := alfabeto[len(path)-1]
			encontrada = true
			copy(TempMBR.Mbr_particion[i].Part_status[:], "1")
			ID := "47" + strconv.Itoa(contadorMontadas) + string(letra)
			fmt.Println("id:" + ID)
			iden = ID

			copy(TempMBR.Mbr_particion[i].Part_id[:], ID)
			particionesMontadasListado += structs.GetPartition(TempMBR.Mbr_particion[i]) + "\n"
			contadorMontadas++
			break

		}
	}

	//Validar si existe una particion extendida
	var EPartition = false
	var EPartitionStart int
	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
		}
	}

	//?EBR verificacion
	if !encontrada && EPartition {
		// Validar que si no existe una particion extendida no se puede crear una logica
		for i := 0; i < 4; i++ {
			//?EBR verificacion
			var x = 0
			for x < 1 {
				var TempEBR structs.EBR
				if err := utilities_test.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if TempEBR.Part_s != 0 {
					if bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
						if bytes.Equal(TempEBR.Part_mount[:], compareMBR.Mbr_particion[0].Part_status[:]) {
							println("Error: La particion ya esta montada")
							config.SetErrorMessage("Error: La particion ya esta montada")
							return
						}
						copy(TempEBR.Part_mount[:], "1") // Cambia a 1 (montada) es estado de la particion
						encontrada = true
						// Escribir el nuevo EBR en el archivo binario
						if err := utilities_test.WriteObject(file, TempEBR, int64(EPartitionStart)); err != nil {
							return
						}
						particionesMontadasListado += structs.GetEBR(TempEBR) + "\n"
					}
					//structs_test.PrintEBR(TempEBR)
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					x = 1
				}
			}
		}
	}
	if encontrada {

		MountedDiskList = append(MountedDiskList, DiskMounted{
			id:    iden,
			PATHH: path,
		})
		saveMountedPartitions()
		fmt.Println("--------------------------------------------------------------------------")
		config.SetGeneralMessage("                        MOUNT: PARTICION " + name + " MONTADA                       \n")
		fmt.Println("--------------------------------------------------------------------------")
		// Overwrite the MBR
		if err := utilities_test.WriteObject(file, TempMBR, 0); err != nil {
			return
		}
		//structs_test.PrintMBR(TempMBR)

	} else {
		println("Error: no se encontro la particion")
		config.SetErrorMessage("Error: no se encontro la particion")
	}
	particionesMontadasListado += "--------------------------------------------------------------------------\n"
	println(particionesMontadasListado)

}
func CreateUserFile(data string) error {
	// Crear o abrir el archivo usuario.txt
	file, err := os.Create("/home/dennis/user.txt")

	if err != nil {
		return fmt.Errorf("error al crear o truncar el archivo: %w", err)
	}
	defer file.Close() // Asegúrate de cerrar el archivo al final

	// Escribir el contenido de data en el archivo
	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %w", err)
	}

	return nil
}
func FDISK2(size int, path string, name string, unit string, type_ string, fit string) {
	if strings.ToLower(unit) == "k" {
		size = size * 1024
	} else if strings.ToLower(unit) == "m" {
		size = size * 1024 * 1024
	}

	file, err := utilities_test.OpenFile(path)
	if err != nil {
		config.ErrorMessage = config.ErrorMessage + err.Error()
		return
	}

	var compareMBR structs.MBR
	copy(compareMBR.Mbr_particion[0].Part_name[:], name)
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")
	var TempMBR structs.MBR

	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}
	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
			config.SetErrorMessage("Error: El nombre de la partición ya está en uso! \n")
			return
		}
	}

	//Validar si existe una particion extendida
	var EPartition = false
	var EPartitionStart int
	var ELimit int32
	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true

			EPartitionStart = int(partition.Part_start)

			if strings.ToLower(type_) == "e" {
				config.SetErrorMessage("EN EL MBR:" + path + " ya existe una particion extendida")

				return
			}
			// println("Tamaño de la particion ", partition.Part_size)
			// println("Tipo de particion ", string(partition.Part_type[:]))
			// println("Start de particion ", partition.Part_start)
			ELimit = partition.Part_start + partition.Part_size
			//println("Fin de particion ", partition.Part_start + partition.Part_size)
			//fmt.Println("¡Existe una particion extendida!")
		}
	}

	var count = 0
	var gap = int32(0)
	// Iterate over the partitions
	for i := 0; i < 4; i++ {

		if TempMBR.Mbr_particion[i].Part_size != 0 {
			count++
			gap = TempMBR.Mbr_particion[i].Part_start + TempMBR.Mbr_particion[i].Part_size
		}
	}
	if count == 4 {
		config.ErrorMessage = config.ErrorMessage + "Ya se excedio el Numero de particiones"

	}
	for i := 0; i < 4; i++ {

		if TempMBR.Mbr_particion[i].Part_size == 0 {
			TempMBR.Mbr_particion[i].Part_size = int32(size)

			if count == 0 {
				TempMBR.Mbr_particion[i].Part_start = int32(binary.Size(TempMBR))
			} else {
				TempMBR.Mbr_particion[i].Part_start = gap
			}

			suma := int32(size) + int32(binary.Size(TempMBR))
			//println("Tamaño del disco:", TempMBR.Mbr_tamano)
			//println("Suma:", suma)
			if suma > TempMBR.Mbr_tamano {
				println("Error: La particion exede el tamaño del disco!")
				config.SetErrorMessage("Error: La particion exede el tamaño del disco!")
				return
			}

			copy(TempMBR.Mbr_particion[i].Part_name[:], name)
			copy(TempMBR.Mbr_particion[i].Part_fit[:], fit)
			copy(TempMBR.Mbr_particion[i].Part_status[:], "0")
			copy(TempMBR.Mbr_particion[i].Part_type[:], type_)
			TempMBR.Mbr_particion[i].Part_correlative = int32(count + 1)
			fmt.Println("--------------------------------------------------------------------------")
			config.SetGeneralMessage("                       FDISK: PARTICION " + type_ + "CREADA                         ")
			fmt.Println("                       FDISK: PARTICION " + type_ + "CREADA                         ")
			fmt.Println("--------------------------------------------------------------------------")
			break
		}
	}

	if EPartition && type_ == "l" {
		//?EBR verificacion
		var x = 0
		for x < 1 {
			var TempEBR structs.EBR
			if err := utilities_test.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
				return
			}

			if TempEBR.Part_s != 0 {
				// Escribir un nuevo EBR en el archivo binario
				var newEBR structs.EBR
				copy(newEBR.Part_mount[:], "0")                                   // Indica si la partición está montada o no
				copy(newEBR.Part_fit[:], fit)                                     // Tipo de ajuste de la partición
				newEBR.Part_start = int32(EPartitionStart) + 1                    // Indica en qué byte del disco inicia la partición
				newEBR.Part_s = TempEBR.Part_s                                    // Contiene el tamaño total de la partición en bytes
				newEBR.Part_next = int32(EPartitionStart) + int32(TempEBR.Part_s) // Byte en el que está el próximo EBR (-1 si no hay siguiente)
				copy(newEBR.Part_name[:], TempEBR.Part_name[:])                   // Nombre de la partición

				// Escribir el nuevo EBR en el archivo binario
				if err := utilities_test.WriteObject(file, newEBR, int64(EPartitionStart)); err != nil {
					return
				}
				EPartitionStart = EPartitionStart + int(TempEBR.Part_s)
				structs.PrintEBR(newEBR)
			} else {
				// Escribir un nuevo EBR en el archivo binario
				var newEBR structs.EBR
				copy(newEBR.Part_mount[:], "0")                // Indica si la partición está montada o no
				copy(newEBR.Part_fit[:], fit)                  // Tipo de ajuste de la partición
				newEBR.Part_start = int32(EPartitionStart) + 1 // Indica en qué byte del disco inicia la partición
				newEBR.Part_s = int32(size)                    // Contiene el tamaño total de la partición en bytes
				newEBR.Part_next = -1                          // Byte en el que está el próximo EBR (-1 si no hay siguiente)
				copy(newEBR.Part_name[:], name)                // Nombre de la partición

				// Escribir el nuevo EBR en el archivo binario
				if err := utilities_test.WriteObject(file, newEBR, int64(EPartitionStart)); err != nil {
					return
				}
				structs.PrintEBR(newEBR)
				suma := newEBR.Part_start + newEBR.Part_s
				if suma > ELimit {
					println("Error: la particion logica supera el tamaño de la particion extendida")
					config.SetErrorMessage("Error: la particion logica supera el tamaño de la particion extendida \n")
					return
				}
				x = 1
			}
		}
		fmt.Println("--------------------------------------------------------------------------")
		config.SetGeneralMessage("                       FDISK: PARTICION " + type_ + "CREADA                         ")
		fmt.Println("--------------------------------------------------------------------------")
		return
	}

	// Overwrite the MBR
	if err := utilities_test.WriteObject(file, TempMBR, 0); err != nil {
		return
	}

	var TempMBR2 structs.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR2, 0); err != nil {
		return
	}

	// Print object
	// fmt.Println(">>>>>DESPUES")
	// structs_test.PrintMBR(TempMBR2)

	// Close bin file
	defer file.Close()

}
func MKDISK1(path string, size int, fit string, unit string) error {
	fmt.Println(strconv.Itoa(size))
	if strings.ToLower(unit) == "k" {
		size = size * 1024
	} else if strings.ToLower(unit) == "m" {
		size = size * 1024 * 1024
	}
	fmt.Println(strconv.Itoa(size))
	fmt.Print(unit)
	err := utilities_test.CreateFile(path)
	if err != nil {
		return err
	}

	// Open bin file
	file, err := utilities_test.OpenFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create buffered writer
	writer := bufio.NewWriter(file)

	// Write 0 binary data to the file using buffer
	zeroBytes := make([]byte, 1024) // 1024 bytes buffer
	for i := 0; i < size; i += len(zeroBytes) {
		remaining := size - i
		if remaining < len(zeroBytes) {
			zeroBytes = make([]byte, remaining)
		}
		_, err := writer.Write(zeroBytes)
		if err != nil {
			return err
		}
	}

	// Flush buffer to ensure all data is written
	if err := writer.Flush(); err != nil {
		return err
	}

	// Obtener la hora actual
	currentTime := time.Now()
	// Formatear la hora actual como una cadena
	timeString := currentTime.Format("2006-01-02 15:04:05")
	//Asignacion de datos al MBR
	var TempMBR structs.MBR
	TempMBR.Mbr_tamano = int32(size)
	copy(TempMBR.Mbr_fecha_creacion[:], []byte(timeString))
	TempMBR.Mbr_dsk_signature = int32(GenerateUniqueID())
	copy(TempMBR.Dsk_fit[:], fit)

	// Write object in bin file
	if err := utilities_test.WriteObject(file, TempMBR, 0); err != nil {
		return err
	}

	// Read object from bin file
	var mbr structs.MBR
	if err := utilities_test.ReadObject(file, &mbr, 0); err != nil {
		return err
	}

	fmt.Println("--------------------------------------------------------------------------")
	config.SetGeneralMessage("               MKDISK:" + path + " DISCO  CREADO CORRECTAMENTE                      \n")
	fmt.Println("--------------------------------------------------------------------------")

	return nil
}
func FDISK(size int, path string, name string, unit string, type_ string, fit string, delete string, add int) {
	if strings.ToLower(unit) == "k" {
		size = size * 1024
	} else if strings.ToLower(unit) == "m" {
		size = size * 1024 * 1024
	}

	file, err := utilities_test.OpenFile(path)
	if err != nil {
		config.ErrorMessage = config.ErrorMessage + err.Error()
		return
	}

	var compareMBR structs.MBR
	copy(compareMBR.Mbr_particion[0].Part_name[:], name)
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")
	var TempMBR structs.MBR

	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}
	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
			config.SetErrorMessage("Error: El nombre de la partición ya está en uso! \n")
			return
		}
	}

	//Validar si existe una particion extendida
	var EPartition = false
	var EPartitionStart int
	var ELimit int32
	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true

			EPartitionStart = int(partition.Part_start)

			if strings.ToLower(type_) == "e" {
				config.SetErrorMessage("EN EL MBR:" + path + " ya existe una particion extendida")

				return
			}
			// println("Tamaño de la particion ", partition.Part_size)
			// println("Tipo de particion ", string(partition.Part_type[:]))
			// println("Start de particion ", partition.Part_start)
			ELimit = partition.Part_start + partition.Part_size
			//println("Fin de particion ", partition.Part_start + partition.Part_size)
			//fmt.Println("¡Existe una particion extendida!")
		}
	}
	if delete == "full" {
		encontrada := false
		// Buscar la partición por nombre y eliminarla
		for i := range TempMBR.Mbr_particion {
			if bytes.Equal(TempMBR.Mbr_particion[i].Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
				//Particiones primarias
				if bytes.Equal(TempMBR.Mbr_particion[i].Part_type[:], compareMBR.Mbr_particion[0].Part_type[:]) {
					TempMBR.Mbr_particion[i].Part_correlative = 0
					copy(TempMBR.Mbr_particion[i].Part_fit[:], "")
					copy(TempMBR.Mbr_particion[i].Part_id[:], "")
					copy(TempMBR.Mbr_particion[i].Part_name[:], "")
					copy(TempMBR.Mbr_particion[i].Part_type[:], "")
					copy(TempMBR.Mbr_particion[i].Part_status[:], "")
					encontrada = true
				}
				//Particiones extendidas
				if bytes.Equal(TempMBR.Mbr_particion[i].Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
					end := TempMBR.Mbr_particion[i].Part_start + TempMBR.Mbr_particion[i].Part_size
					utilities_test.ConvertToZeros(path, int64(TempMBR.Mbr_particion[i].Part_start), int64(end))
					TempMBR.Mbr_particion[i].Part_correlative = 0
					copy(TempMBR.Mbr_particion[i].Part_fit[:], "")
					copy(TempMBR.Mbr_particion[i].Part_id[:], "")
					copy(TempMBR.Mbr_particion[i].Part_name[:], "")
					copy(TempMBR.Mbr_particion[i].Part_type[:], "")
					copy(TempMBR.Mbr_particion[i].Part_status[:], "")
					encontrada = true
				}
				break
			}

		}
		//Particiones logicas
		if !encontrada && EPartition {
			//?EBR verificacion
			var x = 0
			for x < 1 {
				var TempEBR structs.EBR
				if err := utilities_test.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if TempEBR.Part_s != 0 {
					if bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {

						copy(TempEBR.Part_mount[:], "0") // Indica si la partición está montada o no
						copy(TempEBR.Part_fit[:], "")    // Tipo de ajuste de la partición
						TempEBR.Part_s = 0               // Contiene el tamaño total de la partición en bytes
						copy(TempEBR.Part_name[:], "")   // Nombre de la partición
						// Escribir el nuevo EBR en el archivo binario
						if err := utilities_test.WriteObject(file, TempEBR, int64(EPartitionStart)); err != nil {
							return
						}
						encontrada = true
						break
					}
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					x = 1
				}
			}
		}

		if encontrada {
			fmt.Println("--------------------------------------------------------------------------")
			fmt.Printf("                       FDISK: PARTICION %s ELIMINADA                      \n", name)
			config.GeneralMessage = config.GeneralMessage + " FDISK: PARTICION " + name + "ELIMINADA \n"
			fmt.Println("--------------------------------------------------------------------------")
		} else {
			fmt.Println("--------------------------------------------------------------------------")
			fmt.Printf("                    FDISK: NO SE ENCONTRO LA PARTICION %s                 \n", name)
			config.ErrorMessage = config.ErrorMessage + " FDISK: NO SE ENCONTRO LA PARTICION " + name + " \n"
			fmt.Println("--------------------------------------------------------------------------")
		}

		/* -------------------------------------------------------------------------- */
		/*                                     ADD                                    */
		/* -------------------------------------------------------------------------- */

	} else if add != 0 {
		//println("ADD", *add)
		// Añadir o quitar espacio en las particiones
		for i := range TempMBR.Mbr_particion {
			if bytes.Equal(TempMBR.Mbr_particion[i].Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
				// Validar que no queden números negativos en el espacio de las particiones
				if TempMBR.Mbr_particion[i].Part_size+int32(add) < 0 {
					fmt.Println("Error: El espacio de la partición no puede ser negativo")
					config.ErrorMessage = config.ErrorMessage + "Error: El espacio de la partición no puede ser negativo \n"
					return
				}
				// Validar que al añadir no se sobrepase el start de la siguiente partición
				if i < len(TempMBR.Mbr_particion)-1 && TempMBR.Mbr_particion[i+1].Part_start < TempMBR.Mbr_particion[i].Part_start+TempMBR.Mbr_particion[i].Part_size+int32(add) {
					if TempMBR.Mbr_particion[i+1].Part_start != 0 {
						fmt.Println("Error: Al añadir espacio, se sobrepasa el start de la siguiente partición")
						config.ErrorMessage = config.ErrorMessage + "Error: Al añadir espacio, se sobrepasa el start de la siguiente partición \n"
						return
					}
				}
				TempMBR.Mbr_particion[i].Part_size += int32(add)
				if TempMBR.Mbr_particion[i].Part_size > TempMBR.Mbr_tamano {
					println("Error: El tamaño supera el tamaño del disco")
					config.ErrorMessage = config.ErrorMessage + "Error: El tamaño supera el tamaño del disco"
					return
				}
				fmt.Println("--------------------------------------------------------------------------")
				fmt.Printf("                    FDISK: ESPACIO EN %s MODIFICADO                       \n", name)
				config.GeneralMessage = config.GeneralMessage + "  FDISK: ESPACIO EN  MODIFICADO " + name + " \n"
				fmt.Println("--------------------------------------------------------------------------")
				break
			}
		}

		/* -------------------------------------------------------------------------- */
		/*                                   CREATE                                   */
		/* -------------------------------------------------------------------------- */

	} else {
		var count = 0
		var gap = int32(0)
		// Iterate over the partitions
		for i := 0; i < 4; i++ {

			if TempMBR.Mbr_particion[i].Part_size != 0 {
				count++
				gap = TempMBR.Mbr_particion[i].Part_start + TempMBR.Mbr_particion[i].Part_size
			}
		}
		if count == 4 {
			config.ErrorMessage = config.ErrorMessage + "Ya se excedio el Numero de particiones"

		}
		for i := 0; i < 4; i++ {

			if TempMBR.Mbr_particion[i].Part_size == 0 {
				TempMBR.Mbr_particion[i].Part_size = int32(size)

				if count == 0 {
					TempMBR.Mbr_particion[i].Part_start = int32(binary.Size(TempMBR))
				} else {
					TempMBR.Mbr_particion[i].Part_start = gap
				}

				suma := int32(size) + int32(binary.Size(TempMBR))
				//println("Tamaño del disco:", TempMBR.Mbr_tamano)
				//println("Suma:", suma)
				if suma > TempMBR.Mbr_tamano {
					println("Error: La particion exede el tamaño del disco!")
					config.SetErrorMessage("Error: La particion exede el tamaño del disco!")
					return
				}

				copy(TempMBR.Mbr_particion[i].Part_name[:], name)
				copy(TempMBR.Mbr_particion[i].Part_fit[:], fit)
				copy(TempMBR.Mbr_particion[i].Part_status[:], "0")
				copy(TempMBR.Mbr_particion[i].Part_type[:], type_)
				TempMBR.Mbr_particion[i].Part_correlative = int32(count + 1)
				fmt.Println("--------------------------------------------------------------------------")
				config.SetGeneralMessage("                       FDISK: PARTICION " + type_ + "CREADA                         ")
				fmt.Println("                       FDISK: PARTICION " + type_ + "CREADA                         ")
				fmt.Println("--------------------------------------------------------------------------")
				break
			}
		}

		if EPartition && type_ == "l" {
			//?EBR verificacion
			var x = 0
			for x < 1 {
				var TempEBR structs.EBR
				if err := utilities_test.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if TempEBR.Part_s != 0 {
					// Escribir un nuevo EBR en el archivo binario
					var newEBR structs.EBR
					copy(newEBR.Part_mount[:], "0")                                   // Indica si la partición está montada o no
					copy(newEBR.Part_fit[:], fit)                                     // Tipo de ajuste de la partición
					newEBR.Part_start = int32(EPartitionStart) + 1                    // Indica en qué byte del disco inicia la partición
					newEBR.Part_s = TempEBR.Part_s                                    // Contiene el tamaño total de la partición en bytes
					newEBR.Part_next = int32(EPartitionStart) + int32(TempEBR.Part_s) // Byte en el que está el próximo EBR (-1 si no hay siguiente)
					copy(newEBR.Part_name[:], TempEBR.Part_name[:])                   // Nombre de la partición

					// Escribir el nuevo EBR en el archivo binario
					if err := utilities_test.WriteObject(file, newEBR, int64(EPartitionStart)); err != nil {
						return
					}
					EPartitionStart = EPartitionStart + int(TempEBR.Part_s)
					structs.PrintEBR(newEBR)
				} else {
					// Escribir un nuevo EBR en el archivo binario
					var newEBR structs.EBR
					copy(newEBR.Part_mount[:], "0")                // Indica si la partición está montada o no
					copy(newEBR.Part_fit[:], fit)                  // Tipo de ajuste de la partición
					newEBR.Part_start = int32(EPartitionStart) + 1 // Indica en qué byte del disco inicia la partición
					newEBR.Part_s = int32(size)                    // Contiene el tamaño total de la partición en bytes
					newEBR.Part_next = -1                          // Byte en el que está el próximo EBR (-1 si no hay siguiente)
					copy(newEBR.Part_name[:], name)                // Nombre de la partición

					// Escribir el nuevo EBR en el archivo binario
					if err := utilities_test.WriteObject(file, newEBR, int64(EPartitionStart)); err != nil {
						return
					}
					structs.PrintEBR(newEBR)
					suma := newEBR.Part_start + newEBR.Part_s
					if suma > ELimit {
						println("Error: la particion logica supera el tamaño de la particion extendida")
						config.SetErrorMessage("Error: la particion logica supera el tamaño de la particion extendida \n")
						return
					}
					x = 1
				}
			}
			fmt.Println("--------------------------------------------------------------------------")
			config.SetGeneralMessage("                       FDISK: PARTICION " + type_ + "CREADA                         ")
			fmt.Println("--------------------------------------------------------------------------")
			return
		}
	}
	// Overwrite the MBR
	if err := utilities_test.WriteObject(file, TempMBR, 0); err != nil {
		return
	}

	var TempMBR2 structs.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR2, 0); err != nil {
		return
	}

	// Print object
	// fmt.Println(">>>>>DESPUES")
	// structs_test.PrintMBR(TempMBR2)

	// Close bin file
	AgregarParticion(name, path)
	defer file.Close()
	GuardarDatos()
	CargarDatos()

}

func MKDISK(path string, size int, fit string, unit string) error {
	fmt.Println("Size in KB/MB:", size)
	if strings.ToLower(unit) == "k" {
		size = size * 1024
	} else if strings.ToLower(unit) == "m" {
		size = size * 1024 * 1024
	}
	fmt.Println("Size in bytes:", size)
	fmt.Print("Unit:", unit)

	err := utilities_test.CreateFile(path)
	if err != nil {
		return err
	}

	file, err := utilities_test.OpenFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	zeroBytes := make([]byte, 1024) // 1 KB buffer

	for i := 0; i < size; i += len(zeroBytes) {
		remaining := size - i
		if remaining < len(zeroBytes) {
			zeroBytes = make([]byte, remaining) // Last write
		}
		_, err := writer.Write(zeroBytes)
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02 15:04:05")

	var TempMBR structs.MBR
	TempMBR.Mbr_tamano = int32(size)
	copy(TempMBR.Mbr_fecha_creacion[:], []byte(timeString))
	TempMBR.Mbr_dsk_signature = int32(GenerateUniqueID())
	copy(TempMBR.Dsk_fit[:], fit)

	// Ensure file position is at the beginning for writing the MBR
	file.Seek(0, 0)
	if err := utilities_test.WriteObject(file, TempMBR, 0); err != nil {
		return err
	}

	// Optional: Read back the MBR for verification
	var mbr structs.MBR
	file.Seek(0, 0)
	if err := utilities_test.ReadObject(file, &mbr, 0); err != nil {
		return err
	}

	fmt.Println("--------------------------------------------------------------------------")
	config.SetGeneralMessage("               MKDISK:" + path + " DISCO CREADO CORRECTAMENTE                      \n")
	AgregarDisco(path)
	GuardarDatos()
	CargarDatos()
	fmt.Println("--------------------------------------------------------------------------")

	return nil
}

func UNMOUNT_Partition(id *string) {
	loadMountedPartitions()
	var path string
	path = FindPathByID(MountedDiskList, *id)

	correlativo, err := strconv.ParseInt(string((*id)[len(*id)-3]), 10, 32)
	if err != nil {
		fmt.Println("Error: no se logro convertir la cadena a int32:", err)
		config.ErrorMessage = config.ErrorMessage + "Error: no se logro convertir la cadena a int32:" + err.Error() + "\n"
		return
	}
	filepath := path
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR structs.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	var compareMBR structs.MBR
	compareMBR.Mbr_particion[0].Part_correlative = int32(correlativo)

	for i := 0; i < 4; i++ {

		if TempMBR.Mbr_particion[i].Part_correlative == compareMBR.Mbr_particion[0].Part_correlative {
			//println("entro a la igualacion")
			copy(TempMBR.Mbr_particion[i].Part_status[:], "0")
			break
		}
	}

	// Overwrite the MBR
	if err := utilities_test.WriteObject(file, TempMBR, 0); err != nil {
		return
	}
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Printf("          UNMOUNT: SE DESMONTO LA PARTICION CON EL ID %s                  \n", strings.ToUpper(*id))
	config.GeneralMessage = config.GeneralMessage + "UNMOUNT: SE DESMONTO LA PARTICION CON EL ID" + strings.ToUpper(*id) + "\n"
	fmt.Println("--------------------------------------------------------------------------")
	structs.PrintMBR(TempMBR)
}

func RMDISK(filename string) {
	if _, err := os.Stat(filename); err == nil {
		// El archivo existe, intenta eliminarlo

		if err := os.Remove(filename); err != nil {
			fmt.Println("Error al eliminar el archivo:", err)
			return
		}
		fmt.Println("--------------------------------------------------------------------------")
		fmt.Printf("                RMDISK: DISCO %s ELIMINADO CORRECTAMENTE                  \n", strings.ToUpper(filename))
		config.GeneralMessage = config.GeneralMessage + "                RMDISK: DISCO %s ELIMINADO CORRECTAMENTE                  " + strings.ToUpper(filename) + " \n"
		EliminarDisco(filename)
		fmt.Println("--------------------------------------------------------------------------")

	} else if os.IsNotExist(err) {
		// El archivo no existe
		config.SetErrorMessage("Error: Disco No existe")
	} else {
		// Otro error ocurrió
		fmt.Println("Error al verificar la existencia del archivo:", err)
		config.SetErrorMessage("Error: Al verificar la existencia del archivo")
	}

}

func GenerateUniqueID() int {
	// Obtener la marca de tiempo actual
	currentTime := time.Now()
	// Generar un número aleatorio entre 0 y 9999
	randomNumber := rand.Intn(10000)
	// Combinar la marca de tiempo y el número aleatorio para crear un identificador único
	uniqueID := currentTime.UnixNano() * int64(randomNumber) % (1 << 31)
	// Tomar el valor absoluto para asegurarse de que sea positivo
	uniqueID = int64(math.Abs(float64(uniqueID)))
	return int(uniqueID)
}

func FindPathByID(disks []DiskMounted, id string) string {
	for _, disk := range disks {
		if disk.id == id {
			return disk.PATHH
		}
	}
	return ""
}

func FindID2(disks []DiskMounted, id string) string {
	for _, disk := range disks {
		if disk.id == id {
			return disk.id
		}
	}
	return ""
}
