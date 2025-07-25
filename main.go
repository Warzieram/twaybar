package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
	"github.com/warzieram/twaybar/waybar"
)

type Config struct {
	Subs bool `json:"subs"`
	Chat bool `json:"chat"`
	Resubs bool `json:"resubs"`
	Debug bool `json:"debug"`
}

func LoadConfig() (*Config, error) {

	data, err := os.ReadFile("conf.json")
	if err != nil {
		return nil, err
	}

	var conf Config
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, err
}

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Couldn't load .env file")
	}

	output := &waybar.FormatOutput{}


	client := &http.Client{}

	var userToken string

	conf, confErr := LoadConfig()
	if confErr != nil {
		log.Fatal("[ERROR] Couldn't load configuration: ", confErr)
	}

	subs := conf.Subs
	chat := conf.Chat
	resubs := conf.Resubs
	debug := conf.Debug


	if storage, err := LoadToken(); err == nil {
		userToken = storage.UserToken
	} else {
		userToken, err = startOAuthServer(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
		if err != nil {
			log.Fatal("[ERROR] OAuth failed: ", err)
		}
		if err != nil {
			log.Println("Coulnd't store the user token")
		}
	}

	authorizationToken, err := GetAuthorizationToken(client)
	if err != nil {
		log.Fatal("[ERROR] Couldn't retrieve authorization token: ", err)
	}
	helixClient, err := helix.NewClient(&helix.Options{
		ClientID:       os.Getenv("CLIENT_ID"),
		ClientSecret:   os.Getenv("CLIENT_SECRET"),
		AppAccessToken: authorizationToken,
	})
	if err != nil {
		log.Printf("[ERROR] Creating Client: %v", err)
	}

	user, userErr := GetUser(os.Getenv("USER_LOGIN"), helixClient)
	if userErr != nil {
		log.Fatal("[ERROR] Retrieving the user: ", userErr)
	}

	broadcaster, broadcasterErr := GetUser(os.Getenv("BROADCASTER_LOGIN"), helixClient)
	if broadcasterErr != nil {
		log.Fatal("[ERROR] Retrieving the user: ", userErr)
	}

	serverURL := "wss://eventsub.wss.twitch.tv/ws"

	headers := http.Header{}
	headers.Add("Authorization", authorizationToken)

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, headers)
	if err != nil {
		log.Fatalf("Error during connexion : %v", err)
	}
	defer conn.Close()

	log.Println("Connection established successfully")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("Erreur lors de la lecture du message : %v", err)
		}
		//log.Printf("Message reçu : %s", message)

		data := &Message{}
		if jsonErr := json.Unmarshal([]byte(message), data); jsonErr != nil {
			log.Fatalf("Error serializing message data: %s", jsonErr)
		}

		messageType := data.Metadata.MessageType

		if messageType == "session_welcome" {
			//log.Println("Welcome Message Received")
			sessionID := data.Payload.Session.ID
			if chat {
				subErr := CreateChatSubscription(sessionID, userToken, user.ID, broadcaster.ID, *client, debug)
				if subErr != nil {
					log.Fatal("[ERROR] Couldn't create subscription: ", subErr)
				}
			}
			if subs {

				subErr := CreateChannelSubscribeSubscription(sessionID, userToken, user.ID, broadcaster.ID, *client, debug)
				if subErr != nil {
					log.Fatal("[ERROR] Couldn't create subscription: ", subErr)
				}
			}
			if resubs {

				subErr := CreateResubMessageSubscription(sessionID, userToken, user.ID, broadcaster.ID, *client, debug)
				if subErr != nil {
					log.Fatal("[ERROR] Couldn't create subscription: ", subErr)
				}
			}
		}

		if messageType == "notification" {
			subType := data.Payload.Subscription.Type
			output.Text = ""
			if subType == "channel.chat.message" {

				chatMessage := &ChatMessageEvent{}
				json.Unmarshal(message, chatMessage)
				messageText := chatMessage.Payload.Event.Message.Text
				messageSender := chatMessage.Payload.Event.ChatterUserName

				//fmt.Printf("%s : %s\n", messageSender, messageText)
				output.Text += messageSender + ": " + messageText
			}
			if subType == "channel.subscribe" {
				subscribeEvent := &SubcriptionEvent{}
				json.Unmarshal(message, subscribeEvent)
				subUserName := subscribeEvent.Payload.Event.UserName

				//fmt.Printf("[NEW SUBSCRIBER]: %s\n", subUserName)
				output.Tooltip += "New sub: " + subUserName + "\n"
			}

			if subType == "channel.subscription.message" {
				subMessageEvent := &SubscriptionMessageEvent{}
				json.Unmarshal(message, subMessageEvent)
				subUserName := subMessageEvent.Payload.Event.UserName
				subMessage := subMessageEvent.Payload.Event.Message.Text

				//fmt.Printf("[NEW RESUB] : %s : %s", subUserName, subMessage)
				output.Text += "New Resub: " + subUserName + ": " + subMessage + "\n"

			}
			err := output.Print()
			if err != nil {
				log.Fatalf("Error giving the output: %s", err)
			}
		}

	}
}

