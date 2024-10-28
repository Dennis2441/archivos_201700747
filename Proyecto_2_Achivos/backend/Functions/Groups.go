package functions_test

import (
	"backend/Global"
	structs "backend/Structs"
	utilities_test "backend/Utilities"
	"backend/config"
	"bufio"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	session           = false
	UserFound2        bool
	usuario           = Global.UserInfo{}
	groupCounter      = 1
	userCounter       = 1
	InodeIndex        = int32(1)
	blockIndex        = 0
	searchIndex       = 0
	dire              = ""
	DIRRE             = ""
	ID                = ""
	CrrSuperblock     structs.Superblock
	indexSB           = 0
	verificar_usuario = 0
	VERIFICARLOGIN    = false
)
var Usuario_propietarioList []Usuario_propietario

type Usuario_propietario struct {
	id     string
	nombre string
	PATHH  string
}

func LOGIN(user string, pass string, id string) {

	if session == true {
		config.ErrorMessage = config.ErrorMessage + "Un usuario ya esta logeado"

	}
	loadMountedPartitions()
	ID = id

	path := FindPathByID(MountedDiskList, id)
	dire = path
	DIRRE = path
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := path
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}

	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs.Superblock
	indexSB = index
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	CrrSuperblock = tempSuperblock

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	// fmt.Println("Bitmap de bloques del inodo1")
	// fmt.Println(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs.Fileblock
	blockNum := crrInode.I_block[searchIndex]
	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
	// 	fmt.Println("Error reading Fileblock:", err)
	// 	return
	// }
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}
	//fmt.Println("Fileblock " + fmt.Sprint(searchIndex))
	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	userFound := false
	UserFound2 = false
	for _, line := range lines {
		// Imprimir cada línea
		// fmt.Println(line)
		items := strings.Split(line, ",")
		if len(items) > 3 {
			//fmt.Println("items[2]->" + items[2])
			if user == items[len(items)-2] {

				userFound = true
				UserFound2 = true
				usuario.Nombre = items[len(items)-2]
				identificacion, err := strconv.Atoi(items[0])
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				usuario.ID = int32(identificacion)
				session = true
				break
			}
		}
	}

	if !userFound {
		searchIndex++
		if searchIndex <= blockIndex {
			LOGIN(user, pass, id)
			// VERIFICAR SI EL USUARIO ES ROOT
		} else {
			fmt.Println("Error: no se encontro al usuario")
			config.SetErrorMessage("Error: no se encontro al usuario")
			searchIndex = 0
			return
		}
	} else {
		Global.PrintUser(usuario)
		config.SetGeneralMessage(usuario.Nombre + " Logeado")
		VERIFICARLOGIN = true
		if strings.ToLower(usuario.Nombre) == "root" {
			verificar_usuario = 1
		} else {
			verificar_usuario = 2
		}
		searchIndex = 0
		return
	}
}

func LOGOUT() {
	if session {
		// fmt.Println("--------------------------------------------------------------------------")
		// fmt.Println("                        LOGOUT: SESION CERRADA                            ")
		// config.SetGeneralMessage("                        LOGOUT: SESION CERRADA                            ")
		// fmt.Println("--------------------------------------------------------------------------")
		session = false
		searchIndex = 0
		usuario.Nombre = ""
		usuario.ID = -1
		verificar_usuario = 0
		return
	}
	config.SetErrorMessage("Error: no hay una sesion activa")
	println("Error: no hay una sesion activa")
}

func MKGRP(name string) {
	filepath := dire

	if verificar_usuario != 1 {
		if verificar_usuario == 2 {
			config.SetErrorMessage("Solo usuario Root puede hacer esta función")
		} else {
			config.SetErrorMessage("No hay usuarios logeado")
		}
		return
	}

	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	var tempSuperblock structs.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	var Fileblock structs.Fileblock
	blockNum := crrInode.I_block[blockIndex]
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}

	data := string(Fileblock.B_content[:])
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		items := strings.Split(line, ",")
		if len(items) == 3 && name == items[2] {
			fmt.Println("Error: nombre repetido")
			config.SetErrorMessage("Error: Nombre Repetido")
			return
		}
	}

	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	groupCounter++
	nuevoGrupo := fmt.Sprintf("%d,G,%s\n", groupCounter, name)
	newContent := currentContent + nuevoGrupo

	if len(newContent) > len(Fileblock.B_content) {
		if blockIndex >= int(len(crrInode.I_block)) {
			fmt.Println("Error: no hay más bloques disponibles")
			config.SetErrorMessage("Error: no hay más bloques disponibles")
			return
		}
		blockIndex++
		tempSuperblock.S_blocks_count++

		var NEWFileblock structs.Fileblock
		copy(NEWFileblock.B_content[:], newContent)
		if err := utilities_test.WriteObject(file, &NEWFileblock, int64(tempSuperblock.S_block_start+tempSuperblock.S_blocks_count*int32(binary.Size(structs.Fileblock{})))); err != nil {
			fmt.Println("Error writing new Fileblock:", err)
			return
		}

		crrInode.I_block[blockIndex] = tempSuperblock.S_blocks_count
		if err := utilities_test.WriteObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
			fmt.Println("Error writing Inode to disk:", err)
			return
		}

		if err := utilities_test.WriteObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
			fmt.Println("Error writing superblock:", err)
			return
		}
	} else {
		copy(Fileblock.B_content[:], newContent)
		if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
			fmt.Println("Error writing Fileblock to disk:", err)
			return
		}

		// Verificando y mostrando líneas
		fmt.Println("Contenido actualizado de Fileblock:")
		updatedData := string(Fileblock.B_content[:])
		updatedLines := strings.Split(updatedData, "\n")
		for _, line := range updatedLines {
			fmt.Println(line)
		}

		if err := utilities_test.WriteObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
			fmt.Println("Error writing superblock:", err)
			return
		}
	}

	// Imprimir el contenido completo de `users.txt`.
	fmt.Println("Contenido completo de users.txt:")
	for _, line := range strings.Split(newContent, "\n") {
		if line != "" {
			fmt.Println(line)
		}
	}
}
func RMGRP(name string) {
	filepath := dire

	if verificar_usuario != 1 {
		if verificar_usuario == 2 {
			config.SetErrorMessage("Solo usuario Root puede hacer esta función")
		} else {
			config.SetErrorMessage("No hay usuarios logeado")
		}
		return
	}

	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	var tempSuperblock structs.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	searchIndex := 0
	deleted := false

	// Recorre todos los bloques en busca del grupo
	for searchIndex < len(crrInode.I_block) && crrInode.I_block[searchIndex] != 0 {
		var Fileblock structs.Fileblock
		blockNum := crrInode.I_block[searchIndex]

		if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
			fmt.Println("Error reading Fileblock:", err)
			return
		}

		currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
		lines := strings.Split(currentContent, "\n")

		// Buscar específicamente el nombre del grupo exacto
		for i, line := range lines {
			if strings.HasPrefix(line, "1,G,"+name) || strings.HasPrefix(line, "2,G,"+name) {
				lines[i] = strings.Replace(lines[i], "G,", "0,", 1)
				deleted = true
				break
			}
		}

		if deleted {
			newContent := strings.Join(lines, "\n")
			copy(Fileblock.B_content[:], newContent)

			if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
				fmt.Println("Error writing Fileblock to disk:", err)
				return
			}

			fmt.Println("Group deleted successfully. Updated group list:")
			for _, line := range lines {
				if line != "" {
					fmt.Println(line)
				}
			}
			return
		}

		searchIndex++
	}

	if !deleted {
		fmt.Println("Group not found")
	}
}

