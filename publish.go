package main


import (
    "fmt"
    "github.com/datavoc/server-pubsub/processor"
    "github.com/datavoc/server-pubsub/db"
)

func (ps *Pubsub) Publish(topic string, msg string) {
  ps.mu.RLock()
  defer ps.mu.RUnlock()

  if ps.closed {
    return
  }
  
  //TODO: call the processor and wait for a response (ie, processed message)
  processedMsg := processor.Process(msg)

  for _, ch := range ps.subs[topic] {
    //TODO: if the channel "ch" is closed, then
        //delete it from this array of channels (subscribers)
    //else
        ch <- processedMsg
    //end if
  }
  
  //TODO: prepare ProcessingResult object from the processedMsg valid json string
  // 
  
  //persist the procssed maessage to the db for history
  database, err := db.Connect()
  if err != nil {
    fmt.Println(err)
  }else{
    database.Create(&ProcessingResult{
      Date: "2021-10-20", 
      Time: "09:15:02", 
      Sniffer: "DV0897", 
      Disease: "Late blight", 
      PlantStatus: "mild +ve", 
      Recommendation: "Please spray using the recommended chemical immediately or call experts for help"
    })
  }
}

func publishing(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
	//Upgrade to websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true; }
	
	webclient, _ := upgrader.Upgrade(w, r, nil) 
	defer webclient.Close() 
	
	pubsubBroker.Publish(topic, msg)
}



