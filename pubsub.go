package main


import (
    "github.com/gorilla/websocket"
)

type Pubsub struct {
  subs   map[string][]chan string
  closed bool
}

func NewPubsub() *Pubsub {
  ps := &Pubsub{}
  ps.subs = make(map[string][]chan string)
  return ps
}

var pubsubBroker *Pubsub
var upgrader = websocket.Upgrader{} //use default options 

func (ps *Pubsub) Close() {
  if !ps.closed {
    ps.closed = true
    for _, subs := range ps.subs {
      for _, ch := range subs {
        close(ch)
      }
    }
  }
}





