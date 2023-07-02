package middleware

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
	"github.com/radekkrejcirik01/PingMe-backend/services/user/pkg/database"
	"time"
)

func CreateJWT(username string) (string, error) {
	claims := jtoken.MapClaims{
		"username": username,
		"created":  time.Now(),
	}

	// Create token
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS512, claims)

	// Generate encoded token and send it as response.
	return token.SignedString([]byte(database.GetJWTSecret()))
}

func Authorize(c *fiber.Ctx) (string, error) {
	// Get the JWT token from the Authorization header
	authHeader := c.Get("Authorization")
	tokenString := ""

	// Check if the Authorization header is present and starts with "Bearer "
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		// Return an error if the Authorization header is missing or in an invalid format
		return "", fiber.ErrUnauthorized
	}

	// Get the key or public key used for validating the signature
	key := []byte(database.GetJWTSecret())

	// Validate the signature and additional claims
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the key
		if token.Method != jwt.SigningMethodHS512 {
			return nil, fiber.ErrUnauthorized
		}
		return key, nil
	})

	if err != nil {
		// Return an error if the JWT signature is invalid or claims are invalid
		return "", fiber.ErrUnauthorized
	}

	// Validate additional claims as per your requirements
	username, ok := claims["username"].(string)
	if !ok {
		return "", fiber.ErrUnauthorized
	}

	return username, nil
}

// GetHashPassword use sha512 to encrypt the password
func GetHashPassword(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
