package cloudamqp

import (
	"log"
	"math/rand"
	"time"
)

func randomSleep(ms int, name string) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5000)
	log.Printf("[DEBUG] %s sleep for %d ms...\n", name, n)
	time.Sleep(time.Duration(n) * time.Millisecond)
}
