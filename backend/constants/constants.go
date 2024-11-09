package constants

import "os"

var (
	LLM_API_URL = "http://localhost:8000"
	PORT        = "9000"

	EASE_FACTOR_DEFAULT = 2.5

	REVIEW_ITEM_DIFFICULTY_DEFAULT = 3.0
	REVIEW_ITEM_DIFFICULTY_MIN     = 1.0
	REVIEW_ITEM_DIFFICICULT_MAX    = 5.0

	REVIEW_ITEM_STREAK_DEFAULT int32 = 0

	REVIEW_ITEM_INTERVAL_IN_MINUTES_DEFAULT int32 = 60
)

func init() {
	if envLLMApiURL, exists := os.LookupEnv("LLM_API_URL"); exists {
		LLM_API_URL = envLLMApiURL
	}
	if envPort, exists := os.LookupEnv("PORT"); exists {
		PORT = envPort
	}
}
