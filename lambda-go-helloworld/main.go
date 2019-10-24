package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda/messages"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type MyFunction struct {
	lambda.Function
}


func (fn *MyFunction) Invoke(req *messages.InvokeRequest, response *messages.InvokeResponse) error {
	defer func() {
		if err := recover(); err != nil {
			panicInfo := getPanicInfo(err)
			response.Error = &messages.InvokeResponse_Error{
				Message:    panicInfo.Message,
				Type:       getErrorType(err),
				StackTrace: panicInfo.StackTrace,
				ShouldExit: true,
			}
		}
	}()

	deadline := time.Unix(req.Deadline.Seconds, req.Deadline.Nanos).UTC()
	invokeContext, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	lc := &lambdacontext.LambdaContext{
		AwsRequestID:       req.RequestId,
		InvokedFunctionArn: req.InvokedFunctionArn,
		Identity: lambdacontext.CognitoIdentity{
			CognitoIdentityID:     req.CognitoIdentityId,
			CognitoIdentityPoolID: req.CognitoIdentityPoolId,
		},
	}
	if len(req.ClientContext) > 0 {
		if err := json.Unmarshal(req.ClientContext, &lc.ClientContext); err != nil {
			response.Error = lambdaErrorResponse(err)
			return nil
		}
	}
	invokeContext = lambdacontext.NewContext(invokeContext, lc)

	invokeContext = context.WithValue(invokeContext, "x-amzn-trace-id", req.XAmznTraceId)

	payload, err := fn.handler.Invoke(invokeContext, req.Payload)
	if err != nil {
		response.Error = lambdaErrorResponse(err)
		return nil
	}
	response.Payload = payload
	return nil
}

func Handler(ctx context.Context) {
	mlc, _ := lambdacontext.FromContext(ctx)
	log.Printf("lc: %#v\n", mlc)
	log.Printf("Custom: %#v\n", mlc.ClientContext.Custom)
	log.Printf("Hoge: %#v\n", mlc.ClientContext.Custom["hoge"])
}

func main() {
	lambda.Start(Handler)
}
