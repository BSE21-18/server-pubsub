package main


import (
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/datavoc/server-pubsub/processor"
  "github.com/datavoc/server-pubsub/db"
  "github.com/julienschmidt/httprouter"
  "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} //use default options 

func isJSON(str string) bool {
  var jsn json.RawMessage
  return json.Unmarshal([]byte(str), &jsn) == nil
}

func (ps *Pubsub) Publish(topic string, msg string) {
  ps.mu.RLock()
  defer ps.mu.RUnlock()

  if ps.closed {
    return
  }
  
  //call the processor and wait for a response (ie, processed message)
  processedMsg, err := processor.Process(msg)
  if err != nil {
    fmt.Println(err)
    return
  }
  //enforce that processedMsg is a valid json string
  ok := isJSON(processedMsg)
  if !ok {
    fmt.Println("Error: processedMsg from the processor is not a valid json string. Message neighther published nor saved.")
    return
  }

  for _, ch := range ps.subs[topic] {
       ch <- processedMsg
  }
  
  //prepare ProcessedResult object from the processedMsg valid json string
  var processedMsgResponse db.ProcessedResult
  json.Unmarshal([]byte(processedMsg), &processedMsgResponse)
  
  //persist the procssed message to the db for history
  database, err := db.Connect()
  if err != nil {
    fmt.Println(err)
  }else{
    database.Create(&db.ProcessedResult{
      Date: processedMsgResponse.Date, 
      Time: processedMsgResponse.Time, 
      Sniffer: processedMsgResponse.Sniffer, 
      Disease: processedMsgResponse.Disease, 
      PlantStatus: processedMsgResponse.PlantStatus, 
      Recommendation: processedMsgResponse.Recommendation,
    })
  }
  
}

type Publication struct {
    Topic string `json:"topic"`
    Message string `json:"message"`
}

func publishing(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//Upgrade to websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true; }
	
	webclient, _ := upgrader.Upgrade(w, r, nil) 
	defer webclient.Close() 
	
	var pbn Publication
	_ = json.NewDecoder(r.Body).Decode(&pbn)
	pubsubBroker.Publish(pbn.Topic, pbn.Message)
}



