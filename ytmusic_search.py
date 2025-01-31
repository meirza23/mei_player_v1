import sys
import json
from ytmusicapi import YTMusic

def search_ytmusic(query):
    try:
        ytmusic = YTMusic("headers_auth.json")
        results = ytmusic.search(query, filter="songs", limit=5)
        
        formatted = []
        for song in results[:5]:
            formatted.append({
                "title": song.get("title", "Bilinmiyor"),
                "artists": [artist["name"] for artist in song.get("artists", [])],
                "duration": song.get("duration", "Bilinmiyor"),
                "videoId": song.get("videoId", "")
            })
        
        return json.dumps(formatted)
        
    except Exception as e:
        return json.dumps({"error": str(e)})

if __name__ == "__main__":
    if len(sys.argv) > 1:
        print(search_ytmusic(sys.argv[1]))