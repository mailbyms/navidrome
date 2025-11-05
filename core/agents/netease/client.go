package netease

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/navidrome/navidrome/core/agents"
	"github.com/navidrome/navidrome/log"
)

const (
	neteaseAPIBaseURL = "http://netease.api:3000"
)

type client struct {
	hc *http.Client
}

func newClient(hc *http.Client) *client {
	return &client{hc: hc}
}

// Artist search - search for artists by name
func (c *client) artistSearch(ctx context.Context, name string) (*Artist, error) {
	log.Trace(ctx, "Searching artist on Netease", "name", name)

	params := c.buildSearchParams(name, SearchTypeArtist, 10, 0)
	resp, err := c.makeRequest(ctx, "/search", params)
	if err != nil {
		return nil, err
	}

	var searchResp SearchResponse
	if err := c.parseResponse(resp, &searchResp); err != nil {
		return nil, err
	}

	if searchResp.Code != 200 {
		return nil, fmt.Errorf("netease API error: code %d", searchResp.Code)
	}

	if len(searchResp.Result.Artists) == 0 {
		return nil, agents.ErrNotFound
	}

	// Return the first (most relevant) artist
	return &searchResp.Result.Artists[0], nil
}

// Album search - search for albums by name and artist
func (c *client) albumSearch(ctx context.Context, name, artist string) (*Album, error) {
	log.Trace(ctx, "Searching album on Netease", "album", name, "artist", artist)

	// Combine album name and artist for better search results
	keywords := name
	if artist != "" {
		keywords = artist + " " + name
	}

	params := c.buildSearchParams(keywords, SearchTypeAlbum, 10, 0)
	resp, err := c.makeRequest(ctx, "/search", params)
	if err != nil {
		return nil, err
	}

	var searchResp SearchResponse
	if err := c.parseResponse(resp, &searchResp); err != nil {
		return nil, err
	}

	if searchResp.Code != 200 {
		return nil, fmt.Errorf("netease API error: code %d", searchResp.Code)
	}

	if len(searchResp.Result.Albums) == 0 {
		return nil, agents.ErrNotFound
	}

	// Return the first (most relevant) album
	return &searchResp.Result.Albums[0], nil
}

// Artist detail - get detailed artist information by ID
func (c *client) artistDetail(ctx context.Context, artistID string) (*ArtistDetail, error) {
	log.Trace(ctx, "Getting artist detail from Netease", "artistID", artistID)

	params := url.Values{}
	params.Set("id", artistID)

	resp, err := c.makeRequest(ctx, "/artist/detail", params)
	if err != nil {
		return nil, err
	}

	var detailResp ArtistDetailResponse
	if err := c.parseResponse(resp, &detailResp); err != nil {
		return nil, err
	}

	if detailResp.Code != 200 {
		return nil, fmt.Errorf("netease API error: code %d, message: %s", detailResp.Code, detailResp.Message)
	}

	return &detailResp.Data, nil
}

// Album detail - get detailed album information by ID
func (c *client) albumDetail(ctx context.Context, albumID string) (*AlbumDetail, error) {
	// TODO: Implement album detail API call
	log.Trace(ctx, "Getting album detail from Netease", "albumID", albumID)

	// Placeholder implementation
	return nil, agents.ErrNotFound
}

// Similar artists - get similar artists by artist ID
func (c *client) similarArtists(ctx context.Context, artistID string) ([]Artist, error) {
	// TODO: Implement similar artists API call
	log.Trace(ctx, "Getting similar artists from Netease", "artistID", artistID)

	// Placeholder implementation
	return nil, agents.ErrNotFound
}

// Artist top songs - get top songs by artist ID
func (c *client) artistTopSongs(ctx context.Context, artistID string, limit int) ([]Song, error) {
	// TODO: Implement artist top songs API call
	log.Trace(ctx, "Getting artist top songs from Netease", "artistID", artistID, "limit", limit)

	// Placeholder implementation
	return nil, agents.ErrNotFound
}

// Song comments - get comments for a song by song ID
func (c *client) songComments(ctx context.Context, songID string, pageSize int, pageNo int, sortType int) (*CommentsResponse, error) {
	log.Trace(ctx, "Getting song comments from Netease", "songID", songID, "pageSize", pageSize, "pageNo", pageNo, "sortType", sortType)

	params := url.Values{}
	params.Set("id", songID)
	params.Set("type", "0") // 0 for song comments

	// 网易新版评论接口参数
	if pageSize > 0 {
		params.Set("pageSize", strconv.Itoa(pageSize))
	}

	if pageNo > 0 {
		params.Set("pageNo", strconv.Itoa(pageNo))
	}

	// 排序方式 (1:按推荐排序, 2:按热度排序, 3:按时间排序)
	if sortType > 0 && sortType <= 3 {
		params.Set("sortType", strconv.Itoa(sortType))
	} else {
		// 默认按热度排序
		log.Trace(ctx, "set sortType to default value: 2")
		params.Set("sortType", "2")
	}

	resp, err := c.makeRequest(ctx, "/comment/new", params)
	if err != nil {
		return nil, err
	}

	var commentsResp CommentsResponse
	if err := c.parseResponse(resp, &commentsResp); err != nil {
		return nil, err
	}

	if commentsResp.Code != 200 {
		return nil, fmt.Errorf("netease API error: code %d, message: %s", commentsResp.Code, commentsResp.Message)
	}


	return &commentsResp, nil
}

// Helper method to make HTTP requests to Netease API
func (c *client) makeRequest(ctx context.Context, endpoint string, params url.Values) (*http.Response, error) {
	reqURL := fmt.Sprintf("%s%s", neteaseAPIBaseURL, endpoint)

	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	// Set common headers for Netease API
	req.Header.Set("User-Agent", "Navidrome/1.0")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("netease API returned status %d", resp.StatusCode)
	}

	return resp, nil
}

// Helper method to parse JSON response
func (c *client) parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(target)
}

// Helper method to build search parameters
func (c *client) buildSearchParams(keywords string, searchType string, limit int, offset int) url.Values {
	params := url.Values{}
	params.Set("keywords", keywords)
	params.Set("type", searchType)
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))
	return params
}

// Search types for Netease API
const (
	SearchTypeSong     = "1"
	SearchTypeAlbum    = "10"
	SearchTypeArtist   = "100"
	SearchTypePlaylist = "1000"
	SearchTypeMV       = "1004"
	SearchTypeLyric    = "1006"
	SearchTypeUser     = "1002"
)
