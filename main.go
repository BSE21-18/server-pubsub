package main


import (
    "fmt"
    "sync"
    "log"
	"flag"
	"strings"
	"github.com/julienschmidt/httprouter"
)


type Pubsub struct {
  mu     sync.RWMutex
  subs   map[string][]chan string
  closed bool
}

func NewPubsub() *Pubsub {
  ps := &Pubsub{}
  ps.subs = make(map[string][]chan string)
  return ps
}

func (ps *Pubsub) Close() {
  ps.mu.Lock()
  defer ps.mu.Unlock()

  if !ps.closed {
    ps.closed = true
    for _, subs := range ps.subs {
      for _, ch := range subs {
        close(ch)
      }
    }
  }
}

func getRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/pub", publishing)
	router.GET("/sub", subscribing)
	router.POST("/register", registering)
	return router
}

func main() {
    //++++| os.Args |+++++
    wsEndPoint := ":7000" 
    addr := flag.String("addr", wsEndPoint, "websocket API gateway service address") 
    flag.Parse()
    //++++++++++++++++++++
    
    pubsubBroker := NewPubsub()
    defer pubsubBroker.Close()
    
    fmt.Println("DATAVOC Websocket API gateway server listening on port: "+(strings.Split(wsEndPoint,":")[1])) 
    log.Fatal(http.ListenAndServe(*addr, getRouter()))
}








