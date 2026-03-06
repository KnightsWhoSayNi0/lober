package db

import "time"

type Event struct {
	ID      int64
	Command string    `json:"command"`
	User    string    `json:"user"`
	C2      string    `json:"c2"`
	Scope   string    `json:"scope"`
	Time    time.Time `json:"time"`
}

type User struct {
	ID       int64
	Username string
	Email    string
	Team     string // todo; users can be on multiple teams
}

type Team struct {
	ID    int64
	Name  string
	Color string
	Lead  string
}

type C2 struct {
	ID   int64
	Name string
}

type Scope struct {
	ID   int64
	Name string
}
