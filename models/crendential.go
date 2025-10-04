package models

import "time"

type Credential struct {
	ID                int    `json:"id"`
	UserID            int    `json:"user_id"`
	ServiceName       string `json:"service_name" binding:"required"`
	URL               string `json:"url"`
	UserName          string `json:"username" binding:"required"`
	PasswordCipher    string `json:"-"`
	Nonce             string `json:"-"`
    Notes             string `json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CredentialRequest struct{
	ServiceName string `json:"service_name" binding:"required"`
	URL         string `json:"url"`
	Username    string `json:"username" binding:"required"`
    PasswordCipher string `json:"password_cipher" binding:"required"`
    Nonce          string `json:"nonce" binding:"required"`
	Notes       string `json:"notes"`
}

type CredentialResponse struct{
	ID   int  `json:"id"`
	ServiceName string `json:"service_name"`
	URL  string `json:"url"`
	Username    string    `json:"username"`
    Notes string `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SpecifiCrendentialResponse struct{
	CredentialResponse
	PasswordCipher string `json:"password_cipher"`
    Nonce          string `json:"nonce"`


}