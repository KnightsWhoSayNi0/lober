package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Event struct {
	Name    string    `json:"name"`
	Command string    `json:"command"`
	C2      string    `json:"c2"`
	Scope   string    `json:"scope"`
	Time    time.Time `json:"time"`
}

func main() {
	// impl env lookup
	dsn := "postgres://lober:lober@localhost:5432/lober?sslmode=disable"
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		rows, _ := db.Query("select * from events")
		c.JSON(200, gin.H{
			"message": rows,
		})
	})

	router.POST("/events", newEvent)

	router.Run()
}

func newEvent(c *gin.Context) {
	var event Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.Time = time.Now()
	c.JSON(http.StatusOK, event)
}
