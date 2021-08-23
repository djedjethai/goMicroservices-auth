package domain

import (
	"time"
)

type CustomerDomain struct {
	Name        string
	DateOfBirth time.Time
	City        string
	ZipCode     int
	Username    string
	Password    string
}
