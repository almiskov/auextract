package main

import (
	"fmt"
	"io/fs"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	des, _ := os.ReadDir(".")

	videos := make([]fs.DirEntry, 0)

	for _, de := range des {
		if de.IsDir() {
			continue
		}

		mt := mime.TypeByExtension(filepath.Ext(de.Name()))

		if !strings.HasPrefix(mt, "video") {
			continue
		}

		videos = append(videos, de)
	}

	fmt.Println("Выбери номер видео, чтобы извлечь аудио:")

	for i, v := range videos {
		fmt.Printf("%d\t%s\n", i+1, v.Name())
	}

	var selectedIndex int

	_, err := fmt.Scanf("%d", &selectedIndex)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	selectedIndex--

	if selectedIndex < 0 || selectedIndex > (len(videos)-1) {
		fmt.Println("Недопустимый индекс")
		return
	}

	selectedVideo := videos[selectedIndex]

	fmt.Printf("Выбрано видео %s\n", selectedVideo)
	fmt.Println("Извлекаем аудио...")

	outFilename := strings.TrimSuffix(selectedVideo.Name(), filepath.Ext(selectedVideo.Name())) + ".mp3"

	// command := fmt.Sprintf("ffmpeg -hide_banner -loglevel error -i '%s' -vn -acodec copy '%s'", selectedVideo, outFilename)

	input, err := filepath.Abs(selectedVideo.Name())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	stat, err := os.Stat(selectedVideo.Name())
	fmt.Println(stat, err)

	cmd := exec.Command(
		"ffmpeg", "-hide_banner", "-loglevel", "error",
		"-i", "'"+input+"'",
		"-vn", "-acodec", "copy",
		"'"+outFilename+"'",
		"-y",
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println(cmd.String())

	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Аудио сохранено в", outFilename)
}
