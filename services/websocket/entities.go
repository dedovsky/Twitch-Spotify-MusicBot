package websocket

import (
	"encoding/json"
	"time"
)

type Metadata struct {
	MessageID        string `json:"message_id"`
	MessageType      string `json:"message_type"`
	MessageTimestamp string `json:"message_timestamp"`
}

type SessionWelcomeMessage struct {
	Metadata Metadata `json:"metadata"`
	Payload  struct {
		Session struct {
			ID                      string  `json:"id"`
			Status                  string  `json:"status"`
			ConnectedAt             string  `json:"connected_at"`
			KeepaliveTimeoutSeconds int     `json:"keepalive_timeout_seconds"`
			ReconnectURL            *string `json:"reconnect_url"`
			RecoveryURL             *string `json:"recovery_url"`
		} `json:"session"`
	} `json:"payload"`
}

type KeepaliveMessage struct {
	Metadata Metadata `json:"metadata"`
	Payload  struct{} `json:"payload"`
}

type NotificationMessage struct {
	Metadata struct {
		MessageID           string `json:"message_id"`
		MessageType         string `json:"message_type"`
		MessageTimestamp    string `json:"message_timestamp"`
		SubscriptionType    string `json:"subscription_type"`
		SubscriptionVersion string `json:"subscription_version"`
	} `json:"metadata"`
	Payload struct {
		Subscription struct {
			ID        string `json:"id"`
			Status    string `json:"status"`
			Type      string `json:"type"`
			Version   string `json:"version"`
			Cost      int    `json:"cost"`
			Condition struct {
				BroadcasterUserID string `json:"broadcaster_user_id"`
			} `json:"condition"`
			Transport struct {
				Method    string `json:"method"`
				SessionID string `json:"session_id"`
			} `json:"transport"`
			CreatedAt string `json:"created_at"`
		} `json:"subscription"`
		Event struct {
			UserID               string `json:"user_id"`
			UserLogin            string `json:"user_login"`
			UserName             string `json:"user_name"`
			BroadcasterUserID    string `json:"broadcaster_user_id"`
			BroadcasterUserLogin string `json:"broadcaster_user_login"`
			BroadcasterUserName  string `json:"broadcaster_user_name"`
			FollowedAt           string `json:"followed_at"`
		} `json:"event"`
	} `json:"payload"`
}

type TwitchMessage struct {
	Metadata Metadata        `json:"metadata"`
	Payload  json.RawMessage `json:"payload"`
}

type ChatMessage struct {
	Metadata struct {
		MessageId           string    `json:"message_id"`
		MessageType         string    `json:"message_type"`
		MessageTimestamp    time.Time `json:"message_timestamp"`
		SubscriptionType    string    `json:"subscription_type"`
		SubscriptionVersion string    `json:"subscription_version"`
	} `json:"metadata"`
	Payload struct {
		Subscription struct {
			Id        string `json:"id"`
			Status    string `json:"status"`
			Type      string `json:"type"`
			Version   string `json:"version"`
			Condition struct {
				BroadcasterUserId string `json:"broadcaster_user_id"`
				UserId            string `json:"user_id"`
			} `json:"condition"`
			Transport struct {
				Method    string `json:"method"`
				SessionId string `json:"session_id"`
			} `json:"transport"`
			CreatedAt time.Time `json:"created_at"`
			Cost      int       `json:"cost"`
		} `json:"subscription"`
		Event struct {
			BroadcasterUserId    string `json:"broadcaster_user_id"`
			BroadcasterUserLogin string `json:"broadcaster_user_login"`
			BroadcasterUserName  string `json:"broadcaster_user_name"`
			ChatterUserId        string `json:"chatter_user_id"`
			ChatterUserLogin     string `json:"chatter_user_login"`
			ChatterUserName      string `json:"chatter_user_name"`
			MessageId            string `json:"message_id"`
			Message              struct {
				Text      string `json:"text"`
				Fragments []struct {
					Type      string      `json:"type"`
					Text      string      `json:"text"`
					Cheermote interface{} `json:"cheermote"`
					Emote     interface{} `json:"emote"`
					Mention   interface{} `json:"mention"`
				} `json:"fragments"`
			} `json:"message"`
			Color  string `json:"color"`
			Badges []struct {
				SetId string `json:"set_id"`
				Id    string `json:"id"`
				Info  string `json:"info"`
			} `json:"badges"`
			MessageType                 string      `json:"message_type"`
			Cheer                       interface{} `json:"cheer"`
			Reply                       interface{} `json:"reply"`
			ChannelPointsCustomRewardId interface{} `json:"channel_points_custom_reward_id"`
			ChannelPointsAnimationId    interface{} `json:"channel_points_animation_id"`
		} `json:"event"`
	} `json:"payload"`
}

type RewardRedeemedMessage struct {
	Payload struct {
		Event struct {
			Reward struct {
				Title string `json:"title"`
			} `json:"reward"`
		} `json:"event"`
	} `json:"payload"`
}

type TwitchChannelID struct {
	Data []struct {
		Id string `json:"id"`
	} `json:"data"`
}
