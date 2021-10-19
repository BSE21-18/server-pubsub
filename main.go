package main


import (
    "fmt"
    "sync"
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

func main() {
    fmt.Println("DATAVOC server started: waiting for publications and subscriptions")
    //TODO: create a new instance of pubsub
    //TODO: defer close the pubsub
}






