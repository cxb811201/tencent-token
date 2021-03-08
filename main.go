package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tidwall/buntdb"
)

var app = NewApp()

func main() {
	var err error

	var (
		version = flag.Bool("version", false, "version v0.2")
		config  = flag.String("config", "account.json", "config file.")
		port    = flag.Int("port", 8000, "listen port.")
	)

	flag.Parse()

	if *version {
		fmt.Println("v0.2")
		os.Exit(0)
	}

	app.SetAccounts(config)
	app.DB, err = buntdb.Open("tencent.db")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer app.DB.Close()

	app.StartTokenCheckTask()

	InitRoute(app.Web.HttpServer)
	log.Println("Start AccessToken Server on ", *port)
	app.Web.StartServer(*port)
}
