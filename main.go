package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

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

	go statusCheck()

	router := gin.Default()
	router.GET("/wake", wake)
	router.GET("/status", status)

	router.Run("0.0.0.0:8080")
}

// A mutex is used to protect it from concurrent access issues.
var (
	server_up     bool
	pinger_status *probing.Statistics
	mu            sync.Mutex
)

// var pinger_status = make(chan *probing.Statistics)

func statusCheck() {
	// var up = false
	pinger, err := probing.NewPinger(os.Getenv("DESKTOP_IP"))
	if err != nil {
		panic(err)
	}
	pinger.Interval = 1 * time.Minute
	// pinger.Count = 100
	pinger.OnSend = func(p *probing.Packet) {
		mu.Lock()
		server_up = false
		pinger_status = pinger.Statistics()
		mu.Unlock()
	}
	pinger.OnRecv = func(pkt *probing.Packet) {
		mu.Lock()
		server_up = true
		pinger_status = pinger.Statistics()
		mu.Unlock()
	}
	// pinger.OnRecvError = func(err error) {
	// 	mu.Lock()
	// 	server_up = false
	// 	// fmt.Println(err)
	// 	pinger_status = pinger.Statistics()
	// 	mu.Unlock()
	// }
	pinger.Run() // blocks
}

func wake(c *gin.Context) {
	if packet, err := NewMagicPacket(os.Getenv("DESKTOP_MAC")); err == nil {
		packet.Send(os.Getenv("DESKTOP_WOL_BROADCAST")) // send to broadcast
		// packet.Send("192.168.100.255") // send to broadcast
		// packet.SendPort("192.168.100.255", "9") // specify receiving port
	}
	c.Status(http.StatusNoContent)
	// c.JSON(http.StatusNoContent, nil)
}

func status(c *gin.Context) {
	// Acquire the lock to safely read the shared variable.
	mu.Lock()
	stats := pinger_status
	up := server_up
	mu.Unlock()

	// var up = stats.PacketLoss < 0.5
	c.JSON(http.StatusOK, gin.H{
		"up":               up,
		"packet_loss":      stats.PacketLoss,
		"packets_sent":     stats.PacketsSent,
		"packets_received": stats.PacketsRecv,
	})
}
