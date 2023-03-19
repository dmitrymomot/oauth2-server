package validator

import (
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

// Email is a validator that checks if the value is a valid email address.
func ValidateEmail(s string) error {
	email, err := parseEmail(s)
	if err != nil {
		return err
	}

	// if err := email.ValidateHost(); err != nil {
	// 	return err
	// }

	if err := email.ValidateIcanSuffix(); err != nil {
		return err
	}

	return nil
}

var (
	// rfc5322 is a RFC 5322 regex, as per: https://stackoverflow.com/a/201378/5405453.
	// Note that this can't verify that the address is an actual working email address.
	// Use ValidateHost as a starter and/or send them one :-).
	rfc5322          = "(?i)(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"
	validEmailRegexp = regexp.MustCompile(fmt.Sprintf("^%s*$", rfc5322))
)

// EmailAddress is a structure that stores the address local-part@domain parts.
type EmailAddress struct {
	// LocalPart usually the username of an email address.
	LocalPart string

	// Domain is the part of the email address after the last @.
	// This should be DNS resolvable to an email server.
	Domain string
}

func (e EmailAddress) String() string {
	if e.LocalPart == "" || e.Domain == "" {
		return ""
	}

	return fmt.Sprintf("%s@%s", e.LocalPart, e.Domain)
}

// ValidateHost will test if the email address is actually reachable. It will first try to resolve
// the host and then start a mail transaction.
func (e EmailAddress) ValidateHost() error {
	host, err := lookupHost(e.Domain)
	if err != nil {
		return err
	}

	return tryHost(host, e)
}

// ValidateIcanSuffix will test if the public suffix of the domain is managed by ICANN using
// the golang.org/x/net/publicsuffix package. If not it will return an error. Note that if this
// method returns an error it does not necessarily mean that the email address is invalid. Also the
// suffix list in the standard package is embedded and thereby not up to date.
func (e EmailAddress) ValidateIcanSuffix() error {
	d := strings.ToLower(e.Domain)
	if s, icann := publicsuffix.PublicSuffix(d); !icann {
		return fmt.Errorf("public suffix is not managed by ICANN, got %s", s)
	}

	return nil
}

// Parse will parse the input and validate the email locally. If you want to validate the host of
// this email address remotely call the ValidateHost method.
func parseEmail(email string) (*EmailAddress, error) {
	if !validEmailRegexp.MatchString(email) {
		return nil, fmt.Errorf("format is incorrect for %s", email)
	}

	i := strings.LastIndexByte(email, '@')
	e := &EmailAddress{
		LocalPart: email[:i],
		Domain:    email[i+1:],
	}

	return e, nil
}

// lookupHost first checks if any MX records are available and if not, it will check
// if A records are available as they can resolve email server hosts. An error indicates
// that non of the A or MX records are available.
func lookupHost(domain string) (string, error) {
	if mx, err := net.LookupMX(domain); err == nil {
		return mx[0].Host, nil
	}
	if ips, err := net.LookupIP(domain); err == nil {
		return ips[0].String(), nil // randomly returns IPv4 or IPv6 (when available)
	}
	return "", fmt.Errorf("failed finding MX and A records for domain %s", domain)
}

// tryHost will verify if we can start a mail transaction with the host.
func tryHost(host string, e EmailAddress) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, 25), 5*time.Second)
	if err != nil {
		return errors.Wrap(err, "failed to connect to host")
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Hello(e.Domain); err == nil {
		if err = client.Mail(fmt.Sprintf("hello@%s", e.Domain)); err == nil {
			if err = client.Rcpt(e.String()); err == nil {
				client.Reset() // #nosec
				client.Quit()  // #nosec

				return nil
			}
		}
	}

	return err
}
