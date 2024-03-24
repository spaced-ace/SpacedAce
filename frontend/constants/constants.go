package constants

import "os"

var (
	BACKEND_URL = "http://localhost:9000"
	PORT        = "42069"
)

func init() {
	if envBackendURL, exists := os.LookupEnv("BACKEND_URL"); exists {
		BACKEND_URL = envBackendURL
	}
	if envPort, exists := os.LookupEnv("PORT"); exists {
		PORT = envPort
	}
}
