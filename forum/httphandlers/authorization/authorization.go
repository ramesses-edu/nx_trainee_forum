package authorization

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"nx_trainee_forum/forum/application/config"
	"nx_trainee_forum/forum/models"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func generateOauthStateProvider() string {
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return state
}

func buildAuthHeader(cfg *config.Config, method, path string, params map[string]string) string {
	vals := url.Values{}
	vals.Add("oauth_consumer_key", cfg.Twitter.TwitterAPIKey)
	vals.Add("oauth_nonce", generateNonce())
	vals.Add("oauth_signature_method", "HMAC-SHA1")
	vals.Add("oauth_timestamp", strconv.Itoa(int(time.Now().Unix())))
	vals.Add("oauth_token", cfg.Twitter.TwitterTokenKey)
	vals.Add("oauth_version", "1.0")
	for k, v := range params {
		vals.Set(k, v)
	}
	parameterString := strings.Replace(vals.Encode(), "+", "%20", -1)
	signatureBase := strings.ToUpper(method) + "&" + url.QueryEscape(path) + "&" + url.QueryEscape(parameterString)
	signingKey := url.QueryEscape(cfg.Twitter.TwitterAPISecret) + "&" + url.QueryEscape(cfg.Twitter.TwitterTokenSecret)
	signature := CalculateSignature(signatureBase, signingKey)
	vals.Add("oauth_signature", signature)
	returnString := "OAuth"
	for k := range vals {
		returnString += fmt.Sprintf(` %s="%s",`, k, url.QueryEscape(vals.Get(k)))
	}
	return strings.TrimRight(returnString, ",")
}

func generateNonce() string {
	const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 48)
	for i := range b {
		b[i] = allowed[rand.Intn(len(allowed))]
	}
	return string(b)
}

func CalculateSignature(base, key string) string {
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(base))
	signature := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}

func GenerateAccessToken() string {
	b := make([]byte, 64)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return state
}

func GetCurrentUser(cfg *config.Config, db *gorm.DB, r *http.Request) models.User {
	accessToken := ""
	apiKey := ""
	var u models.User = models.User{}
	accessTokenCookie, err := r.Cookie("UAAT")
	if err == nil {
		accessToken = accessTokenCookie.Value
	}
	apiKey = r.Header.Get("APIKey")
	if accessToken != "" {
		hashAccTok := CalculateSignature(accessToken, cfg.HASHKey)
		result := u.GetUser(db, map[string]interface{}{
			"access_token": hashAccTok,
		})
		if result.Error != nil || result.RowsAffected == 0 {
			u = models.User{}
		}
	}
	if apiKey != "" {
		hashApiKey := CalculateSignature(apiKey, cfg.HASHKey)
		result_apiKey := u.GetUser(db, map[string]interface{}{
			"apikey": hashApiKey,
		})
		if result_apiKey.Error != nil || result_apiKey.RowsAffected == 0 {
			u = models.User{}
		}
	}
	return u
}
