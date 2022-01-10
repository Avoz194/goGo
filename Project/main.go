package main

import (
	api "github.com/Avoz194/goGo/Project/APIHandler"
	db "github.com/Avoz194/goGo/Project/DBHandler"
)

func main() {
	db.CreateDatabase()
	println("created DB.")
	api.CreateServer()
	println("server on.")


}
