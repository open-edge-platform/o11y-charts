package main

import (
	"fmt"
	"regexp"
)

// Regular expressions to extract (timestamp, level, caller, and message) from the log
var (
	regexForTime = regexp.MustCompile(`ts=(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z)`)
	regexForLevel = regexp.MustCompile(`level=(\w+)`)
	regexForCaller = regexp.MustCompile(`caller=([\w\.]+:\d+)`)
	regexForMsg = regexp.MustCompile(`msg="([^"]+)"`)
)

func parseTime(log string) string {
	matches := regexForTime.FindStringSubmatch(log)
	if len(matches) > 1 {
		return matches[1] 
	}
	return "" 
}

func parseLevel(log string) string {
	matches := regexForLevel.FindStringSubmatch(log)
	if len(matches) > 1 {
		return matches[1] 
	}
	return "" 
}


func parseCaller(log string) string {
	matches := regexForCaller.FindStringSubmatch(log)
	if len(matches) > 1 {
		return matches[1] 
	}
	return "" 
}


func parseMsg(log string) string {
	matches := regexForMsg.FindStringSubmatch(log)
	if len(matches) > 1 {
		return matches[1] 
	}
	return "" 
}
