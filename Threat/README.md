# ThreatPlaybook Client

> version 3.0.0b1

## Commands

### Install Client

#### Windows

> ThreatPlaybook has not been tested on win* OS

* Download the binary exe from github to a path of your choice

### Configure 

```bash

playbook configure --help

```

> This creates a `.cred` file in the current directory with some config parameters

### Change Default Password first

```bash

playbook change-password --help

```
> If you've set a default password, you need to change that with this command before you can do anything on ThreatPlaybook

### Login

```bash

playbook login --help

```

### Add a Feature/User Story

```bash

playbook apply feature -f <path-to-yaml-file>

```

## TODO Features
* Get Feature/Abuser Story/Threat Scenario from API - Done
* Delete features - Done
* Export Features
* Vulnerability Management Push and Sync Features