package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mcpinger "github.com/Raqbit/mc-pinger"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		// NOW cancel
		fmt.Println("aborting due to interrupt...")
		cancel()
	}()

	// Create new Pinger with 10 seconds Timeout
	pinger := mcpinger.New("mc.herobone.de", 25565, mcpinger.WithContext(ctx))

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
