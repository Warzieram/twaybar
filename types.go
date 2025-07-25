package main

import "time"

// Message represents a WebSocket message from Twitch EventSub
type Message struct {
	Metadata struct {
		MessageID        string `json:"message_id"`
		MessageType      string `json:"message_type"`
		MessageTimestamp string `json:"message_timestamp"`
	} `json:"metadata"`
	Payload struct {
		Subscription struct {
			Type string `json:"type"`
		}
		Session struct {
			ID          string `json:"id"`
			Status      string `json:"connected"`
			ConnectedAt string `json:"connected_at"`
			KeepAlive   int    `json:"keepalive_timeout_seconds"`
		}
	}
}

type MessageMetadata struct {
	MessageID           string `json:"message_id"`
	MessageType         string `json:"message_type"`
	MessageTimestamp    string `json:"message_timestamp"`
	SubscriptionType    string `json:"subscription_type"`
	SubscriptionVersion string `json:"subscription_version"`
}

type SubscriptionMessageEvent struct {
	Metadata MessageMetadata `json:"metadata"`
	Payload  struct {
		Subscription struct {
			ID        string `json:"id"`
			Status    string `json:"status"`
			Type      string `json:"type"`
			Version   string `json:"version"`
			Condition struct {
				BroadcasterUserID string `json:"broadcaster_user_id"`
			} `json:"condition"`
			Transport struct {
				Method    string `json:"method"`
				SessionID string `json:"session_id"`
			} `json:"transport"`
			CreatedAt string `json:"created_at"`
			Cost      int    `json:"cost"`
		} `json:"subscription"`
		Event struct {
			UserID    string `json:"user_id"`
			UserLogin string `json:"user_login"`
			UserName  string `json:"user_name"`
			Tier      string `json:"tier"`
			Message struct {
				Text string `json:"text"`
				Cumul int `json:"cumulative_months"`
			}
		}
	}
}

type SubcriptionEvent struct {
	Metadata MessageMetadata `json:"metadata"`
	Payload  struct {
		Subscription struct {
			ID        string `json:"id"`
			Status    string `json:"status"`
			Type      string `json:"type"`
			Version   string `json:"version"`
			Condition struct {
				BroadcasterUserID string `json:"broadcaster_user_id"`
			} `json:"condition"`
			Transport struct {
				Method    string `json:"method"`
				SessionID string `json:"session_id"`
			} `json:"transport"`
			CreatedAt string `json:"created_at"`
			Cost      int    `json:"cost"`
		} `json:"subscription"`
		Event struct {
			UserID    string `json:"user_id"`
			UserLogin string `json:"user_login"`
			UserName  string `json:"user_name"`
			Tier      string `json:"tier"`
			IsGift    bool   `json:"is_gift"`
		}
	}
}

type ChatMessageEvent struct {
	Metadata MessageMetadata `json:"metadata"`
	Payload  struct {
		Subscription struct {
			ID        string `json:"id"`
			Status    string `json:"status"`
			Type      string `json:"type"`
			Version   string `json:"version"`
			Condition struct {
				BroadcasterUserID string `json:"broadcaster_user_id"`
				UserID            string `json:"user_id"`
			} `json:"condition"`
			Transport struct {
				Method    string `json:"method"`
				SessionID string `json:"session_id"`
			} `json:"transport"`
			CreatedAt string `json:"created_at"`
			Cost      int    `json:"cost"`
		} `json:"subscription"`
		Event struct {
			BroadcasterUserID          string  `json:"broadcaster_user_id"`
			BroadcasterUserLogin       string  `json:"broadcaster_user_login"`
			BroadcasterUserName        string  `json:"broadcaster_user_name"`
			SourceBroadcasterUserID    *string `json:"source_broadcaster_user_id"`
			SourceBroadcasterUserLogin *string `json:"source_broadcaster_user_login"`
			SourceBroadcasterUserName  *string `json:"source_broadcaster_user_name"`
			ChatterUserID              string  `json:"chatter_user_id"`
			ChatterUserLogin           string  `json:"chatter_user_login"`
			ChatterUserName            string  `json:"chatter_user_name"`
			MessageID                  string  `json:"message_id"`
			SourceMessageID            *string `json:"source_message_id"`
			IsSourceOnly               *bool   `json:"is_source_only"`
			Message                    struct {
				Text      string `json:"text"`
				Fragments []struct {
					Type      string    `json:"type"`
					Text      string    `json:"text"`
					Cheermote *struct{} `json:"cheermote"`
					Emote     *struct {
						ID         string   `json:"id"`
						EmoteSetID string   `json:"emote_set_id"`
						OwnerID    string   `json:"owner_id"`
						Format     []string `json:"format"`
					} `json:"emote"`
					Mention *struct{} `json:"mention"`
				} `json:"fragments"`
			} `json:"message"`
			Color  string `json:"color"`
			Badges []struct {
				SetID string `json:"set_id"`
				ID    string `json:"id"`
				Info  string `json:"info"`
			} `json:"badges"`
			SourceBadges                *[]struct{} `json:"source_badges"`
			MessageType                 string      `json:"message_type"`
			Cheer                       *struct{}   `json:"cheer"`
			Reply                       *struct{}   `json:"reply"`
			ChannelPointsCustomRewardID *string     `json:"channel_points_custom_reward_id"`
			ChannelPointsAnimationID    *string     `json:"channel_points_animation_id"`
		} `json:"event"`
	} `json:"payload"`
}

// ClientCredentialsRequestBody represents the request body for client credentials OAuth flow
type ClientCredentialsRequestBody struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

// ClientCredentials represents the response from client credentials OAuth flow
type ClientCredentials struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// CreateSubscriptionRequestBody represents the request body for creating EventSub subscriptions
type CreateSubscriptionRequestBody struct {
}

type TokenStorage struct {
	UserToken string    `json:"access_token"`
	ExpiresIn int       `json:"expires_in"`
	CreatedAt time.Time `json:"created_at"`
}
