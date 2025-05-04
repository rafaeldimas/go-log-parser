package parser

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrInvalidLine         = errors.New("invalid line")
	ErrInvalidNumberMatchs = errors.New("invalid number matchs")
	RegexDate              = `([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})`
	RegexIp                = `((?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})`
	RegexGroup             = `(INFO|WARNING|DANGER)`
	RegexMsg               = `(.*)`
	RegexLine              = fmt.Sprintf(`^%s%s%s%s$`, RegexDate, RegexIp, RegexGroup, RegexMsg)
)

type LogInfo struct {
	CreatedAt string `json:"created_at"`
	Ip        string `json:"ip"`
	Group     string `json:"group"`
	Message   string `json:"message"`
}

type ParserInterface interface {
	Parse(line string) (LogInfo, error)
}

type Parser struct {
	valid bool
}

func NewParser() ParserInterface {
	return &Parser{
		valid: true,
	}
}

func (p *Parser) Parse(line string) (LogInfo, error) {
	regex := regexp.MustCompile(RegexLine)

	if !regex.MatchString(line) {
		p.valid = false
		return LogInfo{}, ErrInvalidLine
	}

	groups := regex.FindStringSubmatch(line)

	if len(groups) != 5 {
		return LogInfo{}, ErrInvalidNumberMatchs
	}

	return LogInfo{
		CreatedAt: groups[1],
		Ip:        groups[2],
		Group:     groups[3],
		Message:   groups[4],
	}, nil
}
