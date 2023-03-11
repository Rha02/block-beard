package main

import "github.com/Rha02/block-beard/src/utils"

func main() {
	msg := utils.IsFoundHost("127.0.0.2", 8080)
	println(msg)

	neighbors := utils.FindNeighbors("127.0.0.1", 8080, 0, 3, 8080, 8083)
	println("Neighbors:")
	for _, neighbor := range neighbors {
		println(neighbor)
	}

}
