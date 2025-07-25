# Netease Agent Configuration Example

## Navidrome Configuration

Add netease to your agents configuration:

### Option 1: Configuration File (navidrome.toml)
```toml
# Enable netease agent along with other agents
Agents = "netease,lastfm,spotify"

# Optional: Set agent priority (netease will be tried first)
Agents = "netease,lastfm,spotify"
```

### Option 2: Environment Variables
```bash
# Enable netease agent
export ND_AGENTS="netease,lastfm,spotify"

# Or set it in your docker-compose.yml
environment:
  - ND_AGENTS=netease,lastfm,spotify
```

### Option 3: Command Line
```bash
./navidrome --agents "netease,lastfm,spotify"
```

## Agent Priority

Agents are tried in the order specified. For example:
- `"netease,lastfm,spotify"` - Try Netease first, then LastFM, then Spotify
- `"lastfm,netease,spotify"` - Try LastFM first, then Netease, then Spotify

## Expected Behavior

Once configured, the Netease agent will:

1. **Artist Information**: Provide Chinese artist names, detailed biographies, and high-quality images
2. **Album Information**: Get album descriptions and cover art from Netease
3. **Artist Biographies**: Rich biographical information in Chinese
4. **High-Quality Images**: Artist covers and avatars in high resolution
5. **Localization**: Provide Chinese metadata for international artists
6. **Artist URLs**: Direct links to Netease Cloud Music artist pages

## Benefits for Chinese Music

The Netease agent is particularly useful for:
- Chinese/Mandarin music collections
- Getting proper Chinese artist and album names
- Accessing rich metadata for Asian artists
- Better recommendations for Chinese music fans

## Fallback Behavior

If Netease agent fails to find information:
1. The next agent in the list (e.g., LastFM) will be tried
2. If all agents fail, local metadata will be used
3. No errors will be shown to the user - it's transparent

## Troubleshooting

### Agent Not Working
1. Check that "netease" is in your Agents configuration
2. Verify network connectivity to music.163.com
3. Check Navidrome logs for any error messages

### Slow Response
1. Netease API might be rate-limited
2. Try moving netease later in the agent priority list
3. Check your network connection to Chinese servers

### Missing Metadata
1. Not all artists/albums are available on Netease
2. Try different agent order
3. Some content might be geo-restricted

## Development Notes

This configuration will work once the API implementation is completed. The agent structure is ready and waiting for the actual API integration code.