func MKUSR(name string, password string, group string) {
	filepath := dire

	if verificar_usuario != 1 {
		if verificar_usuario == 2 {
			config.SetErrorMessage("Solo usuario Root puede hacer esta función")
		} else {
			config.SetErrorMessage("No hay usuarios logeado")
		}
		return
	}

	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	var tempSuperblock structs.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	searchIndex := 0
	userExists := false

	for searchIndex < len(crrInode.I_block) && crrInode.I_block[searchIndex] != 0 {
		var Fileblock structs.Fileblock
		blockNum := crrInode.I_block[searchIndex]

		if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
			fmt.Println("Error reading Fileblock:", err)
			return
		}

		currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
		lines := strings.Split(currentContent, "\n")

		// Verificar si el usuario ya existe y comprobar el grupo
		for _, line := range lines {
			parts := strings.Split(line, ",")
			if len(parts) > 3 && parts[1] == "U" && parts[2] == name && parts[3] == group {
				fmt.Println("Usuario ya existe")
				return
			}
		}

		searchIndex++
	}

	// Si el usuario no existe, lo creamos
	if !userExists {
		for searchIndex = 0; searchIndex < len(crrInode.I_block) && crrInode.I_block[searchIndex] != 0; searchIndex++ {
			var Fileblock structs.Fileblock
			blockNum := crrInode.I_block[searchIndex]

			if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
				fmt.Println("Error reading Fileblock:", err)
				return
			}

			currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
			lines := strings.Split(currentContent, "\n")

			// Agregar nuevo usuario sólo si hay espacio
			for i, line := range lines {
				if line == "" {
					newUserEntry := fmt.Sprintf("1,U,%s,%s,%s", name, group, password) // Ajusta el número de acuerdo a tus necesidades
					lines[i] = newUserEntry
					fmt.Println("Usuario creado:", newUserEntry)

					newContent := strings.Join(lines, "\n")
					copy(Fileblock.B_content[:], newContent)

					if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
						fmt.Println("Error writing Fileblock to disk:", err)
						return
					}

					userExists = true
					break
				}
			}

			if userExists {
				break
			}
		}

		if !userExists {
			fmt.Println("No se pudo crear el usuario, no hay espacio disponible")
		}
	}
}
func MKGRP1(name string) {
	filepath := dire

	if verificar_usuario == 1 {

	} else if verificar_usuario == 2 {
		config.SetErrorMessage("Solo usuario Root puede hacer esta funcion")
		return
	} else {

		config.SetErrorMessage("No hay usuarios logeado")
		return
	}
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs.Fileblock
	blockNum := crrInode.I_block[blockIndex]

	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
	// 	fmt.Println("Error reading Fileblock:", err)
	// 	return
	// }
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}

	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	/* -------------------------------------------------------------------------- */
	/*          ITERAMOS EN CADA LINEA PARA QUE NO HAYAN GRUPOS REPETIDOS         */
	/* -------------------------------------------------------------------------- */
	for _, line := range lines {
		// Imprimir cada línea
		fmt.Println(line)
		items := strings.Split(line, ",")
		if len(items) == 3 {
			if name == items[2] {
				println("Error: nombre repetido")
				config.SetErrorMessage("Error: Nombre Repetido")
				return
			}
		}
	}

	/* -------------------------------------------------------------------------- */
	/*                          PARSEAMOS LA INFORMACION                          */
	/* -------------------------------------------------------------------------- */
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	groupCounter++
	nuevoGrupo := fmt.Sprintf("%d,G,%s\n", groupCounter, name)
	result := strconv.Itoa(groupCounter) + ",G" + name
	config.GeneralMessage = result
	newContent := currentContent + nuevoGrupo

	/* -------------------------------------------------------------------------- */
	/*                 CREAMOS MAS FILEBLOCKS PARA GUARDAR LA INFO                */
	/* -------------------------------------------------------------------------- */
	if len(newContent) > len(Fileblock.B_content) {
		if blockIndex > int(len(crrInode.I_block)) {
			fmt.Println("Error: no hay mas bloques disponibles")
			config.SetErrorMessage("Error: no hay mas bloques disponibles")
			return
		}
		blockIndex++
		//BlockCounter++
		CrrSuperblock.S_blocks_count++

		var NEWFileblock structs.Fileblock
		// if err := utilities_test.WriteObject(file, &NEWFileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
		// 	fmt.Println("Error reading Fileblock:", err)
		// 	return
		// }
		if err := utilities_test.WriteObject(file, &NEWFileblock, int64(CrrSuperblock.S_block_start+CrrSuperblock.S_blocks_count*int32(binary.Size(structs.Fileblock{})))); err != nil {
			fmt.Println("Error reading Fileblock:", err)
			return
		}

		/* -------------------------------------------------------------------------- */
		/*                     ACTUALIZAMOS LOS BLOQUES DEL INODO 1                   */
		/* -------------------------------------------------------------------------- */
		crrInode.I_block[blockIndex] = CrrSuperblock.S_blocks_count
		if err := utilities_test.WriteObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
			fmt.Println("Error writing Inode to disk:", err)
			return
		}
		/* -------------------------------------------------------------------------- */
		/*                         ACTUALIZAMOS EL SUPERBLOQUE                        */
		/* -------------------------------------------------------------------------- */
		if err := utilities_test.WriteObject(file, &CrrSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
			fmt.Println("Error reading superblock:", err)
			return
		}
		MKGRP(name)
	} else {
		/* -------------------------------------------------------------------------- */
		/*                GUARDA LA INFORMACION EN EL FILEBLOCK ACTUAL                */
		/* -------------------------------------------------------------------------- */
		copy(Fileblock.B_content[:], newContent)

		// if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
		// 	fmt.Println("Error writing Fileblock to disk:", err)
		// 	return
		// }
		blockNum := crrInode.I_block[blockIndex]

		if err := utilities_test.WriteObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
			fmt.Println("Error reading Fileblock:", err)
			return
		}

		println("ACTUALIZACION")
		// Mostrar el contenido actualizado del Fileblock
		data := string(Fileblock.B_content[:])
		// Dividir la cadena en líneas
		lines := strings.Split(data, "\n")

		/* -------------------------------------------------------------------------- */
		/*          ITERAMOS EN CADA LINEA PARA QUE NO HAYAN GRUPOS REPETIDOS         */
		/* -------------------------------------------------------------------------- */
		for _, line := range lines {
			// Imprimir cada línea
			fmt.Println(line)
			config.SetGeneralMessage(line)
		}
		/* -------------------------------------------------------------------------- */
		/*                         ACTUALIZAMOS EL SUPERBLOQUE                        */
		/* -------------------------------------------------------------------------- */
		if err := utilities_test.WriteObject(file, &CrrSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
			fmt.Println("Error reading superblock:", err)
			return
		}
	}

}

