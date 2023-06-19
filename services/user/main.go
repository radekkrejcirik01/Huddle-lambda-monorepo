package main

import (
	"context"
	"log"

	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/devices"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/huddles"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/messaging"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/people"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/users"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/rest"
)

var fiberLambda *fiberadapter.FiberLambda

func init() {
	database.Connect()
	if err := database.DB.AutoMigrate(
		&people.Invite{},
		&people.Hide{},
		&people.MutedConversation{},
		&people.MutedHuddle{},
		&huddles.HuddleInteracted{},
		&huddles.Huddle{},
		&huddles.HuddleComment{},
		&huddles.HuddleCommentLike{},
		&messaging.Conversation{},
		&messaging.Message{},
		&messaging.PersonInConversation{},
		&messaging.LastReadMessage{},
		&messaging.ConversationLike{},
		&messaging.MessageReaction{},
		&users.User{},
		&devices.Device{},
	); err != nil {
		log.Fatal(err)
	}

	fiberLambda = fiberadapter.New(rest.Create())
}

// Handler will deal with Fiber working with Lambda
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return fiberLambda.ProxyWithContext(ctx, request)
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(Handler)
}
