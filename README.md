# 🎵 Mei Player 🎵

YouTube Music'ten şarkı arama, indirme ve yerel müzik çalma özelliklerine sahip terminal tabanlı bir müzik oynatıcı.

## ✨ Özellikler
- 🔍 YouTube Music'te şarkı arama
- ⬇️ MP3 formatında şarkı indirme
- 📁 Playlist yönetimi (oluşturma/şarkı ekleme)
- ▶️ Yerel ve çevrimiçi şarkı çalma
- ⏯️ Çalma kontrolü (duraklat/devam et/durdur)

## 🛠️ Kurulum

### Gereksinimler
- [Go](https://golang.org/dl/) (v1.16+)
- [Python 3](https://www.python.org/downloads/)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp#installation)
- [MPV Player](https://mpv.io/installation/)
- Python paketi: `ytmusicapi`

### Adımlar
1. Depoyu klonlayın:
   ```bash
   git clone https://github.com/meirza23/mei_player_v1.git
   cd mei-player
   ```
2. Python bağımlılığını yükleyin:
   ```bash
   pip install ytmusicapi
   ```
3. Go bağımlılıklarını yükleyin:
   ```bash
   go mod init mei-player
   go mod tidy
   ```
4. Gerekli dizinleri oluşturun:
   ```bash
   mkdir -p Songs Playlists/Downloads
   ```

## 🎮 Kullanım

Mei Player'ı kurmak ve çalıştırmak için:
```bash
 go install
 go build
 ./mei_player
sudo mv mei_player /usr/local/bin/
```

### Ana Menü
```
🎵 Mei Player 🎵
0. Çıkış
1. Şarkı Ara
2. Playlistleri Görüntüle
3. Şarkıları Görüntüle
```

### Temel İşlemler

#### Şarkı Arama
- Ana menüden `1` seçeneğini kullanarak arama yapabilirsiniz.
- Format: `d<numara>` (indir) veya `<numara>` (çal)
- Örnek: `d3` → 3. şarkıyı indirir.

#### Playlist Yönetimi
- Şarkı indirirken `E` seçerek bir playlist'e ekleyebilirsiniz.
- Yeni playlist oluşturabilir veya mevcut bir listeye ekleme yapabilirsiniz.

#### Şarkı Çalma
- Yerel şarkılar otomatik olarak `Songs/` klasörüne kaydedilir.
- Ana menüden `3` seçerek yerel şarkıları oynatabilirsiniz.

