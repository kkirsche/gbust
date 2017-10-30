package libgbust

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

// PrepareSignalHandler allows for us to catch the control+c signal and exit
// if we see it. We need to clean up, which is why we add this.
func (a *Attacker) PrepareSignalHandler() {
	a.signalCh = make(chan os.Signal, 1)
	signal.Notify(a.signalCh, os.Interrupt)
	go func() {
		for _ = range a.signalCh {
			// caught CTRL+C
			logrus.Warnln("[!] Keyboard interrupt detected, exiting...")
			logrus.Debugln("[+] cancelling workers...")
			a.cancel()
			logrus.Debugln("[+] waiting for cleanup...")
			a.Wg.Wait()
			logrus.Debugln("[+] exiting...")
			os.Exit(130)
		}
	}()
}
