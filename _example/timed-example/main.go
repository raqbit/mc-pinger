package main

import (
	"fmt"
	"log"
	"time"

	mcpinger "github.com/Raqbit/mc-pinger"
)

func main() {
	// Ping Minecraft server with 10 seconds timeout
	info, err := mcpinger.PingTimeout("mc.hypixel.net", "", 10*time.Second)

	if err != nil {
		log.Println(err)
		return
	}

	// Print server info
	fmt.Printf("Description: \"%s\"\n", info.Description.Text)
	fmt.Printf("Online: %d/%d\n", info.Players.Online, info.Players.Max)
	fmt.Printf("Version: %s\n", info.Version.Name)
}
