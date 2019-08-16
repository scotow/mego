package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Speed uint          `short:"s" long:"speed-limit" description:"Speed limit passed to megadl as --limit-speed" default:"0" value-name:"SPEED"`
	Pipe  bool          `short:"p" long:"pipe-outputs" description:"Pipe megadl's stdout and stderr"`
	Retry time.Duration `short:"r" long:"retry" description:"Interval between two retries" default:"15min" value-name:"INTERVAL"`
}

var (
	opts      options
	linkRegex = regexp.MustCompile(`^(?:https?://)?mega\.nz/#.+$`)
)

var (
	outLogger = log.New(os.Stdout, "", log.LstdFlags)
	errLogger = log.New(os.Stderr, "", log.LstdFlags)
)

func isValidLink(link string) bool {
	return linkRegex.MatchString(link)
}

func isAlreadyDownloadedError(line, link string) bool {
	if strings.HasPrefix(line, "ERROR: File already exists at ") {
		return true
	}
	// Typo in the original program.
	if strings.HasPrefix(line, fmt.Sprintf("ERROR: Download failed for '%s': Can't rename donwloaded temporary file ", link)) {
		return true
	}
	return false
}

func downloadRepeat(link string) {
	for !downloadCommand(link) {
		errLogger.Printf("Download of \"%s\" failed, waiting %s before retrying.\n", link, opts.Retry.String())
		time.Sleep(opts.Retry)
	}

	outLogger.Printf("Download of \"%s\" done.\n", link)
}

func downloadCommand(link string) bool {
	cmd := exec.Command("megadl", fmt.Sprintf("--limit-speed=%d", opts.Speed), link)

	var errBuff bytes.Buffer
	if opts.Pipe {
		cmd.Stdout = os.Stdout
		cmd.Stderr = io.MultiWriter(os.Stderr, &errBuff)
	} else {
		cmd.Stderr = &errBuff
	}

	err := cmd.Run()
	if err != nil {
		logLines := strings.Split(errBuff.String(), "\n")

		for _, line := range logLines {
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			if !isAlreadyDownloadedError(line, link) {
				return false
			}
		}
	}

	return true
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
		if len(line) == 0 {
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
	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = "[-s SPEED] [-p] [-r INTERVAL] LINK... LINK_PATH..."

	args, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if len(args) == 0 {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	for _, arg := range args {
		if isValidLink(arg) {
			downloadRepeat(arg)
		} else {
			downloadFromFilesList(arg)
		}
	}

	outLogger.Println("All download(s) done.")
}
