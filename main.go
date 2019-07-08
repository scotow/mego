package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
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
	speedFlag    = flag.Uint("l", 0, "speed limit passed to megadl as --limit-speed")
	silentFlag   = flag.Bool("s", false, "silent mode. do not pipe megadl's stdout nor stderr")
	intervalFlag = flag.Duration("r", retryInterval, "interval between two retries")

	linkRegex = regexp.MustCompile(`^(?:https?://)?mega\.nz/#.+$`)
)

var (
	outLogger = log.New(os.Stdout, "", log.LstdFlags)
	errLogger = log.New(os.Stderr, "", log.LstdFlags)
)

func isValidLink(link string) bool {
	return linkRegex.MatchString(link)
}

func downloadRepeat(link string) {
	for !downloadCommand(link) {
		errLogger.Printf("Download of \"%s\" failed, waiting %s before retrying.\n", link, retryInterval.String())
		time.Sleep(*intervalFlag)
	}

	outLogger.Printf("Download of \"%s\" done.\n", link)
}

func downloadCommand(link string) bool {
	cmd := exec.Command("megadl", fmt.Sprintf("--limit-speed=%d", *speedFlag), link)

	var errBuff bytes.Buffer
	if *silentFlag {
		cmd.Stderr = &errBuff
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = io.MultiWriter(os.Stderr, &errBuff)
	}

	err := cmd.Run()
	if err != nil {
		logs := strings.Split(strings.TrimSpace(errBuff.String()), "\n")
		if len(logs) == 0 {
			return false
		}

		for _, line := range logs {
			if !strings.HasPrefix(line, "ERROR: File already exists at ") {
				return false
			}
		}

		return true
	}

	return false
}

func downloadFromFilesList(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		errLogger.Printf("Cannot open file \"%s\". Skipping. (%s)\n", path, err.Error())
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
			errLogger.Printf("Skipping \"%s\".\n", link[1:])
			continue
		}

		if !isValidLink(link) {
			errLogger.Printf("Invalid link %s. Skipping.\n", link)
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
		errLogger.Printf("Cannot write file \"%s\". (%s)\n", path, err.Error())
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

	outLogger.Println("All download(s) done.")
}
