package Global

import "fmt"

type UserInfo struct {
	ID      int32
	Nombre  string
}

func PrintUser(usr UserInfo) {
	fmt.Print("ID: ")
	fmt.Println(usr.ID)
	fmt.Println("Nombre: " + usr.Nombre)
}

var Usuario UserInfo