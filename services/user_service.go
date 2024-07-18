package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang-forum-backend/internal/models"
	"golang-forum-backend/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file")
	}
}

func CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
	})

	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.Token = tokenString

	db := utils.GetDB()
	_, err = db.Exec("INSERT INTO users (id, username, email, password, token ) VALUES (?,?,?,?,?)", user.ID, user.Username, user.Email, user.Password, user.Token)
	if err != nil {
		return err
	}
	return nil
}

func CreateUserWithAuthProvider(user *models.User) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
	})

	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		return err
	}
	fmt.Println("User info:", tokenString)
	user.Token = tokenString
	db := utils.GetDB()
	_, err = db.Exec("INSERT INTO users (username, email, token) VALUES (?, ?, ?)", user.Username, user.Email, user.Token)
	if err != nil {
		return err
	}
	return nil
}

func NewAuthProvider(provider *models.AuthProvider) error {
	db := utils.GetDB()
	_, err := db.Exec("INSERT INTO auth_providers (user_id, provider, provider_id) VALUES (?,?,?)", provider.UserID, provider.Provider, provider.ProviderID)
	if err != nil {
		return err
	}
	return nil
}

func GetAuthProviderByUserID(userID int) (*models.AuthProvider, error) {
	db := utils.GetDB()

	var provider models.AuthProvider
	err := db.QueryRow("SELECT id, user_id, provider, provider_id FROM auth_providers WHERE user_id =?", userID).Scan(&provider.ID, &provider.UserID, &provider.Provider, &provider.ProviderID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	db := utils.GetDB()

	var user models.User
	err := db.QueryRow("SELECT id, username, email, password, token FROM users WHERE email =?", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Token)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func LoginUser(user *models.User) (string, error) {
	db := utils.GetDB()

	var dbUser models.User
	err := db.QueryRow("SELECT id, username, email, password, token FROM users WHERE email =?", user.Email).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Password, &dbUser.Token)
	if err != nil {
		return "Kullanıcı bulunamadı", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return "Şifre yanlış", err
	}

	return dbUser.Token, nil
}

func LoginUserWithGoogle(email string) (string, error) {
	db := utils.GetDB()

	var dbUser models.User
	err := db.QueryRow("SELECT id, username, email, token FROM users WHERE email =?", email).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Token)
	if err != nil {
		return "Kullanıcı bulunamadı", err
	}

	return dbUser.Token, nil
}
