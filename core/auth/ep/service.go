package ep

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Mboukhal/SvGoPg/cmd/settings"
	sqlc "github.com/Mboukhal/SvGoPg/internal/db"
	"github.com/golang-jwt/jwt/v5"
)

var PASSWORD_HASH_KEY = ""

func init() {
	PASSWORD_HASH_KEY = os.Getenv("PASSWORD_HASH_KEY")
	if PASSWORD_HASH_KEY == "" {
		log.Fatalln("PASSWORD_HASH_KEY not set in environment variables")
	}
}

// const magicLinkExpiryMinutes = 3

// func sendTokenEmail(to string, token string) error {
// 	// implement email sending logic here
// 	auth_link := os.Getenv("APP_DOMAIN") + "/auth?token=" + token

// 	log.Println("Sending token email to:", to, "with token:", auth_link)
// 	err := email.SendEmailSys(to, "Your Login Link", "Click the following link to log in: \n"+auth_link)

// 	if err != nil {
// 		log.Println("Error sending email:", err)
// 		return err
// 	}

// 	log.Println("Email sent successfully to:", to)

// 	return nil
// }

func checkEmailAuthorization(ctx context.Context, email string) error {

	// log.Println("Checking email authorization for:", email)
	// check if email has profile in db
	// Get queries from context
	queries := settings.GetQueries(ctx)
	// log.Println("Queries in login link handler:", queries)
	if queries == nil {
		return http.ErrServerClosed
	}

	count, err := queries.CheckUserEmailExists(ctx, email)
	// log.Println("Email authorization check count:", count)
	if err != nil {
		return err
	}
	if count == false {
		return errors.New("email not authorized")
	}
	return nil
}

func HashPassword(password string) (string, error) {

	// 1. Convert the key and message to byte slices.
	keyBytes := []byte(PASSWORD_HASH_KEY)
	messageBytes := []byte(password)

	// 2. Create a new HMAC hash using the sha256 algorithm and the key.
	h := hmac.New(sha256.New, keyBytes)

	// 3. Write the message to the hash.
	h.Write(messageBytes)

	// 4. Get the final hash sum as a byte slice.
	// The Sum(nil) method appends the current hash to the provided slice (nil in this case) and returns the result.
	sum := h.Sum(nil)

	// 5. Encode the byte slice to a hex string for easy printing/storage.
	return hex.EncodeToString(sum), nil
}

func VerifyPassword(password string, hash string) (bool, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return false, err
	}
	// println("Verifying password. Hashed input:", hashedPassword, "Stored hash:", hash, "Password:", password, "Match:", hashedPassword == hash)
	return hashedPassword == hash, nil
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func generateJWTToken(userData JWTClaims) (string, error) {
	// Create a new JWT token with user ID and role as claims
	claims := jwt.MapClaims{
		"user_id": userData.UserID,
		"email":   userData.Email,
		"role":    userData.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		return "", errors.New("JWT secret key not set in environment variables")
	}
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateJWTToken(tokenString string) bool {
	_, err := parseJWTToken(tokenString)
	if err != nil {
		return false
	}
	// println("Validated JWT claims:", claims.UserID, claims.Email, claims.Role)
	return true

}

func parseJWTToken(tokenString string) (*JWTClaims, error) {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		return nil, errors.New("JWT secret key not set in environment variables")
	}
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")

}

func CreateUser(ctx context.Context, email string, password string, name string) error {
	q := settings.GetQueries(ctx)
	if q == nil {
		return http.ErrServerClosed
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}

	err = q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:    email,
		Password: passwordHash,
		Username: name,
	})
	if err != nil {
		return err
	}
	return nil
}
