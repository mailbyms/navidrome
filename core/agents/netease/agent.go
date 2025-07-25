package netease

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/core/agents"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
)

const (
	neteaseAgentName = "netease"
)

type neteaseAgent struct {
	ds           model.DataStore
	client       *client
	getInfoMutex sync.Mutex
}

func neteaseConstructor(ds model.DataStore) *neteaseAgent {
	// TODO: Add configuration checks if needed
	// For now, Netease API doesn't require API keys for basic metadata

	n := &neteaseAgent{
		ds: ds,
	}

	hc := &http.Client{
		Timeout: 10 * time.Second,
	}

	n.client = newClient(hc)
	return n
}

func (n *neteaseAgent) AgentName() string {
	return neteaseAgentName
}

// AlbumInfoRetriever interface implementation
func (n *neteaseAgent) GetAlbumInfo(ctx context.Context, name, artist, mbid string) (*agents.AlbumInfo, error) {
	log.Debug(ctx, "Getting album info from Netease", "album", name, "artist", artist)

	album, err := n.callAlbumSearch(ctx, name, artist)
	if err != nil {
		return nil, err
	}

	response := &agents.AlbumInfo{
		Name:        album.Name,
		MBID:        "", // Netease doesn't provide MBID
		Description: album.Description,
		URL:         fmt.Sprintf("https://music.163.com/album?id=%d", album.ID),
		Images:      make([]agents.ExternalImage, 0),
	}

	// Add album cover image if available
	if album.PicURL != "" {
		response.Images = append(response.Images, agents.ExternalImage{
			URL:  album.PicURL,
			Size: 0, // Netease doesn't specify image size
		})
	}

	return response, nil
}

// ArtistMBIDRetriever interface implementation
func (n *neteaseAgent) GetArtistMBID(ctx context.Context, id string, name string) (string, error) {
	log.Debug(ctx, "Getting artist MBID from Netease", "name", name)

	// Netease doesn't provide MBID directly, so we return not found
	// This allows other agents (like LastFM) to handle MBID retrieval
	return "", agents.ErrNotFound
}

// ArtistURLRetriever interface implementation
func (n *neteaseAgent) GetArtistURL(ctx context.Context, id, name, mbid string) (string, error) {
	log.Debug(ctx, "Getting artist URL from Netease", "name", name)

	artist, err := n.callArtistSearch(ctx, name)
	if err != nil {
		return "", err
	}

	// Generate Netease artist URL based on artist ID
	artistURL := fmt.Sprintf("https://music.163.com/artist?id=%d", artist.ID)
	return artistURL, nil
}

// ArtistBiographyRetriever interface implementation
func (n *neteaseAgent) GetArtistBiography(ctx context.Context, id, name, mbid string) (string, error) {
	log.Debug(ctx, "Getting artist biography from Netease", "name", name)

	// First, search for the artist to get their ID
	artist, err := n.callArtistSearch(ctx, name)
	if err != nil {
		return "", err
	}

	// Then get detailed artist information using the ID
	artistID := fmt.Sprintf("%d", artist.ID)
	detail, err := n.callArtistDetail(ctx, artistID)
	if err != nil {
		return "", err
	}

	// Extract biography from briefDesc
	biography := strings.TrimSpace(detail.Artist.BriefDesc)
	if biography == "" {
		return "", agents.ErrNotFound
	}

	return biography, nil
}

// ArtistSimilarRetriever interface implementation
func (n *neteaseAgent) GetSimilarArtists(ctx context.Context, id, name, mbid string, limit int) ([]agents.Artist, error) {
	// TODO: Implement similar artists retrieval
	log.Debug(ctx, "Getting similar artists from Netease", "name", name, "limit", limit)

	// Placeholder implementation
	return nil, agents.ErrNotFound
}

