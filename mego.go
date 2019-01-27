package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	retryInterval = time.Minute
)

var (
	speedFlag 	= flag.Uint("l", 0, "speed limit passed to megadl as --limit-speed")
	silentFlag	= flag.Bool("s", false, "silent mode. do not pipe megadl to stdout")
	linkRegex 	= regexp.MustCompile(`^(?:https?://)?mega\.nz/#.+$`)
)

func isValidLink(link string) bool {
	return linkRegex.MatchString(link)
}

func downloadRepeat(link string) {
	for !downloadCommand(link) {
		log.Printf("Download of \"%s\" failed, waiting %s before retrying.\n", link, retryInterval.String())
		time.Sleep(retryInterval)
	}

	log.Printf("Download of \"%s\" done.\n", link)
}

func downloadCommand(link string) bool {
	cmd := exec.Command("megadl", fmt.Sprintf("--limit-speed=%d", *speedFlag), link)
	if !*silentFlag {
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	}
	return cmd.Run() == nil
}

func downloadFromFilesList(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Cannot open file \"%s\". Skipping. (%s)\n", path, err.Error())
		return
	}

	lines := strings.Split(string(data), "\n")
	links := make([]string, 0, len(lines))

	// Parsing links in file.
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		links = append(links, line)
	}

	// Download each links in list.
	for i, link := range links {
		if link[0] == '#' {
			log.Printf("Skipping \"%s\"\n", link)
			continue
		}

		if !isValidLink(link) {
			log.Printf("Invalid link %s. Skipping.\n", link)
			links[i] = fmt.Sprintf("#-%s", link)
			writeFilesList(path, links)
			continue
		}

		downloadRepeat(link)
		links[i] = fmt.Sprintf("#%s", link)
		writeFilesList(path, links)
	}
}

func writeFilesList(path string, links []string) {
	err := ioutil.WriteFile(path, []byte(strings.Join(links, "\n")), 0664)
	if err != nil {
		log.Printf("Cannot write file \"%s\". Skipping. (%s)\n", path, err.Error())
	}
}

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		if isValidLink(arg) {
			downloadRepeat(arg)
		} else {
			downloadFromFilesList(arg)
		}
	}

	log.Println("All downloads done.")
}
