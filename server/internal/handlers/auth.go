package handlers

import (
	"geekCode/internal/auth"
	"geekCode/internal/models"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	
	"geekCode/internal/config"
)

//authHandler
//login

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	FirstName string `json:"firstname" binding:"required"`
	LastName string `json:"lastname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"` //the tags must never have space after binding
}


func (h *Handler) Login (c *gin.Context) {
	cfg := config.LoadConfig() 


	var req LoginRequest;
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return 
	}
	var user models.User
	if err := h.DB.First(&user, "email = ?", req.Email).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error" : "invalid email or password"})
		return
	}
	
	//checking password

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) ; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := auth.GenerateToken(user.ID, cfg.JWTSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error"  : "failed to create tokebn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token" : token,
		"user": user,
	})
}

//register
func (h *Handler) Register(c *gin.Context) {

	//loading env variables
	cfg := config.LoadConfig();

	var RegisterReq RegisterRequest

	fmt.Println("Reached here")

	if err := c.ShouldBindJSON(&RegisterReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error":  err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, "email = ? ", RegisterReq.Email).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error" : "user already exists"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword(([]byte(RegisterReq.Password)), bcrypt.DefaultCost)
	if err!= nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}
	user.FirstName = RegisterReq.FirstName
	user.LastName = RegisterReq.LastName
	user.Username = RegisterReq.Username
	user.Email = RegisterReq.Email
	user.Password = string(hashed)

	if err := h.DB.Create(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	token, err := auth.GenerateToken(user.ID, cfg.JWTSecret);
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}


	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully", "user": user,"token" : token})
}

//profile (getMe in js)
func(h *Handler) GetProfile(c *gin.Context) {
	//first is to get the user id from the context storage
	uid, exists := c.Get("userId") //the 'exists' is only available to use if we use type assertion

	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user id not found in context"})
	}
	userId := uid.(uint) //type assertion

	var currentUser models.User
	if err := h.DB.First(&currentUser, "id = ?", userId).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		
	}



	c.JSON(http.StatusOK, gin.H{
		"currentUser": gin.H{
			"id":         currentUser.ID,
			"first_name": currentUser.FirstName,
			"last_name": currentUser.LastName,
			"username": currentUser.Username,
			"email": currentUser.Email,
		},
	})
	
}