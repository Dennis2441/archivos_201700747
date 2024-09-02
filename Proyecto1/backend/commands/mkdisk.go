package commands

import "fmt"

func CmdMkdisk(size int, fit string, unit string, path string) {
	fmt.Printf("\n mkdisk \n")
	fmt.Printf("Size: %d\n", size)
	fmt.Printf("Fit: %s\n", fit)
	fmt.Printf("Unit: %s\n", unit)
	fmt.Printf("Path: %s\n", path)
}
