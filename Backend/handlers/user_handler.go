package handlers

import (
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/vihantandon/Coders_Hub/boot"
	"github.com/vihantandon/Coders_Hub/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Validate Name
	if len(strings.TrimSpace(user.Name)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name required"})
		return
	}

	//Validate Email
	if !ValidateMail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email Format"})
		return
	}

	//Validate Password Length
	if len(user.Password) < 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be atleast 7 character long"})
		return
	}

	//Email already exist
	var existing models.User
	if err := boot.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}

	//Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Password = string(hashed)

	dberr := boot.DB.Create(&user).Error
	if dberr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Validate Email format
	if !ValidateMail(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email format"})
		return
	}

	//Validate Password not empty
	if len(input.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password cant be empty"})
		return
	}

	err := boot.DB.Where("email = ?", input.Email).First(&user).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Email or Password"})
		return
	}

	//Compare hash password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	//Issue JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": tokenString})
}

func ValidateMail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
