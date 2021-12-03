package main


import (
  "io"
  "fmt"
  "bytes"
  "net/http"
  "encoding/json"
  "github.com/datavoc/server-pubsub/processor"
  "github.com/datavoc/server-pubsub/db"
  "github.com/gorilla/websocket"
)

func isJSON(str string) bool {
  var jsn json.RawMessage
  return json.Unmarshal([]byte(str), &jsn) == nil
}

func (ps *Pubsub) Publish(topic string, msg string) string {
  if ps.closed {
    return "Error: server channels closed"
  }
  
  //call the processor and wait for a response (ie, processed message)
  processedMsg, err := processor.Process(msg)
  if err != nil {
    fmt.Println(err)
    return "Error: server cant process"
  }
  //enforce that processedMsg is a valid json string
  ok := isJSON(processedMsg)
  if !ok {
    fmt.Println("Error: processedMsg from the processor is not a valid json string. Message neighther published nor saved.")
    return "Error: invalid data"
  }
  
  for _, ch := range ps.subs[topic] {
       ch <- processedMsg
       fmt.Println("Published on topic (",topic,") to a channel")
  }
  
  //prepare ProcessedResult object from the processedMsg valid json string
  var processedMsgResponse db.ProcessedResult
  json.Unmarshal([]byte(processedMsg), &processedMsgResponse)
  
  //persist the procssed message to the db for history
  database.Create(&db.ProcessedResult{
      Date: processedMsgResponse.Date, 
      Time: processedMsgResponse.Time, 
      Sniffer: processedMsgResponse.Sniffer, 
      Disease: processedMsgResponse.Disease, 
      PlantStatus: processedMsgResponse.PlantStatus, 
      Recommendation: processedMsgResponse.Recommendation,
  })
  
  shortMessageForLCD := fmt.Sprintf("plant status: %s", processedMsgResponse.PlantStatus)
  
  return shortMessageForLCD
}

type Publication struct {
    Topic string `json:"topic"`
    Message string `json:"message"`
}

func publishing(w http.ResponseWriter, r *http.Request) {
	//Upgrade to websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true; }
	
	webclient, _ := upgrader.Upgrade(w, r, nil)
	defer webclient.Close() 
	fmt.Println("Sniffer device connetion established.")
	for {
	    _, dataFromClient, err := webclient.ReadMessage() 
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) || err == io.EOF { 
		        break
	        }
	        fmt.Println("!! Error:", err) 
		}else if len(dataFromClient) > 8 {
		    fmt.Println("Data received from Sniffer device.")
		    //use the data submited by the client
		    var pbn Publication
	        _ = json.NewDecoder(bytes.NewReader(dataFromClient)).Decode(&pbn)
	        shortMessageForLCD := pubsubBroker.Publish(pbn.Topic, pbn.Message)
	        _ = webclient.WriteMessage(websocket.TextMessage, []byte(shortMessageForLCD))
	        fmt.Println("Short message for LCD screen => ", shortMessageForLCD)
		    //terminate the infinite loop
		    break
		}
	}
}



