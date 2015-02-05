package conf

import (
	"bitbucket.org/kardianos/osext"
	"github.com/mediocregopher/lever"
	"os"
	"path/filepath"
	"sync"
)

var (
	lev     *lever.Lever
	levLock = sync.RWMutex{}
)

func Parse() error {
	// Change working dir to that of the executable
	exePath, err := osext.Executable()
	if err != nil {
		return err
	}
	exeName := filepath.Base(exePath)
	exeFolder := filepath.Dir(exePath)
	os.Chdir(exeFolder)

	l := lever.New(exeName, nil)
	l.Add(lever.Param{Name: "--host", Default: "127.0.0.1", Description: "host or ip to listen on"})
	l.Add(lever.Param{Name: "--port", Default: "60000", Description: "port to listen on"})
	l.Add(lever.Param{Name: "--log-level", Default: "INFO", Description: "logging level (DEBUG, INFO, WARN, ERROR, FATAL)"})
	l.Add(lever.Param{Name: "--buffer", Default: "52428800", Description: "size in bytes of read and write buffers"})
	l.Add(lever.Param{Name: "--flush", Default: "5", Description: "flush rate in seconds"})
	l.Parse()

	levLock.Lock()
	defer levLock.Unlock()
	lev = l

	return nil
}

func ParamStr(k string) (string, bool) {
	levLock.RLock()
	defer levLock.RUnlock()
	return lev.ParamStr(k)
}

func GetStr(k string) string {
	levLock.RLock()
	defer levLock.RUnlock()
	s, _ := lev.ParamStr(k)
	return s
}

func GetInt(k string) int {
	levLock.RLock()
	defer levLock.RUnlock()
	i, _ := lev.ParamInt(k)
	return i
}
