package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/fatih/color"
)

type dockerImage struct {
	repository string
	tag        string
	id         string
	created    string
	size       string
}

type columnLength struct {
	repository int
	tag        int
	id         int
	created    int
	size       int
}

var (
	images []dockerImage
	column columnLength
)

func main() {
	cmd := exec.Command("docker", "image", "ls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("[docker: %s\n", err.Error())
		os.Exit(1)
	}
	output := strings.Split(string(out), "\n")
	for i, line := range output {
		// don't want the header line
		if i == 0 {
			continue
		}
		lineSlice := strings.Split(line, "  ")
		var strippedLine []string
		for _, piece := range lineSlice {
			piece = strings.TrimSpace(piece)
			if len(piece) > 0 {
				strippedLine = append(strippedLine, piece)
			}
		}
		if len(strippedLine) == 5 {
			column.repository = maxInt(column.repository, len(strippedLine[0]))
			column.tag = maxInt(column.tag, len(strippedLine[1]))
			column.id = maxInt(column.id, len(strippedLine[2]))
			column.created = maxInt(column.created, len(strippedLine[3]))
			column.size = maxInt(column.size, len(strippedLine[4]))
			images = append(images, dockerImage{
				repository: strippedLine[0],
				tag:        strippedLine[1],
				id:         strippedLine[2],
				created:    strippedLine[3],
				size:       strippedLine[4],
			})
		}
	}

	displayImagesByName()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func displayImagesByName() {
	sort.Slice(images, func(i, j int) bool {
		return images[i].repository+images[i].tag < images[j].repository+images[j].tag
	})
	displayImages()
}

func displayImages() {
	fmt.Printf("%*s %*s %*s %*s %*s\n",
		-column.repository, "REPOSITORY",
		-column.tag, "TAG",
		-column.id, "IMAGE ID",
		-column.created, "CREATED",
		-column.size, "SIZE")
	for _, image := range images {
		fmt.Printf("%*s %*s %*s %*s %*s\n",
			-column.repository, image.repository,
			-column.tag, image.tag,
			-column.id, image.id,
			-column.created, image.created,
			-column.size, image.size)
	}
}
