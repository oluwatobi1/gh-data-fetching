package utilities_test

import (
	"testing"

	"github.com/oluwatobi1/gh-api-data-fetch/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseLinkHeader(t *testing.T) {
	linkHeader := `<https://api.github.com/repositories/1300192/issues?page=2>; rel="prev", <https://api.github.com/repositories/1300192/issues?page=4>; rel="next", <https://api.github.com/repositories/1300192/issues?page=515>; rel="last", <https://api.github.com/repositories/1300192/issues?page=1>; rel="first"`
	links := utils.ParseLinkHeader(linkHeader)
	assert.Equal(t, "https://api.github.com/repositories/1300192/issues?page=2", links["prev"])
	assert.Equal(t, "https://api.github.com/repositories/1300192/issues?page=4", links["next"])

}
