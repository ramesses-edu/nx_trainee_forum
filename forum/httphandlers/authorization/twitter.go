package authorization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"nx_trainee_forum/forum/application/config"
	"nx_trainee_forum/forum/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

func AuthTwitter(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	reqTokUrl := cfg.Twitter.ReqTokenURL
	request, err := http.NewRequest(http.MethodPost, reqTokUrl, nil)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	autorize := buildAuthHeader(cfg, http.MethodPost, reqTokUrl, map[string]string{"oauth_callback": cfg.Twitter.RedirectURL})
	request.Header.Set("Authorization", autorize)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	if resp.StatusCode != 200 {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	dataMap := make(map[string]string)
	data := strings.Split(string(respBody), "&")
	for _, v := range data {
		datav := strings.Split(v, "=")
		dataMap[datav[0]] = datav[1]
	}
	stateToken := dataMap["oauth_token"]
	cookie := http.Cookie{Name: "oauthstate", Value: stateToken, Expires: time.Now().Add(5 * time.Minute)}
	http.SetCookie(w, &cookie)
	/////////////////////////////////////////////////////////
	urlAuthUser := cfg.Twitter.AuthURL + "=" + stateToken
	http.Redirect(w, r, urlAuthUser, http.StatusFound)
}

func CallbackTwitter(cfg *config.Config, db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	o_token := r.FormValue("oauth_token")
	o_verifier := r.FormValue("oauth_verifier")
	oauthstate, err := r.Cookie("oauthstate")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	if o_token != (oauthstate.Value) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reqTokUrl := cfg.Twitter.TokenURL
	request, err := http.NewRequest(http.MethodPost, reqTokUrl, bytes.NewBuffer([]byte(fmt.Sprintf("oauth_verifier=%s", o_verifier))))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	autorize := buildAuthHeader(cfg, http.MethodPost, reqTokUrl,
		map[string]string{"oauth_token": o_token, "oauth_verifier": o_verifier})
	request.Header.Set("Authorization", autorize)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	if resp.StatusCode != 200 {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	dataMap := make(map[string]string)
	data := strings.Split(string(respBody), "&")
	for _, v := range data {
		datav := strings.Split(v, "=")
		dataMap[datav[0]] = datav[1]
	}

	reqTokUrl = "https://api.twitter.com/1.1/account/verify_credentials.json"
	request, err = http.NewRequest(http.MethodGet, reqTokUrl, nil)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	autorize = buildAuthHeader(cfg, http.MethodGet, reqTokUrl,
		map[string]string{"oauth_token": dataMap["oauth_token"]})
	request.Header.Set("Authorization", autorize)
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	if resp.StatusCode != 200 {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//decode answer JSON to map
	var respMap map[string]interface{} = make(map[string]interface{})
	err = json.Unmarshal(respBody, &respMap)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//check request error
	if _, ok := respMap["errors"]; ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//generate new accessToken for user
	accessToken := GenerateAccessToken()
	hashAccToken := CalculateSignature(accessToken, "provider")
	//check user registration
	var u models.User
	result := u.GetUser(db, map[string]interface{}{
		"login":    respMap["id_str"].(string),
		"provider": "twitter",
	})
	//if user not found, register new user
	if result.Error != nil || result.RowsAffected == 0 {
		u = models.User{
			Login:       respMap["id_str"].(string),
			Provider:    "twitter",
			Name:        respMap["name"].(string),
			AccessToken: hashAccToken,
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
