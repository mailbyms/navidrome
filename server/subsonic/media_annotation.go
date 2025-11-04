package subsonic

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/navidrome/navidrome/core/agents"
	"github.com/navidrome/navidrome/core/scrobbler"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/model/request"
	"github.com/navidrome/navidrome/server/events"
	"github.com/navidrome/navidrome/server/subsonic/responses"
	"github.com/navidrome/navidrome/utils/req"
)

func (api *Router) SetRating(r *http.Request) (*responses.Subsonic, error) {
	p := req.Params(r)
	id, err := p.String("id")
	if err != nil {
		return nil, err
	}
	rating, err := p.Int("rating")
	if err != nil {
		return nil, err
	}

	log.Debug(r, "Setting rating", "rating", rating, "id", id)
	err = api.setRating(r.Context(), id, rating)
	if err != nil {
		log.Error(r, err)
		return nil, err
	}

	return newResponse(), nil
}

func (api *Router) setRating(ctx context.Context, id string, rating int) error {
	var repo model.AnnotatedRepository
	var resource string

	entity, err := model.GetEntityByID(ctx, api.ds, id)
	if err != nil {
		return err
	}
	switch entity.(type) {
	case *model.Artist:
		repo = api.ds.Artist(ctx)
		resource = "artist"
	case *model.Album:
		repo = api.ds.Album(ctx)
		resource = "album"
	default:
		repo = api.ds.MediaFile(ctx)
		resource = "song"
	}
	err = repo.SetRating(rating, id)
	if err != nil {
		return err
	}
	event := &events.RefreshResource{}
	api.broker.SendMessage(ctx, event.With(resource, id))
	return nil
}

func (api *Router) Star(r *http.Request) (*responses.Subsonic, error) {
	p := req.Params(r)
	ids, _ := p.Strings("id")
	albumIds, _ := p.Strings("albumId")
	artistIds, _ := p.Strings("artistId")
	if len(ids)+len(albumIds)+len(artistIds) == 0 {
		return nil, newError(responses.ErrorMissingParameter, "Required id parameter is missing")
	}
	ids = append(ids, albumIds...)
	ids = append(ids, artistIds...)

	err := api.setStar(r.Context(), true, ids...)
	if err != nil {
		return nil, err
	}

	return newResponse(), nil
}

func (api *Router) Unstar(r *http.Request) (*responses.Subsonic, error) {
	p := req.Params(r)
	ids, _ := p.Strings("id")
	albumIds, _ := p.Strings("albumId")
	artistIds, _ := p.Strings("artistId")
	if len(ids)+len(albumIds)+len(artistIds) == 0 {
		return nil, newError(responses.ErrorMissingParameter, "Required id parameter is missing")
	}
	ids = append(ids, albumIds...)
	ids = append(ids, artistIds...)

	err := api.setStar(r.Context(), false, ids...)
	if err != nil {
		return nil, err
	}

	return newResponse(), nil
}

