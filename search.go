package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func searchPython(query string) ([]Song, error) {
	cmd := exec.Command("python3", "ytmusic_search.py", query)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("python hatası: %v", string(output))
	}

	// Önce genel hata kontrolü
	var errorResp struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(output, &errorResp) == nil && errorResp.Error != "" {
		return nil, fmt.Errorf(errorResp.Error)
	}

	// Normal sonuçları parse et
	var raw []struct {
		Title    string   `json:"title"`
		Artists  []string `json:"artists"`
		Duration string   `json:"duration"`
		VideoID  string   `json:"videoId"`
	}

	if err := json.Unmarshal(output, &raw); err != nil {
		return nil, fmt.Errorf("json parse hatası: %v", err)
	}

	var songs []Song
	for _, item := range raw {
		songs = append(songs, Song{
			Title:    item.Title,
			Artists:  item.Artists,
			Duration: item.Duration,
			VideoID:  item.VideoID,
		})
	}
	return songs, nil
}

func handleSearchResults(songs []Song) {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		fmt.Println("🔍 Arama Sonuçları:")
		for i, song := range songs {
			// Sanatçıları birleştir
			artists := "Bilinmiyor"
			if len(song.Artists) > 0 {
				artists = strings.Join(song.Artists, ", ")
			}
			fmt.Printf("%d. %s - %s (%s)\n", i+1, song.Title, artists, song.Duration)
		}

		fmt.Println("\n0. Geri dön")
		fmt.Println("Oynatmak için numara girin (örn: 1)")
		fmt.Println("İndirmek için 'd' + numara girin (örn: d1)")
		fmt.Print("Seçiminiz: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch {
		case choice == "0":
			return

		case strings.HasPrefix(choice, "d"):
			numStr := strings.TrimPrefix(choice, "d")
			num, err := strconv.Atoi(numStr)
			if err != nil || num < 1 || num > len(songs) {
				fmt.Println("Geçersiz indirme seçimi!")
				time.Sleep(2 * time.Second)
				continue
			}
			song := songs[num-1]
			downloadSong(song.VideoID, song.Title)
			return

		default:
			num, err := strconv.Atoi(choice)
			if err != nil || num < 1 || num > len(songs) {
				fmt.Println("Geçersiz seçim!")
				time.Sleep(2 * time.Second)
				continue
			}
			song := songs[num-1]
			playSong(song.VideoID, song.Title)
			return
		}
	}
}
