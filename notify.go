package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"strings"
)

func sendEmail(to []string, subject, body string) {
	go func() {
		from := ""      // sender
		password := "" // app password or real password

		// SMTP server configuration
		smtpHost := "smtp.gmail.com"
		smtpPort := "587"

		// Message
		msg := "From: " + from + "\n" +
			"To: " + fmt.Sprint(to) + "\n" +
			"Subject: " + subject + "\n\n" +
			body

		// Authentication
		auth := smtp.PlainAuth("", from, password, smtpHost)

		// Send email
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
		if err != nil {
			fmt.Println("Failed to send email:", err)
		} else {
			fmt.Println("Email sent successfully")
		}
	}()
}

const fast2smsAPIKey = ""

func sendSMSAsync(numbers []string, message string) {
	go func() {
		// Join the list into a single comma-separated string
		numbersStr := strings.Join(numbers, ",")

		payload := map[string]string{
			"sender_id": "TXTIND", // must be 6 characters
			"message":   message,
			"language":  "english",
			"route":     "q", // use "q" for transactional messages
			"numbers":   numbersStr,
		}

		data, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "https://www.fast2sms.com/dev/bulkV2", bytes.NewBuffer(data))
		req.Header.Set("authorization", fast2smsAPIKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("SMS Error:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Println("SMS sent to", numbersStr, "Status:", resp.Status, "Response:", string(body))
	}()
}
