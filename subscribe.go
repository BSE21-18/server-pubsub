package main


import (
    "github.com/datavoc/server-pubsub/processor"
)

func (ps *Pubsub) Subscribe(topic string) <-chan string {
  ps.mu.Lock()
  defer ps.mu.Unlock()

  ch := make(chan string, 10)
  ps.subs[topic] = append(ps.subs[topic], ch)
  
  //TODO: persist the details of this subcription into db
  
  return ch
}
