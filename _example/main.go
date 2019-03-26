package main

import (
	"fmt"
	"github.com/Raqbit/mc-pinger"
	"log"
)

func main() {
	// Create new Pinger
	pinger := mcpinger.New("localhost", 25565)

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
