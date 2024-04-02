package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

// Configuration: set the command and the hotkey trigger
var (
	commandName = "open"
	commandArgs = []string{"-na", "Brave Browser", "--args", "--new-window", "--profile-directory=Default"}

	keyModifier = []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModOption, hotkey.ModCmd}
	key         = hotkey.KeySpace
)

func main() { mainthread.Init(fn) }

func fn() {
	hk := hotkey.New(keyModifier, key)
	if err := hk.Register(); err != nil {
		log.Printf("hotkey failed to register: %v", err)
		return
	}
	log.Printf("hotkey %v is registered\n", hk)

	defer func() {
		if err := hk.Unregister(); err != nil {
			log.Printf("hotkey %v failed to unregister: %v\n", hk, err)
			return
		}
		log.Printf("hotkey %v is unregistered\n", hk)
		log.Printf("exiting")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigChan:
			log.Printf("exiting (signal: %v)\n", sig)
			return
		case <-hk.Keyup():
			cmd := exec.Command(commandName, commandArgs...) // #nosec G204
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Printf("error executing command: %v\n", err)
				return
			}
		}
	}
}
