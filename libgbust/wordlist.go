package libgbust

import (
	"fmt"
	"os"
)

func (a *Attacker) wordlistsExists() error {
	if len(a.config.Wordlists) < 1 {
		return fmt.Errorf("[!] at least one wordlist (--wordlist / -w) must be specified")
	}
	for _, wordlist := range a.config.Wordlists {
		_, err := os.Stat(wordlist)
		if err != nil {
			return fmt.Errorf("[!] wordlist %s does not exist", wordlist)
		}
	}

	return nil
}
