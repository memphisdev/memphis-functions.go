<a href="![Github (4)](https://github.com/memphisdev/memphis-terraform/assets/107035359/a5fe5d0f-22e1-4445-957d-5ce4464e61b1)">![Github (4)](https://github.com/memphisdev/memphis-terraform/assets/107035359/a5fe5d0f-22e1-4445-957d-5ce4464e61b1)</a>
<p align="center">
<a href="https://memphis.dev/discord"><img src="https://img.shields.io/discord/963333392844328961?color=6557ff&label=discord" alt="Discord"></a>
<a href="https://github.com/memphisdev/memphis/issues?q=is%3Aissue+is%3Aclosed"><img src="https://img.shields.io/github/issues-closed/memphisdev/memphis?color=6557ff"></a> 
  <img src="https://img.shields.io/npm/dw/memphis-dev?color=ffc633&label=installations">
<a href="https://github.com/memphisdev/memphis/blob/master/CODE_OF_CONDUCT.md"><img src="https://img.shields.io/badge/Code%20of%20Conduct-v1.0-ff69b4.svg?color=ffc633" alt="Code Of Conduct"></a> 
<img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/memphisdev/memphis?color=61dfc6">
<img src="https://img.shields.io/github/last-commit/memphisdev/memphis?color=61dfc6&label=last%20commit">
</p>

<div align="center">
  
  <img width="200" alt="CNCF Silver Member" src="https://github.com/cncf/artwork/raw/master/other/cncf-member/silver/color/cncf-member-silver-color.svg#gh-light-mode-only">
  <img width="200" alt="CNCF Silver Member" src="https://github.com/cncf/artwork/raw/master/other/cncf-member/silver/white/cncf-member-silver-white.svg#gh-dark-mode-only">
  
</div>
 
 <b><p align="center">
  <a href="https://memphis.dev/pricing/">Cloud</a> - <a href="https://memphis.dev/docs/">Docs</a> - <a href="https://twitter.com/Memphis_Dev">X</a> - <a href="https://www.youtube.com/channel/UCVdMDLCSxXOqtgrBaRUHKKg">YouTube</a>
</p></b>

<div align="center">

  <h4>

**[Memphis.dev](https://memphis.dev)** is a highly scalable, painless, and effortless data streaming platform.<br>
Made to enable developers and data teams to collaborate and build<br>
real-time and streaming apps fast.

  </h4>
  
</div>

# Installation
After installing and running memphis broker,<br>
In your project's directory:

```shell
go get github.com/memphisdev/memphis-functions.go
```

# Importing
```go
import "github.com/memphisdev/memphis-functions.go"
```

### Creating a Memphis function
Memphis provides a `CreateFunction` utility for more easily creating Memphis Functions.

The user will write a function which will act as an event handler and will be called for every event that is to be processed by Functions. 

The user created event handler must fulfill the following function signature:
```go
type EventHandlerFunction func([]byte, map[string]string, map[string]string) ([]byte, map[string]string, error)
```

The event handler will take in a `[]byte` representation of an event, also a `map[string]string` of headers that belong to that event, and `map[string]string` representation of the function inputs.

The event handler will then return a modified version of these fields.

If the processing fails, the user function should return `nil, nil, err`. If the user wishes to skip over an event and not send it to the station, return `nil, nil, nil` and the event will be skipped. Events that return an error will be sent to the dead letter station. 

```go
package main

import (
	"encoding/json"
    "github.com/memphisdev/memphis.go"
)

type Event struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func eventHandlerFunc(msgPayload[]byte, msgHeaders[string]string, inputs[string]string) ([]byte, map[string]string, error){
    // Get data from msgPayload
    var event Event
    json.Unmarshal(msgPayload, &event)
    
    // Modify or do something with the payload
    event.Field1 = "modified"
    
    // Return the payload back as []bytes
    eventBytes, _ := json.Marshal(event)
    return eventBytes, msgHeaders, nil
}

func main() {
	memphis.CreateFunction(eventHandlerFunc);
}
```

As mentioned previously, if the user would like to send the message to the dead letter station, simply return an error. The unproccessed payload and headers will be included with the message to the dead letter station.

```go
package main

import (
	"encoding/json"
    "strings"
    "errors"
    "github.com/memphisdev/memphis.go"
)

type Event struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func eventHandlerFunc(msgPayload[]byte, msgHeaders[string]string, inputs[string]string) ([]byte, map[string]string, error){
    // Get data from msgPayload
    var event Event
    json.Unmarshal(msgPayload, &event)
    
    // Modify or do something with the payload
    if strings.Contains(event.Field1, "Bob"){
        return nil, nil, errors.New("String had Bob in it!")
    } 
    
    // Return the payload back as []bytes
    eventBytes, _ := json.Marshal(event)
    return eventBytes, msgHeaders, nil
}

func main() {
	memphis.CreateFunction(eventHandlerFunc);
}
```

If the user would rather this message just be skipped, instead of being sent to the dead letter station, return `nil` for all values:

```go
package main

import (
	"encoding/json"
    "strings"
    "github.com/memphisdev/memphis.go"
)

type Event struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func eventHandlerFunc(msgPayload[]byte, msgHeaders[string]string, inputs[string]string) ([]byte, map[string]string, error){
    // Get data from msgPayload
    var event Event
    json.Unmarshal(msgPayload, &event)
    
    // Modify or do something with the payload
    if strings.Contains(event.Field1, "Bob"){
        return nil, nil, nil
    } 
    
    // Return the payload back as []bytes
    eventBytes, _ := json.Marshal(event)
    return eventBytes, msgHeaders, nil
}

func main() {
	memphis.CreateFunction(eventHandlerFunc);
}
```

Lastly, if the user is using a format like Protocol Buffers, the user may simple decode the msgPayload with their proto function. Assuming that we have a proto definition like: 
```proto
syntax = "proto3";
package protobuf_example;

message Message{
    string data_field = 1;
}
```
That is saved in the directory `user_message/message.pb.go`, we can use this Message object like so:

```go
package main

import (
	"encoding/json"
    "strings"
    "github.com/memphisdev/memphis.go"
    "google.golang.org/protobuf/proto"
    "current_directory/user_message"
)

func eventHandlerFunc(msgPayload[]byte, msgHeaders[string]string, inputs[string]string) ([]byte, map[string]string, error){
    // Get data from msgPayload
    var my_message user_message.Message
    proto.Unmarshal(msgPayload, &user_message)
    
    // Modify or do something with the payload
    my_message.data_field = "new_data"
    
    // Return the payload back as []bytes
    eventBytes, _ := proto.Marshal(&my_message)
    return eventBytes, msgHeaders, nil
}

func main() {
	memphis.CreateFunction(eventHandlerFunc);
}
```