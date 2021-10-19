package db

import (
	"gorm.io/gorm"
)

type Customer struct {
  gorm.Model
  Firstname  string
  Lastname string
}
