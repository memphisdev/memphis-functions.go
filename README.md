[![Github (6)](https://github.com/memphisdev/memphis/assets/107035359/bc2feafc-946c-4569-ab8d-836bc0181890)](https://www.functions.memphis.dev/)
<p align="center">
<a href="https://memphis.dev/discord"><img src="https://img.shields.io/discord/963333392844328961?color=6557ff&label=discord" alt="Discord"></a>
<a href="https://github.com/memphisdev/memphis/issues?q=is%3Aissue+is%3Aclosed"><img src="https://img.shields.io/github/issues-closed/memphisdev/memphis?color=6557ff"></a> 
  <img src="https://img.shields.io/npm/dw/memphis-dev?color=ffc633&label=installations">
<a href="https://github.com/memphisdev/memphis/blob/master/CODE_OF_CONDUCT.md"><img src="https://img.shields.io/badge/Code%20of%20Conduct-v1.0-ff69b4.svg?color=ffc633" alt="Code Of Conduct"></a> 
<img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/memphisdev/memphis?color=61dfc6">
<img src="https://img.shields.io/github/last-commit/memphisdev/memphis?color=61dfc6&label=last%20commit">
</p>

<div align="center">
  
<img width="177" alt="cloud_native 2 (5)" src="https://github.com/memphisdev/memphis/assets/107035359/a20ea11c-d509-42bb-a46c-e388c8424101">
  
</div>
 
 <b><p align="center">
  <a href="https://memphis.dev/pricing/">Cloud</a> - <a href="https://memphis.dev/docs/">Docs</a> - <a href="https://twitter.com/Memphis_Dev">X</a> - <a href="https://www.youtube.com/channel/UCVdMDLCSxXOqtgrBaRUHKKg">YouTube</a>
</p></b>

<div align="center">

  <h4>

**[Memphis.dev](https://memphis.dev)** is more than a broker. It's a new streaming stack.<br>
Memphis.dev is a highly scalable event streaming and processing engine.<br>

  </h4>
  
</div>

## ![20](https://user-images.githubusercontent.com/70286779/220196529-abb958d2-5c58-4c33-b5e0-40f5446515ad.png) About

Before Memphis came along, handling ingestion and processing of events on a large scale took months to adopt and was a capability reserved for the top 20% of mega-companies. Now, Memphis opens the door for the other 80% to unleash their event and data streaming superpowers quickly, easily, and with great cost-effectiveness.

This repository is responsible for the Memphis Functions Golang SDK

# Installation
After installing and running memphis broker,<br>
In your project's directory:

```shell
go get github.com/memphisdev/memphis-functions.go
```

# Importing
```go
import "github.com/memphisdev/memphis-functions.go/memphis"
```

### Creating a Memphis function
Memphis provides a `CreateFunction` utility for more easily creating Memphis Functions.

The user will write a function which will act as an event handler and will be called for every event that is to be processed by Functions. 

The user created event handler must fulfill the following function signature:
```go
type HandlerType func(any, map[string]string, map[string]string) (any, map[string]string, error)

```

The `any` first parameter will always be a pointer to the message payload. 

The event handler will take in a `[]byte` or `*Object` representation of an event. It will also take a `map[string]string` of headers that belong to that event, and `map[string]string` representation of the function inputs.

The event handler will then return a modified version of these fields. The return must be either a `[]byte` or a type that will be able to be used with `json.Marshal`.

If the processing fails, the user function should return `nil, nil, err`. If the user wishes to skip over an event and not send it to the station, return `nil, nil, nil` and the event will be skipped. Events that return an error will be sent to the dead letter station. 

```go
package main

import (
	"encoding/json"
    "github.com/memphisdev/memphis-functions.go/memphis"
)

type Event struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func eventHandlerFunc(msgPayload any, msgHeaders map[string]string, inputs map[string]string) (any, map[string]string, error){
    // Type assert to user type. By default this is a []byte
    as_bytes, ok := msgPayload.([]byte)
	if !ok{
		return nil, nil, fmt.Errorf("object failed type assertion: %v, %v", message, reflect.TypeOf(message))
	}
    
    // Get data from msgPayload
    var event Event
    json.Unmarshal(as_bytes, &event)
    
    // Modify or do something with the payload
    event.Field1 = "modified"
    
    // Return the modified payload
    return event, msgHeaders, nil
}

func main() {
	memphis.CreateFunction(eventHandlerFunc)
}
```

To use a user specified object, create an empty one in main and pass that to the `memphis.PayloadAsJSON` function.

```go
package main

import (
	"encoding/json"
    "github.com/memphisdev/memphis-functions.go/memphis"
)

type Event struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func eventHandlerFunc(msgPayload any, msgHeaders map[string]string, inputs map[string]string) (any, map[string]string, error){
    // Get data from msgPayload
    typedPayload = msgPayload.(*Event)

    // Modify or do something with the payload
    typedPayload.Field1 = "modified"
    
    return typedPayload, msgHeaders, nil
}

func main() {
    var eventObject Event
	memphis.CreateFunction(eventHandlerFunc, memphis.PayloadAsJSON(&eventObject))
}
```

> Note the type assertion is using a pointer to the object, and the address of the object is passed into `memphis.PayloadAsJSON`.

As mentioned previously, if the user would like to send the message to the dead letter station, simply return an error. The unproccessed payload and headers will be included with the message to the dead letter station.

```go
package main

import (
	"encoding/json"
    "strings"
    "errors"
    "github.com/memphisdev/memphis-functions.go/memphis"
)

type Event struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func eventHandlerFunc(msgPayload any, msgHeaders map[string]string, inputs map[string]string) ([]byte, map[string]string, error){
    // Type assert to user type. By default this is a []byte
    as_bytes, ok := msgPayload.([]byte)
	if !ok{
		return nil, nil, fmt.Errorf("object failed type assertion: %v, %v", message, reflect.TypeOf(message))
	}

    // Get data from msgPayload
    var event Event
    json.Unmarshal(as_bytes, &event)
    
    // Modify or do something with the payload
    if strings.Contains(event.Field1, "Bob"){
        return nil, nil, fmt.Errorf("String had Bob in it!")
    } 
    
    // Return the modified payload
    return event, msgHeaders, nil
}

func main() {
	memphis.CreateFunction(eventHandlerFunc)
}
```

If the user would rather this message just be skipped, instead of being sent to the dead letter station, return `nil` for all values:

```go
package main

import (
	"encoding/json"
    "strings"
    "github.com/memphisdev/memphis-functions.go/memphis"
)

type Event struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func eventHandlerFunc(msgPayload any, msgHeaders map[string]string, inputs map[string]string) (any, map[string]string, error){
    // Type assert to user type. By default this is a []byte
    as_bytes, ok := msgPayload.([]byte)
	if !ok{
		return nil, nil, fmt.Errorf("object failed type assertion: %v, %v", message, reflect.TypeOf(message))
	}

    // Get data from msgPayload
    var event Event
    json.Unmarshal(as_bytes, &event)
    
    // Modify or do something with the payload
    if strings.Contains(event.Field1, "Bob"){
        return nil, nil, nil
    } 
    
    // Return the modified payload
    return event, msgHeaders, nil
}

func main() {
	memphis.CreateFunction(eventHandlerFunc)
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
    "github.com/memphisdev/memphis-functions.go/memphis"
    "google.golang.org/protobuf/proto"
    "current_directory/user_message"
)

func eventHandlerFunc(msgPayload any, msgHeaders map[string]string, inputs map[string]string) (any, map[string]string, error){
    // Type assert to user type. By default this is a []byte
    as_bytes, ok := msgPayload.([]byte)
	if !ok{
		return nil, nil, fmt.Errorf("object failed type assertion: %v, %v", message, reflect.TypeOf(message))
	}

    // Get data from msgPayload
    var my_message user_message.Message
    proto.Unmarshal(as_bytes, &user_message)
    
    // Modify or do something with the payload
    my_message.data_field = "new_data"
    
    // Return the payload back as []bytes so we don't json.Marshal the message
    eventBytes, _ := proto.Marshal(&my_message)
    return eventBytes, msgHeaders, nil
}

func main() {
	memphis.CreateFunction(eventHandlerFunc)
}
```
