# go-cmd

Control application by console.

## Install

```console
go get github.com/dainiauskas/go-cmd
```

## Use

```go
import "github.com/dainiauskas/go-cmd"
```

This gives you access to the go-cmd package.

```go
comm := cmd.New(appName, appDescription, appVersion, appBuild)

// Adding app version show
comm.AddVersion()

// Load configuration from yaml. Configuration name is app_name.yaml
comm.LoadConfig(config)

// Adding service control
comm.AddService(fun() {
  fmt.Println("Servic start")
})

// Start
comm.Execute()
```
