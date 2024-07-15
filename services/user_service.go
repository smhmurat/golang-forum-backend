package services

import (
	"github.com/golang-jwt/jwt/v5"
	"golang-forum-backend/internal/models"
	"golang-forum-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

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
	_, err = db.Exec("INSERT INTO users (username, email, password, token ) VALUES (?,?,?,?)", user.Username, user.Email, user.Password, user.Token)
	if err != nil {
		return err
	}
	return nil
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
