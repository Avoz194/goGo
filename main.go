package main

import (
	api "github.com/Avoz194/goGo/APIHandler"
	db "github.com/Avoz194/goGo/DBHandler"
	"os"
)

func main() {

	os.Setenv("MYSQL_DBUSER", "Aviv")
	os.Setenv("MYSQL_DBPASS", "123456goGO")
	db.CreateDatabase()
	print("created DB.")
	api.CreateServer()
	print("server on.")


}
