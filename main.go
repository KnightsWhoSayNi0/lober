package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

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

	// api endpoints
	router.GET("/events", getEvents)
	router.POST("/events", newEvent)

	router.Run()
}

func getEvents(c *gin.Context) {
	c.JSON(http.StatusOK, getEventsSlice(c))
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

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"events": getEventsSlice(c),
	})
}

func usersHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "users.html", gin.H{
		"users": getUsersSlice(c),
	})
}

func teamsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "teams.html", gin.H{
		"teams": getTeamsSlice(c),
	})
}

func c2sHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "c2s.html", gin.H{
		"c2s": getC2sSlice(c),
	})
}

func scopeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "scope.html", gin.H{
		"scope": getScopeSlice(c),
	})
}

func configHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "config.html", gin.H{})
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

func getEventsSlice(c *gin.Context) []Event {
	rows, err := db.Query("select events.command, users.username, c2s.name, scope.name, events.time from events inner join users on events.user_id=users.id inner join c2s on events.c2_id=c2s.id inner join scope on events.scope_id=scope.id order by events.time desc;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	events := make([]Event, 0)
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.Command, &e.User, &e.C2, &e.Scope, &e.Time); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		}
		events = append(events, e)
	}
	return events
}

func getUsersSlice(c *gin.Context) []User {
	rows, err := db.Query("select users.username, users.email, teams.name from users inner join teams on users.team_id=teams.id;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Username, &u.Email, &u.Team); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		}

		users = append(users, u)
	}
	return users
}

func getTeamsSlice(c *gin.Context) []Team {
	rows, err := db.Query("select teams.name, teams.color, users.username from teams inner join users on teams.lead_id=users.id;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	teams := make([]Team, 0)
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.Name, &t.Color, &t.Lead); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		}

		teams = append(teams, t)
	}
	return teams
}

func getC2sSlice(c *gin.Context) []C2 {
	rows, err := db.Query("select c2s.name from c2s;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	c2s := make([]C2, 0)
	for rows.Next() {
		var c2 C2
		if err := rows.Scan(&c2.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		}

		c2s = append(c2s, c2)
	}
	return c2s
}

func getScopeSlice(c *gin.Context) []Scope {
	rows, err := db.Query("select scope.name from scope;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	scope := make([]Scope, 0)
	for rows.Next() {
		var s Scope
		if err := rows.Scan(&s.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		}

		scope = append(scope, s)
	}
	return scope
}
