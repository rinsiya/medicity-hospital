package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"time"

	"golang.org/x/crypto/bcrypt"
)
type OTPService interface {
    SendOTP(phone string) (string, error)
    VerifyOTP(sessionID, otp string,otpExpiresAt time.Time) (time.Time, error)
}

type otpService struct {
    apiKey string
}
func NewOTPService(apiKey string) OTPService {
    return &otpService{
        apiKey: apiKey,
    }
}

type otpResponse struct {
    Status  string `json:"Status"`
    Details string `json:"Details"`
}
func (s *otpService) SendOTP(phone string) (string, error) {

    // Development mode
    if os.Getenv("APP_ENV") == "development" {

        otp := "123456"

        hashedOTP, err := bcrypt.GenerateFromPassword(
            []byte(otp),
            bcrypt.DefaultCost,
        )

        if err != nil {
            return "", err
        }

        return string(hashedOTP), nil
    }

    // Production mode → 2Factor
    url := fmt.Sprintf(
        "https://2factor.in/API/V1/%s/SMS/91%s/AUTOGEN",
        s.apiKey,
        phone,
    )

    response, err := http.Get(url)

    if err != nil {
        return "", err
    }

    defer response.Body.Close()

    var result otpResponse

    err = json.NewDecoder(response.Body).Decode(&result)

    if err != nil {
        return "", err
    }

    if result.Status != "Success" {
        return "", fmt.Errorf("failed to send OTP: %s", result.Details)
    }

    return result.Details, nil
}


func (s *otpService) VerifyOTP(
    sessionID string,
    otp string,
    otpExpiresAt time.Time,
) (time.Time, error) {

    // Check expiry
    if time.Now().After(otpExpiresAt) {
        return otpExpiresAt, fmt.Errorf("OTP expired")
    }

    // Development mode
    if os.Getenv("APP_ENV") == "development" {

        err := bcrypt.CompareHashAndPassword(
            []byte(sessionID),
            []byte(otp),
        )

        if err != nil {
            return otpExpiresAt, fmt.Errorf("invalid OTP")
        }

        return otpExpiresAt, nil
    }

    // Production mode → 2Factor
    url := fmt.Sprintf(
        "https://2factor.in/API/V1/%s/SMS/VERIFY/%s/%s",
        s.apiKey,
        sessionID,
        otp,
    )

    response, err := http.Get(url)

    if err != nil {
        return otpExpiresAt, err
    }

    defer response.Body.Close()

    var result otpResponse

    err = json.NewDecoder(response.Body).Decode(&result)

    if err != nil {
        return otpExpiresAt, err
    }

    if result.Status == "Success" {
        return otpExpiresAt, nil
    }

    return otpExpiresAt, fmt.Errorf("invalid OTP")
}