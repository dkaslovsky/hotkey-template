# hotkey-template
Template for triggering an arbitrary commands with a global hotkey


### plist
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

```
launchctl load /Library/LaunchAgents/com.dkas.hotkeys.plist
launchctl start com.dkas.hotkeys
```
