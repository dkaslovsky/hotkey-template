package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

type hotkeyTrigger struct {
	*hotkey.Hotkey
	name        string
	commandName string
	commandArgs []string
}

func (h *hotkeyTrigger) String() string {
	return fmt.Sprintf("'%s' (%v)", h.name, h.Hotkey)
}

// configuration for each hotkey trigger
type keyConfig struct {
	Name         string   `json:"name"`
	CommandName  string   `json:"command_name"`
	CommandArgs  []string `json:"command_args"`
	Key          string   `json:"key"`
	KeyModifiers []string `json:"key_modifiers"`
}

func parseKeyConfig(r io.Reader) ([]*hotkeyTrigger, error) {
	// Read and unmarshal configuration file
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	keyConfigs := []*keyConfig{}
	err = json.Unmarshal(raw, &keyConfigs)
	if err != nil {
		return nil, err
	}

	hotkeyTriggers := []*hotkeyTrigger{}
	for _, config := range keyConfigs {
		// Map string representations of keys to key constants
		key, ok := keyMap[config.Key]
		if !ok {
			return nil, fmt.Errorf("unknown key: %s", config.Key)
		}
		modifiers := []hotkey.Modifier{}
		for _, m := range config.KeyModifiers {
			modifier, ok := modifierMap[m]
			if !ok {
				return nil, fmt.Errorf("unknown key modifier: %s", m)
			}
			modifiers = append(modifiers, modifier)
		}

		hotkeyTriggers = append(hotkeyTriggers, &hotkeyTrigger{
			Hotkey:      hotkey.New(modifiers, key),
			name:        config.Name,
			commandName: config.CommandName,
			commandArgs: config.CommandArgs,
		})
	}

	return hotkeyTriggers, nil
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "file", "", "path to key configuration file")
	flag.Parse()

	if configFile == "" {
		log.Fatal("no configuration file specified")
	}
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("failed to open key configuration file %s: %v", configFile, err)
	}
	hotkeyTriggers, err := parseKeyConfig(f)
	if err != nil {
		log.Fatalf("failed to parse key configuration: %v", err)
	}

	mainthread.Init(func() {
		wg := &sync.WaitGroup{}
		for _, hkCmd := range hotkeyTriggers {
			wg.Add(1)
			go listener(hkCmd, wg)
		}
		wg.Wait()
		log.Print("done")
	})
}

func listener(trigger *hotkeyTrigger, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := trigger.Register(); err != nil {
		log.Printf("hotkey failed to register: %v", err)
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
			log.Printf("exiting (signal: %v)\n", sig)
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
