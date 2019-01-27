package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func download(link string) bool {
	cmd := exec.Command("megadl", "--limit-speed=500", link)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run() == nil
}

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
		return
	}

	links := strings.Split(string(data), "\n")
	log.Println(len(links))
	for i, link := range links {
		if link[0] == '#' {
			log.Printf("Skipping \"%s\"\n", link)
			continue
		}

		for !download(link) {
			log.Printf("Download of \"%s\" failed, waiting before retrying.\n", link)
			time.Sleep(time.Minute)
		}

		links[i] = "#" + link
		ioutil.WriteFile(os.Args[1], []byte(strings.Join(links, "\n")), 0664)
	}

	log.Println("Downloads done.")
}
