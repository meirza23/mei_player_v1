package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
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
		fmt.Print("Seçiminiz (Çalmak için numara, İndirmek için 'd<numara>', Ana menü için 0):\n")
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
	fmt.Println("Durdurmak için 's', Devam için 'c', Bitir için 'q'")

	if mpvProcess != nil {
		mpvProcess.Kill()
		mpvProcess = nil
	}

	cmd := exec.Command("mpv",
		"--no-video",
		"--ytdl-format=bestaudio",
		"--input-ipc-server=/tmp/mpv-socket",
		"--quiet",
		url,
	)

	// Socket dosyasını temizle
	os.Remove("/tmp/mpv-socket")

	if err := cmd.Start(); err != nil {
		fmt.Printf("Oynatma hatası: %v\n", err)
		return
	}
	mpvProcess = cmd.Process

	// Socket'in hazır olmasını bekle
	time.Sleep(2 * time.Second)

	inputCh := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			input, _ := reader.ReadString('\n')
			inputCh <- strings.TrimSpace(input)
		}
	}()

	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	for {
		select {
		case err := <-done:
			mpvProcess = nil
			if err != nil {
				fmt.Printf("Hata: %v\n", err)
			}
			return
		case input := <-inputCh:
			switch input {
			case "s":
				clearScreen()
				fmt.Printf("🎧 Çalınıyor: %s\n", title)
				fmt.Println("Durdurmak için 's', Devam için 'c', Bitir için 'q'")
				sendMPVCommand([]interface{}{"set_property", "pause", true})
				fmt.Println("⏸️ Duraklatıldı")
			case "c":
				clearScreen()
				fmt.Printf("🎧 Çalınıyor: %s\n", title)
				fmt.Println("Durdurmak için 's', Devam için 'c', Bitir için 'q'")
				sendMPVCommand([]interface{}{"set_property", "pause", false})
				fmt.Println("▶️ Devam ediliyor")
			case "q":
				sendMPVCommand([]interface{}{"stop"})
				fmt.Println("⏹️ Durduruluyor...")
				return
			default:
				fmt.Println("Geçersiz komut!")
			}
		}
	}
}

func sendMPVCommand(args []interface{}) {
	conn, err := net.Dial("unix", "/tmp/mpv-socket")
	if err != nil {
		fmt.Printf("Socket bağlantı hatası: %v\n", err)
		return
	}
	defer conn.Close()

	cmd := map[string]interface{}{
		"command": args,
	}
	jsonCmd, _ := json.Marshal(cmd)
	_, err = conn.Write(append(jsonCmd, '\n'))
	if err != nil {
		fmt.Printf("Komut gönderme hatası: %v\n", err)
	}
}

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

		fmt.Println("Yeni bir playlist oluşturmak ister misin(E/H):\n")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)
		switch input {
		case "e":
			fmt.Println("Oluşturucağınız playlistin adını giriniz: ")
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
		fmt.Println("Tüm Şarkılar:\n")
		for i, file := range songFiles {
			fmt.Printf("%d. %s\n", i+1, file)
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

func playLocalSong(filename string) {
	clearScreen()
	fmt.Printf("🎧 Çalınıyor: %s\n", filename)
	fmt.Println("Durdurmak için 's', Devam için 'c', Bitir için 'q'")

	if mpvProcess != nil {
		mpvProcess.Kill()
		mpvProcess = nil
	}

	// Socket dosyasını temizle
	os.Remove("/tmp/mpv-socket")

	cmd := exec.Command("mpv",
		"--no-video",
		"--quiet",
		"--no-terminal",
		"--input-ipc-server=/tmp/mpv-socket",
		filename,
	)

	if err := cmd.Start(); err != nil {
		fmt.Printf("Oynatma hatası: %v\n", err)
		return
	}
	mpvProcess = cmd.Process

	// Socket'in hazır olmasını bekle
	time.Sleep(1 * time.Second)

	inputCh := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			input, _ := reader.ReadString('\n')
			inputCh <- strings.TrimSpace(input)
		}
	}()

	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	for {
		select {
		case err := <-done:
			mpvProcess = nil
			if err != nil {
				fmt.Printf("Hata: %v\n", err)
			}
			return
		case input := <-inputCh:
			switch input {
			case "s":
				clearScreen()
				fmt.Printf("🎧 Çalınıyor: %s\n", filename)
				fmt.Println("Durdurmak için 's', Devam için 'c', Bitir için 'q'")
				sendMPVCommand([]interface{}{"set_property", "pause", true})
				fmt.Println("⏸️ Duraklatıldı")
			case "c":
				clearScreen()
				fmt.Printf("🎧 Çalınıyor: %s\n", filename)
				fmt.Println("Durdurmak için 's', Devam için 'c', Bitir için 'q'")
				sendMPVCommand([]interface{}{"set_property", "pause", false})
				fmt.Println("▶️ Devam ediliyor")
			case "q":
				sendMPVCommand([]interface{}{"stop"})
				fmt.Println("⏹️ Durduruluyor...")
				return
			default:
				fmt.Println("Geçersiz komut!")
			}
		}
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
			{
				MainSearch()
			}
		/*case secim == 2:
		{
			ShowPlaylists()
		}*/
		case secim == 3:
			{
				ShowSongs()
			}

		default:
			fmt.Println("Geçersiz Seçim")
		}

	}
}
