package main

import (
	"time"
)

// Estructura para un inodo
type Inode struct {
	Uid    int32     // UID del usuario propietario del archivo o carpeta
	Gid    int32     // GID del grupo al que pertenece el archivo o carpeta
	Size   int64     // Tamaño del archivo en bytes
	Atime  time.Time // Última fecha en que se leyó el inodo sin modificarlo
	Ctime  time.Time // Fecha en la que se creó el inodo
	Mtime  time.Time // Última fecha en que se modifica el inodo
	Blocks [15]int64 // Array donde los primeros 12 son bloques directos, el 13 es bloque simple indirecto, el 14 es bloque doble indirecto, el 15 es bloque triple indirecto
	Type   byte      // Indica si es archivo o carpeta (1 = Archivo, 0 = Carpeta)
	Perm   [3]byte   // Permisos del archivo o carpeta (UGO)
}
