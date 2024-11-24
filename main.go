package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	mathRand "math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func init() {
	// make sure the envs are set
	hasSteamApiKey := os.Getenv("STEAM_API_KEY") != ""
	hasSteamAuthToken := os.Getenv("STEAM_AUTH_TOKEN") != ""
	hasSteamId := os.Getenv("STEAM_ID") != ""
	if !hasSteamApiKey || !hasSteamAuthToken || !hasSteamId {
		fmt.Println("The following env variables are missing: ")
		fmt.Printf("STEAM_API_KEY: %t\nSTEAM_AUTH_TOKEN: %t\nSTEAM_ID: %t\n", !hasSteamApiKey, !hasSteamAuthToken, !hasSteamId)
		os.Exit(1)
	}

	info := strings.Split(os.Getenv("STEAM_AUTH_TOKEN"), ".")[1]

	// base64 decode the token
	decoded, err := base64.StdEncoding.DecodeString(info + "==")
	if err != nil {
		fmt.Printf("Error decoding base64 from auth token %v\n", err)
		os.Exit(2)
	}

	// json decode the token
	var token AuthToken
	err = json.Unmarshal(decoded, &token)
	if err != nil {
		fmt.Printf("Error decoding JSON from base64 decoded auth token %v\n", err)
		os.Exit(3)
	}

	// check if the token is expired
	if time.Now().Unix() > int64(token.Exp) {
		fmt.Println("Auth token expired")
		fmt.Println("Go to https://store.steampowered.com/pointssummary/ajaxgetasyncconfig when logged in and grab the value of webapi_token and throw it in the .env as STEAM_AUTH_TOKEN")
		os.Exit(4)
	}
}

func main() {
	games := getSteamGamesList()

	// choose a random game
	randomGame := games[mathRand.Intn(len(games))]
	gameID := strconv.Itoa(randomGame)
	fmt.Println("Found a game to install, starting install now")
	installGame(gameID)

	for {
		installed, wait, percentageStr := gameInstalled(gameID)
		if wait == -1 {
			fmt.Println("Error checking if game is installed")
			os.Exit(1)
		}

		if installed {
			break
		}
		fmt.Print("Install progress: ", percentageStr, "\r")
		time.Sleep(time.Duration(wait) * time.Second)
	}

	fmt.Println("Game installed, starting in 2 seconds")
	time.Sleep(2 * time.Second)
	startGame(gameID)
}
