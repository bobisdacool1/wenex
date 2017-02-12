package wenex

import (
	"log"
	"os"
	"path"
	"sync"
)

func newLogger(name string) (func(string) *log.Logger, error) {
	var err error
	var mutex sync.Mutex

	loggers := make(map[string]*log.Logger)

	f := func(name string) *log.Logger {
		mutex.Lock()
		defer mutex.Unlock()

		if logger, ok := loggers[name]; ok {
			return logger
		}

		var file *os.File

		if err = os.MkdirAll(path.Dir(name), 0755); err == nil {
			if file, err = os.OpenFile(name+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644); err != nil {
				file = os.Stdout
			}
		} else {
			file = os.Stdout
		}

		if loggers[""] != nil {
			loggers[name] = log.New(file, loggers[""].Prefix(), loggers[""].Flags())
		} else {
			loggers[name] = log.New(file, "[!] ", log.LstdFlags)
		}

		return loggers[name]
	}

	loggers[""] = f(name)
	return f, err
}
