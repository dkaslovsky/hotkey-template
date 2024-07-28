# hotkeys
Trigger arbitrary commands with global hotkeys on macOS

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
