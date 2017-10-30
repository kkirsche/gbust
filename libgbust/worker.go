package libgbust

import "github.com/sirupsen/logrus"

// StartWorkers is used to create the number of goroutines we will be doing
// work in
func (a *Attacker) StartWorkers() {
	for i := 0; i < a.config.Goroutines-1; i++ {
		logrus.Debugln("[+] creating check worker...")
		a.Wg.Add(1)
		go a.CheckWorker()
	}
	logrus.Debugln("[+] creating result worker...")
	a.Wg.Add(1)
	go a.ResultWorker()
}

// CheckWorker is the goroutine which manages requests to be made
func (a *Attacker) CheckWorker() {
	for {
		select {
		case word := <-a.workCh:
			a.resultCh <- a.CheckDir(word)
		case <-a.context.Done():
			logrus.Debugln("[+] exiting check worker...")
			a.Wg.Done()
			return
		}
	}
}

// ResultWorker ensures that we have a way to print our results as they come
// in from the workers
func (a *Attacker) ResultWorker() {
	for {
		select {
		case r := <-a.resultCh:
			if r.Err != nil {
				logrus.WithError(r.Err).Errorln(r.Msg)
				continue
			}
			logrus.WithFields(logrus.Fields{
				"statusCode": r.StatusCode,
				"size":       r.Size,
				"url":        r.URL.String(),
			}).Infof("[*] found %s with status %d", r.URL.String(), r.StatusCode)
		case <-a.context.Done():
			logrus.Debugln("[+] exiting result worker...")
			a.Wg.Done()
			return
		}
	}
}
