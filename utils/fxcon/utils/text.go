package utils

import "strings"

const key = "^"

func ColorText(msg string) string {
	converted := msg + "^7"
	converted = strings.ReplaceAll(converted, key+"0", "\x1b[97m") // white
	converted = strings.ReplaceAll(converted, key+"1", "\x1b[91m") // red
	converted = strings.ReplaceAll(converted, key+"2", "\x1b[32m") // dark_green
	converted = strings.ReplaceAll(converted, key+"3", "\x1b[93m") // yellow
	converted = strings.ReplaceAll(converted, key+"4", "\x1b[94m") // blue
	converted = strings.ReplaceAll(converted, key+"5", "\x1b[36m") // aqua
	converted = strings.ReplaceAll(converted, key+"6", "\x1b[35m") // dark_purple
	converted = strings.ReplaceAll(converted, key+"7", "\x1b[0m")  // reset
	converted = strings.ReplaceAll(converted, key+"8", "\x1b[31m") // dark_red
	converted = strings.ReplaceAll(converted, key+"9", "\x1b[34m") // dark_blue

	return converted
}
