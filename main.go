package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type LeaderboardForm struct {
	Name string `form:"name" binding:"required"`
	Mode string `form:"mode" binding:"required"`
}

func main() {
	// Load environment.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database.
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Database connected.")
	}
	defer dbpool.Close()

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
	c.IndentedJSON(http.StatusOK, "Leaderboards")
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
