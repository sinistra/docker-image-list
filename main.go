package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strings"
)

type dockerImage struct {
	Containers   string `json:"Containers"`
	CreatedAt    string `json:"CreatedAt"`
	CreatedSince string `json:"CreatedSince"`
	Digest       string `json:"Digest"`
	ID           string `json:"ID"`
	Repository   string `json:"Repository"`
	SharedSize   string `json:"SharedSize"`
	Size         string `json:"Size"`
	Tag          string `json:"Tag"`
	UniqueSize   string `json:"UniqueSize"`
	VirtualSize  string `json:"VirtualSize"`
	Architecture string `json:"Architecture"`
}

type columnLength struct {
	repository int
	tag        int
	id         int
	created    int
	size       int
	arch       int
}

var column columnLength

func main() {
	out, err := getDockerImagesAsJson()
	if err != nil {
		log.Fatalf("[docker: %s\n", err)
	}

	dockerImages, err := convertJsonListToSlice(out)
	if err != nil {
		log.Fatalf("error converting json to slice: %s", err)
	}

	addArchitectureToImages(dockerImages)

	for _, image := range *dockerImages {
		column.repository = maxInt(column.repository, len(image.Repository))
		column.tag = maxInt(column.tag, len(image.Tag))
		column.id = maxInt(column.id, len(image.ID))
		column.created = maxInt(column.created, len(image.CreatedSince))
		column.size = maxInt(column.size, len(image.Size))
		column.arch = maxInt(column.arch, len(image.Architecture))
	}

	displayImagesByName(*dockerImages)
}

func addArchitectureToImages(images *[]dockerImage) {
	for i, image := range *images {
		arch, err := getImageArchitecture(image.ID)
		if err != nil {
			log.Fatalf("getImageArchitecture failed: %s", err)
		}
		(*images)[i].Architecture = strings.TrimSpace(string(arch))
	}
}

func getImageArchitecture(id string) ([]byte, error) {
	cmd := exec.Command("docker", "image", "inspect", id, "--format", "{{.Architecture}}")
	return cmd.CombinedOutput()
}

func convertJsonListToSlice(out []byte) (*[]dockerImage, error) {
	output := strings.Split(string(out), "\n")
	var dockerImages []dockerImage
	for _, line := range output {
		var image dockerImage
		if len(line) > 0 {
			err := json.Unmarshal([]byte(line), &image)
			if err != nil {
				return nil, err
			}
			dockerImages = append(dockerImages, image)
		}
	}
	return &dockerImages, nil
}

func getDockerImagesAsJson() ([]byte, error) {
	cmd := exec.Command("docker", "image", "ls", "--format", "json")
	return cmd.CombinedOutput()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func displayImagesByName(images []dockerImage) {
	sort.Slice(images, func(i, j int) bool {
		return images[i].Repository+images[i].Tag < images[j].Repository+images[j].Tag
	})
	displayImages(images)
}

func displayImages(images []dockerImage) {
	fmt.Printf("%*s %*s %*s %*s %*s %*s\n",
		-column.repository, "REPOSITORY",
		-20, "TAG",
		-column.arch, "ARCH",
		-column.id, "IMAGE ID",
		-column.created, "CREATED",
		-column.size, "SIZE")

	for _, image := range images {
		var tag string
		if len(image.Tag) > 20 {
			tag = image.Tag[:20]
		} else {
			tag = image.Tag
		}

		fmt.Printf("%*s %*s %*s %*s %*s %*s\n",
			-column.repository, image.Repository,
			-20, tag,
			-column.arch, image.Architecture,
			-column.id, image.ID,
			-column.created, image.CreatedSince,
			-column.size, image.Size)
	}
}
