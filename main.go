package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
	h   *Hub
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func init() {
	// Use environment variables for database connection
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	dbname := os.Getenv("PGDATABASE")
	sslmode := os.Getenv("PGSSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=lober",
		user, password, host, port, dbname, sslmode)

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to postgres")

	// setup default user if needed
	count := 0
	err = db.QueryRow("select count(*) from users").Scan(&count)
	if count == 0 {
		// add default team
		teamID := 0
		err := db.QueryRow("insert into teams (name) values($1) returning id;", "default").Scan(&teamID)
		if err != nil {
			log.Fatal(err)
		}

		pHash, err := hashPassword("admin")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Query("insert into users (username, password, team_id) values ($1, $2, $3);", "admin", pHash, teamID)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// @title Lober API
// @description Centralized Logging Server for Red Teams
// @BasePath /api
func main() {
	router := gin.Default()

	router.Use(tokenAuthMiddleware())

	// websocket
	h = newHub()
	go h.run()
	router.GET("/api/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		h.register <- conn
	})

	// api endpoints
	router.GET("/api/events", func(c *gin.Context) { c.JSON(http.StatusOK, getEventsSlice(c)) })
	router.GET("/api/users", func(c *gin.Context) { c.JSON(http.StatusOK, getUsersSlice(c)) })
	router.GET("/api/teams", func(c *gin.Context) { c.JSON(http.StatusOK, getTeamsSlice(c)) })
	router.GET("/api/c2s", func(c *gin.Context) { c.JSON(http.StatusOK, getC2sSlice(c)) })
	router.GET("/api/scope", func(c *gin.Context) { c.JSON(http.StatusOK, getScopeSlice(c)) })
	router.GET("/api/tokens", func(c *gin.Context) { c.JSON(http.StatusOK, getTokensSlice(c)) })
	router.GET("/api/metrics", getMetrics)
	router.GET("/api/export/events", exportEventsCSV)

	router.GET("/api/users/:name", func(c *gin.Context) { c.JSON(http.StatusOK, getUser(c, c.Param("name"))) })
	router.GET("/api/teams/:name", func(c *gin.Context) { c.JSON(http.StatusOK, getTeam(c, c.Param("name"))) })
	router.GET("/api/c2s/:name", func(c *gin.Context) { c.JSON(http.StatusOK, getC2(c, c.Param("name"))) })
	router.GET("api/scope/:name", func(c *gin.Context) { c.JSON(http.StatusOK, getScope(c, c.Param("name"))) })
	router.GET("/api/tokens/:prefix", func(c *gin.Context) { c.JSON(http.StatusOK, getToken(c, c.Param("prefix"))) })

	router.POST("/api/events", newEvent)
	router.POST("/api/users", newUser)
	router.POST("/api/teams", newTeam)
	router.POST("/api/c2s", newC2)
	router.POST("/api/scope", newScope)
	router.POST("/api/tokens", newToken)

	router.DELETE("/api/users/:name", removeUser)
	router.DELETE("/api/teams/:name", removeTeam)
	router.DELETE("/api/c2s/:name", removeC2)
	router.DELETE("/api/scope/:name", removeScope)
	router.DELETE("/api/tokens/:prefix", removeToken)

	err := router.Run()
	if err != nil {
		return
	}
}

// ID util
func getID(table string, column string, value string) (int64, error) {
	var id int64
	safeTable := pq.QuoteIdentifier(table)
	safeColumn := pq.QuoteIdentifier(column)
	q := fmt.Sprintf("select id from %s where %s = $1", safeTable, safeColumn)

	err := db.QueryRow(q, value).Scan(&id)
	return id, err
}

// GET slices and individual

func getEventsSlice(c *gin.Context) []Event {
	filter := c.Query("filter")
	team := c.Query("team")
	user := c.Query("user")
	c2 := c.Query("c2")
	scope := c.Query("scope")
	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	q := `
select events.command, users.username, teams.color, c2s.name, scope.name, events.time
from events inner join users on events.user_id=users.id
inner join c2s on events.c2_id=c2s.id inner join scope on events.scope_id=scope.id
inner join teams on users.team_id=teams.id
where 1=1
`
	var args []interface{}
	argCount := 1

	if filter != "" {
		q += fmt.Sprintf(" AND (events.command ILIKE $%d OR users.username ILIKE $%d OR c2s.name ILIKE $%d OR scope.name ILIKE $%d)", argCount, argCount, argCount, argCount)
		args = append(args, "%"+filter+"%")
		argCount++
	}
	if team != "" {
		q += fmt.Sprintf(" AND teams.name = $%d", argCount)
		args = append(args, team)
		argCount++
	}
	if user != "" {
		q += fmt.Sprintf(" AND users.username = $%d", argCount)
		args = append(args, user)
		argCount++
	}
	if c2 != "" {
		q += fmt.Sprintf(" AND c2s.name = $%d", argCount)
		args = append(args, c2)
		argCount++
	}
	if scope != "" {
		q += fmt.Sprintf(" AND scope.name = $%d", argCount)
		args = append(args, scope)
		argCount++
	}

	q += fmt.Sprintf(" order by events.time desc limit $%d offset $%d;", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := db.Query(q, args...)
	if err != nil {
		log.Printf("DB error in getEventsSlice: %v", err)
		return []Event{}
	}
	defer rows.Close()

	events := make([]Event, 0)
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.Command, &e.User, &e.TeamColor, &e.C2, &e.Scope, &e.Time); err != nil {
			log.Printf("Scan error in getEventsSlice: %v", err)
			return []Event{}
		}
		events = append(events, e)
	}
	return events
}

