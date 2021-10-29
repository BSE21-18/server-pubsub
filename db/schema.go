package db

import (
	"gorm.io/gorm"
)

type Subscriber struct {
  gorm.Model
  Firstname  string
  Lastname string
  Phone string `gorm:"unique"`
}

type Subscription struct {
  gorm.Model
  Topic string `gorm:"unique"` //sniffer device unique ID
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



