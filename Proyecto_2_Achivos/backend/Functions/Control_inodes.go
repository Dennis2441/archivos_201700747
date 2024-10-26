package functions_test

import (
	structs "backend/Structs"
	utilities_test "backend/Utilities"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func CargarBloque(file *os.File, bloque int32) structs.Folderblock {
	var FolderBlock structs.Folderblock
	//var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &FolderBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Folderblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return structs.Folderblock{}
	}
	return FolderBlock
}
func CargarArbol(tempSuperblock structs.Superblock, file *os.File, num int32) {
	var Inode structs.Inode
	if err := utilities_test.ReadObject(file, &Inode, int64(CrrSuperblock.S_inode_start+num*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                         RECORREMOS CREANDO EL ARBOL                        */
	/* -------------------------------------------------------------------------- */
	inodos += fmt.Sprintf("inodo%d[label=\"{INODO %d}\n|{I_uid|%d}\n|{I_gid|%d}\n|{I_size|%d}\n|{I_atime|%s}\n|{I_ctime|%s}\n|{I_mtime|%s}\n|{I_block|%d}\n|{I_type|%s}\n|{I_perm|%s}\"];\n\n",
		int(num),
		int(num),
		int(Inode.I_uid),
		int(Inode.I_gid),
		int(Inode.I_size),
		Inode.I_atime[:],
		Inode.I_ctime[:],
		Inode.I_mtime[:],
		Inode.I_block[:],
		string(Inode.I_type),
		string(Inode.I_perm[:]))
	for _, i := range Inode.I_block {
		if i == -1 {
			break
		}
		relaciones += fmt.Sprintf("inodo%d -> bloque%d;\n", num, i)
		bloque := CargarBloque(file, i)
		LeerBloque(bloque, i, file)
	}
}

func LeerBloque(bloque structs.Folderblock, j int32, file *os.File) {
	// print("BLOQUE ")
	// println(i)
	bloques += fmt.Sprintf("bloque%d[label=\"{Bloque %d}\n",
		int(j),
		int(j))
	//structs_test.PrintFolderBlock(bloque)

	for i := 0; i < 4; i++ {
		nombre := string(bytes.Trim(bloque.B_content[i].B_name[:], "\x00"))
		comas := strings.Split(nombre, ",")
		punto := strings.Split(nombre, ".")
		if nombre == "" {
			bloques += fmt.Sprintf("|{Inode | %d | Name: | %s}\n", -1, " ")
			//println("nombre vacio")
			continue
		}
		if nombre == "." {
			//println("nombre .")
			bloques += fmt.Sprintf("|{Inode | %d | Name: | %s}\n", bloque.B_content[i].B_inodo, ".")
			continue
		}
		if nombre == ".." {
			//println("nombre ..")
			bloques += fmt.Sprintf("|{Inode | %d | Name: | %s}\n", bloque.B_content[i].B_inodo, "..")
			continue
		}
		if len(comas) != 1 {
			//println("comas")
			//println(len(comas))
			continue
		}
		if len(punto) != 1 {
			bloques += fmt.Sprintf("|{Inode | %d | Name: | %s}\n", bloque.B_content[i].B_inodo, nombre)
			relaciones += fmt.Sprintf("bloque%d -> inodo%d;\n", j, bloque.B_content[i].B_inodo)

			//println("comas")
			//println(len(comas))
			continue
		}
		bloques += fmt.Sprintf("|{Inode | %d | Name: | %s}\n", bloque.B_content[i].B_inodo, nombre)
		relaciones += fmt.Sprintf("bloque%d -> inodo%d;\n", j, bloque.B_content[i].B_inodo)
	}
	bloques += "\"];\n"
	for i := 0; i < 4; i++ {
		nombre := string(bytes.Trim(bloque.B_content[i].B_name[:], "\x00"))
		comas := strings.Split(nombre, ",")
		punto := strings.Split(nombre, ".")

		if nombre == "" {
			//println("nombre vacio")
			continue
		}
		if nombre == "." {
			//println("nombre .")
			continue
		}
		if nombre == ".." {
			//println("nombre ..")
			continue
		}
		if len(comas) != 1 {
			LeerFileblock(bloque.B_content[i].B_inodo, file)
			continue
		}
		if len(punto) != 1 {
			LeerFileblock(bloque.B_content[i].B_inodo, file)
			continue
		}
		LeerInodo(bloque.B_content[i].B_inodo, file)
	}

}
func CargarInodo(file *os.File, inodo int32) structs.Inode {
	var Inode structs.Inode
	if err := utilities_test.ReadObject(file, &Inode, int64(CrrSuperblock.S_inode_start+inodo*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return structs.Inode{}
	}
	return Inode
}
func LeerInodo(i int32, file *os.File) {
	// print("INODO ")
	// println(i)
	Inode := CargarInodo(file, i)
	//structs_test.PrintInode(Inode)
	inodos += fmt.Sprintf("inodo%d[label=\"{INODO %d}\n|{I_uid|%d}\n|{I_gid|%d}\n|{I_size|%d}\n|{I_atime|%s}\n|{I_ctime|%s}\n|{I_mtime|%s}\n|{I_block|%d}\n|{I_type|%s}\n|{I_perm|%s}\"];\n\n",
		int(i),
		int(i),
		int(Inode.I_uid),
		int(Inode.I_gid),
		int(Inode.I_size),
		Inode.I_atime[:],
		Inode.I_ctime[:],
		Inode.I_mtime[:],
		Inode.I_block[:],
		string(Inode.I_type),
		string(Inode.I_perm[:]))
	for _, j := range Inode.I_block {
		if j == -1 {
			break
		}
		relaciones += fmt.Sprintf("inodo%d -> bloque%d;\n", i, j)
		bloque := CargarBloque(file, j)
		LeerBloque(bloque, j, file)
	}
}

func LeerFileblock(i int32, file *os.File) {
	/* -------------------------------------------------------------------------- */
	/*                                CARGAR INODO                                */
	/* -------------------------------------------------------------------------- */
	Inode := CargarInodo(file, i)
	inodos += fmt.Sprintf("inodo%d[label=\"{INODO %d}\n|{I_uid|%d}\n|{I_gid|%d}\n|{I_size|%d}\n|{I_atime|%s}\n|{I_ctime|%s}\n|{I_mtime|%s}\n|{I_block|%d}\n|{I_type|%s}\n|{I_perm|%s}\"];\n\n",
		int(i),
		int(i),
		int(Inode.I_uid),
		int(Inode.I_gid),
		int(Inode.I_size),
		Inode.I_atime[:],
		Inode.I_ctime[:],
		Inode.I_mtime[:],
		Inode.I_block[:],
		string(Inode.I_type),
		string(Inode.I_perm[:]))
	for _, block := range Inode.I_block {
		if block == -1 {
			break
		}
		relaciones += fmt.Sprintf("inodo%d -> bloque%d;\n", i, block)
		LeerFile(block, file)
	}
}
func CargarFileblock(file *os.File, bloque int32) structs.Fileblock {
	var FileBlock structs.Fileblock
	//var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &FileBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return structs.Fileblock{}
	}
	return FileBlock
}
func LeerFile(i int32, file *os.File) {
	bloques += fmt.Sprintf("bloque%d[label=\"{Bloque %d}\n",
		int(i),
		int(i))

	bloque := CargarFileblock(file, i)
	data := string(bytes.Trim(bloque.B_content[:], "\x00"))
	lines := strings.Split(data, "\n")
	content := ""
	for _, line := range lines {
		content += line + "\\n"
	}

	bloques += fmt.Sprintf("|{%s}\n", content)

	bloques += "\"];\n"
}
func BuscarRuta(ruta []string, bloque int32, busqueda int) bool {

	print("BUSCANDO: " + ruta[busqueda] + " en el bloque ")
	println(bloque)
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := DIRRE
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return false
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return true
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
		return true
	}
	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL BLOQUE                             */
	/* -------------------------------------------------------------------------- */
	var FolderBlock structs.Folderblock
	//var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &FolderBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Folderblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return false
	}
	encontrado := false
	for _, content := range FolderBlock.B_content {
		nombre := string(bytes.Trim(content.B_name[:], "\x00")) // Eliminar bytes nulos del final
		println("encontrando: " + nombre)
		//fmt.Printf("Encontrado -> B_inode: %d B_name: %s\n", content.B_inodo, nombre)

		if strings.TrimSpace(nombre) == strings.TrimSpace(ruta[busqueda]) {
			Padre = content
			encontrado = true
			padreBusqueda = busqueda
			//Cargar el inodo
			//recorrerlos bloques del inodo
			//recursividad para cada bloque hasta que se encuentre la otra parte de la ruta
			if busqueda < len(ruta)-1 {
				println("Aun no termina la ruta")
				Inode := CargarInodo(file, content.B_inodo)
				busqueda++
				println("BUSCANDO " + ruta[busqueda])
				fmt.Printf("Encontrado -> B_inode: %d B_name: %s\n", content.B_inodo, nombre)
				print("En los bloques del inodo ")
				println(content.B_inodo)
				structs.PrintInode(Inode)
				for _, i := range Inode.I_block {
					if i == -1 {
						break
					}
					existe := BuscarRuta(ruta, i, busqueda)
					if existe {
						return existe
					}
				}
				return false
			}
			return true // Elemento encontrado
		}
		if encontrado {
			break
		}
	}

	return false
}
func BuscarEspacioEnRoot(carpetaNueva string, bloque int32) bool {
	println("La nueva carpeta a añadir es " + carpetaNueva)
	print("Buscando en el bloque ")
	println(bloque)
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := DIRRE
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return false
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return true
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
		return true
	}
	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL BLOQUE                             */
	/* -------------------------------------------------------------------------- */
	var FolderBlock structs.Folderblock
	//var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &FolderBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Folderblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return false
	}

	//structs_test.PrintFolderBlock(FolderBlock)

	for i, content := range FolderBlock.B_content {
		nombre := string(bytes.Trim(content.B_name[:], "\x00")) // Eliminar bytes nulos del final
		println("encontro: " + nombre)
		if nombre == "" {
			println("Espacio disponible")
			//Escribir la carpeta
			bytes := []byte(carpetaNueva)
			copy(content.B_name[:], bytes)
			// InodeCounter++
			// BlockCounter++
			CrrSuperblock.S_inodes_count++
			//CrrSuperblock.S_blocks_count++
			punto := strings.Split(carpetaNueva, ".")
			content.B_inodo = CrrSuperblock.S_inodes_count
			if len(punto) != 1 {
				CrearInodoFileblock(file, CrrSuperblock.S_inodes_count)
				//CrearBloque(file, CrrSuperblock.S_blocks_count)
				FolderBlock.B_content[i] = content
				println("Actualizo")
				if err := utilities_test.WriteObject(file, &FolderBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Folderblock{})))); err != nil {
					fmt.Println("Error reading Fileblock:", err)
					return false
				}
				return true
			} else {
				CrearInodo(file, CrrSuperblock.S_inodes_count)
				//CrearBloque(file, CrrSuperblock.S_blocks_count)
				FolderBlock.B_content[i] = content
				println("Actualizo")
				if err := utilities_test.WriteObject(file, &FolderBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Folderblock{})))); err != nil {
					fmt.Println("Error reading Fileblock:", err)
					return false
				}
				return true
			}
		}
	}

	return false
}

