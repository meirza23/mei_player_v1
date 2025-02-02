package main

import (
	"fmt"
	"os"
	"time"
)

type SearchResults struct {
	Title    string `json:"title"`
	Duration string `json:"duration"`
	Artists  []struct {
		Name string `json:"name"`
	} `json:"artists"`
	VideoID string `json:"videoId"`
}

var mpvProcess *os.Process

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
func showMainMenu() {
	fmt.Println("🎵 Mei Player 🎵")
	fmt.Println("0. Çıkış")
	fmt.Println("1. Şarkı Ara")
	fmt.Println("2. Playlistleri Görüntüle")
	fmt.Println("3. Şarkıları Görüntüle")
	fmt.Println("🎵 Mei Player 🎵")
}

func main() {
	directories := []string{"./Playlists", "./Songs", "./Playlists/Favourites"}
	for _, dir := range directories {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Klasör Oluşturulamadı!", err)
			return
		}

	}
	for {
		clearScreen()
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
			MainSearch()

		case secim == 2:
			{
				originalDir, err := os.Getwd()
				if err != nil {
					fmt.Println("Dizin alınamadı:", err)
					return
				}
				err = os.Chdir("./Playlists")
				if err != nil {
					fmt.Println("Dizine girilemedi: ", err)
					return
				}
				ShowPlaylists()
				err = os.Chdir(originalDir)
				if err != nil {
					fmt.Println("Dizin değiştirilemedi:", err)
				}
			}
		case secim == 3:
			{
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
				ShowSongs()
				err = os.Chdir(originalDir)
				if err != nil {
					fmt.Println("Orijinal dizine dönülemedi: ", err)
					return
				}
			}

		default:
			fmt.Println("Geçersiz Seçim")
		}

	}
}
