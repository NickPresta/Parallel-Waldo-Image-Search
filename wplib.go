package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type image struct {
	height   int
	width    int
	rotation int
	fileName string
	data     []string
}

func ReadFileContents(file *os.File) (img *image) {
	reader, err := bufio.NewReaderSize(file, 6*1024)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	img = new(image)

	img.fileName = filepath.Base(file.Name())

	line, isPrefix, err := reader.ReadLine()
	// Get first line, for the dimensions
	dimensions := strings.Split(string(line), " ")
	img.height, _ = strconv.Atoi(dimensions[0])
	img.width, _ = strconv.Atoi(dimensions[1])

	line, isPrefix, err = reader.ReadLine()
	for err == nil && !isPrefix {
		s := string(line)
		img.data = append(img.data, s)
		line, isPrefix, err = reader.ReadLine()
	}
	if isPrefix {
		fmt.Println("Buffer was declared to be too small for file (", file.Name(), ")")
		return nil
	}
	if err != os.EOF {
		fmt.Println(err)
		return nil
	}

	return
}

func OpenFile(filePath string, callback func(file *os.File) *image) *image {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	return callback(file)
}

func TraverseDirectory(directory string, callback func(file *os.File) *image) (images []*image) {
	dirContents, err := os.Open(directory)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer dirContents.Close()

	file, err := dirContents.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, file := range file {
		if file.IsRegular() {
			path, _ := filepath.Abs(filepath.Join(directory, file.Name))
			image := OpenFile(path, callback)
			if image != nil {
				images = append(images, image)
			}
		}
	}

	return
}
