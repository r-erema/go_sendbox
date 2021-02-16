package ip

import (
	"strconv"
	"strings"
)

func wrongIPPart(ip string) (result []string) {
	splitIP := strings.Split(ip, ".")
	if len(splitIP) < 4 {
		result = append(result, ip)
		return
	}
	for _, ipPart := range splitIP {
		ipPartInt, err := strconv.Atoi(ipPart)
		if err != nil {
			result = append(result, ipPart)
		}
		if ipPartInt < 0 || ipPartInt > 255 {
			result = append(result, ipPart)
		}
	}
	return
}
