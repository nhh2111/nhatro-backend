package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/smtp"
	"os"
	"strings"
)

func GenerateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func VerifyEmailDomain(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("định dạng email không hợp lệ")
	}

	domain := parts[1]
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return errors.New("tên miền email không tồn tại hoặc không thể nhận thư")
	}

	return nil
}

func SendOTPEmail(toEmail string, otpCode string, subject string, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || password == "" {
		return errors.New("hệ thống chưa được cấu hình địa chỉ email gửi đi")
	}

	message := []byte(fmt.Sprintf("Subject: %s\r\n"+
		"To: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n"+
		"%s %s\r\n", subject, toEmail, body, otpCode))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	smtpAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	errSend := smtp.SendMail(smtpAddr, auth, from, []string{toEmail}, message)

	if errSend != nil {
		return errors.New("không thể gửi email lúc này, vui lòng thử lại sau")
	}

	return nil
}
