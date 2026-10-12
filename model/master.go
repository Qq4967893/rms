package model

import "time"

type Master struct {
	Login    int64  `json:"login"`
	Group    string `json:"group"`
	Passport int64  `json:"passport"`
}

type MasterRecord struct {
	ID        int64     `json:"id"`
	Login     int64     `json:"login"`
	Type      int8      `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
