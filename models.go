package main

import "time"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Team     string `json:"team"` // todo; users can be on multiple teams
}

type Team struct {
	ID    int64    `json:"id"`
	Name  string   `json:"name"`
	Color string   `json:"color"`
	Users []string `json:"users"`
}

type C2 struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Scope struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Event struct {
	ID        int64     `json:"id"`
	Command   string    `json:"command"`
	User      string    `json:"user"`
	TeamColor string    `json:"team_color"`
	C2        string    `json:"c2"`
	Scope     string    `json:"scope"`
	Time      time.Time `json:"time"`
}

type Token struct {
	ID       int64     `json:"id"`
	Prefix   string    `json:"prefix"`
	Username string    `json:"username"`
	C2       string    `json:"c2"`
	Token    string    `json:"token"`
	Created  time.Time `json:"created"`
	Expires  time.Time `json:"expires"`
}
