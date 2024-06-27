package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Message struct {
	Message string `json:"message"`
}

type Goal struct {
	Desc  string `json:"desc"`
	Pos_X int    `json:"pos_x"`
	Pos_Y int    `json:"pos_y"`
}

type User struct {
	Difficulty string `json:"difficulty"`
	Name       string `json:"name"`
	Score      int    `json:"score"`
}

type LeaderboardForm struct {
	Name       string `form:"name" binding:"required"`
	Difficulty string `form:"difficulty" binding:"required"`
	Score      int    `form:"score" binding:"required"`
}

var dbpool *pgxpool.Pool

func main() {
	// Load dev environment.
	env := os.Getenv("ENV_NAME")
	if env != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Connect to database.
	var dbpool_err error
	dbpool, dbpool_err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if dbpool_err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", dbpool_err)
		os.Exit(1)
	}
	defer dbpool.Close()

	router := gin.Default()

	// Allow all origins.
	router.Use(cors.Default())

	router.GET("/goal", getGoal)
	router.GET("/leaderboards", getLeaderboards)
	router.POST("/leaderboards", postLeaderboards)

	router.Run(":80")
}

// Get specific goal data.
// Example query: goal?difficulty=hard&desc=Cowboy%20on%20Horse
func getGoal(c *gin.Context) {
	var goal_data Goal
	difficulty := c.Query("difficulty")
	desc := c.Query("desc")

	err := dbpool.QueryRow(context.Background(), fmt.Sprintf(
		`SELECT goal_desc, goal_pos_x, goal_pos_y 
		FROM goal 
		INNER JOIN game 
		ON goal.game_id = game.game_id 
		WHERE game_name = '%s' 
		AND goal_desc = '%s'`,
		difficulty, desc)).
		Scan(&goal_data.Desc, &goal_data.Pos_X, &goal_data.Pos_Y)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, Message{"Query failed."})
		return
	}

	c.IndentedJSON(http.StatusOK, goal_data)
}

// Get leaderboard data.
func getLeaderboards(c *gin.Context) {
	var users []User
	difficulty := c.Query("difficulty")

	rows, err := dbpool.Query(context.Background(), fmt.Sprintf(
		`SELECT game_name, user_name, user_score 
		FROM "user" 
		INNER JOIN game 
		ON "user".game_id = game.game_id 
		WHERE game_name = '%s'
		ORDER BY user_score ASC
		LIMIT 15`, difficulty))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		return
	}

	// Releases any resources held by the rows no matter how the function returns.
	// Looping all the way through the rows also closes it implicitly,
	// but it is better to use defer to make sure rows is closed no matter what.
	defer rows.Close()

	// Iterate through all rows returned from the query.
	for rows.Next() {
		// Loop through rows, using Scan to assign column data to struct fields.
		var user User
		if err := rows.Scan(&user.Difficulty, &user.Name, &user.Score); err != nil {
			c.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("Query failed: %v", err))
			return
		}
		users = append(users, user)
	}

	// Check if there were any issues when reading rows.
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("Error reading queries: %v", err))
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

// Post new score to leaderboards.
func postLeaderboards(c *gin.Context) {
	var form LeaderboardForm

	// Bind JSON fields to form variable.
	if err := c.BindJSON(&form); err != nil {
		c.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	// Begin transaction.
	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	// Rollback transaction if it doesn't commit successfully.
	defer tx.Rollback(context.Background())

	// Execute insert statement.
	_, err = tx.Exec(context.Background(), fmt.Sprintf(
		`INSERT INTO "user" (game_id, user_name, user_score)
		VALUES (
			(SELECT game_id FROM game WHERE game_name = '%s'),
			'%s',
			%d
		)`,
		form.Difficulty, form.Name, form.Score))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	// Commit transaction.
	err = tx.Commit(context.Background())
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	c.IndentedJSON(http.StatusOK, "Leaderboard posted")
}
