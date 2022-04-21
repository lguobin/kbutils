package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lguobin/kbutils"
	"github.com/lguobin/kbutils/kbfile"
)

func openFile(filepa string, archive bool) (file *os.File, fileName, fileDir string, err error) {
	fullPath := kbfile.RealPath(filepa)
	fileDir, fileName = filepath.Split(fullPath)
	if archive {
		archiveName := kbutils.FormatTime(time.Now(), "Y-m-d")
		ext := filepath.Ext(fileName)
		base := strings.TrimSuffix(fileName, ext)
		fileDir = kbfile.RealPathMkdir(fileDir+base, true)
		fileName = archiveName + ext
		fullPath = fileDir + fileName
	}
	_ = mkdirLog(fileDir)
	if kbfile.FileExist(fullPath) {
		file, err = os.OpenFile(fullPath, os.O_APPEND|os.O_RDWR, 0644)
	} else {
		file, err = os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	}
	if err != nil {
		return nil, "", "", err
	}

	return
}

// SetFile Setting log file output
func (log *Logger) SetFile(filepath string, archive ...bool) {
	log.DisableConsoleColor()
	logArchive := len(archive) > 0 && archive[0]
	if logArchive {
		fmt.Println("暂时不做定时触发日志")
		fmt.Println("暂时不做定时触发日志")
		// c := cron.New()
		// _, _ = c.Add("0 0 * * *", func() {
		// 	log.setLogfile(filepath, logArchive)
		// })
		// c.Run()
	}
	log.setLogfile(filepath, logArchive)
}

func (log *Logger) setLogfile(filepath string, archive bool) {
	fileObj, fileName, fileDir, _ := openFile(filepath, archive)
	log.mu.Lock()
	defer log.mu.Unlock()
	log.CloseFile()
	log.file = fileObj
	log.fileDir = fileDir
	log.fileName = fileName
	if log.fileAndStdout {
		log.out = io.MultiWriter(log.file, os.Stdout)
	} else {
		log.out = fileObj
	}
}

func (log *Logger) Discard() {
	log.mu.Lock()
	log.out = ioutil.Discard
	if log.file != nil {
		_ = log.file.Close()
	}
	log.level = LogNot
	log.mu.Unlock()
}

func (log *Logger) SetSaveFile(filepath string, archive ...bool) {
	log.SetFile(filepath, archive...)
	log.mu.Lock()
	log.fileAndStdout = true
	log.out = io.MultiWriter(log.file, os.Stdout)
	log.mu.Unlock()
}

func (log *Logger) CloseFile() {
	if log.file != nil {
		_ = log.file.Close()
		log.file = nil
		log.out = os.Stderr
	}
}

// func oldLogFile(fileDir, fileName string) string {
// 	ext := path.Ext(fileName)
// 	name := strings.TrimSuffix(fileName, ext)
// 	timeStr := time.Now().Format("2006-01-02")
// 	oldLogFile := fileDir + "/" + name + "." + timeStr + ext
// judge:
// 	for {
// 		if !kbfile.FileExist(oldLogFile) {
// 			break judge
// 		} else {
// 			oldLogFile = fileDir + "/" + name + "." + timeStr + "_" + strconv.Itoa(int(time.Now().UnixNano())) + ext
// 		}
// 	}
//
// 	return oldLogFile
// }

func mkdirLog(dir string) (e error) {
	if kbfile.DirExist(dir) {
		return
	}
	if err := os.MkdirAll(dir, 0775); err != nil && os.IsPermission(err) {
		e = err
	}
	return
}
