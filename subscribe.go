package main


import (
    "encoding/json"
    "github.com/datavoc/server-pubsub/processor"
    "github.com/datavoc/server-pubsub/db"
)

func (ps *Pubsub) Subscribe(topic string) <-chan string {
  ps.mu.Lock()
  defer ps.mu.Unlock()
  ch := make(<-chan string, 10)
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
        json.NewEncoder(w).Encode(struct{ errors: err})
      }
	database.Create(&Subscription{
       Topic: reg.SnifferId,
       Subscriber: &Subscriber{
          Firstname: reg.Firstname,
          Lastname: reg.Lastname,
          Phone: reg.Phone
       }
    })
	json.NewEncoder(w).Encode(reg)
}

func subscribing(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//Upgrade to websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true; }
	
	webclient, _ := upgrader.Upgrade(w, r, nil) 
	defer webclient.Close() 
	
	fmt.Println("> New websocket request: subscribing ...")
	//TODO: retrieve a list of topics this client subscribed to from db
	//
	//TODO: for each of the topics, request for a channel through which 
	//      to receive updates
	myChannel := pubsubBroker.Subscribe(topic)
	//TODO: push the channel to a list of channels
	//TODO ensure to close all channels before leaving ie defer close them
	//TODO: keep watching the channels for any new updates
	//TODO: if a channel is closed, remove it from the list of channels to listen to
	//TODO: if a new update arrives, push it to the front end
}





