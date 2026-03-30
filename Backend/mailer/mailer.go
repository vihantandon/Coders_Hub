package mailer

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type EmailPayload struct {
	To           string
	UserName     string
	ContestName  string
	Platform     string
	ContestStart string
	ContestURL   string
}

func SendReminderEmail(p EmailPayload) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	if from == "" {
		from = user
	}

	auth := smtp.PlainAuth("", user, pass, host)

	subject := fmt.Sprintf("⏰ Reminder: %s starts soon!", p.ContestName)

	body := buildEmailBody(p)

	msg := strings.Join([]string{
		fmt.Sprintf("From: Coders Hub <%s>", from),
		fmt.Sprintf("To: %s", p.To),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, from, []string{p.To}, []byte(msg))
}

func buildEmailBody(p EmailPayload) string {
	contestLink := p.ContestURL
	if contestLink == "" {
		contestLink = platformURL(p.Platform)
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }
    .card { background: #fff; border-radius: 12px; max-width: 520px; margin: 0 auto; padding: 32px; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
    .header { font-size: 22px; font-weight: bold; color: #1a1a1a; margin-bottom: 8px; }
    .platform { display: inline-block; background: #f0f0f0; border-radius: 6px; padding: 4px 12px; font-size: 13px; color: #555; margin-bottom: 20px; }
    .row { margin: 12px 0; font-size: 15px; color: #333; }
    .label { color: #888; font-size: 13px; margin-bottom: 2px; }
    .value { font-weight: 600; }
    .cta { display: inline-block; margin-top: 24px; background: #4f46e5; color: #fff; text-decoration: none; padding: 12px 28px; border-radius: 8px; font-size: 15px; font-weight: 600; }
    .footer { margin-top: 28px; font-size: 12px; color: #aaa; text-align: center; }
  </style>
</head>
<body>
  <div class="card">
    <div class="header">Hey %s, your contest is coming up!</div>
    <div class="platform">%s</div>
 
    <div class="row">
      <div class="label">Contest</div>
      <div class="value">%s</div>
    </div>
    <div class="row">
      <div class="label">Starts at</div>
      <div class="value">%s (UTC)</div>
    </div>
 
    <a class="cta" href="%s">Go to contest &rarr;</a>
 
    <div class="footer">
      You received this because you set a reminder on Coders Hub.<br>
      To manage reminders, log in to your account.
    </div>
  </div>
</body>
</html>
`, p.UserName, p.Platform, p.ContestName, p.ContestStart, contestLink)
}

func platformURL(platform string) string {
	switch platform {
	case "CodeChef":
		return "https://www.codechef.com/contests"
	case "Codeforces":
		return "https://codeforces.com/contests"
	case "Leetcode":
		return "https://leetcode.com/contest"
	default:
		return "#"
	}
}
