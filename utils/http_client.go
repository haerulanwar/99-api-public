package utils

import (
	"99-api-public/config"

	"github.com/go-resty/resty/v2"
)

var Client = resty.New()

var (
	UserServiceURL    string
	ListingServiceURL string
)

func InitURLs() {
	UserServiceURL = config.GetEnv("USER_SERVICE_URL", "http://localhost:8080")
	ListingServiceURL = config.GetEnv("LISTING_SERVICE_URL", "http://localhost:6000")
}
