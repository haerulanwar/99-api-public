package handlers

import (
	"99-api-public/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

func ProxyListingServiceGet(c *gin.Context) {

	msUrl := utils.ListingServiceURL + strings.TrimSuffix(strings.TrimPrefix(c.Request.URL.Path, "/public-api"), "/")
	resp, err := utils.Client.R().
		SetQueryParamsFromValues(c.Request.URL.Query()).
		SetHeader("Content-Type", "application/json").
		Get(msUrl)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Listing Service unavailable"})
		return
	}

	var result struct {
		Result   bool                     `json:"result"`
		Listings []map[string]interface{} `json:"listings"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse listings"})
		return
	}
	listings := result.Listings

	var wg sync.WaitGroup
	errChan := make(chan error, len(listings))
	userInfoChan := make(chan struct {
		index int
		user  map[string]interface{}
	}, len(listings))

	for i, listing := range listings {
		wg.Add(1)
		go func(i int, listing map[string]interface{}) {
			defer wg.Done()
			userID := int(listing["user_id"].(float64))

			resp, err := utils.Client.R().
				Get(utils.UserServiceURL + "/users/" + strconv.Itoa(userID))

			if err != nil {
				errChan <- err
				return
			}

			var userInfo struct {
				Result bool                   `json:"result"`
				User   map[string]interface{} `json:"user"`
			}
			if err := json.Unmarshal(resp.Body(), &userInfo); err != nil {
				errChan <- err
				return
			}
			userInfoChan <- struct {
				index int
				user  map[string]interface{}
			}{index: i, user: userInfo.User}
		}(i, listing)
	}

	go func() {
		wg.Wait()
		close(errChan)
		close(userInfoChan)
	}()

	for i := 0; i < len(listings); i++ {
		select {
		case err := <-errChan:
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
				return
			}
		case userInfo := <-userInfoChan:
			listings[userInfo.index]["user"] = userInfo.user
			delete(listings[userInfo.index], "user_id")
		}
	}

	result.Listings = listings
	c.JSON(http.StatusOK, result)

}

func ProxyListingServicePost(c *gin.Context) {
	msUrl := utils.ListingServiceURL + strings.TrimSuffix(strings.TrimPrefix(c.Request.URL.Path, "/public-api"), "/")

	formData, err := convertJSONToFormData(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(formData.Get("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	_, err = utils.Client.R().
		Get(utils.UserServiceURL + "/users/" + strconv.Itoa(userID))

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "User Service unavailable"})
		return
	}

	resp, err := utils.Client.R().
		SetBody(formData.Encode()).
		SetQueryParamsFromValues(c.Request.URL.Query()).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Post(msUrl)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Listing Service unavailable"})
		return
	}

	var result struct {
		Result  bool                   `json:"result"`
		Listing map[string]interface{} `json:"listing"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		log.Printf("Error when parse result %s", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"listing": result.Listing})

}

func ProxyUserServicePost(c *gin.Context) {
	msUrl := utils.UserServiceURL + strings.TrimSuffix(strings.TrimPrefix(c.Request.URL.Path, "/public-api"), "/")

	formData, err := convertJSONToFormData(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := utils.Client.R().
		SetBody(formData.Encode()).
		SetQueryParamsFromValues(c.Request.URL.Query()).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Post(msUrl)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "User Service unavailable"})
		return
	}

	var result struct {
		Result bool                   `json:"result"`
		User   map[string]interface{} `json:"user"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		log.Printf("Error when parse result %s", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"user": result.User})

}

func convertJSONToFormData(c *gin.Context) (url.Values, error) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	formData := url.Values{}
	for key, value := range jsonData {
		formData.Set(key, fmt.Sprintf("%v", value))
	}

	return formData, nil
}
