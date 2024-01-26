# dogger
Dogger is a simple wrapper for Zerolog that synchronizes logs using Go channels. It also saves logs to a file using a rolling file appender (lumberjack).


## Installation
```bash
go get github.com/douglaspetrin/dogger@latest
```

## Environment variables

````
CORRELATION_KEY: The key used to identify the correlation id in the context. Default: "corrId"
LOG_LEVEL: The log level. Default: 0 -> debug
SERVICE_NAME: The service name. Default: nil
SERVICE_ENV: The service environment. "development" uses zerolog.ConsoleWriter as well for colored logs.
LOG_MAX_SIZE: The max size of the log file in megabytes. Default: 100
LOG_MAX_BACKUPS: The max number of log files to keep. Default: 0
LOG_MAX_AGE: The max age of the log file in days. Default: 0
** If LOG_MAX_BACKUPS and LOG_MAX_AGE are both 0, no old log files will be deleted. **
LOG_COMPRESS: Whether or not to compress the log files. Default: true
USING_GIT_REVISION: In case you want to add git revision field to the logging message. Default: false
USING_GO_VERSION: In case you want to add go version field to the logging message. Default: false
USING_PID: In case you want to add pid field to the logging message. Default: false
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

## Output
```json
{"level":"info","app":"my-app-name","corrId":"my-corrId","event":"my-event-name","data":{"my-key":"my-value"},"time":"2024-01-26T03:33:43.843873865-08:00"}
```
