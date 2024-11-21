package main

import (
	"fmt"
	"os"
)

func checkEnv(name string) bool {
	if os.Getenv(name) == "" {
		return false
	} else {
		return true
	}
}

func getEnv(name string) (variable string) {
	val := os.Getenv(name)
	switch name {
	case "TWITCH_CLIENT_ID":
		return panicVal(val)
	case "TWITCH_CLIENT_SECRET":
		return panicVal(val)
	case "TWITCH_CHANNEL_NAME":
		return panicVal(val)
	case "DISCORD_WEBHOOK_URL":
		return panicVal(val)
	case "DISCORD_LIVE_MESSAGE":
		return backupVal(val, fmt.Sprintf("%s is now live!", getEnv("TWITCH_CHANNEL_NAME")))
	default:
		return ""
	}
}

func backupVal(chk_val string, backup_val string) (val string) {
	if chk_val == "" {
		return backup_val
	} else {
		return chk_val
	}
}

func panicVal(chk_val string) (val string) {
	if chk_val == "" {
		panic(10)
	} else {
		return chk_val
	}
}
