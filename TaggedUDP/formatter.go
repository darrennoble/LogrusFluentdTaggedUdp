package TaggedUDP

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

// Formatter formats a message as JSON for TaggedUDP
// it is a modification of https://github.com/bshuster-repo/logrus-logstash-hook/blob/master/logstash_formatter.go
type Formatter struct {
	TimestampFormat string
	StaticFields    logrus.Fields
}

// Format formats a log messages as JSON
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	fields := make(logrus.Fields)

	for k, v := range f.StaticFields {
		fields[k] = v
	}

	for k, v := range entry.Data {
		switch v := v.(type) {
		case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
			fields[k] = v
		case []byte:
			fields[k] = string(v)
		default:
			fields[k] = fmt.Sprintf("%+v", v)
		}
	}

	timeStampFormat := f.TimestampFormat

	if timeStampFormat == "" {
		timeStampFormat = time.RFC3339
	}

	fields["timestamp"] = entry.Time.Format(timeStampFormat)

	// set message field
	v, ok := entry.Data["message"]
	if ok {
		fields["fields.message"] = v
	}
	fields["message"] = entry.Message

	// set level field
	v, ok = entry.Data["level"]
	if ok {
		fields["fields.level"] = v
	}
	fields["level"] = entry.Level.String()

	serialized, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
