package functions_test

import (
	"backend/Global"
	structs "backend/Structs"
	utilities_test "backend/Utilities"
	"backend/config"
	"bufio"
	"encoding/binary"
	"fmt"
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
func RMGRP(name string) {
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
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

	// Leer el MBR del disco
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                               CARGAMOS EL MBR                              */
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
	/*                   LEEMOS EN INODO 1 DONDE ESTA USERS.TXT                   */
	/* -------------------------------------------------------------------------- */
	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                      LEEMOS EL CONTENIDO DEL FILEBLOCK                     */
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

	/* -------------------------------------------------------------------------- */
	/*                      COLOCAMOS EL STATUS DE ELIMINADO                      */
	/* -------------------------------------------------------------------------- */
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	lines := strings.Split(currentContent, "\n")
	deleted := false
	for i, line := range lines {
		if strings.Contains(line, name) {
			lines[i] = "0,G," + name
			deleted = true
			break
		}
	}

	/* -------------------------------------------------------------------------- */
	/*                   VERIFICAMOS BLOQUES O MENSAJE NOT FOUND                  */
	/* -------------------------------------------------------------------------- */
	if !deleted {
		searchIndex++
		if searchIndex > blockIndex {
			fmt.Println("Group not found")
			searchIndex = 0
			return
		}
		RMGRP(name)

	}

	/* -------------------------------------------------------------------------- */
	/*                          ACTUALIZAMOS EL CONTENIDO                         */
	/* -------------------------------------------------------------------------- */
	newContent := strings.Join(lines, "\n")
	copy(Fileblock.B_content[:], newContent)
	if deleted {
		blockNum := crrInode.I_block[searchIndex]

		// Elimina la línea correspondiente a `name`
		currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
		lines := strings.Split(currentContent, "\n")

		// Reemplazar el nombre con la lógica moral que tenías:
		for i := 0; i < len(lines); i++ {
			if strings.Contains(lines[i], name) {
				// Eliminar la línea
				lines = append(lines[:i], lines[i+1:]...) // Elimina la línea
				deleted = true
				break
			}
		}

		// Si se eliminó una línea, actualizar el contenido del bloque
		if deleted {
			// Reunir el contenido actualizado
			updatedContent := strings.Join(lines, "\n") + "\n"
			// Convertir a un slice de bytes
			copy(Fileblock.B_content[:], []byte(updatedContent))

			// Escribir el bloque actualizado en el disco
			if err := utilities_test.WriteObject(file, &Fileblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
				fmt.Println("Error writing Fileblock to disk:", err)

				return
			}

			config.SetErrorMessage("ELIMINADO")
		}

		// Imprimir las líneas para depuración
		for i := range lines {
			println(lines[i])
		}

		searchIndex = 0
	}
}
func MKUSR(user string, pass string, grp string) {
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

func MKDIR(path *string, r *bool) {
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
func MKFILE(path string, r bool, size int, cont string) {
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
func REMOVE(path *string) {
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
func RENAME(path *string, name *string) {
	/* -------------------------------------------------------------------------- */
	/*                  COMPROBAMOS SI HAY UNA SESSION EXISTENTE                  */
	/* -------------------------------------------------------------------------- */

	if !session {
		fmt.Println("--------------------------------------------------------------------------")
		fmt.Println("                   MKDIR: NO HAY UNA SESION INICIADA                      ")
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

	carpetas := strings.Split(*path, "/")
	carpetas = carpetas[1:]
	nuevaCarpeta := carpetas[len(carpetas)-1]
	/* -------------------------------------------------------------------------- */
	/*                     RECORREMOS LOS BLOQUES DEL INODO 0                     */
	/* -------------------------------------------------------------------------- */
	//println("Bloques del inodo 0:")
	for _, i := range Inode0.I_block {
		if i == -1 {
			break
		}
		//println(i)
		rename := Rename(carpetas, i, 0, *name)
		if rename {
			println(nuevaCarpeta + " renombrada con exito")
		}
	}
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
func MOVE(path *string, destino *string) {
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
		fmt.Printf("                MOVE:  %s  CORRECTAMENTE\n", nuevaCarpeta)
		fmt.Println("--------------------------------------------------------------------------")
	} else {
		println("Error: No se logro mover el elemento")
	}
}
