package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

// Utility shared by several commands


func shorten(x string, maxLen int) string {
	if len(x) > maxLen {
		return x[0:(maxLen-4)] + ".." + x[len(x)-2:]
	} else {
		return x
	}
}


func percentToString(x int) string {
	if x < 0 {
		return "???"
	} else {
		return strconv.Itoa(x) + "%"
	}
}

var factorByUnit = map[string]uint64{
	"b": 1,
	"bi": 1,
	"k": 1000,
	"ki": 1024,
	"m": 1000*1000,
	"mi": 1024*1024,
	"g": 1000*1000*1000,
	"gi": 1024*1024*1024,
	"t": 1000*1000*1000*1000,
	"ti": 1024*1024*1024*1024,
	"p": 1000*1000*1000*1000*1000,
	"pi": 1024*1024*1024*1024*1024,
}

func strBytes2human(x string, unit string) string {
	y, err := strconv.ParseInt(x, 10, 64)
	if err == nil {
		return bytes2human(y, unit)
	} else {
		log.Errorf("while parsing '%s' -> %v", x, err)
		return "???"
	}
}

func bytes2human(x int64, unit string) string {
	if x < 0 {
		return "???"
	} else if x == 0 {
		return "0"
	} else {
		u2 := strings.ToLower(unit)
		if u2 == "a" || u2 == "h" {
			return autoBytes2human(uint64(x))
		} else {
			f, ok := factorByUnit[u2]
			if !ok {
				log.Errorf("Invalid unit '%'", unit)
				return "???"
			}
			return fmt.Sprintf("%d%s", uint64(x) / f, unit)
		}
	}
}


var prefixes = []string{ "Bi", "Ki", "Mi", "Gi", "Ti", "Pi" }

// This will lookup the greatest unit that trigger a lost < 5%
func autoBytes2human(x uint64) string {
	for i := 5; i >= 0; i-- {
		y := x >> (i * 10)
		y2 := y << (i * 10)
		lost := ((x - y2) * 100) / x
		//fmt.Printf("%d%s   y2:%d  (lost:%d%%)\n", y, prefixes[i], y2, lost)
		if lost < 5 {
			return fmt.Sprintf("%d%s", y, prefixes[i])
		}
	}
	return "????1??"
}


