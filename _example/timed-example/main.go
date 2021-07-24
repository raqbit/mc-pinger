package main

import (
	"fmt"
	"log"
	"time"

	mcpinger "github.com/Raqbit/mc-pinger"
)

func main() {
	// Create new Pinger with 10 seconds Timeout
	pinger := mcpinger.New("mc.herobone.de", 25565, mcpinger.WithTimeout(10*time.Second))

	// Get server info
	info, err := pinger.Ping()

	if err != nil {
		log.Println(err)
		return
	}

	// Print server info
	fmt.Printf("Description: \"%s\"\n", info.Description.Text)
	fmt.Printf("Online: %d/%d\n", info.Players.Online, info.Players.Max)
	fmt.Printf("Version: %s\n", info.Version.Name)
}
