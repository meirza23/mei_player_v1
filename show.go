package main

import (
	"bufio"
	"fmt"
	"os"
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
			err := os.Chdir(selected.Name())
			if err != nil {
				fmt.Println("Playlist açılamadı: ", err)
				time.Sleep(1 * time.Second)
				continue
			}
			ShowSongs()
		}
	}

}
