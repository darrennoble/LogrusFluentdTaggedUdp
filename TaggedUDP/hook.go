package TaggedUDP

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"net"
)

const (
	defaultTag = "\t"
)

var defaultLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
	logrus.DebugLevel,
}

// TaggedUDP is a Logrus hook for sending messages to fluentd using the tagged UDP plugin
// Logrus: https://github.com/sirupsen/logrus
// Fluentd: http://www.fluentd.org/
// Tagged UDP Plugin: https://github.com/toyokazu/fluent-plugin-tagged_udp
type TaggedUDP struct {
	Host      string
	Port      int
	Tag       string
	LogLevels []logrus.Level
	Formatter Formatter
	soc       net.Conn
}

// New returns a new TaggedUDP hook
func New(host string, port int) (*TaggedUDP, error) {
	t := TaggedUDP{
		Host:      host,
		Port:      port,
		Tag:       defaultTag,
		LogLevels: defaultLevels,
		Formatter: Formatter{},
	}

	err := t.Connect()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Connect opens a UDP connection to the fluentd tagged UDP host & port
func (h *TaggedUDP) Connect() (err error) {
	h.soc, err = net.Dial("udp", fmt.Sprintf("%s:%d", h.Host, h.Port))
	return err
}

// Close closes the connection to fluentd
func (h *TaggedUDP) Close() error {
	return h.soc.Close()
}

// Fire sends a log messages to fluentd
func (h *TaggedUDP) Fire(entry *logrus.Entry) error {
	msg, err := h.Formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("Error formatting log message: %s", err.Error())
	}

	err = h.send(msg)
	if err != nil {
		return fmt.Errorf("Error sending log message: %s", err.Error())
	}

	return nil
}

func (h *TaggedUDP) send(msg []byte) error {
	sent := 0
	for sent < len(msg) {
		newSent, err := h.soc.Write(msg[sent:])
		if err != nil {
			return err
		}

		sent += newSent
	}

	return nil
}

// Levels returns the supported levels
func (h *TaggedUDP) Levels() []logrus.Level {
	return h.LogLevels
}

// SetLevels sets the log levels used by this TaggedUDP
func (h *TaggedUDP) SetLevels(levels []logrus.Level) {
	h.LogLevels = levels
}
