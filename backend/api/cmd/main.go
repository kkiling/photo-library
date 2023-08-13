package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

type regexToDate struct {
	pattern     *regexp.Regexp
	transformer func([]string) string
}

var rules = []regexToDate{
	{
		regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2}) (\d{2})-(\d{2})-(\d{2})`),
		func(matches []string) string {
			return matches[1] + "-" + matches[2] + "-" + matches[3] + " " + matches[4] + ":" + matches[5] + ":" + matches[6]
		},
	},
	{
		regexp.MustCompile(`^IMG-(\d{4})(\d{2})(\d{2})-WA\d+`),
		func(matches []string) string { return matches[1] + "-" + matches[2] + "-" + matches[3] },
	},
	{
		regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2}) (\d{2})-(\d{2})-(\d{2})_\d+`),
		func(matches []string) string {
			return matches[1] + "-" + matches[2] + "-" + matches[3] + " " + matches[4] + ":" + matches[5] + ":" + matches[6]
		},
	},
	{
		regexp.MustCompile(`^(\d{4})(\d{2})(\d{2})_(\d{2})(\d{2})(\d{2})`),
		func(matches []string) string {
			return matches[1] + "-" + matches[2] + "-" + matches[3] + " " + matches[4] + ":" + matches[5] + ":" + matches[6]
		},
	},
	{
		regexp.MustCompile(`^IMG_(\d{4})(\d{2})(\d{2})_(\d{2})(\d{2})(\d{2})`),
		func(matches []string) string {
			return matches[1] + "-" + matches[2] + "-" + matches[3] + " " + matches[4] + ":" + matches[5] + ":" + matches[6]
		},
	},
	{
		regexp.MustCompile(`^IMG_(\d{4})(\d{2})(\d{2})`),
		func(matches []string) string { return matches[1] + "-" + matches[2] + "-" + matches[3] },
	},
}

func fileNameToTime(path string) (time.Time, error) {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[0 : len(filename)-len(ext)]

	for _, rule := range rules {
		if matches := rule.pattern.FindStringSubmatch(nameWithoutExt); matches != nil {
			formattedDate := rule.transformer(matches)
			if len(formattedDate) == 10 { // "2006-01-02"
				return time.Parse("2006-01-02", formattedDate)
			}
			return time.Parse("2006-01-02 15:04:05", formattedDate)
		}
	}

	// Handle timestamp
	if ts, err := strconv.ParseInt(nameWithoutExt, 10, 64); err == nil {
		return time.Unix(ts/1000, 0), nil // convert from milliseconds
	}

	return time.Time{}, fmt.Errorf("could not parse date from filename")
}

func main() {
	files := []string{
		"/path/to/2021-03-08 04-00-22.jpg",
		"/path/to/IMG-20150127-WA0026.jpg",
		"/path/to/2020-12-25 06-39-08_1608884124807.jpg",
		"/path/to/1488576529616.jpg",
		"/path/to/20150501_133038.jpg",
		"/path/to/IMG_20150429_214541.jpg",
	}

	for _, file := range files {
		date, err := fileNameToTime(file)
		if err != nil {
			fmt.Printf("Error for %s: %s\n", file, err)
		} else {
			fmt.Printf("%s -> %s\n", file, date)
		}
	}
}
