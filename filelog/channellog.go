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
)

func (l *Logger) LogPrivateMessageForChannel(channel string, user twitch.User, message twitch.PrivateMessage) error {
	year := message.Time.Year()
	month := int(message.Time.Month())
	day := message.Time.Day()
	err := os.MkdirAll(fmt.Sprintf(l.logPath+"/%s/%d/%d/%d", message.Tags["room-id"], year, month, day), 0750)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"/%s/%d/%d/%d/channel.txt", message.Tags["room-id"], year, month, day)

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

func (l *Logger) LogClearchatMessageForChannel(channel string, message twitch.ClearChatMessage) error {
	year := message.Time.Year()
	month := int(message.Time.Month())
	day := message.Time.Day()
	err := os.MkdirAll(fmt.Sprintf(l.logPath+"/%s/%d/%d/%d", message.Tags["room-id"], year, month, day), 0750)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"/%s/%d/%d/%d/channel.txt", message.Tags["room-id"], year, month, day)

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
