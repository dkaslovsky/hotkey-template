# hotkeys
Trigger arbitrary commands with global hotkeys on macOS

<br/>

### Quick Start
1. Download the `hotkeys` binary:
   
    ARM:
    ```
    curl -o hotkeys -L https://github.com/dkaslovsky/hotkeys/releases/latest/download/hotkeys_darwin_arm64
    ```
    AMD:
    ```
    curl -o hotkeys -L https://github.com/dkaslovsky/hotkeys/releases/latest/download/hotkeys_darwin_amd64
    ```
2. Ensure the binary is executable:
   ```
   chmod +x hotkeys
   ```
3. Make sure it is in a directory that is included in `$PATH`
4. Create a configuration file for your hotkeys and associated commands (see example provided below)
5. Create a `.plist` file (see the example provided below) that includes the proper path to the binary and the configuration file. Place this file in `/Library/LaunchAgents/` (using `sudo` if necessary)
7. Use `launchctl` to load the .plist file and start the service (see the example commands below)

<br/>

### plist
Add a plist specifying the location of the binary and configuration files to `/Library/LaunchAgents/`
```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.dkas.hotkeys</string>

    <key>UserName</key>
    <string>dkas</string>

    <key>ProgramArguments</key>
    <array>
        <string>/Users/dkas/bin/hotkeys</string>
        <string>-file</string>
        <string>/Users/dkas/.config/hotkeys/hotkeys</string>
    </array>
</dict>
</plist>
```

<br/>

### launchctl
Load the plist and start the service
```
launchctl load /Library/LaunchAgents/com.dkas.hotkeys.plist
launchctl start com.dkas.hotkeys
```
Be sure to restart the service to pick up any changes to the configuration file
```
launchctl stop com.dkas.hotkeys
sleep 2
launchctl start com.dkas.hotkeys
```

<br/>

### Configuration
Bind commands to hotkeys in a configuration file
```
[
  {
    "name":"Brave Browser - open new window for default profile (Main)",
    "command_name":"open",
    "command_args":[
      "-na", "Brave Browser", "--args", "--new-window", "--profile-directory=Default"
    ],
    "key":"KeyN",
    "key_modifiers":[
      "ModCtrl",
      "ModOption",
      "ModCmd"
    ]
  },
  {
    "name":"Brave Browser - open new window for profile 5",
    "command_name":"open",
    "command_args":[
      "-na", "Brave Browser", "--args", "--new-window", "--profile-directory=Profile 5"
    ],
    "key":"KeyM",
    "key_modifiers":[
      "ModCtrl",
      "ModOption",
      "ModCmd"
    ]
  }
]
```

Keys and modifiers are identified in the configuration file by using the `Key` prefix for a letter/number (e.g., use `KeyN` for the `N` key) and the `Mod` prefix for a modifier (e.g., use `ModCtrl` for the control modifier key).
