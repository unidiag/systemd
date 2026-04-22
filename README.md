# systemd (Go package)

Simple helper for managing systemd services in Go.

## Features

- Create service
- Enable / start / restart
- Works with root and user mode
- Zero boilerplate

## Installation

```bash
go get github.com/unidiag/systemd
```

## Usage
```go
svc, _ := systemd.NewFromCurrentBinary("myapp")
svc.InstallAndStart()
```

```
svc := &systemd.Service{
	Name:        "epgserver",
	Description: "EPG Server",
	ExecStart:   "/opt/epg/epgserver",
	WorkingDir:  "/opt/epg",
	Restart:     "always",
	UserMode:    false,
}

err := svc.InstallAndStart()
```

## License
MIT License
Copyright (c) 2026 unidiag# systemd
