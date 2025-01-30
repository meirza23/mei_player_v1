package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func downloadSong(url string, title string) {
	reader := bufio.NewReader(os.Stdin)
	clearScreen()
	fmt.Print("\nŞarkıyı Playliste eklemek ister misiniz(E/H)?,Geri dönmek için '0'\n")

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
		fmt.Printf("📥 %s İndiriliyor...\n", title)
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
		DownToPlaylist(url, title, originalDir)
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

func DownToPlaylist(url string, title string, originalDir string) {
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
		ShowPlToDown(url, title, originalDir)
	case "2":
		fmt.Println("Yeni playlist adı: ")
		playlistName, _ := reader.ReadString('\n')
		playlist := strings.TrimSpace(playlistName)

		if playlist == "" {
			fmt.Println("❌ Geçersiz playlist adı!")
			return
		}

		err := os.Mkdir(playlist, 0755)
		if err != nil {
			fmt.Println("Playlist oluşturulamadı: ", err)
		}
		originDir, _ := os.Getwd()
		err = os.Chdir(playlist)
		if err != nil {
			fmt.Println("Dizine girilemedi: ", err)
			return
		}

		clearScreen()
		fmt.Printf("📥 %s playlistine %s şarkısı indiriliyor...\n", playlist, title)
		cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", url)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("İndirme hatası: %v\nÇıktı: %s\n", err, string(output))
			return
		}

		files, err := os.ReadDir(".")
		if err != nil {
			fmt.Println("Dosyalar okunamadı:", err)
			return
		}

		var mp3File string
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".mp3") {
				mp3File = file.Name()
				break
			}
		}

		if mp3File != "" {
			targetPath := filepath.Join(originalDir, "Songs", mp3File)

			// Dosya zaten var mı kontrol et
			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				// Hard link oluştur
				err = os.Link(mp3File, targetPath)
				if err != nil {
					fmt.Printf("❌ Hard link oluşturulamadı: %v\n", err)
				} else {
					fmt.Printf("✅ %s, Songs'a hard link olarak eklendi!\n", mp3File)
				}
			} else {
				fmt.Printf("ℹ️ %s zaten Songs klasöründe mevcut\n", mp3File)
			}
		}

		err = os.Chdir(originDir)
		if err != nil {
			fmt.Println("Orijinal dizine dönülemedi: ", err)
			return
		}
	default:
		fmt.Println("Geçersiz Seçim!")
	}
}
