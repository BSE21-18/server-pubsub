package main


import (
    "fmt"
    "net/http"
    "encoding/json"
    "github.com/datavoc/server-pubsub/db"
    "github.com/gorilla/websocket"
    "github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{} //use default options 

func (ps *Pubsub) Subscribe(topic string) chan string {
  ps.mu.Lock()
  defer ps.mu.Unlock()
  ch := make(chan string, 10)
  ps.subs[topic] = append(ps.subs[topic], ch)
  return ch
}

type Registration struct {
  Firstname string `json:"firstname"`
  Lastname string  `json:"lastname"`
  Phone string     `gorm:"unique" json:"phone"`
  SnifferId string `json:"snifferId"`
}

func registering(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type","application/json")
	var reg Registration
	_ = json.NewDecoder(r.Body).Decode(&reg)
	//write into db
	database, err := db.Connect()
      if err != nil {
        json.NewEncoder(w).Encode(struct{}{ errors: err})
      }
	database.Create(&db.Subscription{Topic: reg.SnifferId,Subscriber: &db.Subscriber{ Firstname: reg.Firstname, Lastname: reg.Lastname, Phone: reg.Phone }})
	json.NewEncoder(w).Encode(reg)
}

func subscribing(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//Upgrade to websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true; }
	
	webclient, _ := upgrader.Upgrade(w, r, nil) 
	defer webclient.Close() 
	
	fmt.Println("> New websocket request: subscribing ...")
	
	//retrieve a list of topics this client subscribed to from db
	database, err := db.Connect()
      if err != nil {
        json.NewEncoder(w).Encode(struct{}{ errors: err})
      }
      
      //define the shape/format of the records which will come from db
      type Row struct { 
        Topic string 
      }
      
      strChan := make(<-chan string, 10)
      channels := []strChan{}
      defer (func(){for _, ch := range channels {
        close(ch)
      }})()
      
      var listOfTopics []Row
      database.Model(&db.Subscription{}).Select("subscription.topic").Joins("left join subscriber on subscriber.id = subscription.id").Find(&listOfTopics{})
      
	//for each of the topics, 
	for _, topic := range listOfTopics {
	    //request for a channel through which to receive updates
	    myChannel := pubsubBroker.Subscribe(topic)
	    
	    //push the channel to a list of channels
	    channels = append(channels, myChannel)
	}
	
	for {
	    for _, ch := range channels {
            massege, ok := <-ch
            if ok != false {
               fmt.Println("Received msg: ", massege, ok)
               msg := []byte(massege)
               err = webclient.WriteMessage(websocket.TextMessage, msg)
            }
        }
        _, _, err := webclient.ReadMessage()
        if ce, ok := err.(*websocket.CloseError); ok {
            switch ce.Code {
            case websocket.CloseNormalClosure,
                websocket.CloseGoingAway,
                websocket.CloseNoStatusReceived:
                s.env.Statusf("Web socket closed by client: %s", err)
                break
            }
        }
        
    }
}