func CrearBloque(file *os.File, blockNum int32) {
	CrrSuperblock.S_free_blocks_count -= 1
	var Folderblock structs.Folderblock
	//var crrInode structs_test.Inode
	copy(Folderblock.B_content[0].B_name[:], ".")
	Folderblock.B_content[0].B_inodo = 0
	copy(Folderblock.B_content[1].B_name[:], "..")
	Folderblock.B_content[1].B_inodo = 0
	if err := utilities_test.WriteObject(file, &Folderblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Folderblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}
}

func CrearInodo(file *os.File, inodeNum int32) {
	CrrSuperblock.S_free_inodes_count -= 1
	var Inode structs.Inode
	Inode.I_uid = usuario.ID
	Inode.I_gid = inodeNum
	Inode.I_size = int32(binary.Size(structs.Inode{}))
	// Obtener la marca de tiempo actual
	currentTime := time.Now()
	// Formatear la marca de tiempo como una cadena
	date := currentTime.Format("2006-01-02 15:04:05")
	copy(Inode.I_atime[:], date)
	copy(Inode.I_ctime[:], date)
	copy(Inode.I_mtime[:], date)
	Inode.I_type = '0'
	copy(Inode.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode.I_block[i] = -1
	}
	// BlockCounter++
	CrrSuperblock.S_blocks_count++
	Inode.I_block[0] = CrrSuperblock.S_blocks_count
	CrearBloque(file, CrrSuperblock.S_blocks_count)
	if err := utilities_test.WriteObject(file, &Inode, int64(CrrSuperblock.S_inode_start+inodeNum*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}
}

func CrearInodoFileblock(file *os.File, inodeNum int32) {
	CrrSuperblock.S_free_inodes_count -= 1
	var Inode structs.Inode
	Inode.I_uid = usuario.ID
	Inode.I_gid = inodeNum
	Inode.I_size = int32(binary.Size(structs.Inode{}))
	// Obtener la marca de tiempo actual
	currentTime := time.Now()
	// Formatear la marca de tiempo como una cadena
	date := currentTime.Format("2006-01-02 15:04:05")
	copy(Inode.I_atime[:], date)
	copy(Inode.I_ctime[:], date)
	copy(Inode.I_mtime[:], date)
	Inode.I_type = '0'
	copy(Inode.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode.I_block[i] = -1
	}
	// BlockCounter++
	CrrSuperblock.S_blocks_count++
	Inode.I_block[0] = CrrSuperblock.S_blocks_count
	CrearFileBlock(file, CrrSuperblock.S_blocks_count)
	if err := utilities_test.WriteObject(file, &Inode, int64(CrrSuperblock.S_inode_start+inodeNum*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}
}
func CrearFileBlock(file *os.File, blockNum int32) {
	CrrSuperblock.S_free_blocks_count -= 1
	var FileBlock structs.Fileblock
	//var crrInode structs_test.Inode
	if err := utilities_test.WriteObject(file, &FileBlock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}
}

func CrearFolderBlock(file *os.File, blockNum int32, carpeta string) {
	print("Creando el bloque ")
	println(blockNum)
	CrrSuperblock.S_free_blocks_count -= 1
	var Folderblock structs.Folderblock
	//var crrInode structs_test.Inode
	var folder structs.Folderblock
	copy(Folderblock.B_content[0].B_name[:], []byte(carpeta))
	// InodeCounter++
	CrrSuperblock.S_inodes_count++
	Folderblock.B_content[0].B_inodo = CrrSuperblock.S_inodes_count
	punto := strings.Split(carpeta, ".")
	if len(punto) != 1 {
		CrearInodoFileblock(file, CrrSuperblock.S_inodes_count)
	} else {
		CrearInodo(file, CrrSuperblock.S_inodes_count)
	}
	// BlockCounter++
	// CrearBloque(file, CrrSuperblock.S_blocks_count)
	// if err := utilities_test.WriteObject(file, &Folderblock, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs_test.Folderblock{})))); err != nil {
	// 	fmt.Println("Error reading Fileblock:", err)
	// 	return
	// }
	//var folder structs_test.Folderblock
	if err := utilities_test.ReadObject(file, &folder, int64(CrrSuperblock.S_block_start+blockNum*int32(binary.Size(structs.Folderblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}
	structs.PrintFolderBlock(folder)
}

func CreandoCamino(inodo int32, nuevaCarpeta string, file *os.File, ruta []string) {
	terminarRuta := false
	carpeta := nuevaCarpeta
	if padreBusqueda < len(ruta)-1 {
		terminarRuta = true
		padreBusqueda++
		carpeta = ruta[padreBusqueda]
	}
	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL INODO                             */
	/* -------------------------------------------------------------------------- */
	var Inode structs.Inode
	if err := utilities_test.ReadObject(file, &Inode, int64(CrrSuperblock.S_inode_start+inodo*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}
	structs.PrintInode(Inode)
	existe := int32(-1)
	for _, i := range Inode.I_block {
		if i == -1 {
			break
		}
		println(i)
		existe = BuscarEspacio(carpeta, i)
		if existe > 0 {
			break
		}
	}

	if terminarRuta && existe > 0 {
		CreandoCamino(existe, nuevaCarpeta, file, ruta)
	}
}
func BuscarEspacio(carpetaNueva string, bloque int32) int32 {
	println("La nueva carpeta a añadir es " + carpetaNueva)
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := DIRRE
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return -1
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return -1
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
		return -1
	}
	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL BLOQUE                             */
	/* -------------------------------------------------------------------------- */
	var FolderBlock structs.Folderblock
	//var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &FolderBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Folderblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return -1
	}

	for i, content := range FolderBlock.B_content {
		nombre := string(bytes.Trim(content.B_name[:], "\x00")) // Eliminar bytes nulos del final
		println("encontrando: " + nombre)
		if nombre == "" {
			println("Espacio disponible")
			//Escribir la carpeta
			bytes := []byte(carpetaNueva)
			copy(content.B_name[:], bytes)
			// InodeCounter++
			// BlockCounter++
			CrrSuperblock.S_inodes_count++
			//CrrSuperblock.S_blocks_count++
			content.B_inodo = CrrSuperblock.S_inodes_count
			punto := strings.Split(carpetaNueva, ".")
			if len(punto) != 1 {
				CrearInodoFileblock(file, CrrSuperblock.S_inodes_count)
			} else {
				CrearInodo(file, CrrSuperblock.S_inodes_count)
			}
			//CrearBloque(file, CrrSuperblock.S_blocks_count)
			FolderBlock.B_content[i] = content
			if err := utilities_test.WriteObject(file, &FolderBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Folderblock{})))); err != nil {
				fmt.Println("Error reading Fileblock:", err)
				return -1
			}
			return CrrSuperblock.S_inodes_count
		}
	}
	return -1
}
func EliminarRuta(ruta []string, bloque int32, busqueda int) bool {

	print("BUSCANDO: " + ruta[busqueda] + " en el bloque ")
	println(bloque)
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := dire
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return false
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return true
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
		return true
	}
	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL BLOQUE                             */
	/* -------------------------------------------------------------------------- */
	var FolderBlock structs.Folderblock
	//var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &FolderBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Folderblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return false
	}
	encontrado := false
	for i, content := range FolderBlock.B_content {
		nombre := string(bytes.Trim(content.B_name[:], "\x00")) // Eliminar bytes nulos del final
		println("encontrando: " + nombre)
		//fmt.Printf("Encontrado -> B_inode: %d B_name: %s\n", content.B_inodo, nombre)

		if strings.TrimSpace(nombre) == strings.TrimSpace(ruta[busqueda]) {
			encontrado = true
			//Cargar el inodo
			//recorrerlos bloques del inodo
			//recursividad para cada bloque hasta que se encuentre la otra parte de la ruta
			if busqueda < len(ruta)-1 {
				println("Aun no termina la ruta")
				Inode := CargarInodo(file, content.B_inodo)
				busqueda++
				println("BUSCANDO " + ruta[busqueda])
				fmt.Printf("Encontrado -> B_inode: %d B_name: %s\n", content.B_inodo, nombre)
				print("En los bloques del inodo ")
				println(content.B_inodo)
				structs.PrintInode(Inode)
				for _, i := range Inode.I_block {
					if i == -1 {
						break
					}
					existe := EliminarRuta(ruta, i, busqueda)
					if existe {
						return existe
					}
				}
				return false
			} else {
				println("Se prodece a eliminar " + ruta[busqueda])
				FolderBlock.B_content[i].B_name = [12]byte{}
				FolderBlock.B_content[i].B_inodo = -1
				structs.PrintFolderBlock(FolderBlock)
				if err := utilities_test.WriteObject(file, &FolderBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Folderblock{})))); err != nil {
					fmt.Println("Error reading Fileblock:", err)
					return false
				}
			}
			return true // Elemento encontrado
		}
		if encontrado {
			break
		}
	}

	return false
}
func Rename(ruta []string, bloque int32, busqueda int, name string) bool {

	print("BUSCANDO: " + ruta[busqueda] + " en el bloque ")
	println(bloque)
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := dire
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return false
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return true
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
		return true
	}
	/* -------------------------------------------------------------------------- */
	/*                             CARGAMOS EL BLOQUE                             */
	/* -------------------------------------------------------------------------- */
	var FolderBlock structs.Folderblock
	//var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &FolderBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Folderblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return false
	}
	encontrado := false
	for i, content := range FolderBlock.B_content {
		nombre := string(bytes.Trim(content.B_name[:], "\x00")) // Eliminar bytes nulos del final
		println("encontrando: " + nombre)
		//fmt.Printf("Encontrado -> B_inode: %d B_name: %s\n", content.B_inodo, nombre)

		if strings.TrimSpace(nombre) == strings.TrimSpace(ruta[busqueda]) {
			encontrado = true
			//Cargar el inodo
			//recorrerlos bloques del inodo
			//recursividad para cada bloque hasta que se encuentre la otra parte de la ruta
			if busqueda < len(ruta)-1 {
				println("Aun no termina la ruta")
				Inode := CargarInodo(file, content.B_inodo)
				busqueda++
				println("BUSCANDO " + ruta[busqueda])
				fmt.Printf("Encontrado -> B_inode: %d B_name: %s\n", content.B_inodo, nombre)
				print("En los bloques del inodo ")
				println(content.B_inodo)
				structs.PrintInode(Inode)
				for _, i := range Inode.I_block {
					if i == -1 {
						break
					}
					existe := Rename(ruta, i, busqueda, name)
					if existe {
						return existe
					}
				}
				return false
			} else {
				if nombre == name {
					println("Error: el archivo ya existe")
					return false
				}
				copy(FolderBlock.B_content[i].B_name[:], name)

				if err := utilities_test.WriteObject(file, &FolderBlock, int64(CrrSuperblock.S_block_start+bloque*int32(binary.Size(structs.Folderblock{})))); err != nil {
					fmt.Println("Error reading Fileblock:", err)
					return false
				}
			}
			return true // Elemento encontrado
		}
		if encontrado {
			break
		}
	}

	return false
}
