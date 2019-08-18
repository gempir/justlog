package archiver

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func (a *Archiver) scanLogPath() {
	channelFiles, err := ioutil.ReadDir(a.logPath)
	if err != nil {
		log.Error(err)
	}

	for _, channelId := range channelFiles {
		if channelId.IsDir() {
			yearFiles, err := ioutil.ReadDir(a.logPath + "/" + channelId.Name())
			if err != nil {
				log.Error(err)
			}

			yearFiles = a.filterFiles(yearFiles)

			for _, year := range yearFiles {
				if year.IsDir() {
					monthFiles, err := ioutil.ReadDir(a.logPath + "/" + channelId.Name() + "/" + year.Name())
					if err != nil {
						log.Error(err)
					}

					monthFiles = a.filterFiles(monthFiles)

					for _, month := range monthFiles {
						if month.IsDir() {
							dayFiles, err := ioutil.ReadDir(a.logPath + "/" + channelId.Name() + "/" + year.Name() + "/" + month.Name())
							if err != nil {
								log.Error(err)
							}

							dayFiles = a.filterFiles(dayFiles)

							for _, dayOrUserId := range dayFiles {
								if dayOrUserId.IsDir() {
									channelLogFiles, err := ioutil.ReadDir(a.logPath + "/" + channelId.Name() + "/" + year.Name() + "/" + month.Name() + "/" + dayOrUserId.Name())
									if err != nil {
										log.Error(err)
									}

									channelLogFiles = a.filterFiles(channelLogFiles)

									for _, channelLogFile := range channelLogFiles {
										if strings.HasSuffix(channelLogFile.Name(), ".txt") {
											dayInt, err := strconv.Atoi(dayOrUserId.Name())
											if err != nil {
												log.Errorf("Failure converting day to int in scanner %s", err.Error())
												continue
											}

											if dayInt == int(time.Now().Day()) {
												continue
											}

											a.workQueue <- a.logPath + "/" + channelId.Name() + "/" + year.Name() + "/" + month.Name() + "/" + dayOrUserId.Name() + "/" + channelLogFile.Name()
										}
									}

								} else if strings.HasSuffix(dayOrUserId.Name(), ".txt") {
									monthInt, err := strconv.Atoi(month.Name())
									if err != nil {
										log.Errorf("Failure converting month to int in scanner %s", err.Error())
										continue
									}

									if monthInt == int(time.Now().Month()) {
										continue
									}

									a.workQueue <- a.logPath + "/" + channelId.Name() + "/" + year.Name() + "/" + month.Name() + "/" + dayOrUserId.Name()
								}
							}
						}
					}
				}

			}
		}
	}
}

func (a *Archiver) filterFiles(files []os.FileInfo) []os.FileInfo {
	var result []os.FileInfo

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			result = append(result, file)
		}
	}

	return result
}
