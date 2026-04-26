// url_shortner/handlers/handlers.go
package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type URL struct {
	URL string `json:"url" binding:"required,url"`
}

func generateShortCode(url string) string {

	h := sha256.Sum256([]byte(url))
	return hex.EncodeToString(h[:])[:6]
}

type URLStore interface {
	GetEncodedURL(url string) (string, bool)
	Save(url string, code string)
	GetOriginalURL(code string) (string, bool)
}

func GetEncodedURL(db URLStore, url string) string {
	code, exists := db.GetEncodedURL(url)

	if exists {
		return code
	}

	code = generateShortCode(url)
	db.Save(url, code)
	return code
}

func GetOriginalURL(db URLStore, encoded string) (string, error) {
	url, exists := db.GetOriginalURL(encoded)

	if exists {
		return url, nil
	}
	return "", errors.New("URL not found")
}

func HandleEncodingURL(db URLStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req URL
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		val := GetEncodedURL(db, req.URL)

		c.JSON(http.StatusCreated, gin.H{
			"short_url": "http://localhost:8080/" + val,
		})
	}
}

func HandleDecodingURL(db URLStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		url, err := GetOriginalURL(db, code)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}
