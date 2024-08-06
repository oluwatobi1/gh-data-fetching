package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/types"
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

func ParsePaginationParams(pageStr, pageSizeStr string) (*types.Pagination, error) {

	if pageStr == "" {
		pageStr = "1"
	}

	if pageSizeStr == "" {
		pageSizeStr = "10"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return nil, fmt.Errorf("invalid page value %s", err.Error())
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		return nil, fmt.Errorf("invalid page_size value %s", err.Error())
	}

	return &types.Pagination{
		Page:     page,
		PageSize: pageSize,
	}, err

}

func ValidateDates(startDate, endDate string) error {
	const layout = "2006-01-02"

	if (startDate == "" && endDate != "") || (startDate != "" && endDate == "") {
		return fmt.Errorf("both start_date and end_date must be provided together")
	}

	if startDate == "" || endDate == "" {
		return nil
	}

	start, err := time.Parse(layout, startDate)
	if err != nil {
		return fmt.Errorf("invalid start date format: %v", err)
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		return fmt.Errorf("invalid end date format: %v", err)
	}

	if start.After(end) {
		return errors.New("start date must be before end date")
	}

	return nil
}