func getEventsWithFilter(filter string) []Event {
	q := `
select events.command, users.username, teams.color, c2s.name, scope.name, events.time
from events inner join users on events.user_id=users.id
inner join c2s on events.c2_id=c2s.id inner join scope on events.scope_id=scope.id
inner join teams on users.team_id=teams.id
where events.command ILIKE $1 OR users.username ILIKE $1 OR c2s.name ILIKE $1 OR scope.name ILIKE $1
order by events.time desc;
`
	rows, err := db.Query(q, "%"+filter+"%")
	if err != nil {
		log.Printf("DB error in getEventsWithFilter: %v", err)
		return []Event{}
	}
	defer rows.Close()

	events := make([]Event, 0)
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.Command, &e.User, &e.TeamColor, &e.C2, &e.Scope, &e.Time); err != nil {
			log.Printf("Scan error in getEventsWithFilter: %v", err)
			return []Event{}
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
select teams.name, teams.color, COALESCE(array_agg(users.username) FILTER (WHERE users.username IS NOT NULL), '{}') as users
from teams 
left join users on teams.id = users.team_id
group by teams.id, teams.name, teams.color;
`
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	teams := make([]Team, 0)
	for rows.Next() {
		var t Team
		var users pq.StringArray
		if err := rows.Scan(&t.Name, &t.Color, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		}
		t.Users = users
		teams = append(teams, t)
	}
	return teams
}

func getTeam(c *gin.Context, name string) Team {
	var t Team
	q := `
select teams.name, teams.color, COALESCE(array_agg(users.username) FILTER (WHERE users.username IS NOT NULL), '{}') as users
from teams 
left join users on teams.id = users.team_id
where teams.name = $1
group by teams.id, teams.name, teams.color;`
	var users pq.StringArray
	err := db.QueryRow(q, name).Scan(&t.Name, &t.Color, &users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	t.Users = users
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

func getTokensSlice(c *gin.Context) []Token {
	q := `
select tokens.prefix, users.username, c2s.name, tokens.created, tokens.expires
from tokens inner join users on tokens.user_id=users.id
inner join c2s on tokens.c2_id=c2s.id;
`
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	tokens := make([]Token, 0)
	for rows.Next() {
		var t Token
		if err := rows.Scan(&t.Prefix, &t.Username, &t.C2, &t.Created, &t.Expires); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		}

		tokens = append(tokens, t)
	}
	return tokens
}

func getToken(c *gin.Context, prefix string) Token {
	var t Token
	q := `
select tokens.prefix, users.username, c2s.name, tokens.created, tokens.expires
from tokens inner join users on tokens.user_id=users.id
inner join c2s on tokens.c2_id=c2s.id
where tokens.prefix = $1;
`
	err := db.QueryRow(q, prefix).Scan(&t.ID, &t.Prefix, &t.Username, &t.C2, &t.Created, &t.Expires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	return t
}

func getMetrics(c *gin.Context) {
	timeRange := c.DefaultQuery("range", "24 hours")
	metrics := make(map[string]interface{})

	// Validate range to prevent SQL injection
	validRanges := map[string]bool{
		"1 hour":   true,
		"6 hours":  true,
		"24 hours": true,
		"7 days":   true,
		"30 days":  true,
	}
	if !validRanges[timeRange] {
		timeRange = "24 hours"
	}

	// Dynamic truncation based on range
	trunc := "hour"
	if timeRange == "7 days" || timeRange == "30 days" {
		trunc = "day"
	} else if timeRange == "1 hour" {
		trunc = "minute"
	}

	// Events per hour/day for range
	q := fmt.Sprintf(`
		select date_trunc('%s', time) as period, count(*) 
		from events 
		where time > now() - interval '%s'
		group by period order by period;
	`, trunc, timeRange)

	rows, _ := db.Query(q)
	var timeline []map[string]interface{}
	for rows != nil && rows.Next() {
		var t time.Time
		var count int
		rows.Scan(&t, &count)
		timeline = append(timeline, map[string]interface{}{"time": t, "count": count})
	}
	metrics["timeline"] = timeline

	// Group counts for range
	groupQueries := map[string]string{
		"by_team": `
			select teams.name, count(*) 
			from events inner join users on events.user_id=users.id
			inner join teams on users.team_id=teams.id
			where events.time > now() - interval '%s'
			group by teams.name order by count(*) desc;`,
		"by_scope": `
			select scope.name, count(*) 
			from events inner join scope on events.scope_id=scope.id
			where events.time > now() - interval '%s'
			group by scope.name order by count(*) desc;`,
		"by_c2": `
			select c2s.name, count(*) 
			from events inner join c2s on events.c2_id=c2s.id
			where events.time > now() - interval '%s'
			group by c2s.name order by count(*) desc;`,
		"by_user": `
			select users.username, count(*) 
			from events inner join users on events.user_id=users.id
			where events.time > now() - interval '%s'
			group by users.username order by count(*) desc;`,
	}

	for key, query := range groupQueries {
		rows, _ = db.Query(fmt.Sprintf(query, timeRange))
		var list []map[string]interface{}
		for rows != nil && rows.Next() {
			var name string
			var count int
			rows.Scan(&name, &count)
			list = append(list, map[string]interface{}{"label": name, "count": count})
		}
		metrics[key] = list
	}

	c.JSON(http.StatusOK, metrics)
}

func exportEventsCSV(c *gin.Context) {
	q := `
select users.username, c2s.name, scope.name, events.command, events.time
from events inner join users on events.user_id=users.id
inner join c2s on events.c2_id=c2s.id inner join scope on events.scope_id=scope.id
order by events.time desc;
`
	rows, err := db.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=lober_events.csv")
	c.Header("Content-Type", "text/csv")

	// headers
	c.Writer.Write([]byte("User,C2,Scope,Command,Time\n"))

	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.User, &e.C2, &e.Scope, &e.Command, &e.Time); err != nil {
			log.Printf("export error: %v", err)
			continue
		}

		// simple CSV escaping for commas & quotes
		cmd := strings.ReplaceAll(e.Command, "\"", "\"\"")
		line := fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
			e.User, e.C2, e.Scope, cmd, e.Time.Format(time.RFC3339))
		c.Writer.Write([]byte(line))
	}
}

// NEW

func newEvent(c *gin.Context) {
	var event Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use token identities if not defined in JSON
	if event.User == "" {
		if username, exists := c.Get("username"); exists {
			event.User = username.(string)
		}
	}
	if event.C2 == "" {
		if c2, exists := c.Get("c2"); exists {
			event.C2 = c2.(string)
		}
	}

	userID, err := getID("users", "username", event.User)
	c2ID, err := getID("c2s", "name", event.C2)
	scopeID, err := getID("scope", "name", event.Scope)
	event.Time = time.Now()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch team color
	err = db.QueryRow("select teams.color from users inner join teams on users.team_id=teams.id where users.id = $1", userID).Scan(&event.TeamColor)
	if err != nil {
		log.Printf("Error fetching team color for event: %v", err)
	}

	_, err = db.Exec("insert into events (command, user_id, c2_id, scope_id, time) values ($1, $2, $3, $4, $5)",
		event.Command, userID, c2ID, scopeID, event.Time)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// broadcast new event
	eventJSON, _ := json.Marshal(event)
	h.broadcastData(eventJSON)

	c.JSON(http.StatusOK, event)
}

func newTeam(c *gin.Context) {
	var team Team

	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec("insert into teams (name, color) values ($1, $2)", team.Name, team.Color)
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

	_, err = db.Exec("insert into c2s (name) values ($1)", c2.Name)
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

	_, err = db.Exec("insert into scope (name) values ($1)", scope.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, scope)
	}
}

// REMOVE

func removeUser(c *gin.Context) {
	name := c.Param("name")

	if name == "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin is required"})
	}

	_, err := db.Exec("delete from users where username = $1", name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func removeTeam(c *gin.Context) {
	name := c.Param("name")
	_, err := db.Exec("delete from teams where name = $1", name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func removeC2(c *gin.Context) {
	name := c.Param("name")
	_, err := db.Exec("delete from c2s where name = $1", name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func removeScope(c *gin.Context) {
	name := c.Param("name")
	_, err := db.Exec("delete from scope where name = $1", name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func removeToken(c *gin.Context) {
	prefix := c.Param("prefix")
	_, err := db.Exec("delete from tokens where prefix = $1", prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