func MKUSR1(user string, pass string, grp string) {
	if verificar_usuario == 1 {

	} else if verificar_usuario == 2 {
		config.SetErrorMessage("Solo usuario Root puede hacer esta funcion")
		return
	} else {

		config.SetErrorMessage("No hay usuarios logeado")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := dire
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	// fmt.Println("Bitmap de bloques del inodo1")
	// fmt.Println(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	blockNum := crrInode.I_block[blockIndex]
	var Fileblock structs.Fileblock
	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
	// 	fmt.Println("Error reading Fileblock:", err)
	// 	return
	// }
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                          PARSEAMOS LA INFORMACION                          */
	/* -------------------------------------------------------------------------- */
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	groupCounter++
	searchIndex = 0
	var nuevoUsuario = BuscarGrupo(user, pass, grp)
	//fmt.Println("nuevo usuarios: " + nuevoUsuario)
	if nuevoUsuario == "" {
		fmt.Println("Error: No se encontro el grupo")
		return
	}
	newContent := currentContent + nuevoUsuario

	/* -------------------------------------------------------------------------- */
	/*                 CREAMOS MAS FILEBLOCKS PARA GUARDAR LA INFO                */
	/* -------------------------------------------------------------------------- */
	if len(newContent) > len(Fileblock.B_content) {
		if blockIndex > int(len(crrInode.I_block)) {
			fmt.Println("Error: no hay mas bloques disponibles")
			return
		}
		blockIndex++
		//BlockCounter++
		CrrSuperblock.S_blocks_count++

		var NEWFileblock structs.Fileblock
		copy(NEWFileblock.B_content[:], nuevoUsuario)
		// if err := utilities_test.WriteObject(file, &NEWFileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
		// 	fmt.Println("Error reading Fileblock:", err)
		// 	return
		// }

		if err := utilities_test.WriteObject(file, &NEWFileblock, int64(CrrSuperblock.S_block_start+CrrSuperblock.S_blocks_count*int32(binary.Size(structs.Fileblock{})))); err != nil {
			fmt.Println("Error reading Fileblock:", err)
			return
		}
		println("MKUSR EXITOSO")
		config.SetGeneralMessage("MKUSR EXISTOSO")
		// Mostrar el contenido actualizado del Fileblock
		data := string(NEWFileblock.B_content[:])
		// Dividir la cadena en líneas
		lines := strings.Split(data, "\n")

		for _, line := range lines {
			// Imprimir cada línea
			fmt.Println(line)
		}

		/* -------------------------------------------------------------------------- */
		/*                     ACTUALIZAMOS LOS BLOQUES DEL INODO                     */
		/* -------------------------------------------------------------------------- */
		crrInode.I_block[blockIndex] = CrrSuperblock.S_blocks_count

		if err := utilities_test.WriteObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
			fmt.Println("Error writing Inode to disk:", err)
			return
		}
		searchIndex = 0

		/* -------------------------------------------------------------------------- */
		/*                         ACTUALIZAMOS EL SUPERBLOQUE                        */
		/* -------------------------------------------------------------------------- */
		if err := utilities_test.WriteObject(file, &CrrSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
			fmt.Println("Error reading superblock:", err)
			return
		}

	} else {
		config.SetGeneralMessage("MKUSR EXISTOSO")
		println("MKUSR EXITOSO")
		/* -------------------------------------------------------------------------- */
		/*                GUARDA LA INFORMACION EN EL FILEBLOCK ACTUAL                */
		/* -------------------------------------------------------------------------- */
		copy(Fileblock.B_content[:], newContent)

		// if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(blockIndex))); err != nil {
		// 	fmt.Println("Error writing Fileblock to disk:", err)
		// 	return
		// }

		blockNum := crrInode.I_block[blockIndex]

		if err := utilities_test.WriteObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
			fmt.Println("Error reading Fileblock:", err)
			return
		}

		// Mostrar el contenido actualizado del Fileblock
		data := string(Fileblock.B_content[:])
		// Dividir la cadena en líneas
		lines := strings.Split(data, "\n")

		/* -------------------------------------------------------------------------- */
		/*          ITERAMOS EN CADA LINEA PARA QUE NO HAYAN GRUPOS REPETIDOS         */
		/* -------------------------------------------------------------------------- */
		for _, line := range lines {
			// Imprimir cada línea
			fmt.Println(line)
		}
		searchIndex = 0

		/* -------------------------------------------------------------------------- */
		/*                         ACTUALIZAMOS EL SUPERBLOQUE                        */
		/* -------------------------------------------------------------------------- */
		if err := utilities_test.ReadObject(file, &CrrSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
			fmt.Println("Error reading superblock:", err)
			return
		}
	}
}

func RMUSR(user string) {
	filepath := dire
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	var tempSuperblock structs.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	var Fileblock structs.Fileblock
	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
	// 	fmt.Println("Error reading Fileblock:", err)
	// 	return
	// }
	blockNum := crrInode.I_block[searchIndex]

	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}

	data := string(Fileblock.B_content[:])
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		items := strings.Split(line, ",")
		if len(items) > 3 {
			if user == items[len(items)-2] {
				items[0] = "0" // Setear el ID a 0
				newLine := strings.Join(items, ",")
				copy(Fileblock.B_content[:], []byte(newLine))
				// if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
				// 	fmt.Println("Error writing Fileblock to disk:", err)
				// 	return
				// }
				blockNum := crrInode.I_block[searchIndex]

				if err := utilities_test.WriteObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
					fmt.Println("Error reading Fileblock:", err)
					return
				}
				println("RMUSR " + user + " exitoso")
				config.SetGeneralMessage("RMUSR " + user + " exitoso")
				return
			}
		}
	}

	searchIndex++
	if searchIndex <= blockIndex {
		RMUSR(user)
	} else {
		fmt.Println("User not found")
		config.SetErrorMessage("Usuario no encontrado")
	}
}
func CHGRP(user string, grp string) {
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := dire
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs.Superblock
	indexSB = index
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	CrrSuperblock = tempSuperblock

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	// fmt.Println("Bitmap de bloques del inodo1")
	// fmt.Println(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs.Fileblock
	blockNum := crrInode.I_block[searchIndex]
	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
	// 	fmt.Println("Error reading Fileblock:", err)
	// 	return
	// }
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}
	//fmt.Println("Fileblock " + fmt.Sprint(searchIndex))
	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	userFound := false
	for _, line := range lines {
		// Imprimir cada línea
		//fmt.Println(line)
		items := strings.Split(line, ",")
		if len(items) > 3 {
			//fmt.Println("items[2]->" + items[2])
			if user == items[len(items)-2] {
				//print(items[2])
				items[2] = grp // cambiar el grupo
				newLine := strings.Join(items, ",")
				copy(Fileblock.B_content[:], []byte(newLine))
				// if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
				// 	fmt.Println("Error writing Fileblock to disk:", err)
				// 	return
				// }
				blockNum := crrInode.I_block[searchIndex]

				if err := utilities_test.WriteObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
					fmt.Println("Error reading Fileblock:", err)
					return
				}
				println("RMUSR " + user + " exitoso")
				return
			}
		}
	}

	if !userFound {
		searchIndex++
		if searchIndex <= blockIndex {
			CHGRP(user, grp)
		} else {
			fmt.Println("Error: no se encontro al usuario")
			searchIndex = 0
			return
		}
	} else {
		Global.PrintUser(usuario)
		searchIndex = 0
		return
	}
}

