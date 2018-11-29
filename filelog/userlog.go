package filelog

import (
	"fmt"
	"os"

	"github.com/gempir/go-twitch-irc"
)

// Logger logger struct
type Logger struct {
	logPath string
}

// NewFileLogger create file logger
func NewFileLogger(logPath string) Logger {
	return Logger{
		logPath: logPath,
	}
}

// LogMessageForUser log in file
func (l *Logger) LogMessageForUser(channel string, user twitch.User, message twitch.Message) error {
	year := message.Time.Year()
	month := message.Time.Month()
	err := os.MkdirAll(fmt.Sprintf(l.logPath+"/%s/%d/%s/", channel, year, month), 0755)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"/%s/%d/%s/%s.txt", channel, year, month, user.Username)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	contents := fmt.Sprintf("[%s] %s: %s\r\n", message.Time.Format("2006-01-2 15:04:05"), user.Username, message.Text)
	if _, err = file.WriteString(contents); err != nil {
		return err
	}
	return nil
}
