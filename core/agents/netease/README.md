# Netease Cloud Music Agent

This agent integrates with Netease Cloud Music (网易云音乐) API to retrieve music metadata for Navidrome.

## Features

The Netease agent implements the following interfaces:

- **AlbumInfoRetriever**: Get album information including description and images
- **ArtistMBIDRetriever**: Get MusicBrainz ID for artists (if available)
- **ArtistURLRetriever**: Get artist page URLs on Netease Cloud Music
- **ArtistBiographyRetriever**: Get artist biographies and descriptions
- **ArtistSimilarRetriever**: Get similar artists recommendations
- **ArtistImageRetriever**: Get artist images and photos
- **ArtistTopSongsRetriever**: Get artist's top/popular songs

## Configuration

To enable the Netease agent, add it to your Navidrome configuration:

```toml
# In navidrome.toml
Agents = "netease,lastfm,spotify"
```

Or via environment variable:
```bash
export ND_AGENTS="netease,lastfm,spotify"
```

## API Endpoints

The agent uses the following Netease Cloud Music API endpoints:

### Search APIs
- `/search` - Search for artists, albums, songs
- `/cloudsearch` - Enhanced search with more results

### Artist APIs
- `/artist/detail` - Get detailed artist information
- `/artist/desc` - Get artist description/biography
- `/artist/songs` - Get artist's songs
- `/artist/album` - Get artist's albums
- `/simi/artist` - Get similar artists

### Album APIs
- `/album/detail` - Get detailed album information
- `/album/desc` - Get album description

## Implementation Status

### ✅ Completed
- Basic project structure
- Response data models
- HTTP client setup
- Test framework setup
- Artist search implementation
- Album search implementation
- Basic album info retrieval
- Artist URL generation
- Artist image retrieval (with high-quality fallback)
- Artist detail retrieval (/artist/detail)
- Artist biography extraction
- Error handling and rate limiting

### 🚧 TODO (Waiting for additional API details)
- Album detail retrieval (full metadata)
- Similar artists retrieval
- Top songs retrieval
- Enhanced error handling for specific API errors

## API Response Examples

### Artist Search Response
```json
{
  "code": 200,
  "result": {
    "artists": [{
      "id": 6452,
      "name": "周杰伦",
      "picUrl": "http://p1.music.126.net/...",
      "alias": ["Jay Chou"],
      "briefDesc": "华语流行音乐天王"
    }],
    "artistCount": 1
  }
}
```

### Album Detail Response
```json
{
  "code": 200,
  "album": {
    "id": 32311,
    "name": "范特西",
    "artist": {
      "id": 6452,
      "name": "周杰伦"
    },
    "publishTime": 1000000000000,
    "description": "周杰伦第二张专辑",
    "songs": [...]
  }
}
```

## Testing

Run the tests with:
```bash
go test ./core/agents/netease/...
```

Note: Most tests are currently skipped pending API implementation.

## Contributing

When implementing API methods:

1. Update the corresponding method in `client.go`
2. Add proper error handling
3. Update the agent method in `agent.go`
4. Add/update tests
5. Update this README

## Rate Limiting

Netease Cloud Music API has rate limiting. The client includes:
- Request mutex to prevent concurrent API abuse
- 10-second HTTP timeout
- Proper User-Agent and Referer headers

## Error Handling

The agent handles various error scenarios:
- API rate limiting
- Network timeouts
- Invalid responses
- Missing data (returns `agents.ErrNotFound`)

## Notes

- Netease Cloud Music primarily serves Chinese music content
- Some content may be geo-restricted
- API endpoints may change without notice
- Consider implementing caching for frequently requested data