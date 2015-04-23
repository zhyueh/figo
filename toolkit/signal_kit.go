package toolkit

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var signalFuncs map[os.Signal][]func() = nil
var lock sync.Mutex

func processSignal(s os.Signal) {
	for sig, runFuncs := range signalFuncs {
		if sig == s {
			for _, runFunc := range runFuncs {
				log.Println(sig, runFunc)
				runFunc()
			}
			break
		}
	}
}

func SignalWatchRegister(runFunc func(), sig ...os.Signal) {
	lock.Lock()
	defer lock.Unlock()

	if signalFuncs == nil {
		signalFuncs = make(map[os.Signal][]func())
	}

	for _, s := range sig {
		if _, present := signalFuncs[s]; !present {
			signalFuncs[s] = make([]func(), 0)
		}
		signalFuncs[s] = append(signalFuncs[s], runFunc)
	}
}

func SignalWatchRun() {
	ch := make(chan os.Signal, 1)
	for sig, _ := range signalFuncs {
		signal.Notify(ch, sig)
	}
	go func() {
		s := <-ch
		processSignal(s)
	}()
}
