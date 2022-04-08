aws iam create-role --role-name lambda-ex --assume-role-policy-document file://trust-policy.json

aws iam attach-role-policy --role-name lambda-ex --policy-arn arn:aws:iam:aws:policy/service-role/AWSLambdaBasicExecutionRole

go build
go build main.go

zip function.zip main

aws lambda create-function --function-name go-lambda --zip-file fileb://function.zip --handler main --runtime go1.x --role arn:aws:iam:123456789:role/lambda-ex

aws lambda invoke --function-name go-lambda --cli-binary-format raw-in-base64-out --payload '{"what is your name?": "Jim", "How old are you?": 33}' output.txt