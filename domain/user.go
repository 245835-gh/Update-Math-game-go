package domain

import "time"

type User struct {
	Id        uint64
	Name      string
	TimeSpent time.Duration
}
