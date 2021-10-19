package db

import (
	"gorm.io/gorm"
)

type Subscriber struct {
  gorm.Model
  Firstname  string
  Lastname string
  Phone string `gorm:"unique"`
  Subscriptions []Subscription
}

type Subscription struct {
  gorm.Model
  Topic string  //sniffer device unique ID
  Subscriber Subscriber
  SubscriberID uint
}

type ProcessingResult struct {
  gorm.Model
  Date string
  Time string
  Sniffer string  //the topic ie sniffer device unique ID
  Disease string  //default "Late blight"
  PlantStatus string  //healthy, mild +ve, moderate +ve, severe +ve
  Recommendation string
}



