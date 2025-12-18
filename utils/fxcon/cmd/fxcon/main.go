package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lumina-io/fxcon/rcon"
	"github.com/lumina-io/fxcon/utils"
)

func main() {
	host := utils.Getenv("RCON_ADDRESS", "localhost")
	prt, _ := strconv.Atoi(utils.Getenv("RCON_PORT", "30120"))
	port := prt
	password := utils.Getenv("RCON_PASSWORD", "")

	rcon, err := rcon.New(host, port, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		// panic(err)
		return
	}
	defer rcon.Close()

	if len(os.Args) <= 1 {
	} else {
		resp, err := rcon.Send(strings.Join(os.Args[1:], " "))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send: %v\n", err)
			// panic(err)
			return
		}
		fmt.Print(utils.ColorText(resp))
	}
}
