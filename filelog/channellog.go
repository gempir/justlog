package filelog

import (
	"github.com/gempir/gempbotgo/twitch"
	"strings"
	"os"
	"fmt"
)

func (l *Logger) LogMessageForChannel(msg twitch.Message) error {
	year := msg.Time.Year()
	month := msg.Time.Month()
	channel := strings.Replace(msg.Channel.Name, "#", "", 1)
	err := os.MkdirAll(fmt.Sprintf(l.logPath+"%s/%d/%s/", channel, year, month), 0755)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf(l.logPath+"%s/%d/%s/channel.txt", channel, year, month)

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

