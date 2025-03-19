package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// to redirect to gihub authentication page to get permission of the user
func redirectToGithubAuth(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()

	baseUrl := "https://github.com/login/oauth/authorize?"
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	authUrl := baseUrl + "client_id=" + clientID
	http.Redirect(w, r, authUrl, http.StatusSeeOther)
}

// to get the creadentials of the user from github
func handleGithubAuthCallback(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	godotenv.Load()

	//creating base for accessing token url and body data
	baseUrl := "https://github.com/login/oauth/access_token"
	rawData := map[string]string{
		"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
		"client_secret": os.Getenv("GITHUB_CLIENT_SECRET"),
		"code":          r.URL.Query().Get("code"),
		"redirect_uri":  "http://localhost:8000/auth/callback/github",
	}

	//marshaling into json
	jsonData, err := json.Marshal(rawData)
	if err != nil {
		log.Fatal("Error marshalling", err)
		return
	}

	//put request to get access token from github
	tokenPostReq, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error creating POST request:", err)
	}
	tokenPostReq.Header.Set("content-type", "application/json")
	tokenPostReq.Header.Set("Accept", "application/json")

	tokenPostRes, err := client.Do(tokenPostReq)
	if err != nil {
		log.Fatal("error making put request", err)
	}

	accessJsonData, err := io.ReadAll(tokenPostRes.Body)
	if err != nil {
		log.Fatal("error reading put response", err)
	}

	// Unmarshalling access token
	accessData := map[string]string{}
	err = json.Unmarshal(accessJsonData, &accessData)
	if err != nil {
		log.Fatal("error unmarshalling access data :", err)
	}

	//get user data using access tocken
	githubUserApiReq, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Fatal("Error creating GET request:", err)
	}
	githubUserApiReq.Header.Set("Authorization", "token "+accessData["access_token"])
	githubUserApiReq.Header.Set("Accept", "application/json")

	githubUserApiRes, err := client.Do(githubUserApiReq)
	if err != nil {
		log.Fatal("error making get request", err)
	}

	userJsonData, err := io.ReadAll(githubUserApiRes.Body)
	if err != nil {
		log.Fatal("error reading get response", err)
	}

	//unmarshalling user data
	userData := map[string]interface{}{}
	err = json.Unmarshal(userJsonData, &userData)
	if err != nil {
		log.Fatal("error unmarshalling userdata :", err)
	}
	fmt.Println(userData)
	fmt.Fprintln(w, userData)

	porfile := map[string]string{
		"url":userData["avatar_url"].(string),
		"email":userData["email"].(string),
		"name":userData["name"].(string),
	}

}

// all routes for github login
func Github_Routes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/github", redirectToGithubAuth)
	mux.HandleFunc("/auth/callback/github", handleGithubAuthCallback)
}
