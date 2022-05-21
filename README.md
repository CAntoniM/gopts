# gopts

    A Go options parsing library

## About

gopts goals are are to provide a simple minimal code command line arguments parser based on CLAP arguments parser for rust. To offer
this simplicity we allow you to parse command line args via annotating either annotating a struct or via using the builder api

## Examples 

### Tag API

This example uses the tag api to parse some simple command line arguments. Notice how for the argument to be added the struct field has to be exported otherwise it will be ignored.

```go
package main

import (
    "fmt"
	"os"

	"github.com/CAntoniM/spack/gopts"
)

type Cli struct {
	Config_file string `gopts:"optional,desc=path to configuration file.,"`
	Address     string `gopts:"name=addr,flag,desc=The address the service will listen on"`
	Port        uint16 `gopts:"name=port,flag,desc=The port the service will listen on."`
	Verbose     bool   `gopts:"name=verbose,flag,desc=Enables verbose logging."`
}

func main() {
	var opts Cli = Cli{"spackd.json", "0.0.0.0", 15000, false}

	gopts.Parse(&opts, os.Args, "The simple, self hosted package service")
	fmt.Printf("file: %s, Address: %s, port: %d, verbose: %t \n", opts.Config_file, opts.Address, opts.Port, opts.Verbose)
}

```

### Builder API

The alternative API that is available for parsing your arguments is the builder api this api requires more work on the part of the user however it is more performant. However for larger applications this should be less of a concern as this should be treated as a one time event but what your needs are in terms of performance and ease of use/maintainability should be considered on a application by application basis

```go

package main 

import {
    "os"
    "fmt"

    "github.com/CAntoniM/spack/gopts"
}

func main() {
    _, optional, flags := gopts.Opts("The simple, self-hosted package service").
          Optional("Config_file","the path to the json configuration file.").
          Flag("port","The port the application will listen on").
          Flag("addr","The address the application will listen on").
          Short_flag("verbose","enables verbose logging")

        for key, value := range optional {
            fmt.Printf("Optional value %s=%s",key,value)
        }
        for key, value := range required {
            fmt.Printf("Required value %s=%s",key,value)
        }
        for key, value := range flags {
              fmt.Printf("flag value %s=%s",key,value)
          }
}

```

## Future Development

1. Implement did you mean functionality
2. Add a subcommand category that allows for different options based on the main command specified
