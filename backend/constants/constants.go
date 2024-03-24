package constants

import "os"

var (
	LLM_API_URL = "http://localhost:8000"
	PORT        = "9000"
)

func init() {
	if envLLMApiURL, exists := os.LookupEnv("LLM_API_URL"); exists {
		LLM_API_URL = envLLMApiURL
	}
	if envPort, exists := os.LookupEnv("PORT"); exists {
		PORT = envPort
	}
}
