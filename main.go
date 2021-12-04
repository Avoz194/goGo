package main

import (
	api "github.com/Avoz194/goGo/APIHandler"
	db "github.com/Avoz194/goGo/DBHandler"
)

func main() {
	db.CreateDatabase()
	print("created DB.")
	api.CreateServer()
	print("server on.")


}
