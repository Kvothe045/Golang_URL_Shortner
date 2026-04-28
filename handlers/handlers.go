// url_shortner/handlers/handlers.go
package handlers

import (
	"crypto/sha256"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type URL struct {
	URL string `json:"url" binding:"required,url"`
}

func generateShortCode(url string) string {
	//base62
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hash := sha256.Sum256([]byte(url))

	var num int64
	for i := 0; i < 8; i++ {
		num = (num << 8) | int64(hash[i])
	}

	if num < 0 {
		num = -num
	}

	shortCode := ""
	for num > 0 {
		shortCode = string(charset[num%62]) + shortCode
		num /= 62
	}
	shortCode = shortCode[:6]
	return shortCode

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
		c.Redirect(http.StatusFound, url)
	}
}
