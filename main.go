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
	router.GET("/api/events", func(c *gin.Context) { c.JSON(http.StatusOK, getEventsSlice(c)) })
	router.GET("/api/users", func(c *gin.Context) { c.JSON(http.StatusOK, getUsersSlice(c)) })
	router.GET("/api/teams", func(c *gin.Context) { c.JSON(http.StatusOK, getTeamsSlice(c)) })
	router.GET("/api/c2s", func(c *gin.Context) { c.JSON(http.StatusOK, getC2sSlice(c)) })
	router.GET("/api/scope", func(c *gin.Context) { c.JSON(http.StatusOK, getScopeSlice(c)) })

	router.GET("/api/users/:name", func(c *gin.Context) { c.JSON(http.StatusOK, getUser(c, c.Param("name"))) })
	router.GET("/api/teams/:name", func(c *gin.Context) { c.JSON(http.StatusOK, getTeam(c, c.Param("name"))) })
	router.GET("/api/c2s/:name", func(c *gin.Context) { c.JSON(http.StatusOK, getC2(c, c.Param("name"))) })
	router.GET("api/scope/:name", func(c *gin.Context) { c.JSON(http.StatusOK, getScope(c, c.Param("name"))) })

	router.POST("/api/events", newEvent)
	router.POST("/api/users", newUser)
	router.POST("/api/teams", newTeam)
	router.POST("/api/c2s", newC2)
	router.POST("/api/scope", newScope)

	router.Run()
}

func getEventsSlice(c *gin.Context) []Event {
	q := `
select events.command, users.username, c2s.name, scope.name, events.time 
from events inner join users on events.user_id=users.id 
inner join c2s on events.c2_id=c2s.id inner join scope on events.scope_id=scope.id 
order by events.time desc;
`
	rows, err := db.Query(q)
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
	q := `
select users.id, users.username, users.email, teams.name 
from users inner join teams on users.team_id=teams.id;`
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Team); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		}

		users = append(users, u)
	}
	return users
}

func getUser(c *gin.Context, username string) User {
	var user User
	q := `
select users.id, users.username, users.email, teams.name 
from users inner join teams on users.team_id=teams.id 
where users.username = $1;`
	err := db.QueryRow(q, username).Scan(&user.ID, &user.Username, &user.Email, &user.Team)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	return user
}

func getTeamsSlice(c *gin.Context) []Team {
	q := `
select teams.name, teams.color, users.username 
from teams inner join users on teams.lead_id=users.id;
`
	rows, err := db.Query(q)
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

func getTeam(c *gin.Context, name string) Team {
	var t Team
	q := `
select teams.name, teams.color, users.username
from teams inner join users on teams.lead_id=users.id 
where teams.name = $1;`
	err := db.QueryRow(q, name).Scan(&t.Name, &t.Color, &t.Lead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	return t
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

func getC2(c *gin.Context, name string) C2 {
	var c2 C2
	q := `
select c2s.name from c2s
where c2s.name = $1;`
	err := db.QueryRow(q, name).Scan(&c2.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	return c2
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

func getScope(c *gin.Context, name string) Scope {
	var s Scope
	q := `
select scope.name from scope
where scope.name = $1;`
	err := db.QueryRow(q, name).Scan(&s.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	return s
}

func getUserId(username string, ptr *int) error {
	return db.QueryRow("select id from users where username = $1", username).Scan(ptr)
}
func getC2Id(c2 string, ptr *int) error {
	return db.QueryRow("select id from c2s where name = $1", c2).Scan(ptr)
}
func getScopeId(scope string, ptr *int) error {
	return db.QueryRow("select id from scope where name = $1", scope).Scan(ptr)
}
func getTeamId(team string, ptr *int) error {
	return db.QueryRow("select id from teams where name = $1", team).Scan(ptr)
}

func newEvent(c *gin.Context) {
	var event Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, c2Id, scopeId := 0, 0, 0
	err = getUserId(event.User, &userId)
	err = getC2Id(event.C2, &c2Id)
	err = getScopeId(event.Scope, &scopeId)
	event.Time = time.Now()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

func newUser(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamID := 0
	if err := getTeamId(user.Team, &teamID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(pq.CopyIn("users", "username", "email", "team_id"))
	stmt.Exec(user.Username, user.Email, teamID)
	stmt.Exec()
	stmt.Close()
	tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, user)
}

func newTeam(c *gin.Context) {
	var team Team

	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leadID := 0
	if err := getUserId(team.Lead, &leadID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(pq.CopyIn("teams", "name", "color", "lead_id"))
	stmt.Exec(team.Name, team.Color, leadID)
	stmt.Exec()
	stmt.Close()
	tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, team)
}

func newC2(c *gin.Context) {
	var c2 C2

	if err := c.ShouldBindJSON(&c2); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(pq.CopyIn("c2s", "name"))
	stmt.Exec(c2.Name)
	stmt.Exec()
	stmt.Close()
	tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, c2)
}

func newScope(c *gin.Context) {
	var scope Scope

	if err := c.ShouldBindJSON(&scope); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(pq.CopyIn("scopes", "name"))
	stmt.Exec(scope.Name)
	stmt.Exec()
	stmt.Close()
	tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, scope)
}
