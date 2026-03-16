package main

import "time"

type Event struct {
	ID      int64     `json:"id"`
	Command string    `json:"command"`
	User    string    `json:"user"`
	C2      string    `json:"c2"`
	Scope   string    `json:"scope"`
	Time    time.Time `json:"time"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Team     string `json:"team"` // todo; users can be on multiple teams
}

type Team struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Lead  string `json:"lead"`
}

type C2 struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Scope struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
