package filelog

import (
	"github.com/gempir/gempbotgo/twitch"
	"strings"
	"os"
	"fmt"
)

type Logger struct {
	logPath string
}

func NewFileLogger(logPath string) Logger {
	return Logger{
		logPath: logPath,
	}
}

func (l *Logger) LogMessage(msg twitch.Message) error {
	year := msg.Timestamp.Year()
	month := msg.Timestamp.Month()
	channel := strings.Replace(msg.Channel, "#", "", 1)
	err := os.MkdirAll(fmt.Sprintf(l.logPath+"%s/%d/%s/", channel, year, month), 0755)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"%s/%d/%s/%s.txt", channel, year, month, msg.User.Username)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	contents := fmt.Sprintf("[%s] %s: %s\r\n", msg.Timestamp.Format("2006-01-2 15:04:05"), msg.User.Username, msg.Text)
	if _, err = file.WriteString(contents); err != nil {
		return err
	}
	return nil
}