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

	doneChan := make(chan struct{}, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for s := range sigChan {
			log.Printf("signal handler got signal: %v\n", s)
			doneChan <- struct{}{}
		}
		log.Printf("signal handler done")
	}()

hotkeyListener:
	for {
		select {
		case <-doneChan:
			log.Print("exiting main loop due to interrupt signal\n")
			break hotkeyListener
		case <-hk.Keyup():
			cmd := exec.Command(commandName, commandArgs...) // #nosec G204
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Printf("error executing command: %v\n", err)
				break hotkeyListener
			}
		}
	}

	close(sigChan)

	if err := hk.Unregister(); err != nil {
		log.Printf("hotkey %v failed to unregister: %v\n", hk, err)
		return
	}
	log.Printf("hotkey %v is unregistered\n", hk)

	log.Printf("exiting")
}
