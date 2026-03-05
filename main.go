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
		x, y, z := 0, 0, 0
		if err := rows.Scan(&e.ID, &e.Command, &x, &y, &z, &e.Time); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		getUsername(x, &e.User)
		getC2Name(y, &e.C2)
		getScopeName(z, &e.Scope)
		fmt.Printf("%s ran by %s using %s on %s at %s\n", e.Command, e.User, e.C2, e.Scope, e.Time)
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

	userId, c2Id, scopeId := 0, 0, 0
	getUserId(event.User, &userId)
	getC2Id(event.C2, &c2Id)
	getScopeId(event.Scope, &scopeId)
	event.Time = time.Now()

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(pq.CopyIn("events", "command", "user_id", "c2_id", "scope_id", "time"))
	stmt.Exec(event.Command, userId, c2Id, scopeId, event.Time)
	stmt.Exec()
	stmt.Close()
	tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, event)
}

func getUserId(username string, ptr *int) {
	db.QueryRow("select id from users where username = $1", username).Scan(ptr)
}
func getC2Id(c2 string, ptr *int) {
	db.QueryRow("select id from c2s where name = $1", c2).Scan(ptr)
}
func getScopeId(scope string, ptr *int) {
	db.QueryRow("select id from scope where name = $1", scope).Scan(ptr)
}

func getUsername(id int, ptr *string) {
	db.QueryRow("select username from users where id = $1", id).Scan(ptr)
}
func getC2Name(id int, ptr *string) {
	db.QueryRow("select name from c2s where id = $1", id).Scan(ptr)
}
func getScopeName(id int, ptr *string) {
	db.QueryRow("select name from scope where id = $1", id).Scan(ptr)
}
