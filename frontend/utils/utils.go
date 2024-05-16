package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func StringInArray(target string, arr []string) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}

func HashToColor(input string) string {
	shaHash := sha256.New()
	shaHash.Write([]byte(input))
	hash := hex.EncodeToString(shaHash.Sum(nil))

	// Convert the first 6 characters of the hash into an integer
	colorInt, _ := strconv.ParseInt(hash[:6], 16, 64)

	// Map the integer to a color
	colors := []string{
		"red-300", "red-400", "red-500", "red-600",
		"orange-300", "orange-400", "orange-500", "orange-600",
		"amber-300", "amber-400", "amber-500", "amber-600",
		"yellow-300", "yellow-400", "yellow-500", "yellow-600",
		"green-300", "green-400", "green-500", "green-600",
		"blue-300", "blue-400", "blue-500", "blue-600",
		"purple-300", "purple-400", "purple-500", "purple-600",
		"pink-300", "pink-400", "pink-500", "pink-600",
		"emerald-300", "emerald-400", "emerald-500", "emerald-600",
		"teal-300", "teal-400", "teal-500", "teal-600",
		"cyan-300", "cyan-400", "cyan-500", "cyan-600",
		"indigo-300", "indigo-400", "indigo-500", "indigo-600",
		"violet-300", "violet-400", "violet-500", "violet-600",
		"fuchsia-300", "fuchsia-400", "fuchsia-500", "fuchsia-600",
		"rose-300", "rose-400", "rose-500", "rose-600",
	}
	color := colors[colorInt%int64(len(colors))]

	return color
}

func GenerateColors(title string, id string) (string, string) {
	fromColor := HashToColor(title + id)
	toColor := HashToColor(id)

	return fromColor, toColor
}