func (api *Router) setStar(ctx context.Context, star bool, ids ...string) error {
	if len(ids) == 0 {
		return nil
	}
	log.Debug(ctx, "Changing starred", "ids", ids, "starred", star)
	if len(ids) == 0 {
		log.Warn(ctx, "Cannot star/unstar an empty list of ids")
		return nil
	}
	event := &events.RefreshResource{}
	err := api.ds.WithTxImmediate(func(tx model.DataStore) error {
		for _, id := range ids {
			exist, err := tx.Album(ctx).Exists(id)
			if err != nil {
				return err
			}
			if exist {
				err = tx.Album(ctx).SetStar(star, id)
				if err != nil {
					return err
				}
				event = event.With("album", id)
				continue
			}
			exist, err = tx.Artist(ctx).Exists(id)
			if err != nil {
				return err
			}
			if exist {
				err = tx.Artist(ctx).SetStar(star, id)
				if err != nil {
					return err
				}
				event = event.With("artist", id)
				continue
			}
			err = tx.MediaFile(ctx).SetStar(star, id)
			if err != nil {
				return err
			}
			event = event.With("song", id)
		}
		api.broker.SendMessage(ctx, event)
		return nil
	})
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

func (api *Router) Scrobble(r *http.Request) (*responses.Subsonic, error) {
	p := req.Params(r)
	ids, err := p.Strings("id")
	if err != nil {
		return nil, err
	}
	times, _ := p.Times("time")
	if len(times) > 0 && len(times) != len(ids) {
		return nil, newError(responses.ErrorGeneric, "Wrong number of timestamps: %d, should be %d", len(times), len(ids))
	}
	submission := p.BoolOr("submission", true)
	ctx := r.Context()

	if submission {
		err := api.scrobblerSubmit(ctx, ids, times)
		if err != nil {
			log.Error(ctx, "Error registering scrobbles", "ids", ids, "times", times, err)
		}
	} else {
		err := api.scrobblerNowPlaying(ctx, ids[0])
		if err != nil {
			log.Error(ctx, "Error setting NowPlaying", "id", ids[0], err)
		}
	}

	return newResponse(), nil
}

func (api *Router) scrobblerSubmit(ctx context.Context, ids []string, times []time.Time) error {
	var submissions []scrobbler.Submission
	log.Debug(ctx, "Scrobbling tracks", "ids", ids, "times", times)
	for i, id := range ids {
		var t time.Time
		if len(times) > 0 {
			t = times[i]
		} else {
			t = time.Now()
		}
		submissions = append(submissions, scrobbler.Submission{TrackID: id, Timestamp: t})
	}

	return api.scrobbler.Submit(ctx, submissions)
}

func (api *Router) scrobblerNowPlaying(ctx context.Context, trackId string) error {
	mf, err := api.ds.MediaFile(ctx).Get(trackId)
	if err != nil {
		return err
	}
	if mf == nil {
		return fmt.Errorf(`ID "%s" not found`, trackId)
	}

	player, _ := request.PlayerFrom(ctx)
	username, _ := request.UsernameFrom(ctx)
	client, _ := request.ClientFrom(ctx)
	clientId, ok := request.ClientUniqueIdFrom(ctx)
	if !ok {
		clientId = player.ID
	}

	log.Info(ctx, "Now Playing", "title", mf.Title, "artist", mf.Artist, "user", username, "player", player.Name)
	err = api.scrobbler.NowPlaying(ctx, clientId, client, trackId)
	return err
}

func (api *Router) GetSongComments(r *http.Request) (*responses.Subsonic, error) {
	p := req.Params(r)
	songID, err := p.String("id")
	if err != nil {
		return nil, err
	}

	// 使用新的分页参数格式
	pageSize, _ := p.Int("pageSize")
	if pageSize == 0 {
		pageSize = 10
	}
	pageNo, _ := p.Int("pageNo")
	if pageNo == 0 {
		pageNo = 1
	}
	sortType, _ := p.Int("sortType")
	if sortType == 0 {
		sortType = 2
	}

	log.Debug(r, "Getting song comments", "id", songID, "pageSize", pageSize, "pageNo", pageNo)

	// 首先获取歌曲信息，用于匹配网易云音乐的歌曲
	mf, err := api.ds.MediaFile(r.Context()).Get(songID)
	if err != nil {
		log.Error(r, "Error getting song", "id", songID, err)
		return nil, err
	}
	if mf == nil {
		return nil, newError(responses.ErrorDataNotFound, "Song not found")
	}

	// 获取歌曲评论
	songComments, err := api.getSongComments(r.Context(), r, mf, pageSize, pageNo, sortType)
	if err != nil {
		log.Error(r, "Error getting song comments", "id", songID, err)
		return nil, err
	}

	response := newResponse()
	response.SongComments = songComments

	return response, nil
}

func (api *Router) getSongComments(ctx context.Context, r *http.Request, mf *model.MediaFile, pageSize int, pageNo int, sortType int) (*responses.SongComments, error) {
	// 获取agents实例
	agentsInstance := agents.GetAgents(api.ds)

	// 使用agents接口获取评论，传入分页和排序参数
	comments, total, err := agentsInstance.GetSongComments(ctx, mf.Title, mf.Artist, pageSize, pageNo, sortType)
	if err != nil {
		log.Warn(ctx, "Failed to get comments from agents", "title", mf.Title, "artist", mf.Artist, "sortType", sortType, err)
		return nil, err
	}

	// 转换为Subsonic格式的评论
	var result []responses.SongComment
	for _, comment := range comments {
		result = append(result, responses.SongComment{
			ID:          comment.ID,
			User:        comment.User,
			AvatarURL:   comment.AvatarURL,
			Content:     comment.Content,
			Timestamp:   comment.Timestamp,
			LikedCount:  comment.LikedCount,
			Liked:       comment.Liked,
		})
	}

	return &responses.SongComments{
		Comments:     result,
		Total:        total,
		CommentCount: total,
	}, nil
}
