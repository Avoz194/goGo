package main

import (
	api "github.com/Avoz194/goGo/APIHandler"
	db "github.com/Avoz194/goGo/DBHandler"
)

func main() {
	db.CreateDatabase()
	println("created DB.")
	api.CreateServer()
	println("server on.")


}
