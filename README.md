# Remote-Restarter 

## Description
This is a simple app that i use to remote restart some apps in my windows 10 HTPC. I wanted to make it mainly to restart HyperHDR Screen Capture which sometimes glitches out when i connect to my Sunshine server fullscreen with Moonlight so i just trigger the endpoint from my phone. 

## Installation
For now build the application 

```powershell
go build -o remote-restarter.exe
```

Create\Modify config file in the same directory (RestartConfigs.yml)
```yaml
serverconfig:
  port: "4949"

apps:
  - appid: hyperhdr
    filepath: "C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\HyperHDR\\HyperHDR.lnk"
    waittime: 2
  - appid: notepad++
    filepath: "C:\\Program Files\\Notepad++\\notepad++.exe"
    waittime: 1
```

> [!IMPORTANT]
> Dont forget to add a firewall rule that allows requests to the server.

run it (which starts the web server)
```powershell
.\remote-restarter.exe
```

In the future i am going to create an nssm config to run it as a windows service for auto-start and maybe even a actions workflow for ease of installation.

## Endpoints

### GET /
Returns a simple embedded html with all the apps and restart links in the configuration file.

### GET /health
Returns 200 and OK if server is running.

### GET /restart/$appid
Restarts $appid from the config.

## Home Assistant Configuration
At the end of your Home Assistant's configuration.yaml add the following:

```yaml
rest_command:
  restart_hyperhdr:
    url: "http://<your-server-ip>:4949/restart/hyperhdr"
    method: GET
```

Create a automation:
```yaml
alias: Restart HyperHDR Button
description: ""
triggers:
  - entity_id: button.restart_hyperhdr
    trigger: state
actions:
  - action: rest_command.restart_hyperhdr
mode: single
```

Add the button to your dashboard:

```yaml
views:
  - title: Home
    cards:
      - type: button
        name: Restart HyperHDR
        icon: mdi:restart
        tap_action:
          action: call-service
          service: rest_command.restart_hyperhdr
```
