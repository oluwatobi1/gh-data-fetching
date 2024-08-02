package utils

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ResponseType struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func InfoResponse(ctx *gin.Context, message string, data interface{}, statusCode int) {
	ctx.JSON(statusCode, ResponseType{
		Code:    0,
		Data:    data,
		Message: message,
	})
}

// HandleRateLimit handles GitHub rate limiting by waiting until the rate limit resets
func HandleRateLimit(resp *http.Response) error {
	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter != "" {
		retrySeconds, err := strconv.Atoi(retryAfter)
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(retrySeconds) * time.Second)
		return nil
	}

	rateLimitRemaining := resp.Header.Get("X-RateLimit-Remaining")
	if rateLimitRemaining == "0" {
		rateLimitReset := resp.Header.Get("X-RateLimit-Reset")
		resetTime, err := strconv.ParseInt(rateLimitReset, 10, 64)
		if err != nil {
			return err
		}
		waitDuration := time.Until(time.Unix(resetTime, 0))
		time.Sleep(waitDuration)
		return nil
	}

	return errors.New("unknown rate limit issue")
}

// ParseLinkHeader parses the GitHub link header for pagination
func ParseLinkHeader(header string) map[string]string {
	links := make(map[string]string)
	for _, part := range strings.Split(header, ",") {
		section := strings.Split(strings.TrimSpace(part), ";")
		if len(section) < 2 {
			continue
		}
		url := strings.Trim(section[0], "<>")
		rel := strings.Trim(strings.Split(section[1], "=")[1], "\"")
		links[rel] = url
	}
	return links
}

type DBError struct {
	code int
	err  error
}
