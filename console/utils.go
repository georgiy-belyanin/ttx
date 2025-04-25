package console

import (
	"regexp"
)

type ConnectInfo struct {
	UriOrInstanceName string
	User              string
	Password          string
}

func ParseConnectInfo(s string) ConnectInfo {
	connectInfo := ConnectInfo{
		UriOrInstanceName: s,
		User:              "guest",
		Password:          "",
	}

	connectInfoRe, err := regexp.Compile(`(?:([^:]+)(?::(.+))?@)?(.+)`)
	if err != nil {
		return connectInfo
	}

	matches := connectInfoRe.FindStringSubmatch(s)

	user := matches[1]
	if len(user) > 0 {
		connectInfo.User = user
	}

	connectInfo.Password = matches[2]
	connectInfo.UriOrInstanceName = matches[3]

	return connectInfo
}
