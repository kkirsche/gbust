package libgbust

import (
	"bufio"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Attacker represents the client we use to
type Attacker struct {
	client  *http.Client
	config  *Config
	scanner *bufio.Scanner
}

// NewAttacker creates an instance of attacker
func NewAttacker(c *Config) (*Attacker, error) {
	logrus.Debugln("[+] parsing url...")
	u, err := url.Parse(c.RawURL)
	if err != nil {
		return nil, err
	}

	c.URL = u

	logrus.WithFields(logrus.Fields{
		"cookies":   c.Cookies,
		"timeout":   c.Timeout,
		"url":       c.URL.String(),
		"verbose":   c.Verbose,
		"wordlists": strings.Join(c.Wordlists, ", "),
	}).Debugln("[+] creating attacker...")

	a := &Attacker{
		client: &http.Client{
			Timeout: time.Duration(c.Timeout),
		},
		config: c,
	}

	return a, nil
}

// Attack is used to begin brute forcing
func (a *Attacker) Attack() {
	logrus.Debugln("[+] beginning attack...")
	err := a.wordlistsExists()
	if err != nil {
		logrus.Errorln(err)
		return
	}

	for _, wp := range a.config.Wordlists {
		wordlist, err := os.Open(wp)
		if err != nil {
			logrus.WithError(err).Fatalln("[!] failed to open wordlist")
			return
		}
		defer wordlist.Close()

		a.scanner = bufio.NewScanner(wordlist)
		for a.scanner.Scan() {
			word := strings.TrimSpace(a.scanner.Text())
			if !strings.HasPrefix(word, "#") && !strings.HasPrefix(word, "//") && len(word) > 0 {
				go a.Check(word)
			}
		}
	}
}

func (a *Attacker) Check(word string, resultCh chan string) {

}
