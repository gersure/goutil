package signal

import (
"os"
"os/signal"
"runtime"
"syscall"

log "github.com/sirupsen/logrus"
)


func SetupSignalHandler(shudownFunc func(bool)) {
	usrDefSignalChan := make(chan os.Signal, 1)

	signal.Notify(usrDefSignalChan, syscall.SIGUSR1)
	go func() {
		buf := make([]byte, 1<<16)
		for {
			sig := <-usrDefSignalChan
			if sig == syscall.SIGUSR1 {
				stackLen := runtime.Stack(buf, true)
				log.Printf("\n=== Got signal [%s] to dump goroutine stack. ===\n%s\n=== Finished dumping goroutine stack. ===\n", sig, buf[:stackLen])
			}
		}
	}()

	closeSignalChan := make(chan os.Signal, 1)
	signal.Notify(closeSignalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-closeSignalChan
		log.Infof("Got signal [%s] to exit.", sig)
		shudownFunc(sig == syscall.SIGQUIT)
	}()
}
