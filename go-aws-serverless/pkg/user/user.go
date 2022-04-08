package user

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/tonystrawberry/go-aws-serverless/pkg/validators"
)

var ErrorFailedToFetchRecord = "Failed to fetch record"
var ErrorFailedToUnmarshalRecord = "Failed to unmarshal record"
var ErrorInvalidUserData = "Invalid user data"
var ErrorInvalidEmail = "Invalid email"
var ErrorCouldNotMarshalItem = "Could not marshal item"
var ErrorCouldNotDeleteItem = "Could not delete item"
var ErrorCouldNotPutItem = "Could not put item"
var ErrorUserDoesNotExist = "User does not exist"
var ErrorUserAlreadyExist = "User already exists"

type User struct {
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	res, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	user := new(User)
	err = dynamodbattribute.UnmarshalMap(res.Item, user)

	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return user, nil
}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*[]User, error){
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	res, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	users := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, users)

	return users, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error){
	var u User

	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(u.Email){
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExist)
	}

	av, err := dynamodbattribute.MarshalMap(u)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	log.Println(av)

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)

	if err != nil {
		log.Fatal(err.Error())
		return nil, errors.New(ErrorCouldNotPutItem)
	}

	return &u, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error){
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)
	if currentUser == nil || len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)

	if err != nil {
		log.Fatal(err.Error())
		return nil, errors.New(ErrorCouldNotPutItem)
	}

	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error{
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}
