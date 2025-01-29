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
	originalDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Dizin alınamadı:", err)
		return
	}

	err = os.Chdir("./Songs")
	if err != nil {
		fmt.Println("Dizine girilemedi: ", err)
		return
	}
	defer os.Chdir(originalDir)

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

		fmt.Println("\nÇalmak için numara girin (Ana menü için 0): ")
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
	originalDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Dizin alınamadı:", err)
		return
	}
	err = os.Chdir(originalDir)
	if err != nil {
		fmt.Println("Dizin değiştirilemedi:", err)
	}
}
