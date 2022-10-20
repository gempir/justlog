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

	"github.com/gempir/go-twitch-irc/v3"
	log "github.com/sirupsen/logrus"
)

type Logger interface {
	LogPrivateMessageForUser(user twitch.User, message twitch.PrivateMessage) error
	LogClearchatMessageForUser(userID string, message twitch.ClearChatMessage) error
	LogUserNoticeMessageForUser(userID string, message twitch.UserNoticeMessage) error
	GetLastLogYearAndMonthForUser(channelID, userID string) (int, int, error)
	GetAvailableLogsForUser(channelID, userID string) ([]UserLogFile, error)
	ReadLogForUser(channelID, userID string, year string, month string) ([]string, error)
	ReadRandomMessageForUser(channelID, userID string) (string, error)

	LogPrivateMessageForChannel(message twitch.PrivateMessage) error
	LogClearchatMessageForChannel(message twitch.ClearChatMessage) error
	LogUserNoticeMessageForChannel(message twitch.UserNoticeMessage) error
	ReadLogForChannel(channelID string, year int, month int, day int) ([]string, error)
	ReadRandomMessageForChannel(channelID string) (string, error)
	GetAvailableLogsForChannel(channelID string) ([]ChannelLogFile, error)
}

type FileLogger struct {
	logPath string
}

func NewFileLogger(logPath string) FileLogger {
	return FileLogger{
		logPath: logPath,
	}
}

