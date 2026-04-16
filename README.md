# Remote-Restarter 

## Description
This is a simple web server that can be used to remote restart some apps on a windows 10 HTPC. It was mainly created to restart HyperHDR Direct Screen Capture which sometimes glitches out when i connect to my Sunshine server fullscreen with Moonlight so i trigger the endpoint from my phone. 

## Features
- Restart an application with a filepath and check if it is running with a wait time
- Run a command with arguments

## Installation
Get the latest release from the [releases](https://github.com/mgenc2077/Go-Remote-Restarter/releases) page.

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
  - appid: k3d-stop
    command:
      - "wsl k3d cluster stop mg-homelab"
  - appid: k3d-start
    command:
      - "wsl k3d cluster start mg-homelab"
```

> [!IMPORTANT]
> Dont forget to add a firewall rule that allows requests to the server.

Run the executable (which starts the web server in the system tray)
```powershell
.\remote-restarter.exe
```

## Endpoints

### GET /
Returns a simple embedded html with all the apps and restart links in the configuration file.

### GET /health
Returns 200 and OK if server is running.

### GET /restart/$appid
Restarts or runs the command associated with $appid from the config.

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
