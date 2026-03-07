package models

import "time"

type Notification struct {
	ID        int
	Role      Role
	Message   string
	CreatedAt time.Time
}
