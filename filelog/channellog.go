package filelog

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
	log "github.com/sirupsen/logrus"
)

func (l *Logger) LogPrivateMessageForChannel(message twitch.PrivateMessage) error {
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

func (l *Logger) LogClearchatMessageForChannel(message twitch.ClearChatMessage) error {
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

func (l *Logger) LogUserNoticeMessageForChannel(message twitch.UserNoticeMessage) error {
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

func (l *Logger) ReadLogForChannel(channelID string, year int, month int, day int) ([]string, error) {
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

func (l *Logger) ReadRandomMessageForChannel(channelID) (string, error) {
	var days []string
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
			path := fmt.Sprintf("%s/%s/%s/%s/%s.txt", l.logPath, channelID, year, month, userID)
			if _, err := os.Stat(path); err == nil {
				days = append(days, path)
			} else if _, err := os.Stat(path + ".gz"); err == nil {
				days = append(days, path+".gz")
			}
		}
	}

	if len(days) < 1 {
		return "", errors.New("no log found")
	}

	randomDayIndex := rand.Intn(len(days))
	randomDayPath := days[randomDayIndex]

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

	randomLineNumber := rand.Intn(len(lines))
	return lines[randomLineNumber], nil
}
