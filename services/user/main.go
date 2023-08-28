package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/model/devices"
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
		&people.Blocked{},
		&messaging.Conversation{},
		&messaging.ConversationLike{},
		&devices.Device{},
		&huddles.Huddle{},
		&huddles.HuddleComment{},
		&huddles.HuddleCommentLike{},
		&huddles.HuddleLike{},
		&people.Invite{},
		&huddles.LastSeenHuddle{},
		&messaging.LastSeenMessage{},
		&messaging.Message{},
		&messaging.MessageReaction{},
		&people.MutedConversation{},
		&messaging.PersonInConversation{},
		&users.User{},
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
