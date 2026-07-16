package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

// Konfigurasi SMTP
const (
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
	senderEmail = "adminmantraa@gmail.com"
	senderPass  = "csmttruonthoifgv"
)

// Template HTML untuk email OTP
const otpEmailTemplate = `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Password - Mantra App</title>
    <style>
        body {
            font-family: 'Inter', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f7f7f7;
            margin: 0;
            padding: 0;
            color: #333333;
        }
        .container {
            max-width: 600px;
            margin: 40px auto;
            background-color: #ffffff;
            border-radius: 16px;
            box-shadow: 0 4px 20px rgba(0,0,0,0.05);
            overflow: hidden;
            border: 1px solid #eeeeee;
        }
        .header {
            background-color: #AF510C;
            padding: 30px 20px;
            text-align: center;
        }
        .header h1 {
            color: #ffffff;
            margin: 0;
            font-size: 28px;
            letter-spacing: 1px;
        }
        .content {
            padding: 40px 30px;
            text-align: center;
        }
        .greeting {
            font-size: 20px;
            font-weight: 600;
            margin-bottom: 10px;
        }
        .message {
            font-size: 15px;
            line-height: 1.6;
            color: #666666;
            margin-bottom: 30px;
        }
        .otp-box {
            background-color: #FAEDE4;
            border: 2px dashed #AF510C;
            border-radius: 12px;
            padding: 20px;
            margin: 0 auto 30px auto;
            max-width: 300px;
        }
        .otp-code {
            font-size: 36px;
            font-weight: 700;
            letter-spacing: 6px;
            color: #AF510C;
            margin: 0;
        }
        .warning {
            font-size: 13px;
            color: #999999;
            margin-bottom: 0;
        }
        .footer {
            background-color: #fafafa;
            padding: 20px;
            text-align: center;
            font-size: 12px;
            color: #aaaaaa;
            border-top: 1px solid #eeeeee;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Mantra App</h1>
        </div>
        <div class="content">
            <div class="greeting">Halo, {{.Nama}}!</div>
            <div class="message">
                Kami menerima permintaan untuk mereset password akun Anda di Mantra App. 
                Berikut adalah kode OTP Anda. Kode ini berlaku selama 10 menit.
            </div>
            <div class="otp-box">
                <p class="otp-code">{{.OTP}}</p>
            </div>
            <p class="warning">
                Jika Anda tidak meminta reset password ini, abaikan email ini. 
                Jangan pernah membagikan kode OTP kepada siapapun.
            </p>
        </div>
        <div class="footer">
            &copy; 2026 Mantra App. All rights reserved.
        </div>
    </div>
</body>
</html>
`

// SendOTPEmail mengirimkan kode OTP ke email tujuan dengan template HTML
func SendOTPEmail(toEmail string, otp string, namaLengkap string) error {
	// Jika nama lengkap kosong, fallback
	if namaLengkap == "" {
		namaLengkap = "Pengguna"
	}

	// Menyiapkan template data
	data := struct {
		Nama string
		OTP  string
	}{
		Nama: namaLengkap,
		OTP:  otp,
	}

	// Parse template
	t, err := template.New("email").Parse(otpEmailTemplate)
	if err != nil {
		return fmt.Errorf("gagal parse email template: %w", err)
	}

	// Eksekusi template ke dalam buffer
	var bodyBuffer bytes.Buffer
	if err := t.Execute(&bodyBuffer, data); err != nil {
		return fmt.Errorf("gagal execute email template: %w", err)
	}

	// Menyusun headers dan body email
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	subject := "Subject: Reset Password - Mantra App\r\n"
	to := "To: " + toEmail + "\r\n"
	msg := []byte(to + subject + headers + "\r\n" + bodyBuffer.String())

	// Autentikasi SMTP
	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)

	// Kirim email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("gagal kirim email SMTP: %w", err)
	}

	return nil
}
