## Deployment of Golang app to AWS's Lambda

1. export AWS_PROFILE=admin
2. GOARCH=amd64 GOOS=linux go build main.go
3. zip -r messages.zip .

## Create
aws lambda create-function --function-name messages --zip-file fileb://messages.zip --handler main --runtime go1.x --role "arn:aws:iam::832685173872:role/lambda-basic-execution"

## Update
aws lambda update-function-code --function-name messages --zip-file fileb://messages.zip

## Invoke
aws lambda invoke --function-name PingMe-messages --invocation-type "RequestResponse" response.txt