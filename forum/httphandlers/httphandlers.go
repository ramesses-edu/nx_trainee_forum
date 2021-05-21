package httphandlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"nx_trainee_forum/forum/application/config"
	"nx_trainee_forum/forum/httphandlers/authorization"
	"path/filepath"
	"regexp"

	"gorm.io/gorm"
)

var (
	reNum *regexp.Regexp = regexp.MustCompile(`\d+`)
)

func Authentification(cfg *config.Config, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rPath := r.URL.Path
		reGoogleProvider := regexp.MustCompile(`\/auth\/google(\/)??`)
		reFacebookProvider := regexp.MustCompile(`\/auth\/facebook(\/)??`)
		reTwitterProvider := regexp.MustCompile(`\/auth\/twitter(\/)??`)
		reCallback := regexp.MustCompile(`\/auth\/callback(\/)??\w+`)
		switch {
		case reGoogleProvider.Match([]byte(rPath)):
			authorization.AuthGoogle(cfg, w, r)
		case reFacebookProvider.Match([]byte(rPath)):
			authorization.AuthFacebook(cfg, w, r)
		case reTwitterProvider.Match([]byte(rPath)):
			authorization.AuthTwitter(cfg, w, r)
		case reCallback.Match([]byte(rPath)):
			oauthCallback(cfg, db, w, r)
		}
	})
}

func oauthCallback(cfg *config.Config, db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	reProviderGoogle := regexp.MustCompile(`\/auth\/callback\/google(\/)??`)
	reProviderFacebook := regexp.MustCompile(`\/auth\/callback\/facebook(\/)??`)
	reProviderTwitter := regexp.MustCompile(`\/auth\/callback\/twitter(\/)??`)
	switch {
	case reProviderGoogle.Match([]byte(r.URL.Path)):
		authorization.CallbackGoogle(cfg, db, w, r)
	case reProviderFacebook.Match([]byte(r.URL.Path)):
		authorization.CallbackFacebook(cfg, db, w, r)
	case reProviderTwitter.Match([]byte(r.URL.Path)):
		authorization.CallbackTwitter(cfg, db, w, r)
	}
}

func MainHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := authorization.GetCurrentUser(db, r)
		t, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			fmt.Println(err)
		}
		t.Execute(w, u)
	})
}

type myFileSystem struct {
	fs http.FileSystem
}

func (nfs myFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

func PublicHandler() http.Handler {
	return http.StripPrefix("/public/", http.FileServer(myFileSystem{fs: http.Dir("./static")}))
}

func LogoutHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := authorization.GetCurrentUser(db, r)
		if u.ID == 0 {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		u.AccessToken = authorization.CalculateSignature(authorization.GenerateAccessToken(), "provider")
		u.UpdateAccessToken(db)
		cookie := http.Cookie{Name: "UAAT", Path: "/", MaxAge: -1}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func responseXML(r *http.Request) bool {
	if _, ok := r.Form["xml"]; ok {
		return true
	}
	return false
}

func xmlWrite(w http.ResponseWriter, data interface{}) error {
	xmlB, err := xml.MarshalIndent(data, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":""}`))
		return err
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(xmlB))
	return nil
}

func jsonWrite(w http.ResponseWriter, data interface{}) error {
	jsonB, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":""}`))
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonB))
	return nil
}
