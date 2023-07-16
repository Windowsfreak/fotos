package main

import "fotos/fotos"

func main() {
	err := fotos.DirMainTest()
	if err != nil {
		println(err.Error)
	}
}
