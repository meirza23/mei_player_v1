package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

func playSong(url string, title string) {
	songUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", url)
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
		songUrl,
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
