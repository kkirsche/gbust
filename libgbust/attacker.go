package libgbust

import (
	"bufio"
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Attacker represents the client we use to
type Attacker struct {
	client   *http.Client
	config   *Config
	scanner  *bufio.Scanner
	context  context.Context
	cancel   context.CancelFunc
	workCh   chan string
	resultCh chan *Result
	signalCh chan os.Signal
	wg       sync.WaitGroup
}

// Result is a struct that wraps the details of the check. We're using this so
// we can add or remove with minimal refactoring
type Result struct {
	StatusCode int
	Result     string
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
	a.cancel = cancel

	a.StartWorkers()

	for _, wp := range a.config.Wordlists {
		wordlist, err := os.Open(wp)
		if err != nil {
			logrus.WithError(err).Fatalln("[!] failed to open wordlist")
			a.Exit()
			return
		}
		defer wordlist.Close()

		a.scanner = bufio.NewScanner(wordlist)
		for a.scanner.Scan() {
			word := strings.TrimSpace(a.scanner.Text())
			if !strings.HasPrefix(word, "#") && !strings.HasPrefix(word, "//") && len(word) > 0 {
				a.workCh <- word
			}
		}
	}
	a.Exit()
}

// Exit is used to exit and cleanup
func (a *Attacker) Exit() {
	defer func() { //catch any errors
		if err := recover(); err != nil { //catch
			return
		}
	}()
	// Tell the workers to exit
	a.cancel()
	// Wait for them to acknowledge they are done
	a.wg.Wait()
}
