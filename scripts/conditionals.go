package scripts

import (
	"regexp"
	"sync"
)

var (
	ipRegex         *regexp.Regexp
	ipPortRegex     *regexp.Regexp
	domainRegex     *regexp.Regexp
	domainPortRegex *regexp.Regexp
	regexOnce       sync.Once
)

func initRegex() {
	regexOnce.Do(func() {

		ipRegex = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)

		ipPortRegex = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?):(6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5][0-9]{4}|[1-9][0-9]{0,3})$`)

		domainRegex = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

		domainPortRegex = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}:(6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5][0-9]{4}|[1-9][0-9]{0,3})$`)
	})
}

func IsIPAddress(str string) bool {
	if str == "" {
		return false
	}
	initRegex()
	return ipRegex.MatchString(str)
}

func IsIPAddressWithPort(str string) bool {
	if str == "" {
		return false
	}
	initRegex()
	return ipPortRegex.MatchString(str)
}

func IsDomainName(str string) bool {
	if str == "" {
		return false
	}
	initRegex()
	return domainRegex.MatchString(str)
}

func IsDomainNameWithPort(str string) bool {
	if str == "" {
		return false
	}
	initRegex()
	return domainPortRegex.MatchString(str)
}