func BuscarGrupo(user string, pass string, grp string) string {
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := dire
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return ""
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return ""
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return ""
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return ""
	}

	/* -------------------------------------------------------------------------- */
	/*                   LEEMOS EL INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return ""
	}

	// fmt.Println("Bitmap de bloques del inodo1")
	// fmt.Println(crrInode.I_block)

	/* -------------------------------------------------------------------------- */
	/*                             LEEMOS EL FILEBLOCK                            */
	/* -------------------------------------------------------------------------- */
	var Fileblock structs.Fileblock
	blockNum := crrInode.I_block[searchIndex]

	// if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{}))*int32(searchIndex))); err != nil {
	// 	fmt.Println("Error reading Fileblock:", err)
	// 	return ""
	// }
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return ""
	}
	//fmt.Println("Fileblock " + fmt.Sprint(searchIndex))
	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	groupFound := false
	var newUserLine string
	for _, line := range lines {
		// Imprimir cada línea
		//fmt.Println(line)
		items := strings.Split(line, ",")
		if len(items) == 3 {
			//fmt.Println("items[2]->" + items[2])
			if grp == items[2] {
				groupFound = true
				newUserLine = fmt.Sprintf("%d,G,%s,%s,%s\n", userCounter, grp, user, pass)
				userCounter++
				break
			}
		}
	}

	if !groupFound {
		searchIndex++
		if searchIndex <= blockIndex {
			return BuscarGrupo(user, pass, grp)
		}
	} else {
		return newUserLine
	}
	return ""
}

func Cat(filename string) error {
	if session == false {
		config.SetErrorMessage("Necesitas logiarte")
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %w", err)
		config.SetErrorMessage("error al abrir el archivo")
	}
	defer file.Close() // Asegúrate de cerrar el archivo al final

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Imprimir cada línea del archivo
		config.SetGeneralMessage(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		config.SetErrorMessage("error al leer el archivo")
		return nil
	}

	return nil
}
func GetDir(path string) string {
	// Trim trailing slashes
	path = strings.TrimRight(path, "/")
	// Find the last occurrence of the separator
	lastSeparator := strings.LastIndex(path, "/")
	if lastSeparator == -1 {
		return "." // No directory, return current directory
	}
	return path[:lastSeparator]
}
func MKFILE1(path string, r bool) {
	if verificar_usuario == 1 {

	} else if verificar_usuario == 2 {
		config.SetErrorMessage("Solo usuario Root puede hacer esta funcion")
		return
	} else {

		config.SetErrorMessage("No hay usuarios logeado")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepaths := dire
	file, err := utilities_test.OpenFile(filepaths)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */

	var Inode0 structs.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	structs.PrintInode(Inode0)

	/* -------------------------------------------------------------------------- */
	/*                           OBTENEMOS LA RUTA PADRE                          */
	/* -------------------------------------------------------------------------- */
	tieneComillas := strings.Split(path, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			path = tieneComillas[0]
		} else {
			path = tieneComillas[1]
		}
	}
	rutaPadre := filepath.Dir(path)
	println("Ruta original")
	println(path)
	println("Ruta padre")
	println(rutaPadre)
	Carpetas := strings.Split(path, "/")
	nuevaCarpeta := Carpetas[len(Carpetas)-1]
	partes := strings.Split(rutaPadre, "/")
	partes = partes[1:]
	//println("Elementos en ruta padre")
	//println(len(partes) - 1)
	carpetaCreada := false
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//println("Bloques del inodo 0:")
	ultimo := 0
	root := false
	padreExiste := false
	for cont, i := range Inode0.I_block {
		if len(partes)-1 == 0 {
			//println("root es true")
			root = true
		}
		if i == -1 {
			ultimo = int(cont - 1)
			break
		}
		//println(i)
		if !root {
			existe := BuscarRuta(partes, i, 0)
			if existe {
				println("Existe la ruta padre")
				padreExiste = true
			}
		}
	}

	if root {
		existe := false
		for _, i := range Inode0.I_block {
			if i == -1 {
				break
			}
			// print("Buscando en el inodo 0 el bloque ")
			// println(i)
			existe = BuscarEspacioEnRoot(nuevaCarpeta, i)
			println("hay espacio")
			println(existe)
			if existe {
				break
			}
		}
		if !existe {
			println("Creando nuevo inodo y bloque")
			// BlockCounter++
			CrrSuperblock.S_blocks_count++
			Inode0.I_block[ultimo+1] = CrrSuperblock.S_blocks_count
			CrearFolderBlock(file, CrrSuperblock.S_blocks_count, nuevaCarpeta)
			println("Actualizando inodo 0")
			if err := utilities_test.WriteObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs.Inode{})))); err != nil {
				fmt.Println("Error reading inode:", err)
				return
			}
			structs.PrintInode(Inode0)

		}
		carpetaCreada = true
	}

	if padreExiste && !carpetaCreada {
		println("Creando carpeta desde padre")
		CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		carpetaCreada = true
	}

	if r && !carpetaCreada {
		if string(Padre.B_name[:]) != "" {
			println("creando a partir de carpetas existentes")
			fmt.Printf("Encontrado -> B_inode: %d B_name: %s\n", Padre.B_inodo, Padre.B_name)
			CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		} else {
			println("Creando todas las carpetas")
			CreandoCamino(0, nuevaCarpeta, file, partes)
		}
		carpetaCreada = true
	}
	if carpetaCreada {
		fmt.Println("--------------------------------------------------------------------------")
		config.SetGeneralMessage("                MKDIR: CARPETA " + nuevaCarpeta + " CREADA CORRECTAMENTE\n")
		fmt.Println("--------------------------------------------------------------------------")
	} else {
		config.SetErrorMessage("Error: No se logro crear la carpeta")

	}
}

