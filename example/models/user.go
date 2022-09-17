package models

import (
	"time"
)

type User struct {
	Username  string
	CreatedAt time.Time
}
