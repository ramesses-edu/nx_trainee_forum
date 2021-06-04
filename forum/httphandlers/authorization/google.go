package authorization

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"nx_trainee_forum/forum/application/config"
	"nx_trainee_forum/forum/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

func AuthGoogle(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	URL, err := url.Parse(cfg.Google.Config.Endpoint.AuthURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//create stateToken for CSFR protect
	oauthStateGoogle := generateOauthStateProvider()
	cookie := http.Cookie{Name: "oauthstate", Value: oauthStateGoogle, Expires: time.Now().Add(5 * time.Minute)}
	http.SetCookie(w, &cookie)
	/////////////////////////////////////////////////////////////////////
	parameters := url.Values{}
	parameters.Add("client_id", cfg.Google.Config.ClientID)
	parameters.Add("scope", strings.Join(cfg.Google.Config.Scopes, " "))
	parameters.Add("redirect_uri", cfg.Google.Config.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateGoogle)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	//redirect to provider Authentification
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackGoogle(cfg *config.Config, db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	oauthstate, err := r.Cookie("oauthstate")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//verify stateTokens
	if state != (oauthstate.Value) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//exchange code to provider Access&Refresh tokens
	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	token, err := cfg.Google.Config.Exchange(context.Background(), code)
	if err != nil {
		return
	}
	//get userinfo on provider resource
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//decode answer JSON to map
	var respMap map[string]interface{} = make(map[string]interface{})
	err = json.Unmarshal(response, &respMap)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//check request error
	if _, ok := respMap["error"]; ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//generate new accessToken for user
	accessToken := GenerateAccessToken()
	hashAccToken := CalculateSignature(accessToken, cfg.HASHKey)
	//check user registration
	var u models.User
	result := u.GetUser(db, map[string]interface{}{
		"login":    respMap["id"],
		"provider": "google",
	})
	//if user not found, register new user
	if result.Error != nil || result.RowsAffected == 0 {
		u = models.User{
			Login:       respMap["id"].(string),
			Provider:    "google",
			Name:        respMap["name"].(string),
			AccessToken: hashAccToken,
			APIKey:      CalculateSignature(GenerateAccessToken(), cfg.HASHKey),
		}
		result = u.CreateUser(db)
	} else {
		u.AccessToken = hashAccToken
		u.UpdateAccessToken(db)
	}
	//write cookies
	if result.Error == nil {
		var expiration = time.Now().Add(30 * 24 * time.Hour)
		cookieUID := http.Cookie{Name: "UAAT", Value: accessToken, Expires: expiration, Path: "/"}
		http.SetCookie(w, &cookieUID)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
