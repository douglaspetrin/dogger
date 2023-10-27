# dogger
Dogger is a simple wrapper for Zerolog that synchronizes logs using Go channels.


## Installation
```bash
go get github.com/douglaspetrin/dogger@latest
```

## Environment variables

````
CORRELATION_KEY: The key used to identify the correlation id in the context. Default: "correlation-id"
LOG_LEVEL: The log level. Default: 0 -> debug
SERVICE_NAME: The service name. Default: nil
SERVICE_ENV: The service environment. "development" uses zerolog.ConsoleWriter as well for colored logs.
LOG_MAX_SIZE: The max size of the log file in megabytes. Default: 100
USING_GIT_REVISION: In case you want to add git revision field to the logging message. Default: false
USING_GO_VERSION: In case you want to add go version field to the logging message. Default: false
````

## Usage
```go
package main

import (
    "github.com/douglaspetrin/dogger"
)   

func main() {
	
	if err := myFunc(); err != nil {
		dogger.LogError("my-correlation-id", "Error when calling myFunc", nil, err)
    }
	
	var data = map[string]any{"foo": "bar"}
	
	dogger.LogInfo("my-correlation-id", "My event message", data)
}
``` 
