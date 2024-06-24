package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// type Goal struct {
// 	goal_id    int
// 	game_id    int
// 	goal_desc  string
// 	goal_pos_x int
// 	goal_pos_y int
// }

type User struct {
	Name  string
	Score int
}

type LeaderboardForm struct {
	Name string `form:"name" binding:"required"`
	Mode string `form:"mode" binding:"required"`
}

var dbpool *pgxpool.Pool

func main() {
	// Load environment.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database.
	var dbpool_err error
	dbpool, dbpool_err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if dbpool_err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// Example of a database query.
	// var game_id int
	// err = dbpool.QueryRow(context.Background(), "SELECT game_id FROM game WHERE game_name = 'hard'").Scan(&game_id)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(game_id)

	router := gin.Default()
	router.GET("/:difficulty/:goal", getGoals)
	// router.POST("/:difficulty/:goal", postGoals)
	router.GET("/leaderboards", getLeaderboards)
	router.POST("/leaderboards", postLeaderboards)

	router.Run("localhost:8080")
}

// Get all goals.
func getGoals(c *gin.Context) {
	difficulty := c.Param("difficulty")
	goal := c.Param("goal")
	c.IndentedJSON(http.StatusOK, difficulty+" "+goal)
}

func getLeaderboards(c *gin.Context) {
	var users []User
	rows, err := dbpool.Query(context.Background(), `SELECT user_name, user_score FROM "user" INNER JOIN game ON "user".game_id = game.game_id`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	// Iterate through all rows returned from the query.
	for rows.Next() {
		// Loop through rows, using Scan to assign column data to struct fields.
		var user User
		rows.Scan(&user.Name, &user.Score)
		users = append(users, user)
	}

	// Convert struct to JSON.
	usersJson, err := json.Marshal(users)
	// fmt.Println(string(usersJson))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}
	c.IndentedJSON(http.StatusOK, usersJson)
}

// Post new score to leaderboards.
func postLeaderboards(c *gin.Context) {
	var form LeaderboardForm
	if err := c.BindJSON(&form); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}
	// c.String(200, "Hello %s", form.Name)
	c.IndentedJSON(http.StatusOK, "Leaderboard posted")
}
