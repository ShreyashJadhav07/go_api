package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"auth_api/database"
	"auth_api/models"
	"auth_api/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


func SignUp(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
	
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format. Check email and ensure password is at least 8 characters."})
		return
	}


	if err:= user.Validate();err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":err.Error(),
		})
		return

	}

	if user.Password !=user.ConfirmPassword{
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "Password and Confirm Password do not match"})
		return
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password securely."})
		return
	}


	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`


	db := database.GetDB()

	
	var newID int
	err = db.QueryRow(query, user.Email, hashedPassword).Scan(&newID)
	
	if err != nil {
		
		log.Printf("Database insertion error: %v", err)
	
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered. Please use a different email."})
		return
	}


	tokenString,err :=utils.GenerateToken(user.Email)
	if err !=nil{
		log.Printf("Error Generating JWT on signup: %v",err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to create session token after registration."})
			return
	}
	c.SetCookie("token", tokenString, 60*60*24, "/", "localhost" ,false ,true)
	


	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id":      newID,
		"email":   user.Email,
	})
}


func Login(c *gin.Context){

	var loginReq models.LoginRequest

	if err :=c.ShouldBindJSON(&loginReq);
	err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid input. Provide valid email and password"})
			return
	}

	db :=database.GetDB()

	var storedHash string
	var userID int
	query:= "SELECT id, password_hash FROM users WHERE email=$1 LIMIT 1"
	err := db.QueryRow(query,loginReq.Email).Scan(&userID, &storedHash)

	if err!=nil{
		log.Printf("Login Failed:Email not found (%s)",loginReq.Email)
		c.JSON(http.StatusUnauthorized ,gin.H{
			"error": "Invalid email or password"})

			return
	}

	if err:=bcrypt.CompareHashAndPassword([] byte(storedHash), []byte(loginReq.Password));
	err!=nil{
		log.Printf("Login failed: wrong password for %s",loginReq.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid and password"})
		return
	}

       tokenString,err :=utils.GenerateToken(loginReq.Email)
	    if err !=nil{
		log.Printf("Error Generating JWT on signup: %v",err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to create session token after registration."})
			return
	     }
	    c.SetCookie("token", tokenString, 60*60*24, "/", "localhost" ,false ,true)

		c.JSON(http.StatusOK,gin.H{
			"message": "Login successful",
			"user_id":userID,
			"email":loginReq.Email,
			"token": tokenString,
		})

	

}





func ForgotPassword(c *gin.Context){

	var req models.ForgotPasswordRequest
	if err :=c.ShouldBindJSON(&req) ; err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid email format."})
			return
	}

	db:=database.GetDB()
	var userID int
	var storedEmail string

	query:="SELECT id, email FROM users WHERE email=$1 LIMIT 1"
	err:=db.QueryRow(query,req.Email).Scan(&userID, &storedEmail)

	if err!=nil{
		if err == sql.ErrNoRows{
			c.JSON(http.StatusOK,gin.H{"message":"If the email is registered,an otp has been sent."})
			return
		}
		log.Printf("DB error checking user for forgot Password: %v",err)
		c.JSON(http.StatusInternalServerError,gin.H{"error":"An internal error occurred"})
		return
	}

	otpCode:=utils.GenerateOTP()
	otpExpiry:=utils.ClaculateOTPExpiry()


	updateQuery := `
	UPDATE users
	SET otp_code =$1,otp_expires_at = $2
	WHERE id =$3
	`
	_,err =db.Exec(updateQuery,otpCode,otpExpiry,userID)
	if err!=nil{
		log.Printf("DB error saving OTP for user %d: %v",userID,err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to generate reset code."})
			return
	}

	go func ()  {
		if emailErr:=utils.SendEmail(storedEmail,otpCode);
		emailErr !=nil {
			log.Printf("CRITICAL:Failed to send OTP email to %s: %V",storedEmail,emailErr)
		}
		
	}()

		c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully. Check your email.",
		
		 "otp": otpCode, 
	})


}



func Profile(c * gin.Context){
	userEmail,exists:=c.Get("userEmail")
	if !exists{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Failed to get user from context"})
			return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"profile fetched succesfully",
		"email": userEmail,
	})
}

func Logout(c *gin.Context) {
    
    c.SetCookie("token", "", -1, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

	
