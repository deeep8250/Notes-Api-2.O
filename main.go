package main

import (
	"pr01/config"
	"pr01/db"
)

func main() {
	cfg := config.Load()
	db.DbInit(*cfg)

}

//ererrrr
