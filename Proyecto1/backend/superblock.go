package main

import (
	"time"
)

// Estructura para el superbloque
type SuperBlock struct {
	FilesystemType  int32     // Identificación del sistema de archivos (por ejemplo, EXT2)
	InodesCount     int32     // Número total de inodos
	BlocksCount     int32     // Número total de bloques
	FreeBlocksCount int32     // Número de bloques libres
	FreeInodesCount int32     // Número de inodos libres
	Mtime           time.Time // Última fecha en la que el sistema fue montado
	Umtime          time.Time // Última fecha en que el sistema fue desmontado
	MntCount        int32     // Número de veces que se ha montado el sistema
	Magic           int32     // Valor que identifica al sistema de archivos (0xEF53)
	InodeSize       int32     // Tamaño del inodo
	BlockSize       int32     // Tamaño del bloque
	FirstInode      int32     // Primer inodo libre (dirección del inodo)
	FirstBlock      int32     // Primer bloque libre (dirección del bloque)
	BmInodeStart    int64     // Inicio del bitmap de inodos
	BmBlockStart    int64     // Inicio del bitmap de bloques
	InodeStart      int64     // Inicio de la tabla de inodos
	BlockStart      int64     // Inicio de la tabla de bloques
}
