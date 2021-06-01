package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var guidRegexp = regexp.MustCompile(`{"value":"(.{22}==)"}`)

func dump(output interface{}, currentCall uint64) {
	data, err := json.Marshal(output)
	if err != nil {
		logrus.WithError(err).Warn("Failed to convert dump output to JSON")
		return
	}

	replaced := replaceAllGuids(string(data))

	fmt.Println(replaced)
}

func replaceAllGuids(data string) string {
	replacements := make([][2]string, 0)
	for _, match := range guidRegexp.FindAllStringSubmatch(data, -1) {
		bytes, err := base64.StdEncoding.DecodeString(match[1])
		if err != nil {
			continue
		}
		swapped := make([]byte, 16)
		swapped[0] = bytes[3]
		swapped[1] = bytes[2]
		swapped[2] = bytes[1]
		swapped[3] = bytes[0]
		swapped[4] = bytes[5]
		swapped[5] = bytes[4]
		swapped[6] = bytes[7]
		swapped[7] = bytes[6]
		swapped[8] = bytes[8]
		swapped[9] = bytes[9]
		swapped[10] = bytes[10]
		swapped[11] = bytes[11]
		swapped[12] = bytes[12]
		swapped[13] = bytes[13]
		swapped[14] = bytes[14]
		swapped[15] = bytes[15]
		guid, err := uuid.FromBytes(swapped)
		if err != nil {
			continue
		}
		replacements = append(replacements, [2]string{match[0], `"` + guid.String() + `"`})
	}
	replaced := data
	for _, replacement := range replacements {
		replaced = strings.Replace(replaced, replacement[0], replacement[1], -1)
	}
	return replaced
}
