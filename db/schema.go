package db

import (
	"gorm.io/gorm"
)

type Customer struct {
  gorm.Model
  Firstname  string
  Lastname string
  Phone string
}

type Subscription struct {
  gorm.Model
  Topic string  //sniffer device unique ID
  Client string //mobile client, identified by phone number of customer
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



