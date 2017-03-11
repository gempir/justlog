package filelog

import (
	"fmt"
	"github.com/gempir/gempbotgo/twitch"
	"os"
	"strings"
)

type Logger struct {
	logPath string
}

func NewFileLogger(logPath string) Logger {
	return Logger{
		logPath: logPath,
	}
}

func (l *Logger) LogMessageForUser(msg twitch.Message) error {
	year := msg.Time.Year()
	month := msg.Time.Month()
	channel := strings.Replace(msg.Channel.Name, "#", "", 1)
	err := os.MkdirAll(fmt.Sprintf(l.logPath+"%s/%d/%s/", channel, year, month), 0755)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"%s/%d/%s/%s.txt", channel, year, month, msg.Username)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	contents := fmt.Sprintf("[%s] %s: %s\r\n", msg.Time.Format("2006-01-2 15:04:05"), msg.Username, msg.Text)
	if _, err = file.WriteString(contents); err != nil {
		return err
	}
	return nil
}
