package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func MainSearch() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Aranacak şarkı: ")
	songName, _ := reader.ReadString('\n')
	songName = strings.TrimSpace(songName)

	cmd := exec.Command(
		"yt-dlp",
		"--dump-json",
		"--default-search", "ytmsearch",
		"ytsearch5:"+songName,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Hata oluştu: %v\nÇıktı: %s\n", err, string(output))
		return
	}

	var results []SearchResults
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var item SearchResults
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			continue
		}
		results = append(results, item)
	}
	for {
		clearScreen()
		fmt.Printf("🎵 YouTube Music Sonuçları (%d adet):\n\n", len(results))
		for i, item := range results {
			artistInfo := ""
			if len(item.Artists) > 0 {
				artistInfo = " - " + item.Artists[0].Name
			}
			fmt.Printf("%d. [%s] %s%s\n\n",
				i+1,
				formatTime(item.Duration),
				item.Title,
				artistInfo,
			)
		}
		fmt.Print("Seçiminiz (Çalmak için numara, İndirmek için 'd<numara>', Ana menü için 0):\n")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "0" {
			clearScreen()
			return
		}

		if strings.HasPrefix(input, "d") {
			numStr := strings.TrimPrefix(input, "d")
			num, err := strconv.Atoi(numStr)
			if err != nil || num < 1 || num > len(results) {
				fmt.Println("Geçersiz numara!")
				continue
			}
			selected := results[num-1]
			downloadSong(selected.URL, selected.Title)
		} else {
			num, err := strconv.Atoi(input)
			if err != nil || num < 1 || num > len(results) {
				fmt.Println("Geçersiz numara!")
				continue
			}
			selected := results[num-1]
			playSong(selected.URL, selected.Title)
		}

	}
}

func formatTime(seconds int) string {
	mins := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}
