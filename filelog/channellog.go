package filelog

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
	log "github.com/sirupsen/logrus"
)

func (l *FileLogger) LogPrivateMessageForChannel(message twitch.PrivateMessage) error {
	year := message.Time.Year()
	month := int(message.Time.Month())
	day := message.Time.Day()
	err := os.MkdirAll(fmt.Sprintf(l.logPath+"/%s/%d/%d/%d", message.RoomID, year, month, day), 0750)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"/%s/%d/%d/%d/channel.txt", message.RoomID, year, month, day)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(message.Raw + "\n"); err != nil {
		return err
	}

	return nil
}

func (l *FileLogger) LogClearchatMessageForChannel(message twitch.ClearChatMessage) error {
	year := message.Time.Year()
	month := int(message.Time.Month())
	day := message.Time.Day()
	err := os.MkdirAll(fmt.Sprintf(l.logPath+"/%s/%d/%d/%d", message.RoomID, year, month, day), 0750)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"/%s/%d/%d/%d/channel.txt", message.RoomID, year, month, day)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(message.Raw + "\n"); err != nil {
		return err
	}

	return nil
}

func (l *FileLogger) LogUserNoticeMessageForChannel(message twitch.UserNoticeMessage) error {
	year := message.Time.Year()
	month := int(message.Time.Month())
	day := message.Time.Day()
	err := os.MkdirAll(fmt.Sprintf(l.logPath+"/%s/%d/%d/%d", message.RoomID, year, month, day), 0750)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"/%s/%d/%d/%d/channel.txt", message.RoomID, year, month, day)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(message.Raw + "\n"); err != nil {
		return err
	}

	return nil
}

func (l *FileLogger) ReadLogForChannel(channelID string, year int, month int, day int) ([]string, error) {
	filename := fmt.Sprintf(l.logPath+"/%s/%d/%d/%d/channel.txt", channelID, year, month, day)

	if _, err := os.Stat(filename); err != nil {
		filename = filename + ".gz"
	}

	f, err := os.Open(filename)
	if err != nil {
		return []string{}, errors.New("file not found: " + filename)
	}
	defer f.Close()

	var reader io.Reader

	if strings.HasSuffix(filename, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			log.Error(err)
			return []string{}, errors.New("file gzip not readable")
		}
		reader = gz
	} else {
		reader = f
	}

	scanner := bufio.NewScanner(reader)
	if err != nil {
		log.Error(err)
		return []string{}, errors.New("file not readable")
	}

	content := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		content = append(content, line)
	}

	return content, nil
}

func (l *FileLogger) ReadRandomMessageForChannel(channelID string) (string, error) {
	var dayFileList []string
	var lines []string

	if channelID == "" {
		return "", errors.New("missing channelID")
	}

	years, _ := ioutil.ReadDir(l.logPath + "/" + channelID)

	for _, yearDir := range years {
		year := yearDir.Name()
		months, _ := ioutil.ReadDir(l.logPath + "/" + channelID + "/" + year + "/")
		for _, monthDir := range months {
			month := monthDir.Name()

			possibleDays := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

			for _, day := range possibleDays {
				dayDirPath := l.logPath + "/" + channelID + "/" + year + "/" + month + "/" + fmt.Sprint(day)
				logFiles, err := ioutil.ReadDir(dayDirPath)
				if err != nil {
					continue
				}

				for _, logFile := range logFiles {
					logFilePath := dayDirPath + "/" + logFile.Name()
					dayFileList = append(dayFileList, logFilePath)
				}
			}
		}
	}

	if len(dayFileList) < 1 {
		return "", errors.New("no log found")
	}

	randomDayIndex := rand.Intn(len(dayFileList))
	randomDayPath := dayFileList[randomDayIndex]

	f, _ := os.Open(randomDayPath)
	scanner := bufio.NewScanner(f)

	if strings.HasSuffix(randomDayPath, ".gz") {
		gz, _ := gzip.NewReader(f)
		scanner = bufio.NewScanner(gz)
	}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	f.Close()

	if len(lines) < 1 {
		log.Infof("path %s", randomDayPath)
		return "", errors.New("no lines found")
	}

	randomLineNumber := rand.Intn(len(lines))
	return lines[randomLineNumber], nil
}
