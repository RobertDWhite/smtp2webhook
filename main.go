package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"html"

	"github.com/emersion/go-smtp"
	"github.com/namsral/flag"
)

const (
	name      = "smtp2webhook"
	envPrefix = "SMTP2WEBHOOK"
	version   = "2.0"
)

var (
	fs           *flag.FlagSet
	domain       string
	code         string
	tlsCertPath  string
	tlsKeyPath   string
	healthcheck  bool
	printVersion bool
	defaultWebhookURL string
)

var webhooks = make(map[string]string)

func main() {
	fs = flag.NewFlagSetWithEnvPrefix(name, envPrefix, flag.ExitOnError)
	fs.StringVar(&domain, "domain", "localhost", "domain")
	fs.StringVar(&code, "code", "", "secret code")
	fs.StringVar(&tlsCertPath, "tls-cert", "", "TLS certificate path")
	fs.StringVar(&tlsKeyPath, "tls-key", "", "TLS key path")
	fs.BoolVar(&healthcheck, "healthcheck", false, "run healthcheck")
	fs.BoolVar(&printVersion, "version", false, "print version")
	fs.Parse(os.Args[1:])

	if printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	defaultWebhookURL = os.Getenv("WEBHOOK_URL")
	if defaultWebhookURL == "" {
		log.Fatal("WEBHOOK_URL environment variable is not set")
	}

	s := smtp.NewServer(&Backend{})

	if tlsCertPath != "" && tlsKeyPath != "" {
		cert, err := tls.LoadX509KeyPair(tlsCertPath, tlsKeyPath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		s.TLSConfig = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		}
	}
	s.Domain = domain
	s.AllowInsecureAuth = true
	s.AuthDisabled = true
	s.EnableSMTPUTF8 = false

	go func() {
		if s.TLSConfig != nil {
			log.Printf("Listening on :465")
			s.Addr = "[::]:465"
			if err := s.ListenAndServeTLS(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	log.Printf("Listening on :25")
	s.Addr = "[::]:25"
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

type Backend struct{}

func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &Session{}, nil
}

func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return &Session{}, nil
}

type Session struct {
	From       string
	To         string
	WebhookURL string
	Debug      bool
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	s.From = from
	return nil
}

func sanitizeEmailContent(content string) string {
    sanitizedContent := html.EscapeString(content)
    return sanitizedContent
}

func (s *Session) Data(r io.Reader) error {
    log.Println(s.From, "->", s.To)
	
    buf, err := ioutil.ReadAll(r)
    if err != nil {
        log.Println(err)
        return err
    }

    if s.Debug {
        log.Println(string(buf))
    }

    if s.WebhookURL == "" {
        return nil
    }

    resp, err := http.Post(s.WebhookURL, "message/rfc822", bytes.NewReader(buf))
	log.Println("Received email content:", string(buf))
	log.Println("Forwarding email to:", s.WebhookURL)

    if err != nil {
        log.Println("POST", s.WebhookURL, err)
        return err
    }

    log.Println("POST", s.WebhookURL, resp.StatusCode)

    if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
        return nil
    } else {
        return &smtp.SMTPError{
            Code:         450,
            EnhancedCode: smtp.EnhancedCode{4, 5, 0},
            Message:      "Failed to relay message",
        }
    }
}


func (s *Session) Rcpt(to string) error {
	s.To = to

	// Set the forwarding endpoint for all incoming emails
	s.WebhookURL = defaultWebhookURL

	log.Printf("Forwarding email to: %s", s.WebhookURL)

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
