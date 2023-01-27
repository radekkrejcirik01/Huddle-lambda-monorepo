## Deployment of Golang app to AWS's Lambda

1. export AWS_PROFILE=admin
2. GOARCH=amd64 GOOS=linux go build main.go
3. zip -r pushnotifications.zip .

## Create
aws lambda create-function --function-name pushnotifications --zip-file fileb://pushnotifications.zip --handler main --runtime go1.x --role "arn:aws:iam::832685173872:role/lambda-basic-execution"

## Update
aws lambda update-function-code --function-name pushnotifications --zip-file fileb://pushnotifications.zip

## Invoke
aws lambda invoke --function-name pushnotifications --invocation-type "RequestResponse" response.txt