# Fastbase Logger Client
> Allows you to send a log to the fastbase logger microservice

--- 
## :rocket: How to use

#### Install
```go
    go get gitlab.devops.ukfast.co.uk/fastbase/logger-client
```

#### Use

```go
var LoggerClient = loggerclient.NewClient("log_set_name", "log_type_name")
LoggerClient.Log(map[string]interface{}{"message": msg, "error": err})
```