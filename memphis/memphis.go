package memphis

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aws/aws-lambda-go/lambda"
)

type MemphisMsg struct {
	Headers map[string]string `json:"headers"`
	Payload string            `json:"payload"`
}

type MemphisMsgWithError struct {
	Headers map[string]string `json:"headers"`
	Payload string            `json:"payload"`
	Error   string            `json:"error"`
}

type MemphisEvent struct {
	Inputs   map[string]string `json:"inputs"`
	Messages []MemphisMsg      `json:"messages"`
}

type MemphisOutput struct {
	Messages       []MemphisMsg          `json:"messages"`
	FailedMessages []MemphisMsgWithError `json:"failed_messages"`
}

// HandlerType functions get the message payload as []byte (or any), message headers as map[string]string and inputs as map[string]string and should return the modified payload and headers.
// error should be returned if the message should be considered failed and go into the dead-letter station.
// if all returned values are nil the message will be filtered out of the station.
type HandlerType func(any, map[string]string, map[string]string) (any, map[string]string, error)

type HandlerOption func(*HandlerOptions) error

type HandlerOptions struct {
	Handler           HandlerType
	UserObject        any
}

func ObjectOption(schema any) HandlerOption {
	return func(handlerOptions *HandlerOptions) error {
		handlerOptions.UserObject = schema
		return nil
	}
}

func UnmarshalIntoStruct(data []byte, userStruct any) error {
	// Unmarshal JSON data into the struct
	err := json.Unmarshal(data, userStruct)
	if err != nil {
		return err
	}

	return nil
}

func checkValidStruct(userStruct any) error {
	// Check if the provided variable is a pointer to a struct
	valueOf := reflect.ValueOf(userStruct)
	if valueOf.Kind() != reflect.Ptr || valueOf.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("input parameter must be a pointer to a struct")
	}

	return nil
}

// This function creates a Memphis function and processes events with the passed-in eventHandler function.
// eventHandler gets the message payload as []byte or as the user specified type, 
// message headers as map[string]string and inputs as map[string]string and should return the modified payload and headers.
// The modified payload type will either be the user type, or []byte depending on user requirements.
// error should be returned if the message should be considered failed and go into the dead-letter station.
// if all returned values are nil the message will be filtered out from the station.
func CreateFunction(eventHandler HandlerType, options ...HandlerOption) {
	LambdaHandler := func(ctx context.Context, event *MemphisEvent) (*MemphisOutput, error) {
		params := HandlerOptions{
			Handler:           eventHandler,
			UserObject:        nil,
		}

		for _, option := range options {
			if option != nil {
				if err := option(&params); err != nil {
					return nil, err
				}
			}
		}

		if params.UserObject != nil {
			err := checkValidStruct(params.UserObject)
			if err != nil {
				return nil, err
			}
		}

		if len(options) > 1 {
			return nil, fmt.Errorf("the user passed in too many options. Functions only supports one handler option")
		}

		var processedEvent MemphisOutput
		for _, msg := range event.Messages {
			payload, err := base64.StdEncoding.DecodeString(msg.Payload)
			if err != nil {
				processedEvent.FailedMessages = append(processedEvent.FailedMessages, MemphisMsgWithError{
					Headers: msg.Headers,
					Payload: msg.Payload,
					Error:   "couldn't decode message: " + err.Error(),
				})
				continue
			}

			var handlerInput any
			if params.UserObject != nil{
				UnmarshalIntoStruct(payload, params.UserObject)
				handlerInput = params.UserObject
			}else{
				handlerInput = payload
			}

			modifiedPayload, modifiedHeaders, err := params.Handler(handlerInput, msg.Headers, event.Inputs)
			_, ok := modifiedPayload.([]byte)

			if err == nil && !ok{
				modifiedPayload, err = json.Marshal(modifiedPayload) // err will proagate to next if
			}

			if err != nil {
				processedEvent.FailedMessages = append(processedEvent.FailedMessages, MemphisMsgWithError{
					Headers: msg.Headers,
					Payload: msg.Payload,
					Error:   err.Error(),
				})
				continue
			}

			if modifiedPayload != nil && modifiedHeaders != nil {
				modifiedPayloadStr := base64.StdEncoding.EncodeToString(modifiedPayload.([]byte))
				processedEvent.Messages = append(processedEvent.Messages, MemphisMsg{
					Headers: modifiedHeaders,
					Payload: modifiedPayloadStr,
				})
			}
		}

		return &processedEvent, nil
	}

	lambda.Start(LambdaHandler)
}