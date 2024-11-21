package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func makeRequest(url string, auth_token string) (data *[]byte, status_code *int) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", auth_token)
	req.Header.Set("Client-Id", getEnv("TWITCH_CLIENT_ID"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("makeRequest Request Error: %s - %s\n", url, err.Error())
		return nil, &resp.StatusCode
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("makeRequest Body Error: %s - %s\n", url, err.Error())
		return nil, nil
	}

	return &body, &resp.StatusCode
}

type Users struct {
	Data []User `json:"data"`
}

type User struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	ProfileImageUrl string `json:"profile_image_url"`
}

func getUserID(name string, auth_token string) (user *User) {
	body, _ := makeRequest(TWITCH_URL_USERS+name, auth_token)
	if body == nil {
		return nil
	}

	var users Users
	json.Unmarshal(*body, &users)

	find := -1
	for i := 0; i < len(users.Data); i++ {
		if users.Data[i].Login != getEnv("TWITCH_CHANNEL_NAME") {
			continue
		}
		find = i
		break
	}

	if find == -1 {
		return nil
	}

	return &users.Data[find]
}

type Streams struct {
	Data []Stream `json:"data"`
}

type Stream struct {
	UserId       string `json:"user_id"`
	UserLogin    string `json:"user_login"`
	GameName     string `json:"game_name"`
	Title        string `json:"title"`
	ViewerCount  int    `json:"viewer_count"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

func getStream(user_id string, auth_token string) (stream *Stream, status_code *int) {
	body, status_code := makeRequest(TWITCH_URL_STREAMS+user_id, auth_token)
	if body == nil {
		return nil, status_code
	}

	var streams Streams
	json.Unmarshal(*body, &streams)

	find := -1
	for i := 0; i < len(streams.Data); i++ {
		if streams.Data[i].UserLogin != getEnv("TWITCH_CHANNEL_NAME") {
			continue
		}
		find = i
		break
	}

	if find == -1 {
		return nil, status_code
	}

	return &streams.Data[find], nil
}

type AuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func getAuthToken() (auth_token *string) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", getEnv("TWITCH_CLIENT_ID"))
	data.Set("client_secret", getEnv("TWITCH_CLIENT_SECRET"))

	req, _ := http.NewRequest(http.MethodPost, TWITCH_URL_OAUTH_TOKEN, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("getAuthToken Request Error: %s\n", err.Error())
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("makeRequest Body Error: %s\n", err.Error())
		return nil
	}

	// fmt.Println(string(body))

	var token_data AuthToken
	json.Unmarshal(body, &token_data)

	// fmt.Println(token_data)

	token := fmt.Sprintf("Bearer %s", token_data.AccessToken)

	return &token
}

type MessageEmbed struct {
	Content   string  `json:"content"`
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Embeds    []Embed `json:"embeds"`
}

type Embed struct {
	Title  string       `json:"title"`
	URL    string       `json:"url"`
	Color  int          `json:"color"`
	Fields []EmbedField `json:"fields"`
	Author struct {
		Name    string `json:"name"`
		IconURL string `json:"icon_url"`
	} `json:"author"`
	Image struct {
		URL string `json:"url"`
	} `json:"image"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type MessageReturnData struct {
	Id string `json:"id"`
}

func sendWebhookMessage(user User, stream Stream) (message_id *string) {
	message := getEnv("DISCORD_LIVE_MESSAGE")
	message = strings.ReplaceAll(message, "{CHANNEL_NAME}", user.DisplayName)
	message = strings.ReplaceAll(message, "{GAME}", stream.GameName)
	message = strings.ReplaceAll(message, "{TITLE}", stream.Title)

	var message_embed MessageEmbed

	message_embed.Content = message
	message_embed.Username = user.DisplayName
	message_embed.AvatarURL = user.ProfileImageUrl
	embed := &Embed{}
	embed.Title = stream.Title
	embed.URL = fmt.Sprintf("http://twitch.tv/%s", user.Login)
	color, _ := strconv.ParseInt("9148ff", 16, 32)
	embed.Color = int(color)
	embed.Author.Name = user.DisplayName
	embed.Author.IconURL = user.ProfileImageUrl
	embed.Fields = append(embed.Fields, EmbedField{Name: "Game", Value: stream.GameName, Inline: true})
	embed.Fields = append(embed.Fields, EmbedField{Name: "Viewers", Value: fmt.Sprintf("%d", stream.ViewerCount), Inline: true})
	embed.Fields = append(embed.Fields, EmbedField{Name: "Title", Value: stream.Title, Inline: false})
	embed.Image.URL = strings.ReplaceAll(strings.ReplaceAll(stream.ThumbnailUrl, "{width}", "1920"), "{height}", "1080") + fmt.Sprintf("?t=%d", time.Now().Minute()/5)
	message_embed.Embeds = append(message_embed.Embeds, *embed)

	req_body, _ := json.Marshal(message_embed)
	req, _ := http.NewRequest(http.MethodPost, getEnv("DISCORD_WEBHOOK_URL")+"?wait=true", strings.NewReader(string(req_body)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("sendWebhookMessage Request Error: %s\n", err.Error())
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("sendWebhookMessage Body Error: %s\n", err.Error())
		return nil
	}

	var message_data MessageReturnData
	json.Unmarshal(body, &message_data)

	return &message_data.Id
}

func updateWebhookMessage(message_id string, user User, stream Stream) (success bool) {
	message := getEnv("DISCORD_LIVE_MESSAGE")
	message = strings.ReplaceAll(message, "{CHANNEL_NAME}", user.DisplayName)
	message = strings.ReplaceAll(message, "{GAME}", stream.GameName)
	message = strings.ReplaceAll(message, "{TITLE}", stream.Title)

	var message_embed MessageEmbed

	message_embed.Content = message
	message_embed.Username = user.DisplayName
	message_embed.AvatarURL = user.ProfileImageUrl
	embed := &Embed{}
	embed.Title = stream.Title
	embed.URL = fmt.Sprintf("http://twitch.tv/%s", user.Login)
	color, _ := strconv.ParseInt("9148ff", 16, 32)
	embed.Color = int(color)
	embed.Author.Name = user.DisplayName
	embed.Author.IconURL = user.ProfileImageUrl
	embed.Fields = append(embed.Fields, EmbedField{Name: "Game", Value: stream.GameName, Inline: true})
	embed.Fields = append(embed.Fields, EmbedField{Name: "Viewers", Value: fmt.Sprintf("%d", stream.ViewerCount), Inline: true})
	embed.Fields = append(embed.Fields, EmbedField{Name: "Title", Value: stream.Title, Inline: false})
	embed.Image.URL = strings.ReplaceAll(strings.ReplaceAll(stream.ThumbnailUrl, "{width}", "1920"), "{height}", "1080") + fmt.Sprintf("?t=%d", time.Now().Minute()/5)
	message_embed.Embeds = append(message_embed.Embeds, *embed)

	req_body, _ := json.Marshal(message_embed)
	req, _ := http.NewRequest(http.MethodPatch, getEnv("DISCORD_WEBHOOK_URL")+"/messages/"+message_id, strings.NewReader(string(req_body)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("updateWebhookMessage Request Error: %s\n", err.Error())
		return false
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("updateWebhookMessage Body Error: %s\n", err.Error())
		return false
	}

	fmt.Printf("updateWebhookMessage Status: %d\n", resp.StatusCode)
	fmt.Printf("updateWebhookMessage Body: %s\n", string(body))

	return true
}
