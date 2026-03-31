package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

// @title Lober API
// @description Centralized Logging Server for Red Teams
// @host localhost:8080
// @BasePath /api
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
select users.username, teams.name
from users inner join teams on users.team_id=teams.id;`
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Username, &u.Team); err != nil {
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
select users.username, teams.name
from users inner join teams on users.team_id=teams.id
where users.username = $1;`
	err := db.QueryRow(q, username).Scan(&user.Username, &user.Team)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	return user
}

func getTeamsSlice(c *gin.Context) []Team {
	q := `
select teams.name, teams.color
from teams;
`
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	teams := make([]Team, 0)
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.Name, &t.Color); err != nil {
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
select teams.name, teams.color
from teams
where teams.name = $1;`
	err := db.QueryRow(q, name).Scan(&t.Name, &t.Color)
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

func getID(table string, column string, value string) (int64, error) {
	var id int64
	safeTable := pq.QuoteIdentifier(table)
	safeColumn := pq.QuoteIdentifier(column)
	q := fmt.Sprintf("select id from %s where %s = $1", safeTable, safeColumn)

	err := db.QueryRow(q, value).Scan(&id)
	return id, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashToken(rawToken string) string {
	bytes := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(bytes[:])
}

func newEvent(c *gin.Context) {
	var event Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getID("users", "username", event.User)
	c2ID, err := getID("c2s", "name", event.C2)
	scopeID, err := getID("scope", "name", event.Scope)
	event.Time = time.Now()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(pq.CopyIn("events", "command", "user_id", "c2_id", "scope_id", "time"))
	stmt.Exec(event.Command, userID, c2ID, scopeID, event.Time)
	stmt.Exec()
	stmt.Close()
	tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, event)
	}
}

func newUser(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if user already exists
	exists := "false"
	err := db.QueryRow("select exists(select 1 from users where username = $1);", user.Username).Scan(&exists)

	if exists == "true" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	teamID, err := getID("teams", "name", user.Team)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(pq.CopyIn("users", "username", "password", "team_id"))
	stmt.Exec(user.Username, hash, teamID)
	stmt.Exec()
	stmt.Close()
	tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func newTeam(c *gin.Context) {
	var team Team

	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(pq.CopyIn("teams", "name", "color"))
	stmt.Exec(team.Name, team.Color)
	stmt.Exec()
	stmt.Close()
	tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, team)
	}
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
	} else {
		c.JSON(http.StatusOK, c2)
	}
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
	} else {
		c.JSON(http.StatusOK, scope)
	}
}

func newToken(c *gin.Context) {
	var token Token
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userID, err := getID("users", "username", token.Username)
	c2ID, err := getID("c2s", "name", token.C2)

	if token.Created.IsZero() || token.Expires.IsZero() {
		token.Created = time.Now()
		token.Expires = time.Now().Add(time.Hour * 24 * 7)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token.Token = hashToken(uuid.NewString()[:32])

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(pq.CopyIn("tokens", "token", "user_id", "c2_id", "created", "expires"))
	stmt.Exec(token.Token, userID, c2ID, token.Created, token.Expires)
	stmt.Exec()
	stmt.Close()
	tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, token)
	}
}

func verifyToken(token string) (bool, Token) {
	var rv Token
	q := `
select tokens.user_id, tokens.c2_id, tokens.created, tokens.expires
from tokens inner join users on tokens.user_id=users.id
inner join c2s on tokens.c2_id=c2s.id
where tokens.token = $1;
`

	if len(token) < 1 {
		return false, rv
	}

	tokenHash := hashToken(token)
	err := db.QueryRow(q, tokenHash).Scan(&rv.Username, &rv.C2, &rv.Created, &rv.Expires)
	if err != nil {
		return false, rv
	}

	if rv.Expires.Before(time.Now()) {
		return false, rv
	}

	return true, rv
}