func MKDIR2(path *string, r *bool) {
	/* -------------------------------------------------------------------------- */

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepaths := dire
	file, err := utilities_test.OpenFile(filepaths)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */

	var Inode0 structs.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	structs.PrintInode(Inode0)

	/* -------------------------------------------------------------------------- */
	/*                           OBTENEMOS LA RUTA PADRE                          */
	/* -------------------------------------------------------------------------- */
	tieneComillas := strings.Split(*path, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			*path = tieneComillas[0]
		} else {
			*path = tieneComillas[1]
		}
	}
	rutaPadre := filepath.Dir(*path)
	println("Ruta original")
	println(*path)
	println("Ruta padre")
	println(rutaPadre)
	Carpetas := strings.Split(*path, "/")
	tieneArchivo := strings.Split(Carpetas[len(Carpetas)-1], ".")
	if (len(tieneArchivo) - 1) != 0 {
		fmt.Println("Error: para crear archivos debes usar MKFILE")
		return
	}
	nuevaCarpeta := Carpetas[len(Carpetas)-1]
	partes := strings.Split(rutaPadre, "/")
	partes = partes[1:]
	//println("Elementos en ruta padre")
	//println(len(partes) - 1)
	carpetaCreada := false
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//println("Bloques del inodo 0:")
	ultimo := 0
	root := false
	padreExiste := false
	for cont, i := range Inode0.I_block {
		if len(partes)-1 == 0 {
			//println("root es true")
			root = true
		}
		if i == -1 {
			ultimo = int(cont - 1)
			break
		}
		//println(i)
		if !root {
			existe := BuscarRuta(partes, i, 0)
			if existe {
				println("Existe la ruta padre")
				padreExiste = true
			}
		}
	}

	if root {
		existe := false
		for _, i := range Inode0.I_block {
			if i == -1 {
				break
			}
			// print("Buscando en el inodo 0 el bloque ")
			// println(i)
			existe = BuscarEspacioEnRoot(nuevaCarpeta, i)
			println("hay espacio")
			println(existe)
			if existe {
				break
			}
		}
		if !existe {
			println("Creando nuevo inodo y bloque")
			// BlockCounter++
			CrrSuperblock.S_blocks_count++
			Inode0.I_block[ultimo+1] = CrrSuperblock.S_blocks_count
			CrearFolderBlock(file, CrrSuperblock.S_blocks_count, nuevaCarpeta)
			println("Actualizando inodo 0")
			if err := utilities_test.WriteObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs.Inode{})))); err != nil {
				fmt.Println("Error reading inode:", err)
				return
			}
			structs.PrintInode(Inode0)

		}
		carpetaCreada = true
	}

	if padreExiste && !carpetaCreada {
		println("Creando carpeta desde padre")
		CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		carpetaCreada = true
	}

	if *r && !carpetaCreada {
		if string(Padre.B_name[:]) != "" {
			println("creando a partir de carpetas existentes")
			fmt.Printf("Encontrado -> B_inode: %d B_name: %s\n", Padre.B_inodo, Padre.B_name)
			CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		} else {
			println("Creando todas las carpetas")
			CreandoCamino(0, nuevaCarpeta, file, partes)
		}
		carpetaCreada = true
	}
	if carpetaCreada {
		fmt.Println("--------------------------------------------------------------------------")
		fmt.Printf("                MKDIR: CARPETA %s CREADA CORRECTAMENTE\n", nuevaCarpeta)
		fmt.Println("--------------------------------------------------------------------------")
	} else {
		println("Error: No se logro crear la carpeta")
	}
	/* -------------------------------------------------------------------------- */
	/*                         ACTUALIZAMOS EL SUPERBLOQUE                        */
	/* -------------------------------------------------------------------------- */
	if err := utilities_test.WriteObject(file, &CrrSuperblock, int64(TempMBR.Mbr_particion[indexSB].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}
}
func MKFILE2(path string, r bool, size int, cont string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESIÓN EXISTENTE                  */

	/* -------------------------------------------------------------------------- */
	/*                             PROCESAMOS EL PATH                           */
	/* -------------------------------------------------------------------------- */
	tieneComillas := strings.Split(path, "\"")
	if len(tieneComillas) > 1 {
		path = tieneComillas[1] // Elimina comillas si existen
	}

	/* -------------------------------------------------------------------------- */
	/*                          VERIFICAMOS SI EL ARCHIVO EXISTE                 */
	// Aquí deberías implementar la lógica para verificar si el archivo ya existe
	// Por ejemplo, podrías usar una función similar a BuscarRuta para verificar
	// si el archivo ya está presente en el sistema de archivos.

	/* -------------------------------------------------------------------------- */
	/*                          VERIFICAMOS EL TAMAÑO                            */
	if size < 0 {
		fmt.Println("Error: El tamaño no puede ser negativo.")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                    VERIFICAMOS LA EXISTENCIA DE LAS CARPETAS PADRES      */
	// Aquí deberías implementar la lógica para verificar si las carpetas padres existen
	// y crear las carpetas si es necesario.

	/* -------------------------------------------------------------------------- */
	/*                           CREANDO EL ARCHIVO                              */
	/* -------------------------------------------------------------------------- */
	filepath := DIRRE
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error abriendo el archivo del disco:", err)
		return
	}
	defer file.Close()

	// Buscar espacio en la raíz para el archivo
	ruta := strings.Split(path, "/")
	nombreArchivo := ruta[len(ruta)-1] // Obtener el nombre del archivo
	var bloque int32 = 0               // Supongamos que comenzamos en la raíz

	// Busca espacio para el archivo en la ruta especificada
	existe := BuscarEspacio(nombreArchivo, bloque)
	if existe == -1 {
		fmt.Println("Error: No se pudo encontrar espacio en la ruta indicada.")
		return
	}

	// Crear inode y bloque para el archivo
	CrearInodoFileblock(file, existe)

	// Completar la creación del archivo
	if err := utilities_test.WriteObject(file, &CrrSuperblock, 0); err != nil {
		fmt.Println("Error actualizando el superbloque:", err)
		return
	}

	repus := "MKFILE: Archivo " + path + " creado correctamente\n"
	config.GeneralMessage = config.GeneralMessage + repus

	// Actualizamos el superbloque después de crear el archivo
	if err := utilities_test.WriteObject(file, &CrrSuperblock, 0); err != nil {
		fmt.Println("Error actualizando el superbloque:", err)
		config.ErrorMessage = config.ErrorMessage + "Error actualizando el superbloque \n"
		return
	}
}
func REMOVE2(path *string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */
	if verificar_usuario == 1 {

	} else if verificar_usuario == 2 {
		config.SetErrorMessage("Solo usuario Root puede hacer esta funcion")
		return
	} else {

		config.SetErrorMessage("No hay usuarios logeado")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepaths := dire
	file, err := utilities_test.OpenFile(filepaths)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */

	var Inode0 structs.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	structs.PrintInode(Inode0)

	carpetas := strings.Split(*path, "/")
	carpetas = carpetas[1:]
	nuevaCarpeta := carpetas[len(carpetas)-1]
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//println("Bloques del inodo 0:")
	println("eliminando " + nuevaCarpeta)
	deleted := false
	for _, i := range Inode0.I_block {
		if i == -1 {
			break
		}
		//println(i)
		deleted = EliminarRuta(carpetas, i, 0)
		if deleted {
			println(nuevaCarpeta + " eliminado con exito")
			config.GeneralMessage = config.GeneralMessage + " " + nuevaCarpeta + " eliminado con exito \n"
			break
		}
	}
	if !deleted {
		println("No se logro eliminar " + nuevaCarpeta)
		config.ErrorMessage = config.ErrorMessage + " No se logro eliminar " + nuevaCarpeta + " \n"
	}
}
func REMOVE(path string) {
	// Verificar si la ruta está vacía
	if path == "" {
		config.ErrorMessage = "Error: Se requiere una ruta."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Verificar si la ruta existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		config.ErrorMessage = "Error: La ruta especificada no existe."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Verificar si el usuario tiene permiso de escritura en la ruta
	if !hasWritePermission(path, usuario.Nombre) {
		config.ErrorMessage = "Error: No tiene permiso para eliminar la ruta especificada."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Si es un directorio, verificar permisos para todos los contenidos
	if isDirectory(path) {
		if !canDeleteDirectoryContents(path) {
			config.ErrorMessage = "Error: No se pueden eliminar algunos contenidos del directorio debido a problemas de permisos."
			fmt.Println(config.ErrorMessage)
			return
		}

		// Eliminar el directorio y su contenido
		err := os.RemoveAll(path)
		if err != nil {
			config.ErrorMessage = "Error al eliminar el directorio: " + err.Error()
			fmt.Println(config.ErrorMessage)
			return
		}
	} else {
		// Si es un archivo, intentar eliminarlo
		err := os.Remove(path)
		if err != nil {
			config.ErrorMessage = "Error al eliminar el archivo: " + err.Error()
			fmt.Println(config.ErrorMessage)
			return
		}
	}

	// Eliminar permisos del archivo o directorio de permissions.txt
	removePermission(path, usuario.Nombre)

	config.GeneralMessage = "Eliminado con éxito: " + path
	fmt.Println(config.GeneralMessage)
}

// Helper function to check if a path is a directory
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Helper function to check if the user can delete all contents of a directory
func canDeleteDirectoryContents(dirPath string) bool {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Verificar si el usuario tiene permiso de escritura para cada elemento
		if !hasWritePermission(path, usuario.Nombre) {
			return fmt.Errorf("sin permiso para %s", path)
		}
		return nil
	})

	return err == nil
}

func removePermission(path string, username string) {
	file, err := os.Open("permissions.txt")
	if err != nil {
		fmt.Println("Error al abrir el archivo de permisos:", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) == 4 {
			// Solo agregar líneas que no coincidan con el archivo o directorio eliminado
			if !(parts[0] == username && parts[1] == path) {
				lines = append(lines, line)
			}
		}
	}

	// Escribir de nuevo el archivo sin las líneas eliminadas
	file, err = os.Create("permissions.txt")
	if err != nil {
		fmt.Println("Error al crear el archivo de permisos:", err)
		return
	}
	defer file.Close()

	for _, line := range lines {
		file.WriteString(line + "\n")
	}
}

// EDIT edita el contenido de un archivo basado en los permisos del usuario
func EDIT(path string, contenido string) {
	// Verificar si la ruta del archivo está vacía
	if path == "" {
		config.ErrorMessage = "Error: Se requiere una ruta de archivo."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Verificar si la ruta del archivo existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		config.ErrorMessage = "Error: La ruta del archivo especificado no existe."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Verificar si el usuario tiene permiso de escritura en el archivo
	if !hasWritePermission(path, usuario.Nombre) {
		config.ErrorMessage = "Error: No tiene permiso para editar el archivo especificado."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Leer el contenido del archivo que se usará para la edición
	newContent, err := ioutil.ReadFile(contenido)
	if err != nil {
		config.ErrorMessage = "Error al leer el archivo de contenido: " + err.Error()
		fmt.Println(config.ErrorMessage)
		return
	}

	// Escribir el nuevo contenido en el archivo especificado
	err = ioutil.WriteFile(path, newContent, 0644)
	if err != nil {
		config.ErrorMessage = "Error al escribir en el archivo: " + err.Error()
		fmt.Println(config.ErrorMessage)
		return
	}

	config.GeneralMessage = "Archivo editado con éxito: " + path
	fmt.Println(config.GeneralMessage)
}

// RENAME cambia el nombre de un archivo o carpeta basado en los permisos del usuario
func RENAME(path string, newName string) {
	// Verificar si la ruta del archivo está vacía
	if path == "" {
		config.ErrorMessage = "Error: Se requiere una ruta de archivo."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Verificar si la ruta del archivo existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		config.ErrorMessage = "Error: La ruta del archivo o carpeta especificada no existe."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Verificar si el usuario tiene permiso de escritura en el archivo
	if !hasWritePermission(path, usuario.Nombre) {
		config.ErrorMessage = "Error: No tiene permiso para renombrar el archivo o carpeta especificada."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Obtener el directorio del archivo original
	dir := filepath.Dir(path)
	// Crear la nueva ruta completa
	newPath := filepath.Join(dir, newName)

	// Verificar si ya existe un archivo o carpeta con el nuevo nombre
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		config.ErrorMessage = "Error: Ya existe un archivo o carpeta con el nombre especificado."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Renombrar el archivo o carpeta
	err := os.Rename(path, newPath)
	if err != nil {
		config.ErrorMessage = "Error al renombrar: " + err.Error()
		fmt.Println(config.ErrorMessage)
		return
	}

	config.GeneralMessage = "Renombrado con éxito: " + path + " a " + newPath
	fmt.Println(config.GeneralMessage)
}
func COPY(path *string, destino *string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */
	if !session {
		fmt.Println("--------------------------------------------------------------------------")
		fmt.Println("                   MKFILE: NO HAY UNA SESION INICIADA                     ")
		fmt.Println("--------------------------------------------------------------------------")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepaths := dire
	file, err := utilities_test.OpenFile(filepaths)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */

	var Inode0 structs.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	structs.PrintInode(Inode0)

	/* -------------------------------------------------------------------------- */
	/*                           OBTENEMOS LA RUTA PADRE                          */
	/* -------------------------------------------------------------------------- */
	tieneComillas := strings.Split(*path, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			*path = tieneComillas[0]
		} else {
			*path = tieneComillas[1]
		}
	}
	tieneComillas = strings.Split(*destino, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			*destino = tieneComillas[0]
		} else {
			*destino = tieneComillas[1]
		}
	}
	Carpetas := strings.Split(*path, "/")
	nuevaCarpeta := Carpetas[len(Carpetas)-1]
	partes := strings.Split(*destino, "/")
	partes = partes[1:]
	//println("Elementos en ruta padre")
	//println(len(partes) - 1)
	carpetaCreada := false
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//println("Bloques del inodo 0:")
	ultimo := 0
	root := false
	padreExiste := false
	for cont, i := range Inode0.I_block {
		if len(partes)-1 == 0 {
			//println("root es true")
			root = true
		}
		if i == -1 {
			ultimo = int(cont - 1)
			break
		}
		//println(i)
		if !root {
			existe := BuscarRuta(partes, i, 0)
			if existe {
				println("Existe la ruta padre")
				padreExiste = true
			}
		}
	}

	if root {
		existe := false
		for _, i := range Inode0.I_block {
			if i == -1 {
				break
			}
			// print("Buscando en el inodo 0 el bloque ")
			// println(i)
			existe = BuscarEspacioEnRoot(nuevaCarpeta, i)
			println("hay espacio")
			println(existe)
			if existe {
				break
			}
		}
		if !existe {
			println("Creando nuevo inodo y bloque")
			// BlockCounter++
			CrrSuperblock.S_blocks_count++
			Inode0.I_block[ultimo+1] = CrrSuperblock.S_blocks_count
			CrearFolderBlock(file, CrrSuperblock.S_blocks_count, nuevaCarpeta)
			println("Actualizando inodo 0")
			if err := utilities_test.WriteObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs.Inode{})))); err != nil {
				fmt.Println("Error reading inode:", err)
				return
			}
			structs.PrintInode(Inode0)

		}
		carpetaCreada = true
	}

	if padreExiste && !carpetaCreada {
		println("Creando carpeta desde padre")
		CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		carpetaCreada = true
	}

	if !carpetaCreada {
		if string(Padre.B_name[:]) != "" {
			println("creando a partir de carpetas existentes")
			fmt.Printf("Encontrado -> B_inode: %d B_name: %s\n", Padre.B_inodo, Padre.B_name)
			CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		} else {
			println("Creando todas las carpetas")
			CreandoCamino(0, nuevaCarpeta, file, partes)
		}
		carpetaCreada = true
	}
	if carpetaCreada {
		fmt.Println("--------------------------------------------------------------------------")
		fmt.Printf("                COPY:  %s COPIADO CORRECTAMENTE\n", nuevaCarpeta)
		fmt.Println("--------------------------------------------------------------------------")
	} else {
		println("Error: No se logro copiar el elemento")
	}
}
func MOVE(path string, destino *string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */

	if !session {
		fmt.Println("--------------------------------------------------------------------------")
		fmt.Println("                   MKDIR: NO HAY UNA SESION INICIADA                      ")
		fmt.Println("--------------------------------------------------------------------------")
		return
	}

	REMOVE(path)
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepaths := dire
	file, err := utilities_test.OpenFile(filepaths)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL INODO 0                            */
	/* -------------------------------------------------------------------------- */

	var Inode0 structs.Inode
	if err := utilities_test.ReadObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	structs.PrintInode(Inode0)

	/* -------------------------------------------------------------------------- */
	/*                           OBTENEMOS LA RUTA PADRE                          */
	/* -------------------------------------------------------------------------- */
	tieneComillas := strings.Split(path, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			path = tieneComillas[0]
		} else {
			path = tieneComillas[1]
		}
	}
	tieneComillas = strings.Split(*destino, "\"")
	if len(tieneComillas)-1 != 0 {
		if len(tieneComillas)-1 == 1 {
			*destino = tieneComillas[0]
		} else {
			*destino = tieneComillas[1]
		}
	}
	Carpetas := strings.Split(path, "/")
	nuevaCarpeta := Carpetas[len(Carpetas)-1]
	partes := strings.Split(*destino, "/")
	partes = partes[1:]
	//println("Elementos en ruta padre")
	//println(len(partes) - 1)
	carpetaCreada := false
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//println("Bloques del inodo 0:")
	ultimo := 0
	root := false
	padreExiste := false
	for cont, i := range Inode0.I_block {
		if len(partes)-1 == 0 {
			//println("root es true")
			root = true
		}
		if i == -1 {
			ultimo = int(cont - 1)
			break
		}
		//println(i)
		if !root {
			existe := BuscarRuta(partes, i, 0)
			if existe {
				println("Existe la ruta padre")
				padreExiste = true
			}
		}
	}

	if root {
		existe := false
		for _, i := range Inode0.I_block {
			if i == -1 {
				break
			}
			// print("Buscando en el inodo 0 el bloque ")
			// println(i)
			existe = BuscarEspacioEnRoot(nuevaCarpeta, i)
			println("hay espacio")
			println(existe)
			if existe {
				break
			}
		}
		if !existe {
			println("Creando nuevo inodo y bloque")
			// BlockCounter++
			CrrSuperblock.S_blocks_count++
			Inode0.I_block[ultimo+1] = CrrSuperblock.S_blocks_count
			CrearFolderBlock(file, CrrSuperblock.S_blocks_count, nuevaCarpeta)
			println("Actualizando inodo 0")
			if err := utilities_test.WriteObject(file, &Inode0, int64(CrrSuperblock.S_inode_start+0*int32(binary.Size(structs.Inode{})))); err != nil {
				fmt.Println("Error reading inode:", err)
				return
			}
			structs.PrintInode(Inode0)

		}
		carpetaCreada = true
	}

	if padreExiste && !carpetaCreada {
		println("Creando carpeta desde padre")
		CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		carpetaCreada = true
	}

	if !carpetaCreada {
		if string(Padre.B_name[:]) != "" {
			println("creando a partir de carpetas existentes")
			fmt.Printf("Encontrado -> B_inode: %d B_name: %s\n", Padre.B_inodo, Padre.B_name)
			CreandoCamino(Padre.B_inodo, nuevaCarpeta, file, partes)
		} else {
			println("Creando todas las carpetas")
			CreandoCamino(0, nuevaCarpeta, file, partes)
		}
		carpetaCreada = true
	}
	if carpetaCreada {
		fmt.Println("--------------------------------------------------------------------------")
		fmt.Printf("                MOVE:  %s  CORRECTAMENTE\n", nuevaCarpeta)
		fmt.Println("--------------------------------------------------------------------------")
	} else {
		println("Error: No se logro mover el elemento")
	}
}

