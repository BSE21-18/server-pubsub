package db

import (
	"gorm.io/gorm"
)

type Subscriber struct {
  gorm.Model
  Firstname  string `json:"firstname"`
  Lastname string `json:"lastname"`
  Phone string `json:"phone" gorm:"unique"`
}

type Subscription struct {
  gorm.Model
  Topic string `json:"topic" gorm:"unique"` //sniffer device unique ID
  Subscriber Subscriber
  SubscriberID uint
}

type ProcessedResult struct {
  gorm.Model
  Date string `json:"date"`
  Time string `json:"time"`
  Sniffer string `json:"sniffer"` //the topic ie sniffer device unique ID
  Disease string `json:"disease" default:"Late blight"`
  PlantStatus string `json:"plantStatus"` //healthy, mild +ve, moderate +ve, severe +ve
  Recommendation string `json:"recommendation"`
}



