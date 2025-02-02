from ytmusicapi import YTMusic
import sys
import json

def search_song(query):
    ytmusic = YTMusic()
    results = ytmusic.search(query, filter="songs")[:5]  # İlk 5 sonucu al
    print(json.dumps(results, ensure_ascii=False))

if __name__ == "__main__":
    query = " ".join(sys.argv[1:])  # Komut satırından argüman al
    search_song(query)
