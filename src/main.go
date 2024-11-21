package main

import (
	"fmt"
	"time"
)

func main() {
	var auth_token *string

	for ok := true; ok; ok = (auth_token == nil) {
		auth_token = getAuthToken()
		if auth_token == nil {
			time.Sleep(time.Second * 15)
		}
	}

	// fmt.Println(*auth_token)

	var user *User

	for ok := true; ok; ok = (user == nil) {
		user = getUserID(getEnv("TWITCH_CHANNEL_NAME"), *auth_token)
		if user == nil {
			time.Sleep(time.Second * 15)
		}
	}

	var message_id *string

	cooldown := 0
	message_id = nil

	fmt.Println(*user)

	for {
		var stream *Stream
		var status *int
		user_copy := *user

		zero := 0
		status = &zero

		for ok := true; ok; ok = (status != nil) {

			stream, status = getStream(user_copy.Id, *auth_token)

			if stream == nil && status != nil && *status != 401 {
				fmt.Println("Error Retrieving Stream Info / No Stream Info To Get")
				time.Sleep(time.Second * 15)
				continue
			}

			if stream == nil && status != nil && *status == 401 {
				fmt.Println("Token Expired Getting New One")
				auth_token = nil
				for ok := true; ok; ok = (auth_token == nil) {
					auth_token = getAuthToken()
					if auth_token == nil {
						time.Sleep(time.Second * 15)
						continue
					}
				}
				continue
			}
		}

		if stream != nil {
			fmt.Println(*stream)
		} else {
			fmt.Println(stream)
		}

		if stream == nil {
			if cooldown > 0 {
				cooldown--
				time.Sleep(time.Second * 15)
				continue
			}

			cooldown = 0
			message_id = nil
			time.Sleep(time.Second * 15)
			continue
		} else {
			if cooldown == 0 {
				message_id = sendWebhookMessage(user_copy, *stream)
				if message_id == nil {
					cooldown = 0
					time.Sleep(time.Second * 15)
					continue
				}
			}

			if message_id == nil {
				message_id = sendWebhookMessage(user_copy, *stream)
				if message_id == nil {
					cooldown = 0
					time.Sleep(time.Second * 15)
					continue
				}
			} else {
				if !updateWebhookMessage(*message_id, user_copy, *stream) {
					cooldown = 0
					time.Sleep(time.Second * 15)
					continue
				}
			}

			cooldown = 8
			time.Sleep(time.Second * 15)
			continue
		}
	}
}