// MKFILE creates a file with the specified properties
// MKFILE creates a file with the specified properties

// MKFILE creates a file with specified properties
func MKFILE(path string, size int, cont string, createParents bool) {
	// Check if the path is empty
	if path == "" {
		config.ErrorMessage = "Error: Path is required."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Handle parent directory creation
	if createParents {
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			config.ErrorMessage = "Error creating parent directories: " + err.Error()
			fmt.Println(config.ErrorMessage)
			return
		}
	} else {
		// Check that the parent directory exists
		parentDir := filepath.Dir(path)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			config.ErrorMessage = "Error: Parent directory does not exist."
			fmt.Println(config.ErrorMessage)
			return
		}
	}

	// Check if the file exists
	if _, err := os.Stat(path); err == nil {
		// File exists
		if !hasWritePermission(path, usuario.Nombre) {
			config.ErrorMessage = "Error: No permission to overwrite the existing file."
			fmt.Println(config.ErrorMessage)
			return
		}
	}

	// Create the file
	file, err := os.Create(path)
	if err != nil {
		config.ErrorMessage = "Error creating file: " + err.Error()
		fmt.Println(config.ErrorMessage)
		return
	}
	defer file.Close() // Ensure the file is always closed after opening

	// Handle size parameter
	if size > 0 {
		content := generateContent(size)
		_, err = file.WriteString(content)
		if err != nil {
			config.ErrorMessage = "Error writing content to file: " + err.Error()
		}
	} else if cont != "" {
		// Handle content from existing file
		contentBytes, err := os.ReadFile(cont)
		if err != nil {
			config.ErrorMessage = "Error reading content from file: " + err.Error()
			fmt.Println(config.ErrorMessage)
			return
		}
		_, err = file.Write(contentBytes)
		if err != nil {
			config.ErrorMessage = "Error writing content to file: " + err.Error()
		}
	}

	// Log the permission for the created file with ID
	logPermission(usuario.Nombre, path, "File", ID)
	config.GeneralMessage = "File created successfully: " + path
	fmt.Println(config.GeneralMessage)
}

