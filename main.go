package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	probing "github.com/prometheus-community/pro-bing"
)

func main() {
	// load dotenv file
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Error loading .env file")
		// panic(err)
	}

	router := gin.Default()
	router.GET("/wake", wake)
	router.GET("/status", status)

	router.Run("0.0.0.0:8080")
}

func wake(c *gin.Context) {
	if packet, err := NewMagicPacket(os.Getenv("DESKTOP_MAC")); err == nil {
		packet.Send("192.168.100.255")          // send to broadcast
		// packet.SendPort("192.168.100.255", "9") // specify receiving port
	}
	c.Status(http.StatusNoContent)
	// c.JSON(http.StatusNoContent, nil)
}

func status(c *gin.Context) {
	pinger, err := probing.NewPinger(os.Getenv("DESKTOP_IP"))
	if err != nil {
		panic(err)
	}
	pinger.Count = 4
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats
	var up = stats.PacketLoss < 0.5
	c.JSON(http.StatusOK, gin.H{
		"up":          up,
		"packet_loss": stats.PacketLoss,
	})
}
