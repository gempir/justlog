package filelog

import (
	"fmt"
	"github.com/gempir/go-twitch-irc"
	"os"
)

func (l *Logger) LogMessageForChannel(channel string, user twitch.User, message twitch.Message) error {
	year := message.Time.Year()
	month := message.Time.Month()
	day := message.Time.Day()
	err := os.MkdirAll(fmt.Sprintf(l.logPath+"%s/%d/%s/%d", channel, year, month, day), 0755)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"%s/%d/%s/%d/channel.txt", channel, year, month, day)

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
