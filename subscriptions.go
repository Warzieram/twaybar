package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func CreateChannelSubscribeSubscription(sessionID, authorization, userID, broadcasterUserID string, client http.Client, debug bool) error {
	requestURL := "https://api.twitch.tv/helix/eventsub/subscriptions"
	body := map[string]any{
		"type":    "channel.subscribe",
		"version": "1",
		"condition": map[string]any{
			"broadcaster_user_id": broadcasterUserID,
		},
		"transport": map[string]any{
			"method":     "websocket",
			"session_id": sessionID,
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+authorization)
	req.Header.Set("Client-id", os.Getenv("CLIENT_ID"))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func CreateChatSubscription(sessionID string, authorization string, userID string, broadcasterUserID string, client http.Client, debug bool) error {
	requestURL := "https://api.twitch.tv/helix/eventsub/subscriptions"
	body := map[string]any{
		"type":    "channel.chat.message",
		"version": "1",
		"condition": map[string]any{
			"broadcaster_user_id": broadcasterUserID,
			"user_id":             userID,
		},
		"transport": map[string]any{
			"method":     "websocket",
			"session_id": sessionID,
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+authorization)
	req.Header.Set("Client-id", os.Getenv("CLIENT_ID"))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func CreateResubMessageSubscription(sessionID, authorization, userID, broadcasterUserID string, client http.Client, debug bool) error {

	requestURL := "https://api.twitch.tv/helix/eventsub/subscriptions"
	body := map[string]any{
		"type":    "channel.subscription.message",
		"version": "1",
		"condition": map[string]any{
			"broadcaster_user_id": broadcasterUserID,
		},
		"transport": map[string]any{
			"method":     "websocket",
			"session_id": sessionID,
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+authorization)
	req.Header.Set("Client-id", os.Getenv("CLIENT_ID"))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	if debug {
		if debug {
			reqDump, err := httputil.DumpRequest(req, true)
			if err != nil {
				log.Printf("[WARNING] Could not dump request: %v", err)
			} else {
				log.Println("CREATE RESUB MESSAGE SUBSCRIPTION REQUEST: ", string(reqDump))
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if debug {
			respDump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				log.Printf("[WARNING] Could not dump response: %v", err)
			} else {
				log.Println("CREATE RESUB MESSAGE SUBSCRIPTION RESPONSE: ", string(respDump))
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if debug {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Printf("[WARNING] Could not dump response: %v", err)
		} else {
			log.Println("CREATE CHAT SUBSCRIPTION RESPONSE: ", string(respDump))
		}
	}

	return nil
}
