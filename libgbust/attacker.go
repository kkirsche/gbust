package libgbust

import (
	"bufio"
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// Attacker represents the client we use to
type Attacker struct {
	client   *http.Client
	config   *Config
	scanner  *bufio.Scanner
	context  context.Context
	workCh   chan string
	resultCh chan *Result
	signalCh chan os.Signal
	words    sync.WaitGroup
	Wg       sync.WaitGroup
}

// Result is a struct that wraps the details of the check. We're using this so
// we can add or remove with minimal refactoring
type Result struct {
	URL        *url.URL
	StatusCode int
	Size       *int64
	Msg        string
	Err        error
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
		"url":       c.URL.String(),
		"verbose":   c.Verbose,
		"wordlists": strings.Join(c.Wordlists, ", "),
	}).Debugln("[+] creating attacker...")

	a := &Attacker{
		client:   &http.Client{},
		config:   c,
		resultCh: make(chan *Result),
		workCh:   make(chan string),
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

	logrus.Debugln("[+] creating work context...")
	ctx, cancel := context.WithCancel(context.Background())
	a.context = ctx
	defer cancel()

	a.StartWorkers()

	for _, wp := range a.config.Wordlists {
		logrus.WithField("wordlist", wp).Debugln("[+] beginning wordlist...")
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
				a.words.Add(1)
				a.workCh <- word
			}
		}
	}
	a.words.Wait()
	logrus.Debugln("[+] exiting...")
}
