package apiwg

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/google/uuid"
)

var (
	APIKeyName = "name-value"
)

func SessionAPIWG() *apigateway.APIGateway {
	mySession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a APIGateway client from just a session.
	// svc := apigateway.New(mySession)

	// Create a APIGateway client with additional configuration
	apiwg := apigateway.New(mySession, aws.NewConfig().WithRegion("ap-southeast-1"))
	return apiwg
}

func CreateAPIKey() {
	var (
		err error
	)
	apiwg := SessionAPIWG()
	key := uuid.NewString()
	out, err := apiwg.CreateApiKey(&apigateway.CreateApiKeyInput{
		Value: aws.String(key),
		Name:  aws.String(APIKeyName),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}

func DeleteAPIKey() {
	var (
		err error
	)
	apiwg := SessionAPIWG()
	out, err := apiwg.DeleteApiKey(&apigateway.DeleteApiKeyInput{
		ApiKey: aws.String(APIKeyName),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
