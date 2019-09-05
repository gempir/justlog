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
	"sort"
	"strconv"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	logPath string
}

func NewFileLogger(logPath string) Logger {
	return Logger{
		logPath: logPath,
	}
}

func (l *Logger) LogPrivateMessageForUser(user twitch.User, message twitch.PrivateMessage) error {
	year := message.Time.Year()
	month := int(message.Time.Month())

	err := os.MkdirAll(fmt.Sprintf(l.logPath+"/%s/%d/%d/", message.RoomID, year, month), 0750)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"/%s/%d/%d/%s.txt", message.RoomID, year, month, user.ID)

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

func (l *Logger) LogClearchatMessageForUser(userID string, message twitch.ClearChatMessage) error {
	year := message.Time.Year()
	month := int(message.Time.Month())

	err := os.MkdirAll(fmt.Sprintf(l.logPath+"/%s/%d/%d/", message.RoomID, year, month), 0750)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"/%s/%d/%d/%s.txt", message.RoomID, year, month, userID)

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

func (l *Logger) LogUserNoticeMessageForUser(userID string, message twitch.UserNoticeMessage) error {
	year := message.Time.Year()
	month := int(message.Time.Month())

	err := os.MkdirAll(fmt.Sprintf(l.logPath+"/%s/%d/%d/", message.RoomID, year, month), 0750)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"/%s/%d/%d/%s.txt", message.RoomID, year, month, userID)

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

type userLogFile struct {
	path  string
	year  string
	month string
}

func (l *Logger) GetLastLogYearAndMonthForUser(channelID, userID string) (int, int, error) {
	if channelID == "" || userID == "" {
		return 0, 0, fmt.Errorf("Invalid channelID: %s or userID: %s", channelID, userID)
	}

	logFiles := []userLogFile{}

	years, _ := ioutil.ReadDir(l.logPath + "/" + channelID)

	for _, yearDir := range years {
		year := yearDir.Name()
		months, _ := ioutil.ReadDir(l.logPath + "/" + channelID + "/" + year + "/")
		for _, monthDir := range months {
			month := monthDir.Name()
			path := fmt.Sprintf("%s/%s/%s/%s/%s.txt", l.logPath, channelID, year, month, userID)
			if _, err := os.Stat(path); err == nil {

				logFile := userLogFile{path, year, month}
				logFiles = append(logFiles, logFile)
			} else if _, err := os.Stat(path + ".gz"); err == nil {
				logFile := userLogFile{path + ".gz", year, month}
				logFiles = append(logFiles, logFile)
			}
		}
	}

	sort.Slice(logFiles, func(i, j int) bool {
		yearA, _ := strconv.Atoi(logFiles[i].year)
		yearB, _ := strconv.Atoi(logFiles[j].year)
		monthA, _ := strconv.Atoi(logFiles[i].month)
		monthB, _ := strconv.Atoi(logFiles[j].month)

		if yearA == yearB {
			return monthA > monthB
		} else {
			return yearA > yearB
		}
	})

	if len(logFiles) > 0 {
		year, _ := strconv.Atoi(logFiles[0].year)
		month, _ := strconv.Atoi(logFiles[0].month)

		return year, month, nil
	}

	return 0, 0, errors.New("No logs file")
}

func (l *Logger) ReadLogForUser(channelID, userID string, year int, month int) ([]string, error) {
	if channelID == "" || userID == "" {
		return []string{}, fmt.Errorf("Invalid channelID: %s or userID: %s", channelID, userID)
	}

	filename := fmt.Sprintf(l.logPath+"/%s/%d/%d/%s.txt", channelID, year, month, userID)

	if _, err := os.Stat(filename); err != nil {
		filename = filename + ".gz"
	}

	log.Debug("Opening " + filename)
	f, err := os.Open(filename)
	if err != nil {
		return []string{}, errors.New("file not found: " + filename)
	}
	defer f.Close()

	var reader io.Reader

	if strings.HasSuffix(filename, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return []string{}, errors.New("file gzip not readable")
		}
		reader = gz
	} else {
		reader = f
	}

	scanner := bufio.NewScanner(reader)
	if err != nil {
		return []string{}, errors.New("file not readable")
	}

	content := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		content = append(content, line)
	}

	return content, nil
}

func (l *Logger) ReadRandomMessageForUser(channelID, userID string) (string, error) {
	var userLogs []string
	var lines []string

	if channelID == "" || userID == "" {
		return "", errors.New("missing channelID or userID")
	}

	years, _ := ioutil.ReadDir(l.logPath + "/" + channelID)
	for _, yearDir := range years {
		year := yearDir.Name()
		months, _ := ioutil.ReadDir(l.logPath + "/" + channelID + "/" + year + "/")
		for _, monthDir := range months {
			month := monthDir.Name()
			path := fmt.Sprintf("%s/%s/%s/%s/%s.txt", l.logPath, channelID, year, month, userID)
			if _, err := os.Stat(path); err == nil {
				userLogs = append(userLogs, path)
			} else if _, err := os.Stat(path + ".gz"); err == nil {
				userLogs = append(userLogs, path+".gz")
			}
		}
	}

	if len(userLogs) < 1 {
		return "", errors.New("no log found")
	}

	for _, logFile := range userLogs {
		f, _ := os.Open(logFile)

		scanner := bufio.NewScanner(f)

		if strings.HasSuffix(logFile, ".gz") {
			gz, _ := gzip.NewReader(f)
			scanner = bufio.NewScanner(gz)
		}

		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}
		f.Close()
	}

	ranNum := rand.Intn(len(lines))

	return lines[ranNum], nil
}
