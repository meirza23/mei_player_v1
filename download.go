package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func downloadSong(url string, title string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Şarkıyı Playliste eklemek ister misiniz(E/H)? ")

	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "h":
		originalDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Dizin alınamadı:", err)
			return
		}

		err = os.Chdir("./Songs")
		if err != nil {
			fmt.Println("Dizine girilemedi:", err)
			return
		}

		fmt.Printf("📥 %s İndiriliyor...\n", title) // println yerine printf
		cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", url)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("İndirme hatası: %v\nÇıktı: %s\n", err, string(output))

			os.Chdir(originalDir)
			return
		}

		fmt.Println("✅ İndirme tamamlandı!")

		err = os.Chdir(originalDir)
		if err != nil {
			fmt.Println("Dizin değiştirilemedi:", err)
		}
	case "e":
		originalDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Dizin alınamadı:", err)
			return
		}

		err = os.Chdir("./Playlists")
		if err != nil {
			fmt.Println("Dizine girilemedi:", err)
			return
		}

		playlists, _ := os.ReadDir(".")
		fmt.Println("Mevcut Playlist'ler:")
		for i, file := range playlists {
			fmt.Printf("%d. %s\n", i+1, file.Name())
		}

		fmt.Println("\nYeni playlist için 0, Seçmek için numara girin:")
		input, _ := reader.ReadString('\n')
		choice, _ := strconv.Atoi(strings.TrimSpace(input))

		var targetDir string

		switch {
		case choice == 0:
			fmt.Println("Playlist adı girin:")
			filename, _ := reader.ReadString('\n')
			filename = strings.TrimSpace(filename)

			err = os.Mkdir(filename, 0755)
			if err != nil {
				fmt.Println("Oluşturulamadı:", err)
				return
			}
			targetDir = filename

		case choice > 0 && choice <= len(playlists):
			targetDir = playlists[choice-1].Name()

		default:
			fmt.Println("Geçersiz seçim!")
			return
		}

		err = os.Chdir(targetDir)
		if err != nil {
			fmt.Println("Playlist'e girilemedi:", err)
			return
		}

		fmt.Printf("📥 %s İndiriliyor...\n", title)
		cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", url)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("İndirme hatası: %v\nÇıktı: %s\n", err, string(output))
		} else {
			fmt.Println("✅ Playlist'e eklendi!")
		}

		os.Chdir(originalDir)
	default:
		fmt.Println("❌ Geçersiz seçim! Lütfen sadece E veya H giriniz.")
		time.Sleep(1 * time.Second)
		return
	}

}
