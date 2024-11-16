package constants

import (
	"os"
	"strconv"
)

var (
	BACKEND_URL = "http://localhost:9000"
	PORT        = "42069"

	REVIEW_ITEM_PAGE_SIZE = 10
)

func init() {
	if envBackendURL, exists := os.LookupEnv("BACKEND_URL"); exists {
		BACKEND_URL = envBackendURL
	}
	if envPort, exists := os.LookupEnv("PORT"); exists {
		PORT = envPort
	}
	if envReviewItemPageSize, exist := os.LookupEnv("REVIEW_ITEM_PAGE_SIZE"); exist {
		parsedPageSize, err := strconv.Atoi(envReviewItemPageSize)
		if err != nil {
			envReviewItemPageSize = strconv.Itoa(parsedPageSize)
		}
	}
}
