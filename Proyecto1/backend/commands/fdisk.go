package commands

import "fmt"

func CmdFdiskk(size int, fit string, unit string, path string, tipo string, name string) {
	fmt.Printf("\n fdisk \n")
	fmt.Printf("Size: %d\n", size)
	fmt.Printf("Fit: %s\n", fit)
	fmt.Printf("Unit: %s\n", unit)
	fmt.Printf("Path: %s\n", path)
	fmt.Printf("tipo: %s\n", tipo)
	fmt.Printf("nombre: %s\n", name)

}
