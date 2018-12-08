package filelog

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gempir/go-twitch-irc"
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

func (l *Logger) LogMessageForUser(channel string, user twitch.User, message twitch.Message) error {
	year := message.Time.Year()
	month := int(message.Time.Month())

	err := os.MkdirAll(fmt.Sprintf(l.logPath+"/%s/%d/%d/", message.Tags["room-id"], year, month), 0750)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"/%s/%d/%d/%s.txt", message.Tags["room-id"], year, month, user.UserID)

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

func (l *Logger) ReadLogForUser(channelID string, userID string, year int, month int) ([]string, error) {
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
