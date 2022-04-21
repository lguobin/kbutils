package kbutils

import (
	"bytes"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// killChannel
func KillSignal() bool {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	return true
}

//time
//time
//time
//东八区
var (
	TimeTp          = "2006-01-02 15:04:05"
	GetLocationName = func(zone int) string {
		switch zone {
		case 8:
			return "Asia/Shanghai"
		}
		return "UTC"
	}
	formatKeyTpl = map[byte]string{
		'd': "02",
		'D': "Mon",
		'w': "Monday",
		'N': "Monday",
		'S': "02",
		'l': "Monday",
		'F': "January",
		'm': "01",
		'M': "Jan",
		'n': "1",
		'Y': "2006",
		'y': "06",
		'a': "pm",
		'A': "PM",
		'g': "3",
		'h': "03",
		'H': "15",
		'i': "04",
		's': "05",
		'O': "-0700",
		'P': "-07:00",
		'T': "MST",
		'c': "2006-01-02T15:04:05-07:00",
		'r': "Mon, 02 Jan 06 15:04 MST",
	}
)

// FormatTime format time
func FormatTime(t time.Time, format ...string) string {
	return New().formatTime(t, format...)
}

type TimeEngine struct{ zone *time.Location }

// New new timeEngine
func New(zone ...int) *TimeEngine {
	e := &TimeEngine{}
	timezone := Zone(zone...)
	if timezone != time.Local {
		e.zone = timezone
	}
	return e
}

// Zone eastEightTimeZone
func Zone(zone ...int) *time.Location {
	if len(zone) > 0 {
		return time.FixedZone(GetLocationName(zone[0]), zone[0]*3600)
	}
	return time.Local
}

func (e *TimeEngine) in(t time.Time) time.Time {
	if e.zone == nil {
		return t
	}
	return t.In(e.zone)
}

// formatTime string format of return time
func (e *TimeEngine) formatTime(t time.Time, format ...string) string {
	t = e.in(t)
	tpl := TimeTp
	if len(format) > 0 {
		tpl = FormatTlp(format[0])
	}
	return t.Format(tpl)
}

// FormatTlp FormatTlp
func FormatTlp(format string) string {
	runes := []rune(format)
	buffer := bytes.NewBuffer(nil)
	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case '\\':
			if i < len(runes)-1 {
				buffer.WriteRune(runes[i+1])
				i += 1
				continue
			} else {
				return buffer.String()
			}
		default:
			if runes[i] > 255 {
				buffer.WriteRune(runes[i])
				break
			}
			if f, ok := formatKeyTpl[byte(runes[i])]; ok {
				buffer.WriteString(f)
			} else {
				buffer.WriteRune(runes[i])
			}
		}
	}
	return buffer.String()
}
