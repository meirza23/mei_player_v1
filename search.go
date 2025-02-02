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

	cmd := exec.Command("python3", "search.py", songName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Hata oluştu: %v\nÇıktı: %s\n", err, string(output))
		return
	}

	var results []SearchResults
	if err := json.Unmarshal(output, &results); err != nil {
		fmt.Printf("JSON parse hatası: %v\n", err)
		return
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
				item.Duration,
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
			url := "https://www.youtube.com/watch?v=" + selected.VideoID
			downloadSong(url, selected.Title)
		} else {
			num, err := strconv.Atoi(input)
			if err != nil || num < 1 || num > len(results) {
				fmt.Println("Geçersiz numara!")
				continue
			}
			selected := results[num-1]
			url := "https://www.youtube.com/watch?v=" + selected.VideoID
			playSong(url, selected.Title)
		}
	}
}