func GetAuthorizationToken(httpClient *http.Client) (string, error) {
	requestURL := "https://id.twitch.tv/oauth2/token"
	data := url.Values{}
	data.Set("client_id", os.Getenv("CLIENT_ID"))
	data.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//reqDump, err := httputil.DumpRequest(req, true)
	//if err != nil {
	//	return "", err
	//}

	//log.Println("REQUEST: ", string(reqDump))

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	//respDump, err := httputil.DumpResponse(resp, true)
	//if err != nil {
	//	return "", nil
	//}

	//log.Println("RESPONSE: ", string(respDump))

	credentials := &ClientCredentials{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	jsonErr := json.Unmarshal(body, credentials)
	if jsonErr != nil {
		return "", err
	}

	return credentials.AccessToken, nil
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func startOAuthServer(clientID, clientSecret string) (string, error) {
	code := make(chan string, 1)

	port := os.Getenv("PORT")
	server := &http.Server{Addr: ":"+port}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		//log.Println("[RECIEVED]", r.URL.Query())
		authCode := r.URL.Query().Get("code")
		if authCode != "" {
			code <- authCode
			fmt.Fprintf(w, "Authorization successful! You can close this window.")
		} else {
			fmt.Fprintf(w, "Authorization failed!")
		}
	})

	go server.ListenAndServe()
	scopesList := []string{
		"channel:read:subscriptions",
		"user:read:chat",
	}

	scopeString := strings.Join(scopesList, "+")

	// Open browser to Twitch OAuth page
	authURL := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=http://localhost:8080/callback&response_type=code&scope=%s", clientID, scopeString)
	fmt.Printf("Opening OAuth URL in browser: %s\n", authURL)

	// Automatically open the URL in the default browser
	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Failed to open browser automatically. Please open this URL manually: %s\n", authURL)
	}

	// Wait for authorization code
	select {
	case authCode := <-code:
		server.Shutdown(context.Background())
		return exchangeCodeForToken(authCode, clientID, clientSecret)
	case <-time.After(time.Minute):
		server.Shutdown(context.Background())
		return "", fmt.Errorf("authorization timeout")
	}
}

func exchangeCodeForToken(code, clientID, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", "http://localhost:8080/callback")

	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp = &TokenStorage{}

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &tokenResp)
	// log.Println(tokenResp)
	tokenResp.CreatedAt = time.Now()

	err = SaveTokenToFile(tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.UserToken, nil
}

func GetUser(login string, client *helix.Client) (helix.User, error) {

	resp, err := client.GetUsers(&helix.UsersParams{
		Logins: []string{login},
	})
	if err != nil {
		log.Println("[ERROR] Retrieving user failed: ", err)
		return helix.User{}, err
	}

	return resp.Data.Users[0], nil
}

func SaveTokenToFile(userToken *TokenStorage) error {

	data, _ := json.Marshal(userToken)
	return os.WriteFile("token.json", data, 0600)
}

func LoadToken() (*TokenStorage, error) {
	data, err := os.ReadFile("token.json")
	if err != nil {
		return nil, err
	}

	var token TokenStorage
	err = json.Unmarshal(data, &token)

	if token.UserToken == "" {
		return nil, errors.New("no token stored")
	}
	if token.CreatedAt.Add(time.Duration(token.ExpiresIn) * time.Second).Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return &token, err
}
