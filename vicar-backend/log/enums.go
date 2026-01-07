package log

import "github.com/fatih/color"

type Level int

const (
	DEBUG Level = iota
	INFO
	SUCCESS
	WARNING
	ERROR
)

func (level Level) ColoredString(emote string) string {
	switch level {
	case DEBUG:
		return color.New(color.FgCyan, color.Bold).Sprintf(" [%s DEBUG  ] ", emote)
	case INFO:
		return color.New(color.FgBlue, color.Bold).Sprintf(" [%s INFO   ] ", emote)
	case SUCCESS:
		return color.New(color.FgGreen, color.Bold).Sprintf(" [%s SUCCESS] ", emote)
	case WARNING:
		return color.New(color.FgYellow, color.Bold).Sprintf(" [%s WARNING] ", emote)
	case ERROR:
		return color.New(color.FgRed, color.Bold).Sprintf(" [%s ERROR  ] ", emote)
	default:
		return color.New(color.FgWhite, color.Bold).Sprintf(" [%s LOGGING] ", emote)
	}
}

type LogArea struct {
	Name string
	Bg   color.Attribute
	Fg   color.Attribute
}

func newLogArea(name string, bg, fg color.Attribute) LogArea {
	return LogArea{
		Name: name,
		Bg:   bg,
		Fg:   fg,
	}
}

func (area LogArea) ColoredString() string {
	return color.New(area.Bg, area.Fg).Sprintf(" %s ", area.Name)
}

var (
	Auth     = newLogArea("AUTH", color.BgRed, color.FgWhite)
	Database = newLogArea("DABA", color.BgMagenta, color.FgWhite)
	Server   = newLogArea("SERV", color.BgYellow, color.FgWhite)
	Cache    = newLogArea("CACH", color.BgGreen, color.FgWhite)
	FS       = newLogArea("FILE", color.BgCyan, color.FgWhite)
	Sync     = newLogArea("SYNC", color.BgBlue, color.FgWhite)
)
