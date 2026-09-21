package service
//delete all Pro when replce
import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"time"
)
type OTPServicePro interface {
    SendOTP(phone string) (string, error)
    VerifyOTP(sessionID, otp string,otpExpiresAt time.Time) (time.Time, error)
}

type otpServicePro struct {
    apiKey string
}
func NewOTPServicePro(apiKey string) OTPServicePro {
    return &otpServicePro{
        apiKey: apiKey,
    }
}




type otpResponsePro struct {
    Status  string `json:"Status"`
    Details string `json:"Details"`
}

func (s *otpServicePro) SendOTP(phone string) (string, error) {

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

    var result otpResponsePro

    err = json.NewDecoder(response.Body).Decode(&result)

    if err != nil {
        return "", err
    }

    if result.Status != "Success" {
        return "", fmt.Errorf("failed to send OTP: %s", result.Details)
    }

    return result.Details, nil
}

func (s *otpServicePro) VerifyOTP(sessionID string, otp string, otpExpiresAt time.Time) (time.Time, error) {

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

    return otpExpiresAt, nil
}

