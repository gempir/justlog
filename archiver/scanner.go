package archiver

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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

			for _, year := range yearFiles {
				monthFiles, err := ioutil.ReadDir(a.logPath + "/" + channelId.Name() + "/" + year.Name())
				if err != nil {
					log.Error(err)
				}

				for _, month := range monthFiles {
					dayFiles, err := ioutil.ReadDir(a.logPath + "/" + channelId.Name() + "/" + year.Name() + "/" + month.Name())
					if err != nil {
						log.Error(err)
					}

					for _, dayOrUserId := range dayFiles {
						if dayOrUserId.IsDir() {
							channelLogFiles, err := ioutil.ReadDir(a.logPath + "/" + channelId.Name() + "/" + year.Name() + "/" + month.Name() + "/" + dayOrUserId.Name())
							if err != nil {
								log.Error(err)
							}

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
