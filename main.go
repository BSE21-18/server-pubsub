package main


import (
    "fmt"
    "log"
	"flag"
	"embed"
	"io/fs"
	"strings"
	"net/http"
	"gorm.io/gorm"
	"github.com/gorilla/mux"
	"github.com/datavoc/server-pubsub/db"
)

//go:embed client/web/*
var website embed.FS

//go:embed client/sniffer/*
var sniffer embed.FS

var database *gorm.DB

func getRouter() *mux.Router {
	client_web, _ := fs.Sub(website, "client/web")
	client_sniffer, _ := fs.Sub(sniffer, "client/sniffer")
	router := mux.NewRouter()
	router.HandleFunc("/publish", publishing).Methods("GET")
	router.HandleFunc("/subscribe", registering).Methods("POST")
	router.HandleFunc("/getupdates", subscribing).Methods("GET")
	router.PathPrefix("/client/web").Handler(http.StripPrefix("/client/web/", http.FileServer(http.FS(client_web))) ).Methods("GET")
	router.PathPrefix("/client/sniffer").Handler(http.StripPrefix("/client/sniffer/", http.FileServer(http.FS(client_sniffer))) ).Methods("GET")
	return router
}

func main() {
    //++++| os.Args |+++++
    wsEndPoint := ":7000" 
    addr := flag.String("addr", wsEndPoint, "websocket API gateway service address") 
    flag.Parse()
    //++++++++++++++++++++
    database, _ = db.Connect()
    
    pubsubBroker = NewPubsub()
    //defer pubsubBroker.Close()
    
    fmt.Println("DATAVOC Websocket API gateway server listening on port: "+(strings.Split(wsEndPoint,":")[1])) 
    log.Fatal(http.ListenAndServe(*addr, getRouter()))
}