// MKDIR creates directories with the specified properties
// MKDIR creates directories with the specified properties
// MKDIR creates directories with the specified properties
// MKDIR creates directories with the specified properties
func MKDIR(path string, createParents bool) {
	// Check if the path is empty
	if path == "" {
		config.ErrorMessage = "Error: Path is required."
		fmt.Println(config.ErrorMessage)
		return
	}

	// If createParents is true, create all necessary parent directories
	if createParents {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			config.ErrorMessage = "Error creating directory: " + err.Error()
			fmt.Println(config.ErrorMessage)
			return
		}
	} else {
		// Check that the parent directory exists
		parentDir := filepath.Dir(path)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			config.ErrorMessage = "Error: Parent directory does not exist."
			fmt.Println(config.ErrorMessage)
			return
		}

		// Attempt to create the directory
		err := os.Mkdir(path, 0755)
		if err != nil {
			config.ErrorMessage = "Error creating directory: " + err.Error()
			fmt.Println(config.ErrorMessage)
			return
		}
	}

	// Log the permission for the created directory with ID
	logPermission(usuario.Nombre, path, "DIR", ID)
	config.GeneralMessage = "Directory created successfully: " + path
	fmt.Println(config.GeneralMessage)
}

// generateContent generates a string with numeric characters based on the specified size
func generateContent(size int) string {
	content := strings.Repeat("0123456789", size/10)
	if len(content) < size {
		content += "0123456789"[:size-len(content)]
	}
	return content
}

