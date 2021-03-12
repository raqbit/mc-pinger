package main

import (
	"fmt"
	"log"

	mcpinger "github.com/Raqbit/mc-pinger"
)

func main() {
	// Ping Minecraft server
	info, err := mcpinger.Ping("mc.hypixel.net", "")

	if err != nil {
		log.Println(err)
		return
	}

	// Print server info
	fmt.Printf("Description: \"%s\"\n", info.Description.Text)
	fmt.Printf("Online: %d/%d\n", info.Players.Online, info.Players.Max)
	fmt.Printf("Version: %s\n", info.Version.Name)
}
