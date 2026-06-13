package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

// SendEnquiryNotificationEmail sends an email notifying the admin about a new enquiry.
// It runs in a separate goroutine to prevent slowing down the request.
func SendEnquiryNotificationEmail(
	host string,
	port int,
	username, password, from, to string,
	useSSL bool,
	enquiryName, enquiryEmail, enquiryPhone, enquiryCompany, enquiryProduct, enquiryMsg string,
	enquiryQty int,
) {
	go func() {
		// Log start of process
		log.Printf("Starting email notification delivery to %s...", to)

		if host == "" || to == "" || from == "" {
			log.Println("WARN: Email notification skipped due to missing SMTP host, from, or to fields.")
			return
		}

		subject := fmt.Sprintf("New Enquiry Received: %s - %s", enquiryName, enquiryProduct)

		// Create email body in HTML format
		body := fmt.Sprintf(`
			<h2>New Business Enquiry</h2>
			<p>A new enquiry has been submitted on the Golden Success website.</p>
			<table border="1" cellpadding="8" style="border-collapse: collapse; border-color: #e2e8f0;">
				<tr style="background-color: #f7fafc;">
					<td><strong>Field</strong></td>
					<td><strong>Value</strong></td>
				</tr>
				<tr>
					<td><strong>Name</strong></td>
					<td>%s</td>
				</tr>
				<tr>
					<td><strong>Email</strong></td>
					<td>%s</td>
				</tr>
				<tr>
					<td><strong>Phone</strong></td>
					<td>%s</td>
				</tr>
				<tr>
					<td><strong>Company</strong></td>
					<td>%s</td>
				</tr>
				<tr>
					<td><strong>Product Requested</strong></td>
					<td>%s</td>
				</tr>
				<tr>
					<td><strong>Quantity</strong></td>
					<td>%d</td>
				</tr>
				<tr>
					<td><strong>Message</strong></td>
					<td>%s</td>
				</tr>
			</table>
			<p style="color: #718096; font-size: 12px; margin-top: 20px;">
				This is an automated notification from the Golden Success Admin Panel. You can manage this enquiry directly in the admin dashboard.
			</p>
		`, enquiryName, enquiryEmail, enquiryPhone, enquiryCompany, enquiryProduct, enquiryQty, enquiryMsg)

		// Message formatting
		msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
		msg += fmt.Sprintf("From: %s\n", from)
		msg += fmt.Sprintf("To: %s\n", to)
		msg += fmt.Sprintf("Subject: %s\n\n", subject)
		msg += body

		addr := fmt.Sprintf("%s:%d", host, port)

		var err error
		if useSSL || port == 465 {
			// TLS Config
			tlsConfig := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         host,
			}

			conn, connErr := tls.Dial("tcp", addr, tlsConfig)
			if connErr != nil {
				log.Printf("ERROR: Failed to dial SMTP server via SSL: %v", connErr)
				return
			}
			defer conn.Close()

			client, clientErr := smtp.NewClient(conn, host)
			if clientErr != nil {
				log.Printf("ERROR: Failed to create SMTP client: %v", clientErr)
				return
			}
			defer client.Quit()

			if username != "" && password != "" {
				auth := smtp.PlainAuth("", username, password, host)
				if authErr := client.Auth(auth); authErr != nil {
					log.Printf("ERROR: SMTP authentication failed: %v", authErr)
					return
				}
			}

			if mailErr := client.Mail(from); mailErr != nil {
				log.Printf("ERROR: SMTP mail command failed: %v", mailErr)
				return
			}

			if rcptErr := client.Rcpt(to); rcptErr != nil {
				log.Printf("ERROR: SMTP rcpt command failed: %v", rcptErr)
				return
			}

			w, wErr := client.Data()
			if wErr != nil {
				log.Printf("ERROR: SMTP data command failed: %v", wErr)
				return
			}

			_, writeErr := w.Write([]byte(msg))
			if writeErr != nil {
				log.Printf("ERROR: Failed to write email body: %v", writeErr)
				return
			}

			w.Close()
			log.Println("Email notification sent successfully via SSL/TLS")
		} else {
			// Plain SMTP or STARTTLS (normally port 587/25)
			var auth smtp.Auth
			if username != "" && password != "" {
				auth = smtp.PlainAuth("", username, password, host)
			}

			err = smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
			if err != nil {
				log.Printf("ERROR: Failed to send email via standard SMTP: %v", err)
				return
			}
			log.Println("Email notification sent successfully via standard SMTP")
		}
	}()
}
