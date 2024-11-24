package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func getSteamGamesList() []int {
	games := []int{}
	// check for the games list file
	if _, err := os.Stat("games.txt"); err == nil {
		// games list file exists, check if it's older than 3 days
		fileInfo, err := os.Stat("games.txt")
		if err != nil {
			panic(err)
		}
		if time.Since(fileInfo.ModTime()).Hours() < 72 {
			// games list file is less than 3 days old, don't update it and instead open it into an array of ints
			file, err := os.Open("games.txt")
			if err != nil {
				panic(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				game, err := strconv.Atoi(scanner.Text())
				if err != nil {
					panic(err)
				}
				games = append(games, game)
			}
			return games
		}
	}
	// grab the list of games from steam's api
	resp, err := http.Get(fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/?key=%s&steamid=%s&include_played_free_games=1&include_appinfo=1&format=json", os.Getenv("STEAM_API_KEY"), os.Getenv("STEAM_ID")))
	if err != nil {
		panic(err)
	}

	var response SteamResponse

	// Parse the response
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	for _, game := range response.Response.Games {
		games = append(games, game.Appid)
	}

	// Write the games list to a file (newline separated)
	err = os.WriteFile("games.txt", []byte(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(games)), "\n"), "[]")), 0644)
	if err != nil {
		panic(err)
	}

	return games
}

func installGame(appId string) bool {
	uri := fmt.Sprintf("https://api.steampowered.com/IClientCommService/InstallClientApp/v1/?access_token=%s&appid=%s", os.Getenv("STEAM_AUTH_TOKEN"), appId)
	resp, err := http.Post(uri, "application/json", nil)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode > 204 {
		fmt.Println("Failed to start game installation")
		return false
	}
	return true
}

func gameInstalled(appId string) (bool, int, string) {
	uri := fmt.Sprintf("https://api.steampowered.com/IClientCommService/GetClientAppList/v1?access_token=%s&filters=changing&filter_appids[0]=%s", os.Getenv("STEAM_AUTH_TOKEN"), appId)
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Failed to check if game is installed")
		return false, -1, "???"
	}

	var response SteamClientResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	if len(response.Response.Apps) == 0 {
		return false, -1, "???"
	}

	if response.Response.Apps[0].BytesToDownload == "" || response.Response.Apps[0].BytesToStage == "" {
		return false, response.Response.RefetchIntervalSecUpdating, "0%"
	}

	// convert the strings to ints
	btd, err := strconv.Atoi(response.Response.Apps[0].BytesToDownload)
	if err != nil {
		panic(err)
	}

	bd, err := strconv.Atoi(response.Response.Apps[0].BytesDownloaded)
	if err != nil {
		panic(err)
	}

	bs, err := strconv.Atoi(response.Response.Apps[0].BytesStaged)
	if err != nil {
		panic(err)
	}

	bts, err := strconv.Atoi(response.Response.Apps[0].BytesToStage)
	if err != nil {
		panic(err)
	}

	// calculate the percentage
	percentage := (float64(bd) + float64(bs)) / (float64(btd) + float64(bts)) * 100

	// floor the percentage
	percentageStr := fmt.Sprintf("%.2f%%", percentage)

	if response.Response.Apps[0].BytesToDownload != response.Response.Apps[0].BytesDownloaded {

		return false, response.Response.RefetchIntervalSecUpdating, percentageStr
	}

	if response.Response.Apps[0].BytesToStage != response.Response.Apps[0].BytesStaged {
		return false, response.Response.RefetchIntervalSecUpdating, percentageStr
	}

	return true, response.Response.RefetchIntervalSecUpdating, percentageStr
}

func startGame(appId string) {
	uri := fmt.Sprintf("https://api.steampowered.com/IClientCommService/LaunchClientApp/v1?access_token=%s&appid=%s", os.Getenv("STEAM_AUTH_TOKEN"), appId)

	resp, err := http.Post(uri, "application/json", nil)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Failed to start game")
		return
	}

	fmt.Println("Game started")
}
