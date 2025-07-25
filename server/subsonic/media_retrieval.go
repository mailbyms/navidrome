package subsonic

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/consts"
	"github.com/navidrome/navidrome/core/lyrics"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/resources"
	"github.com/navidrome/navidrome/server/subsonic/filter"
	"github.com/navidrome/navidrome/server/subsonic/responses"
	"github.com/navidrome/navidrome/utils/gravatar"
	"github.com/navidrome/navidrome/utils/req"
)

func (api *Router) GetAvatar(w http.ResponseWriter, r *http.Request) (*responses.Subsonic, error) {
	if !conf.Server.EnableGravatar {
		return api.getPlaceHolderAvatar(w, r)
	}
	p := req.Params(r)
	username, err := p.String("username")
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	u, err := api.ds.User(ctx).FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if u.Email == "" {
		log.Warn(ctx, "User needs an email for gravatar to work", "username", username)
		return api.getPlaceHolderAvatar(w, r)
	}
	http.Redirect(w, r, gravatar.Url(u.Email, 0), http.StatusFound)
	return nil, nil
}

func (api *Router) getPlaceHolderAvatar(w http.ResponseWriter, r *http.Request) (*responses.Subsonic, error) {
	f, err := resources.FS().Open(consts.PlaceholderAvatar)
	if err != nil {
		log.Error(r, "Image not found", err)
		return nil, newError(responses.ErrorDataNotFound, "Avatar image not found")
	}
	defer f.Close()
	_, _ = io.Copy(w, f)

	return nil, nil
}

func (api *Router) GetCoverArt(w http.ResponseWriter, r *http.Request) (*responses.Subsonic, error) {
	// If context is already canceled, discard request without further processing
	if r.Context().Err() != nil {
		return nil, nil //nolint:nilerr
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	p := req.Params(r)
	id, _ := p.String("id")
	size := p.IntOr("size", 0)
	square := p.BoolOr("square", false)

	imgReader, lastUpdate, err := api.artwork.GetOrPlaceholder(ctx, id, size, square)
	switch {
	case errors.Is(err, context.Canceled):
		return nil, nil
	case errors.Is(err, model.ErrNotFound):
		log.Warn(r, "Couldn't find coverArt", "id", id, err)
		return nil, newError(responses.ErrorDataNotFound, "Artwork not found")
	case err != nil:
		log.Error(r, "Error retrieving coverArt", "id", id, err)
		return nil, err
	}

	defer imgReader.Close()
	w.Header().Set("cache-control", "public, max-age=315360000")
	w.Header().Set("last-modified", lastUpdate.Format(time.RFC1123))

	cnt, err := io.Copy(w, imgReader)
	if err != nil {
		log.Warn(ctx, "Error sending image", "count", cnt, err)
	}

	return nil, err
}

func (api *Router) GetLyrics(r *http.Request) (*responses.Subsonic, error) {
	p := req.Params(r)
	artist, _ := p.String("artist")
	title, _ := p.String("title")
	response := newResponse()
	lyricsResponse := responses.Lyrics{}
	response.Lyrics = &lyricsResponse
	mediaFiles, err := api.ds.MediaFile(r.Context()).GetAll(filter.SongsByArtistTitleWithLyricsFirst(artist, title))

	if err != nil {
		return nil, err
	}

	if len(mediaFiles) == 0 {
		return response, nil
	}

	structuredLyrics, err := lyrics.GetLyrics(r.Context(), &mediaFiles[0])
	if err != nil {
		return nil, err
	}

	if len(structuredLyrics) == 0 {
		return response, nil
	}

	lyricsResponse.Artist = artist
	lyricsResponse.Title = title

	lyricsText := ""
	for _, line := range structuredLyrics[0].Line {
		lyricsText += line.Value + "\n"
	}

	lyricsResponse.Value = lyricsText

	return response, nil
}

func (api *Router) GetLyricsBySongId(r *http.Request) (*responses.Subsonic, error) {
	id, err := req.Params(r).String("id")
	if err != nil {
		return nil, err
	}

	mediaFile, err := api.ds.MediaFile(r.Context()).Get(id)
	if err != nil {
		return nil, err
	}

	structuredLyrics, err := lyrics.GetLyrics(r.Context(), mediaFile)
	if err != nil {
		return nil, err
	}

	response := newResponse()
	response.LyricsList = buildLyricsList(mediaFile, structuredLyrics)

	return response, nil
}
