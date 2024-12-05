package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Song struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Genre  string `json:"genre"`
}

var (
	songs  = make(map[string]Song)
	mutex  = &sync.Mutex{}
	nextID = 1
)

func createSong(song Song) Song {
	mutex.Lock()
	defer mutex.Unlock()
	song.ID = strconv.Itoa(nextID)
	nextID++
	songs[song.ID] = song
	return song
}

func createSongHandler(w http.ResponseWriter, r *http.Request) {
	var newSong Song
	if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil || newSong.Title == "" || newSong.Artist == "" || newSong.Genre == "" {
		http.Error(w, `{"message": "Invalid input"}`, http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(createSong(newSong))
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var newSong Song
	if err := json.Unmarshal([]byte(req.Body), &newSong); err != nil || newSong.Title == "" || newSong.Artist == "" || newSong.Genre == "" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: `{"message": "Invalid input"}`}, nil
	}
	response, _ := json.Marshal(createSong(newSong))
	return events.APIGatewayProxyResponse{StatusCode: http.StatusCreated, Body: string(response)}, nil
}

func main() {
	lambda.Start(lambdaHandler)
}
