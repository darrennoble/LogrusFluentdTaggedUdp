package TaggedUDP

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"net"
	"strings"
)

const (
	defaultSeparator   = "\t"
	defaultReplacement = "    "
	defaultTag         = "applog"
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
	Host           string
	Port           int
	LogLevels      []logrus.Level
	Formatter      Formatter
	Seperator      string
	SepReplacement string
	Tag            string
	TagField       string
	soc            net.Conn
}

// New returns a new TaggedUDP hook
func New(host string, port int, tag string) (*TaggedUDP, error) {
	t := TaggedUDP{
		Host:           host,
		Port:           port,
		LogLevels:      defaultLevels,
		Formatter:      Formatter{},
		Seperator:      defaultSeparator,
		SepReplacement: defaultReplacement,
		Tag:            tag,
	}

	t.Formatter.StaticFields = logrus.Fields{}

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
	tag := h.Tag

	if h.TagField != "" {
		if newTag, ok := entry.Data[h.TagField]; ok {
			if newTag, ok := newTag.(string); ok {
				tag = newTag
			}
		}
	}

	msg, err := h.Formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("Error formatting log message: %s", err.Error())
	}

	err = h.send(tag, msg)
	if err != nil {
		return fmt.Errorf("Error sending log message: %s", err.Error())
	}

	return nil
}

func (h *TaggedUDP) send(tag string, msg []byte) error {
	b := []byte(fmt.Sprintf("%s%s", tag, h.Seperator))

	msg = []byte(strings.Replace(string(msg), h.Seperator, h.SepReplacement, -1))

	b = append(b, msg...)

	sent := 0
	for sent < len(b) {
		newSent, err := h.soc.Write(b[sent:])
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