func (l *FileLogger) LogPrivateMessageForUser(user twitch.User, message twitch.PrivateMessage) error {
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

func (l *FileLogger) LogClearchatMessageForUser(userID string, message twitch.ClearChatMessage) error {
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

func (l *FileLogger) LogUserNoticeMessageForUser(userID string, message twitch.UserNoticeMessage) error {
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

type UserLogFile struct {
	path  string
	Year  string `json:"year"`
	Month string `json:"month"`
}

func (l *FileLogger) GetLastLogYearAndMonthForUser(channelID, userID string) (int, int, error) {
	if channelID == "" || userID == "" {
		return 0, 0, fmt.Errorf("Invalid channelID: %s or userID: %s", channelID, userID)
	}

	logFiles := []UserLogFile{}

	years, _ := ioutil.ReadDir(l.logPath + "/" + channelID)

	for _, yearDir := range years {
		year := yearDir.Name()
		months, _ := ioutil.ReadDir(l.logPath + "/" + channelID + "/" + year + "/")
		for _, monthDir := range months {
			month := monthDir.Name()
			path := fmt.Sprintf("%s/%s/%s/%s/%s.txt", l.logPath, channelID, year, month, userID)
			if _, err := os.Stat(path); err == nil {

				logFile := UserLogFile{path, year, month}
				logFiles = append(logFiles, logFile)
			} else if _, err := os.Stat(path + ".gz"); err == nil {
				logFile := UserLogFile{path + ".gz", year, month}
				logFiles = append(logFiles, logFile)
			}
		}
	}

	sort.Slice(logFiles, func(i, j int) bool {
		yearA, _ := strconv.Atoi(logFiles[i].Year)
		yearB, _ := strconv.Atoi(logFiles[j].Year)
		monthA, _ := strconv.Atoi(logFiles[i].Month)
		monthB, _ := strconv.Atoi(logFiles[j].Month)

		if yearA == yearB {
			return monthA > monthB
		}

		return yearA > yearB
	})

	if len(logFiles) > 0 {
		year, _ := strconv.Atoi(logFiles[0].Year)
		month, _ := strconv.Atoi(logFiles[0].Month)

		return year, month, nil
	}

	return 0, 0, errors.New("No logs file")
}

func (l *FileLogger) GetAvailableLogsForUser(channelID, userID string) ([]UserLogFile, error) {
	if channelID == "" || userID == "" {
		return []UserLogFile{}, fmt.Errorf("Invalid channelID: %s or userID: %s", channelID, userID)
	}

	logFiles := []UserLogFile{}

	years, _ := ioutil.ReadDir(l.logPath + "/" + channelID)

	for _, yearDir := range years {
		year := yearDir.Name()
		months, _ := ioutil.ReadDir(l.logPath + "/" + channelID + "/" + year + "/")
		for _, monthDir := range months {
			month := monthDir.Name()
			path := fmt.Sprintf("%s/%s/%s/%s/%s.txt", l.logPath, channelID, year, month, userID)
			if _, err := os.Stat(path); err == nil {

				logFile := UserLogFile{path, year, month}
				logFiles = append(logFiles, logFile)
			} else if _, err := os.Stat(path + ".gz"); err == nil {
				logFile := UserLogFile{path + ".gz", year, month}
				logFiles = append(logFiles, logFile)
			}
		}
	}

	sort.Slice(logFiles, func(i, j int) bool {
		yearA, _ := strconv.Atoi(logFiles[i].Year)
		yearB, _ := strconv.Atoi(logFiles[j].Year)
		monthA, _ := strconv.Atoi(logFiles[i].Month)
		monthB, _ := strconv.Atoi(logFiles[j].Month)

		if yearA == yearB {
			return monthA > monthB
		}

		return yearA > yearB
	})

	if len(logFiles) > 0 {
		return logFiles, nil
	}

	return logFiles, errors.New("No logs file")
}

type ChannelLogFile struct {
	path  string
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
}

func (l *FileLogger) GetAvailableLogsForChannel(channelID string) ([]ChannelLogFile, error) {
	if channelID == "" {
		return []ChannelLogFile{}, fmt.Errorf("Invalid channelID: %s", channelID)
	}

	logFiles := []ChannelLogFile{}

	years, _ := ioutil.ReadDir(l.logPath + "/" + channelID)

	for _, yearDir := range years {
		year := yearDir.Name()
		months, _ := ioutil.ReadDir(l.logPath + "/" + channelID + "/" + year + "/")
		for _, monthDir := range months {
			month := monthDir.Name()

			days, _ := ioutil.ReadDir(l.logPath + "/" + channelID + "/" + year + "/" + month + "/")
			for _, dayDir := range days {
				day := dayDir.Name()
				path := fmt.Sprintf("%s/%s/%s/%s/%s/channel.txt", l.logPath, channelID, year, month, day)

				if _, err := os.Stat(path); err == nil {
					logFile := ChannelLogFile{path, year, month, day}
					logFiles = append(logFiles, logFile)
				} else if _, err := os.Stat(path + ".gz"); err == nil {
					logFile := ChannelLogFile{path + ".gz", year, month, day}
					logFiles = append(logFiles, logFile)
				}
			}
		}
	}

	sort.Slice(logFiles, func(i, j int) bool {
		yearA, _ := strconv.Atoi(logFiles[i].Year)
		yearB, _ := strconv.Atoi(logFiles[j].Year)
		monthA, _ := strconv.Atoi(logFiles[i].Month)
		monthB, _ := strconv.Atoi(logFiles[j].Month)
		dayA, _ := strconv.Atoi(logFiles[j].Day)
		dayB, _ := strconv.Atoi(logFiles[j].Day)

		if yearA == yearB {
			if monthA == monthB {
				return dayA > dayB
			}

			return monthA > monthB
		}

		if monthA == monthB {
			return dayA > dayB
		}

		return yearA > yearB
	})

	if len(logFiles) > 0 {
		return logFiles, nil
	}

	return logFiles, errors.New("No logs file")
}

// ReadLogForUser fetch logs
func (l *FileLogger) ReadLogForUser(channelID, userID string, year string, month string) ([]string, error) {
	if channelID == "" || userID == "" {
		return []string{}, fmt.Errorf("Invalid channelID: %s or userID: %s", channelID, userID)
	}

	filename := fmt.Sprintf(l.logPath+"/%s/%s/%s/%s.txt", channelID, year, month, userID)

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

func (l *FileLogger) ReadRandomMessageForUser(channelID, userID string) (string, error) {
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
