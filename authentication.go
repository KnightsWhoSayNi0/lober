package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

func newUser(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if user already exists
	exists := false
	err := db.QueryRow("select exists(select 1 from users where username = $1);", user.Username).Scan(&exists)

	if exists {
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
		return
	}

	_, err = db.Exec("insert into users (username, password, team_id) values ($1, $2, $3)", user.Username, hash, teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, user)
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

	_, err = db.Exec("insert into tokens (token, user_id, c2_id, created, expires) values ($1, $2, $3, $4, $5)",
		token.Token, userID, c2ID, token.Created, token.Expires)

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