// ArtistImageRetriever interface implementation
func (n *neteaseAgent) GetArtistImages(ctx context.Context, id, name, mbid string) ([]agents.ExternalImage, error) {
	log.Debug(ctx, "Getting artist images from Netease", "name", name)

	// First get basic artist info from search
	artist, err := n.callArtistSearch(ctx, name)
	if err != nil {
		return nil, err
	}

	var images []agents.ExternalImage

	// Try to get detailed artist info for higher quality images
	artistID := fmt.Sprintf("%d", artist.ID)
	detail, detailErr := n.callArtistDetail(ctx, artistID)

	if detailErr == nil {
		// Use high-quality images from artist detail if available
		if detail.Artist.Cover != "" {
			images = append(images, agents.ExternalImage{
				URL:  detail.Artist.Cover,
				Size: 0, // High quality cover image
			})
		}

		if detail.Artist.Avatar != "" && detail.Artist.Avatar != detail.Artist.Cover {
			images = append(images, agents.ExternalImage{
				URL:  detail.Artist.Avatar,
				Size: 0, // High quality avatar image
			})
		}
	}

	// Fallback to search results if detail failed or no images found
	if len(images) == 0 {
		// Add main artist image if available
		if artist.PicURL != "" {
			images = append(images, agents.ExternalImage{
				URL:  artist.PicURL,
				Size: 0,
			})
		}

		// Add img1v1 image if different and available
		if artist.Img1v1URL != "" && artist.Img1v1URL != artist.PicURL {
			images = append(images, agents.ExternalImage{
				URL:  artist.Img1v1URL,
				Size: 0,
			})
		}
	}

	if len(images) == 0 {
		return nil, agents.ErrNotFound
	}

	return images, nil
}

// ArtistTopSongsRetriever interface implementation
func (n *neteaseAgent) GetArtistTopSongs(ctx context.Context, id, artistName, mbid string, count int) ([]agents.Song, error) {
	// TODO: Implement artist top songs retrieval
	log.Debug(ctx, "Getting artist top songs from Netease", "artist", artistName, "count", count)

	// Placeholder implementation
	return nil, agents.ErrNotFound
}

// Helper methods for API calls
func (n *neteaseAgent) callArtistSearch(ctx context.Context, name string) (*Artist, error) {
	n.getInfoMutex.Lock()
	defer n.getInfoMutex.Unlock()

	// TODO: Implement artist search API call
	return n.client.artistSearch(ctx, name)
}

func (n *neteaseAgent) callAlbumSearch(ctx context.Context, name, artist string) (*Album, error) {
	n.getInfoMutex.Lock()
	defer n.getInfoMutex.Unlock()

	// TODO: Implement album search API call
	return n.client.albumSearch(ctx, name, artist)
}

func (n *neteaseAgent) callArtistDetail(ctx context.Context, artistID string) (*ArtistDetail, error) {
	n.getInfoMutex.Lock()
	defer n.getInfoMutex.Unlock()

	// TODO: Implement artist detail API call
	return n.client.artistDetail(ctx, artistID)
}

func (n *neteaseAgent) callAlbumDetail(ctx context.Context, albumID string) (*AlbumDetail, error) {
	n.getInfoMutex.Lock()
	defer n.getInfoMutex.Unlock()

	// TODO: Implement album detail API call
	return n.client.albumDetail(ctx, albumID)
}

func (n *neteaseAgent) callSimilarArtists(ctx context.Context, artistID string) ([]Artist, error) {
	n.getInfoMutex.Lock()
	defer n.getInfoMutex.Unlock()

	// TODO: Implement similar artists API call
	return n.client.similarArtists(ctx, artistID)
}

func (n *neteaseAgent) callArtistTopSongs(ctx context.Context, artistID string, limit int) ([]Song, error) {
	n.getInfoMutex.Lock()
	defer n.getInfoMutex.Unlock()

	// TODO: Implement artist top songs API call
	return n.client.artistTopSongs(ctx, artistID, limit)
}

func init() {
	conf.AddHook(func() {
		agents.Register(neteaseAgentName, func(ds model.DataStore) agents.Interface {
			a := neteaseConstructor(ds)
			if a != nil {
				return a
			}
			return nil
		})
	})
}
