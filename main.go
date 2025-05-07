package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"systementor.se/cloudgolangapi/data"
)

var config Config
var theRandom *rand.Rand

func start(c *gin.Context) {
	c.Data(http.StatusOK, "application/text", []byte("Tjena"))
}

func enableCors(c *gin.Context) {
	(*c).Header("Access-Control-Allow-Origin", "*")
}

func apiStats(c *gin.Context) {
	enableCors(c)
	totalGames, wins := data.Stats()
	c.JSON(http.StatusOK, gin.H{"totalGames": totalGames, "wins": wins})
}

func apiPlay(c *gin.Context) {
	enableCors(c)
	yourSelection := c.Query("yourSelection")
	mySelection := randomizeSelection()
	winner := "Tie"
	if yourSelection == "STONE" && mySelection == "SCISSOR" {
		winner = "You"
	}
	if yourSelection == "SCISSOR" && mySelection == "BAG" {
		winner = "You"
	}
	if yourSelection == "BAG" && mySelection == "STONE" {
		winner = "You"
	}
	if mySelection == "STONE" && yourSelection == "SCISSOR" {
		winner = "Computer"
	}
	if mySelection == "SCISSOR" && yourSelection == "BAG" {
		winner = "Computer"
	}
	if mySelection == "BAG" && yourSelection == "STONE" {
		winner = "Computer"
	}
	data.SaveGame(yourSelection, mySelection, winner)
	c.JSON(http.StatusOK, gin.H{"winner": winner, "yourSelection": yourSelection, "computerSelection": mySelection})
}

func randomizeSelection() string {
	val := theRandom.Intn(3) + 1
	if val == 1 {
		return "STONE"
	}
	if val == 2 {
		return "SCISSOR"
	}
	if val == 3 {
		return "BAG"
	}
	return "ERROR"

}

func main() {
	theRandom = rand.New(rand.NewSource(time.Now().UnixNano()))
	readConfig(&config)

	data.InitDatabase(config.Database.File,
		config.Database.Server,
		config.Database.Database,
		config.Database.Username,
		config.Database.Password,
		config.Database.Port)

	router := gin.Default()
	router.GET("/", start)
	router.GET("/api/play", apiPlay)
	router.GET("/api/stats", apiStats)
	// router.GET("/api/employee/:id", apiEmployeeById)
	// router.PUT("/api/employee/:id", apiEmployeeUpdateById)
	// router.DELETE("/api/employee/:id", apiEmployeeDeleteById)
	// router.POST("/api/employee", apiEmployeeAdd)

	// router.GET("/api/employees", employeesJson)
	// router.GET("/api/addemployee", addEmployee)
	// router.GET("/api/addmanyemployees", addManyEmployees)
	router.Run(":8080")

}
