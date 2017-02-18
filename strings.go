package vokiri

import (
	"log"
	"strings"
)

var _ = log.Print

func hasPrefixAnyAry(s string, prefixes []string) (found bool, length int) {
	if len(prefixes) > 0 {
		for _, p := range prefixes {
			if strings.HasPrefix(s, p) {
				return true, len(p)
			}
		}
	}
	return false, 0
}

func hasPrefixAny(s, prefixes string) (found bool, length int) {
	return hasPrefixAnyAry(s, strings.Split(prefixes, ""))
}

func index(s, sep string) (pos, length int) {
	pos = strings.Index(s, sep)
	if pos != -1 {
		return pos, len(sep)
	}
	return -1, 0
}

func indexAnyAry(s string, seps []string) (pos, length int) {
	poses := make(map[string]int)

	if len(seps) > 0 {
		for _, sep := range seps {
			if pos, length = index(s, sep); pos != -1 {
				poses[sep] = pos
			}
		}

		//log.Printf("poses=%#v\n", poses)
		var min_sep string
		min_pos := -1
		for sep, pos := range poses {
			//log.Printf("sep, pos=%v, %v\n", sep, pos)
			if min_pos == -1 || pos < min_pos {
				min_pos = pos
				min_sep = sep
			}
		}
		return min_pos, len(min_sep)
	}
	return -1, 0
}

func indexAny(s, chars string) (pos, length int) {
	return indexAnyAry(s, strings.Split(chars, ""))
}

func trimSuffixAnyAry(s string, suffixes []string) string {
	if len(suffixes) > 0 {
		for _, suffix := range suffixes {
			if strings.HasPrefix(s, suffix) {
				return strings.TrimSuffix(s, suffix)
			}
		}
	}
	return s
}

func trimSuffixAny(s, suffixes string) string {
	return trimSuffixAnyAry(s, strings.Split(suffixes, ""))
}
