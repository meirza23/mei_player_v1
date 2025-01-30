package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func downloadSong(url string, title string) {
	reader := bufio.NewReader(os.Stdin)
	clearScreen()
	fmt.Print("\nŞarkıyı Playliste eklemek ister misiniz(E/H)?,Geri dönmek için '0'")

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
		clearScreen()
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
		DownToPlaylist(url, title)
		err = os.Chdir(originalDir)
		if err != nil {
			fmt.Println("Dizin değiştirilemedi:", err)
		}
	case "0":
		return
	default:
		fmt.Println("❌ Geçersiz seçim! Lütfen sadece E , H veya 0 giriniz.")
		time.Sleep(1 * time.Second)
		return
	}

}

func DownToPlaylist(url string, title string) {
	clearScreen()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("0. Önceki Sayfaya Geri Dön")
	fmt.Println("1. Varolan Playlist'e Ekle")
	fmt.Println("2. Yeni Playlist Oluştur")
	fmt.Print("Seçiminiz: ")

	input, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(input)

	switch choice {
	case "0":
		return
	case "1":
		ShowPlToDown(url, title)
	case "2":
		fmt.Println("Yeni playlist adı: ")

	default:
		fmt.Println("Geçersiz Seçim!")
	}
}
