package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Event struct {
	ID      int64
	Command string    `json:"command"`
	User    string    `json:"user"`
	C2      string    `json:"c2"`
	Scope   string    `json:"scope"`
	Time    time.Time `json:"time"`
}

var db *sql.DB
var err error

func init() {
	// impl env lookup
	dsn := "postgres://lober:lober@db:5432/lober?sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to postgres")
}

func main() {
	router := gin.Default()

	router.GET("/events", getEvents)
	router.POST("/events", newEvent)

	router.Run()
}

func getEvents(c *gin.Context) {
	rows, err := db.Query("select * from events")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	events := make([]Event, 0)

	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Command, &e.User, &e.C2, &e.Scope, &e.Time); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func newEvent(c *gin.Context) {
	var event Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.Time = time.Now()

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(pq.CopyIn("events", "command", "user_id", "c2_id", "scope_id", "time"))
	stmt.Exec(event.Command, 1, 2, 3, event.Time)
	stmt.Exec()
	stmt.Close()
	tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, event)
}
