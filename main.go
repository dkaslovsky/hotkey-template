package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

// Application appName
const appName = "hotkeys"

// Version is set with ldflags
var version string

func main() {
	flags := appFlags{}
	attachFlags(&flags)

	if flags.versionAndExit {
		displayVersion()
		return
	}

	// Construct hotkey triggers from configuration file
	if flags.configFile == "" {
		log.Fatal("must provide path to key configuration file")
	}
	configReader, err := os.Open(flags.configFile)
	if err != nil {
		log.Fatalf("failed to open key configuration file %s: %v", flags.configFile, err)
	}
	hotkeyTriggers, err := parseKeyConfig(configReader)
	if err != nil {
		log.Fatalf("failed to parse key configuration: %v", err)
	}

	// Run a listener for each hotkey trigger
	mainthread.Init(func() {
		wg := &sync.WaitGroup{}
		for _, h := range hotkeyTriggers {
			wg.Add(1)
			go listener(wg, h)
		}
		wg.Wait()
		log.Printf("%s exiting", appName)
	})
}

type hotkeyTrigger struct {
	*hotkey.Hotkey
	name        string
	commandName string
	commandArgs []string
}

func (h *hotkeyTrigger) String() string {
	return fmt.Sprintf("'%s' (%v)", h.name, h.Hotkey)
}

func listener(wg *sync.WaitGroup, trigger *hotkeyTrigger) {
	defer wg.Done()

	if err := trigger.Register(); err != nil {
		log.Printf("hotkey %s failed to register: %v", trigger, err)
		return
	}
	log.Printf("hotkey %s is registered\n", trigger)

	defer func() {
		err := trigger.Unregister()
		if err != nil {
			log.Printf("hotkey %s failed to unregister: %v\n", trigger, err)
		} else {
			log.Printf("hotkey %s is unregistered\n", trigger)
		}
		log.Printf("listener for hotkey %s exiting", trigger)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigChan:
			log.Printf("listener for hotkey %s exiting (received signal: %v)\n", trigger, sig)
			return
		case <-trigger.Keyup():
			cmd := exec.Command(trigger.commandName, trigger.commandArgs...) // #nosec G204
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Printf("error executing command for hotkey %s: %v\n", trigger, err)
				return
			}
		}
	}
}

func displayVersion() {
	vStr := "%s: version %s\n"
	if version == "" {
		fmt.Printf(vStr, appName, "(development)")
		return
	}
	fmt.Printf(vStr, appName, version)
}
