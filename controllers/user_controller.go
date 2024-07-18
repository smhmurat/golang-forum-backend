package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang-forum-backend/internal/models"
	"golang-forum-backend/services"
	"golang-forum-backend/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8082/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8082/auth/github/callback",
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes: []string{
			"user:email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
	oauthStateString = "randomstring"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := services.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func SignUpWithGoogle(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func SignUpWithGoogleCallback(w http.ResponseWriter, r *http.Request) {
	var user models.User
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Kod eksik", http.StatusBadRequest)
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token alma hatasi "+err.Error(), http.StatusBadRequest)
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Kullanıcı bilgileri alma hatası", http.StatusInternalServerError)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		http.Error(w, "Google OAuth ile giriş yapılamadı", http.StatusInternalServerError)
		return
	}

	userInfo := struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Verified bool   `json:"verified_email"`
		Picture  string `json:"picture"`
		Name     string `json:"name"`
	}{}
	user.ID = uuid.New().String()
	user.Email = userInfo.Email
	user.Username = userInfo.Name
	//user.Token = token.AccessToken

	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Kullanıcı bilgileri alma hatası", http.StatusInternalServerError)
		return
	}

	var existingUser string

	//NEW USER
	if existingUser == "" {
		user.Email = userInfo.Email
		user.Username = utils.MergeLowercaseAndAddRandomNumber(userInfo.Name)
		user.Password = "" // TODO: Hash password

		err = services.CreateUser(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newUser, err := services.GetUserByEmail(userInfo.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var authProvider models.AuthProvider
		authProvider.UserID = newUser.ID
		authProvider.Provider = "google"
		authProvider.ProviderID = userInfo.ID
		err = services.NewAuthProvider(&authProvider)
		if err != nil {
			http.Error(w, "Auth provider ekleme hatası 1"+err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := services.LoginUserWithGoogle(userInfo.Email)
		if err != nil {
			http.Error(w, "Email address or password is wrong.", http.StatusUnauthorized)
			return
		}

		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{
			Name:     "forum_session",
			Value:    t,
			Expires:  expiration,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "http://localhost:8081", http.StatusTemporaryRedirect)
	} else {
		fmt.Println("user email", user.Email)
	}
	//else if existingUser != "" {
	//	newUser, err := services.GetUserByEmail(userInfo.Email)
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//
	//	authProvider, err := services.GetAuthProviderByUserID(newUser.ID)
	//	if err != nil || errors.Is(err, sql.ErrNoRows) {
	//		http.Error(w, "Auth provider bulunamadı 2"+err.Error(), http.StatusInternalServerError)
	//		fmt.Println(authProvider)
	//		return
	//	}
	//	//if authProvider == "" {
	//	//	var authProvider models.AuthProvider
	//	//	authProvider.UserID = newUser.ID
	//	//	authProvider.Provider = "google"
	//	//	authProvider.ProviderID = userInfo.ID
	//	//	err = services.NewAuthProvider(&authProvider)
	//	//	if err != nil {
	//	//		http.Error(w, "Auth provider ekleme hatası 1"+err.Error(), http.StatusInternalServerError)
	//	//		return
	//	//	}
	//	//}
	//	if authProvider.Provider == "google" {
	//		t, err := services.LoginUserWithGoogle(userInfo.Email)
	//		if err != nil {
	//			http.Error(w, "Email address or password is wrong.", http.StatusUnauthorized)
	//			return
	//		}
	//
	//		expiration := time.Now().Add(24 * time.Hour)
	//		cookie := http.Cookie{
	//			Name:     "forum_session",
	//			Value:    t + "Deneme", // Veya gerçek oturum belirteci
	//			Expires:  expiration,
	//			HttpOnly: true,
	//			Path:     "/",
	//		}
	//		http.SetCookie(w, &cookie)
	//
	//		http.Redirect(w, r, "http://localhost:8081", http.StatusTemporaryRedirect)
	//	} else {
	//		var authProvider models.AuthProvider
	//		authProvider.UserID = newUser.ID
	//		authProvider.Provider = "google"
	//		authProvider.ProviderID = userInfo.ID
	//		err = services.NewAuthProvider(&authProvider)
	//		if err != nil {
	//			http.Error(w, "Auth provider ekleme hatası"+err.Error(), http.StatusInternalServerError)
	//			return
	//		}
	//	}
	//}

}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := services.LoginUser(&user)
	if err != nil {
		http.Error(w, "Email address or password is wrong.", http.StatusUnauthorized)
		return
	}

	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    token, // Veya gerçek oturum belirteci
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	// Yanıt gönderme
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if r.FormValue("state") != oauthStateString {
		log.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("code exchange failed: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := githubOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Println("could not get user info: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Println("could not decode user info: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	userData, err := json.Marshal(user)
	if err != nil {
		log.Println("could not marshal user info: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	u.ID = uuid.New().String()
	u.Email = user["email"].(string)
	u.Username = utils.MergeLowercaseAndAddRandomNumber(user["login"].(string))
	u.Password = ""
	//user.Token = token.AccessToken

	var existingUser string

	//NEW USER
	if existingUser == "" {
		err = services.CreateUser(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newUser, err := services.GetUserByEmail(u.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var authProvider models.AuthProvider
		authProvider.UserID = newUser.ID
		authProvider.Provider = "github"
		authProvider.ProviderID = u.ID
		err = services.NewAuthProvider(&authProvider)
		if err != nil {
			http.Error(w, "Auth provider ekleme hatası 1"+err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := services.LoginUserWithGoogle(u.Email)
		if err != nil {
			http.Error(w, "Email address or password is wrong.", http.StatusUnauthorized)
			return
		}

		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{
			Name:     "forum_session",
			Value:    t,
			Expires:  expiration,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "http://localhost:8081", http.StatusTemporaryRedirect)

		// Redirect or respond to the client
		w.Header().Set("Content-Type", "application/json")
		w.Write(userData)
	}
}
