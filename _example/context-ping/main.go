package main

import (
	"context"
	"fmt"
	mcpinger "github.com/Raqbit/mc-pinger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		// NOW cancel
		fmt.Println("Aborting due to signal...")
		cancel()
	}()

	// Ping Minecraft server with context
	info, err := mcpinger.PingContext(ctx, "mc.herobone.de:25565")

	if err != nil {
		log.Println(err)
		return
	}

	// Print server info
	fmt.Printf("Description: \"%s\"\n", info.Description.Text)
	fmt.Printf("Online: %d/%d\n", info.Players.Online, info.Players.Max)
	fmt.Printf("Version: %s\n", info.Version.Name)
}
