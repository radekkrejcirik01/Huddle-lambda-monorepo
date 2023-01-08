## Deployment of Golang app to AWS's Lambda

1. export AWS_PROFILE=admin
2. GOARCH=amd64 GOOS=linux go build main.go
3. zip -r login-registration.zip .

## Create
aws lambda create-function \                                                                    
--function-name login-registration \
--zip-file fileb://login-registration.zip \  
--handler login-registration \
--runtime go1.x \
--role "arn:aws:iam::409186456204:role/lambda-basic-execution"

## Update
aws lambda update-function-code --function-name login-registration --zip-file fileb://login-registration.zip

## Invoke
aws lambda invoke --function-name login-registration --invocation-type "RequestResponse" response.txt