// hasWritePermission checks write permissions for the given directory path
func hasWritePermission(dir string, username string) bool {
	file, err := os.Open("permissions.txt") // Directly using "permissions.txt"
	if err != nil {
		fmt.Println("Error opening permissions file:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) == 4 {
			if parts[0] == username && parts[1] == dir && parts[2] == "File" {
				return true
			}
		}
	}

	// Check for directory permissions
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) == 4 {
			if parts[0] == username && strings.HasPrefix(dir, parts[1]) && parts[2] == "DIR" {
				return true
			}
		}
	}

	return false
}

// logPermission logs the permission for the created file or directory
func logPermission(username string, path string, itemType string, id string) {
	// Read existing permissions to check for duplicates
	existingPermissions := make(map[string]bool)

	file, err := os.Open("permissions.txt") // Directly using "permissions.txt"
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			existingPermissions[line] = true
		}
	}

	// Create the new permission entry
	newEntry := fmt.Sprintf("%s|%s|%s|%s", username, path, itemType, id)

	// Check if the entry already exists
	if _, exists := existingPermissions[newEntry]; !exists {
		// Append the new entry to the permissions file
		file, err := os.OpenFile("permissions.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // Directly using "permissions.txt"
		if err != nil {
			fmt.Println("Error opening permissions file:", err)
			return
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(newEntry + "\n")
		if err != nil {
			fmt.Println("Error writing to permissions file:", err)
			return
		}
		writer.Flush()
	}
}

// updatePermissions actualiza el archivo permissions.txt para cambiar el propietario de un archivo o directorio
func updatePermissions(path string, newOwner string) bool {
	file, err := os.Open("permissions.txt")
	if err != nil {
		fmt.Println("Error al abrir el archivo de permisos:", err)
		return false
	}
	defer file.Close()

	var lines []string
	found := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) == 4 {
			// Verificamos si la línea corresponde al archivo o directorio que se está cambiando
			if parts[1] == path {
				parts[0] = newOwner // Cambiamos el propietario
				found = true
			}
			lines = append(lines, strings.Join(parts, "|"))
		}
	}

	if !found {
		// Si no encontramos la línea, agregamos una nueva entrada para el nuevo propietario
		lines = append(lines, fmt.Sprintf("%s|%s|File|other_permissions", newOwner, path)) // Ajusta 'other_permissions' según sea necesario
	}

	// Reescribimos el archivo con las nuevas líneas
	if err := os.WriteFile("permissions.txt", []byte(strings.Join(lines, "\n")), 0644); err != nil {
		fmt.Println("Error al escribir en el archivo de permisos:", err)
		return false
	}

	return true
}

// userExists verifica si el usuario existe en el sistema
// userExists verifica si el usuario existe en el sistema
func userExists(username string) bool {
	// Obtener la ruta del disco montado
	filepath := FindPathByID(MountedDiskList, ID)
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error abriendo el archivo de usuarios:", err)
		return false
	}
	defer file.Close()

	// Leer el MBR para poder acceder a la partición que contiene users.txt
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error leyendo el MBR:", err)
		return false
	}

	// Buscar la partición con el ID
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partición no encontrada")
		return false
	}

	// Cargar el superblock
	var tempSuperblock structs.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error leyendo el superblock:", err)
		return false
	}
	CrrSuperblock = tempSuperblock

	// Leer el inodo donde está users.txt
	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error leyendo el inodo:", err)
		return false
	}

	// Leer el Fileblock de users.txt
	var Fileblock structs.Fileblock
	blockNum := crrInode.I_block[0] // Asumimos que está en el primer bloque
	if err := utilities_test.ReadObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
		fmt.Println("Error leyendo el Fileblock:", err)
		return false
	}

	data := string(Fileblock.B_content[:])
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		items := strings.Split(line, ",")
		if len(items) > 3 && items[len(items)-2] == username {

			return true // Usuario encontrado
		}
	}
	return false // Usuario no encontrado
}

// CHOWN cambia el propietario de uno o varios archivos o carpetas
func CHOWN(path string, newOwner string, recursive bool) {
	// Verificar si la ruta del archivo está vacía
	if path == "" {
		config.ErrorMessage = "Error: Se requiere una ruta."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Verificar si la ruta del archivo existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		config.ErrorMessage = "Error: La ruta especificada no existe."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Verificar si el nuevo propietario existe
	if !userExists(newOwner) {
		config.ErrorMessage = "Error: El usuario especificado no existe."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Verificar si el usuario logueado tiene permisos para cambiar el propietario
	if !hasWritePermission(path, usuario.Nombre) {
		config.ErrorMessage = "Error: El usuario no tiene permisos para cambiar el propietario de este archivo."
		fmt.Println(config.ErrorMessage)
		return
	}

	// Cambiar el propietario del archivo o carpeta
	if recursive {
		err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Actualiza los permisos después de cambiar el propietario
			if !updatePermissions(currentPath, newOwner) {
				fmt.Println("Error al actualizar los permisos para:", currentPath)
			}
			return nil
		})
		if err != nil {
			config.ErrorMessage = err.Error()
			fmt.Println(config.ErrorMessage)
			return
		}
	} else {
		// Actualiza los permisos del archivo o carpeta especificada
		if !updatePermissions(path, newOwner) {
			fmt.Println("Error al actualizar los permisos para:", path)
		}
	}

	config.GeneralMessage = "Propietario cambiado con éxito: " + path + " a " + newOwner
	fmt.Println(config.GeneralMessage)
}
