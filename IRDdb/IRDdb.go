package IRDdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Callan05/IRacingData"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DBSession struct {
	RaceID int
	PK     string
	Data   IRacingData.Session
	SK     string
}

type DBLeague_Season_Standings struct {
	PK   string
	Data IRacingData.League_Season_Standings
	SK   string
}

var tableName string
var client *dynamodb.DynamoDB

func Init(tname string) {
	tableName = tname
	fmt.Println("Starting DynamoDB client for IRacingData")
	fmt.Println("table: " + tableName)
}

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	client = dynamodb.New(sess)

}

func GetSession(sessionID string) (IRacingData.Session, error) {
	query := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String("session"),
			},
			"SK": {
				S: aws.String(sessionID),
			},
		},
	}

	result, err := client.GetItem(query)
	if err != nil {
		return IRacingData.Session{}, err
	}

	if result.Item == nil {
		return IRacingData.Session{}, errors.New("item not found in DB")
	}

	var ret DBSession
	dynamodbattribute.UnmarshalMap(result.Item, &ret)

	j, _ := json.MarshalIndent(ret, "", "  ")
	fmt.Println(string(j))

	return ret.Data, nil
}

func AddSession(session IRacingData.Session, raceID int) {
	session.Results.RaceID = raceID
	item := DBSession{
		PK:   "session",
		SK:   strconv.Itoa(session.Session_id),
		Data: session,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = client.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Success")
}

func FindSession(sessionID int) (bool, error) {

	result, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String("session"),
			},
			"SK": {
				S: aws.String(strconv.Itoa(sessionID)),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func GetLeagueStandings(leagueID string, seasonID string) (IRacingData.League_Season_Standings, error) {
	query := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String("league_standings"),
			},
			"SK": {
				S: aws.String(leagueID + "|" + seasonID),
			},
		},
	}

	result, err := client.GetItem(query)
	if err != nil {
		return IRacingData.League_Season_Standings{}, err
	}

	if result.Item == nil {
		return IRacingData.League_Season_Standings{}, errors.New("item not found in DB")
	}

	var ret DBLeague_Season_Standings
	dynamodbattribute.UnmarshalMap(result.Item, &ret)

	j, _ := json.MarshalIndent(ret, "", "  ")
	fmt.Println(string(j))

	return ret.Data, nil
}

func AddLeagueStandings(session IRacingData.League_Season_Standings, raceID int) {

	item := DBLeague_Season_Standings{
		PK:   "league_standings",
		SK:   IRacingData.League_Season_Standings.LeagueID + "|" + IRacingData.League_Season_Standings.Season_id,
		Data: session,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = client.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Success")
}
