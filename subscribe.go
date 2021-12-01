package main


import (
    "fmt"
    "bytes"
    "net/http"
    "encoding/json"
    "github.com/datavoc/server-pubsub/db"
    "github.com/gorilla/websocket"
)

func (ps *Pubsub) Subscribe(topic string) chan string {
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

func registering(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	var reg Registration
	_ = json.NewDecoder(r.Body).Decode(&reg)
	//write into db
	database.Create(db.Subscription{Topic: reg.SnifferId, Subscriber: db.Subscriber{ Firstname: reg.Firstname, Lastname: reg.Lastname, Phone: reg.Phone }})
	json.NewEncoder(w).Encode(reg)
}

type UpdatesRequest struct {
    Device string `json:"device"`
    Phone string `json:"phone"`
}

func subscribing(w http.ResponseWriter, r *http.Request) {
	//Upgrade to websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true; }
	
	webclient, _ := upgrader.Upgrade(w, r, nil) 
	defer webclient.Close() 
	
	fmt.Println("> New websocket request: subscribing ...")
	
	//retrieve a list of topics this client subscribed to from db
	var channels []chan string
    defer (func(){for _, ch := range channels {
       close(ch)
    }})()
      
    for {
        _, dataFromClient, err := webclient.ReadMessage() 
        if err != nil {
            if ce, ok := err.(*websocket.CloseError); ok {
                switch ce.Code {
                case websocket.CloseNormalClosure,
                    websocket.CloseGoingAway,
                    websocket.CloseNoStatusReceived:
                    break
                }
            }
        }else if len(dataFromClient) > 8 {
	        fmt.Println("len(dataFromClient) > 8 = true")
	        //use the data submited by the client
	        var updreq UpdatesRequest
            _ = json.NewDecoder(bytes.NewReader(dataFromClient)).Decode(&updreq)
            fmt.Println("updreq.Device =", updreq.Device, ", updreq.Phone=", updreq.Phone)
            if updreq.Device == "All" {
                fmt.Println("updreq.Device == All = true")
                //getting updates from all devices to which this user subscribed 
                //(ie, from all devices registered by the same user)
                rows, err := database.Table("subscriptions").Select("subscriptions.topic as topic").Joins("left join subscribers on subscribers.id = subscriptions.subscriber_id").Where("subscribers.phone = ?", updreq.Phone).Rows()
                if err != nil {
                    fmt.Println("!!Error while reterieving subscriptions. ", err.Error())
                }else{
                    fmt.Println("< rows.Next() >")
                    for rows.Next() {
                        var topic string
                        rows.Scan(&topic)
                        fmt.Println("retrieved topic =", topic)
                        //get a channel through which to receive updates
                        myChannel := pubsubBroker.Subscribe(topic)
                        
                        //push the channel to a list of channels
                        channels = append(channels, myChannel)
                    }
                    fmt.Println("</ rows.Next() >")
                }
           }else{
              //filtering to receive updates from one soecific device (topic)
              //first empty the channels map in case it's not empty
              channels = nil
              topic := updreq.Device
              myChannel := pubsubBroker.Subscribe(topic)
              channels = append(channels, myChannel)
           }
	    }
	    
	    for _, ch := range channels {
            fmt.Println("getting here means the user already has atleast one channel")
            //getting here means the user already has atleast one channel.
            massege, ok := <-ch
            if ok != false {
               fmt.Println("Received msg: ", massege, ok)
               msg := []byte(massege)
               err = webclient.WriteMessage(websocket.TextMessage, msg)
            }
            fmt.Println(ok, "Received msg from channel: ", massege)
        }
    }
    fmt.Println("Loop exited because client closed the connection.")
}





