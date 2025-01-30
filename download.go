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

		fmt.Println("Yeni bir playlist oluşturmak ister misin(E/H):")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)
		switch input {
		case "e":
			fmt.Println("\nOluşturucağınız playlistin adını giriniz: ")
			input, _ := reader.ReadString('\n')
			err := os.Mkdir(input, 0755)
			if err != nil {
				fmt.Println("Playlist oluşturulamadı: ", err)
				return
			}
			ShowPlaylists()
		case "h":
		}
		err = os.Chdir(originalDir)
		if err != nil {
			fmt.Println("Dizin değiştirilemedi:", err)
		}
	default:
		fmt.Println("❌ Geçersiz seçim! Lütfen sadece E veya H giriniz.")
		time.Sleep(1 * time.Second)
		return
	}

}
