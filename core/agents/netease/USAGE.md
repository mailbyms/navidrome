# Netease Agent Usage Guide

## Quick Start

The Netease agent is now functional for basic music metadata retrieval. Here's how to use it:

### 1. Configuration

Add netease to your Navidrome configuration:

```toml
# navidrome.toml
Agents = "netease,lastfm,spotify"
```

### 2. What Works Now

#### ✅ Artist URL Retrieval
```go
// Returns: https://music.163.com/artist?id=6452
url, err := agent.GetArtistURL(ctx, "", "周杰伦", "")
```

#### ✅ Album Information
```go
// Returns album info with description and cover image
info, err := agent.GetAlbumInfo(ctx, "范特西", "周杰伦", "")
// info.Name = "范特西"
// info.URL = "https://music.163.com/album?id=32311"
// info.Images[0].URL = "http://p2.music.126.net/..."
```

#### ✅ Artist Images
```go
// Returns artist profile images (high-quality cover and avatar)
images, err := agent.GetArtistImages(ctx, "", "周杰伦", "")
// images[0].URL = "http://p2.music.126.net/.../cover.jpg" (high-quality cover)
// images[1].URL = "http://p2.music.126.net/.../avatar.jpg" (high-quality avatar)
```

#### ✅ Artist Biography
```go
// Returns detailed artist biography
bio, err := agent.GetArtistBiography(ctx, "", "周杰伦", "")
// Returns: "周杰伦（Jay Chou），1979年1月18日出生于台湾省新北市..."
```

### 3. API Endpoints Used

The agent currently uses these Netease API endpoints:

- `GET /search?keywords={name}&type=100` - Artist search
- `GET /search?keywords={artist} {album}&type=10` - Album search
- `GET /artist/detail?id={artistId}` - Artist detailed information

### 4. Response Examples

#### Artist Search Response
```json
{
  "code": 200,
  "result": {
    "artists": [{
      "id": 6452,
      "name": "周杰伦",
      "picUrl": "https://p1.music.126.net/...",
      "alias": ["Jay Chou"],
      "albumSize": 15,
      "mvSize": 31
    }]
  }
}
```

#### Artist Detail Response
```json
{
  "code": 200,
  "data": {
    "artist": {
      "id": 6452,
      "name": "周杰伦",
      "cover": "http://p2.music.126.net/.../cover.jpg",
      "avatar": "http://p2.music.126.net/.../avatar.jpg",
      "alias": ["Jay Chou", "周董"],
      "briefDesc": "周杰伦（Jay Chou），1979年1月18日出生于台湾省新北市...",
      "albumSize": 44,
      "musicSize": 569
    }
  }
}
```

#### Album Search Response
```json
{
  "code": 200,
  "result": {
    "albums": [{
      "id": 32311,
      "name": "范特西",
      "artist": {
        "id": 6452,
        "name": "周杰伦"
      },
      "picUrl": "http://p2.music.126.net/...",
      "description": "周杰伦第二张专辑"
    }]
  }
}
```

### 5. Error Handling

The agent handles various error scenarios:

- **No results found**: Returns `agents.ErrNotFound`
- **API errors**: Returns formatted error with status code
- **Network timeouts**: 10-second timeout with proper error handling
- **Invalid responses**: JSON parsing errors are handled gracefully

### 6. Performance Considerations

- **Request throttling**: Uses mutex to prevent concurrent API abuse
- **Timeout handling**: 10-second timeout for all requests
- **Efficient search**: Combines artist and album names for better results
- **Image optimization**: Returns multiple image sizes when available

### 7. Chinese Music Support

The Netease agent is particularly effective for:

- **Chinese artists**: 周杰伦, 邓紫棋, 林俊杰, etc.
- **Chinese albums**: Proper Chinese titles and descriptions
- **Mandarin/Cantonese content**: Native language metadata
- **Asian pop music**: K-pop, J-pop artists available on Netease

### 8. Integration with Other Agents

The agent works well in combination:

```toml
# Recommended configuration for Chinese music
Agents = "netease,lastfm,spotify"
```

- **Netease**: Chinese metadata and images
- **LastFM**: International metadata and scrobbling
- **Spotify**: Additional metadata and recommendations

### 9. Troubleshooting

#### Common Issues

1. **Empty results for Western artists**
   - Netease focuses on Chinese/Asian music
   - Try different agent order: `"lastfm,netease,spotify"`

2. **Slow response times**
   - Check network connectivity to Chinese servers
   - Consider VPN if accessing from restricted regions

3. **Missing album descriptions**
   - Some albums may not have descriptions in Netease
   - Other agents will be tried automatically

#### Debug Logging

Enable debug logging to see API calls:

```toml
LogLevel = "debug"
```

Look for log entries like:
```
Getting album info from Netease album=范特西 artist=周杰伦
Searching album on Netease album=范特西 artist=周杰伦
```

### 10. Next Steps

To extend functionality, you can provide additional API endpoints for:

- Artist biographies (`/artist/desc`)
- Similar artists (`/simi/artist`)
- Top songs (`/artist/top/song`)
- Detailed album info (`/album/detail`)

The framework is ready to accept these implementations!