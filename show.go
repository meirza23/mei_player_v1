package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ShowSongs() {

	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Dizin okunamadı: ", err)
		return
	}

	songFiles := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mp3") {
			songFiles = append(songFiles, file.Name())
		}
	}

	if len(songFiles) == 0 {
		fmt.Println("Kütüphanede şarkı bulunamadı!")
		time.Sleep(2 * time.Second)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		fmt.Println("Tüm Şarkılar:")
		for i, file := range songFiles {
			fmt.Printf("\n%d. %s", i+1, file)
		}

		fmt.Println("\nÇalmak için numara girin (Geri dönmek için 0): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "0" {
			return
		}

		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > len(songFiles) {
			fmt.Printf("Geçersiz seçim! (1-%d arası değer girin)\n", len(songFiles))
			time.Sleep(1 * time.Second)
			continue
		}

		selectedSong := songFiles[num-1]
		playLocalSong(selectedSong)
	}
}

func ShowPlaylists() {
	playlists, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Dizin okunamadı: ", err)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		clearScreen()
		fmt.Println("Mevcut Playlist'ler:")
		for i, playlistFiles := range playlists {
			fmt.Printf("\n%d. %s", i+1, playlistFiles.Name())
		}
		fmt.Println("\nSeçim yapın(Geri dönmek için 0):")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "0" {
			return
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(playlists) {
			fmt.Printf("Geçersiz seçim! (1-%d arası girin)\n", len(playlists))
			time.Sleep(1 * time.Second)
			continue
		}

		selected := playlists[choice-1]
		if selected.IsDir() {

			currentDir, _ := os.Getwd()
			targetPath := filepath.Join(currentDir, selected.Name())

			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				fmt.Printf("❌ Dizin bulunamadı: %s\n", targetPath)
				time.Sleep(2 * time.Second)
				continue
			}

			if err := os.Chdir(targetPath); err != nil {
				fmt.Printf("❌ Dizin açılamadı [%s]: %v\n", targetPath, err)
				time.Sleep(2 * time.Second)
				continue
			}

			ShowSongs()

			os.Chdir("..")
		}
	}
}

func ShowPlToDown(url string, title string) {
	playlists, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Dizin okunamadı: ", err)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		clearScreen()
		fmt.Println("Mevcut Playlist'ler:")
		for i, playlistFiles := range playlists {
			fmt.Printf("\n%d. %s", i+1, playlistFiles.Name())
		}
		fmt.Println("\nSeçim yapın(Geri dönmek için 0):")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "0" {
			return
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(playlists) {
			fmt.Printf("Geçersiz seçim! (1-%d arası girin)\n", len(playlists))
			time.Sleep(1 * time.Second)
			continue
		}

		selected := playlists[choice-1]
		if selected.IsDir() {

			currentDir, _ := os.Getwd()
			targetPath := filepath.Join(currentDir, selected.Name())

			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				fmt.Printf("❌ Dizin bulunamadı: %s\n", targetPath)
				time.Sleep(2 * time.Second)
				continue
			}

			if err := os.Chdir(targetPath); err != nil {
				fmt.Printf("❌ Dizin açılamadı [%s]: %v\n", targetPath, err)
				time.Sleep(2 * time.Second)
				continue
			}

			clearScreen()
			fmt.Printf("📥 %s İndiriliyor...\n", title)
			cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", url)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("İndirme hatası: %v\nÇıktı: %s\n", err, string(output))
				return
			}

			fmt.Println("✅ İndirme tamamlandı!")

			os.Chdir("..")
		}
	}

}
