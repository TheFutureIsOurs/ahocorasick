/*
 * @Author: Daiming Liu (xingrufeng)
 * @Copyright (C) Daiming Liu (xingrufeng)
 */
package main

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/TheFutureIsOurs/ahocorasick"
)

var (
	memprofile = "memprofile"
)

func readFile(input string) []rune {
	file, err := os.Open(input)
	if err != nil {
		return nil
	}
	defer file.Close()
	content := make([]rune, 0)
	bufReader := bufio.NewReader(file)
	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			break
		}
		keyword := strings.TrimSpace(string(line))
		if keyword == "" {
			continue
		}
		content = append(content, []rune(keyword)...)
	}
	return content

}

func main() {
	f, err := os.Create(memprofile)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close()

	ac, _ := ahocorasick.BuildFromFile("../dictionary.txt")

	content := readFile(".../text.txt")
	runtime.GC() // get up-to-date statistics
	for i := 0; i < 100; i++ {
		ac.MultiPatternIndexes(content)
	}
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

}
