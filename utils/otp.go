package utils

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const OTPExpirationMinutes = 10

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())

	otp:=rand.Intn(900000) + 100000

	return fmt.Sprintf("%d",otp)

}


func ClaculateOTPExpiry() time.Time {
	return time.Now().Add(OTPExpirationMinutes * time.Minute)
}

func SendEmail(recipientEmail,otpCode string) error{
		// --- EMAIL SIMULATION LOG ---
	log.Printf("------------------------")
	log.Printf("SIMULATION: Sending Password Reset OTP Email")
	log.Printf("To: %s", recipientEmail)
	log.Printf("Code: %s", otpCode)
	log.Printf("Expires: %d minutes", OTPExpirationMinutes)
	log.Printf("------------------------")
	
	// **TODO:** Replace this with actual email sending using a library (e.g., gomail)
    // or a service API (SendGrid, Mailgun, etc.).
	
	return nil // Return nil to simulate success
}
