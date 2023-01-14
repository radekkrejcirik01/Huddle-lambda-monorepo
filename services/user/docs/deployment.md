## Deployment of Golang app to AWS's Lambda

1. export AWS_PROFILE=admin
2. GOARCH=amd64 GOOS=linux go build main.go
3. zip -r user.zip .

## Create
aws lambda create-function --function-name PingMe-user --zip-file fileb://user.zip --handler main --runtime go1.x --role "arn:aws:iam::409186456204:role/lambda-basic-execution"

## Update
aws lambda update-function-code --function-name PingMe-user --zip-file fileb://user.zip

## Invoke
aws lambda invoke --function-name user --invocation-type "RequestResponse" response.txt