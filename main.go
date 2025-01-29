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

type SearchResults struct {
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	Artists  []struct {
		Name string `json:"name"`
	} `json:"artists"`
	URL string `json:"webpage_url"`
}

var mpvProcess *os.Process

func main() {
	clearScreen()
	for {
		showMainMenu()
		fmt.Println("\nSeçiminizi yapınız: ")
		var secim int
		_, err := fmt.Scanln(&secim)

		if err != nil {
			fmt.Println("Lütfen sayı girin!")
			var discard string
			fmt.Scanln(&discard)
			time.Sleep(1500 * time.Millisecond)
			clearScreen()
			continue
		}
		switch {
		case secim == 0:
			{
				clearScreen()
				fmt.Println("Çıkış yapılıyor... Güle Güle 👋👋")
				os.Exit(0)
			}
		case secim == 1:
			{
				MainSearch()
			}
		case secim == 2:
			{
				StopSong()
			}
		/*case secim == 3:
			{
				ShowPlaylists()
			}
		case secim == 4:
			{
				ShowSongs()
			}*/

		default:
			fmt.Println("Geçersiz Seçim")
		}

	}
}
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
func showMainMenu() {
	fmt.Println("🎵 Mei Player 🎵")
	fmt.Println("0. Çıkış")
	fmt.Println("1. Şarkı Oynat")
	fmt.Println("2. Şarkıyı Duraklat")
	fmt.Println("3. Playlistleri Görüntüle")
	fmt.Println("4. Favorileri Görüntüle")
	fmt.Println("🎵 Mei Player 🎵")
}

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
		fmt.Print("Seçiminiz (Çalmak için numara, İndirmek için 'd<numara>', Ana menü için 0): ")
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

func playSong(url string, title string) {
	clearScreen()
	fmt.Printf("🎧 Çalınıyor: %s\n", title)
	if mpvProcess != nil {
		mpvProcess.Kill()
	}

	cmd := exec.Command("mpv",
		"--no-video",
		"--ytdl-format=bestaudio",
		"--no-terminal",
		"--quiet",
		url,
	)

	if err := cmd.Start(); err != nil {
		fmt.Printf("Oynatma hatası: %v\n", err)
		return
	}

	mpvProcess = cmd.Process

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Oynatma hatası: %v\n", err)
	}

	clearScreen()
}
func StopSong() {
	if mpvProcess != nil {
		mpvProcess.Kill()
	}
}

func downloadSong(url string, title string) {

}
