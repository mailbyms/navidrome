package netease

// Base response structure for Netease API
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// Artist represents an artist from Netease API
type Artist struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	PicURL      string   `json:"picUrl"` // 图片链接
	Alias       []string `json:"alias"`
	AlbumSize   int      `json:"albumSize"`
	PicID       int64    `json:"picId"`
	// FansGroup   *string  `json:"fansGroup"`
	Img1v1URL   string   `json:"img1v1Url"`  // 方形图片
	Img1v1      int64    `json:"img1v1"`
	MVSize      int      `json:"mvSize"`
	// Followed    bool     `json:"followed"`
	// Alg         string   `json:"alg"`
	// Trans       *string  `json:"trans"`
	// BriefDesc   string   `json:"briefDesc,omitempty"`
	// MusicSize   int      `json:"musicSize,omitempty"`
	// TopicPerson int      `json:"topicPerson,omitempty"`
}

// ArtistDetail represents detailed artist information from /artist/detail
type ArtistDetail struct {
	// VideoCount int `json:"videoCount"`
	// Identify   struct {
	// 	ImageURL  *string `json:"imageUrl"`
	// 	ImageDesc string  `json:"imageDesc"`
	// 	ActionURL string  `json:"actionUrl"`
	// } `json:"identify"`
	Artist struct {
		ID          int64    `json:"id"`
		Cover       string   `json:"cover"`
		Avatar      string   `json:"avatar"`
		Name        string   `json:"name"`
		TransNames  []string `json:"transNames"`
		Alias       []string `json:"alias"`
		Identities  []string `json:"identities"`
		//IdentifyTag []string  `json:"identifyTag"`
		BriefDesc   string   `json:"briefDesc"`
		// Rank        struct {
		// 	Rank int `json:"rank"`
		// 	Type int `json:"type"`
		// } `json:"rank"`
		AlbumSize int `json:"albumSize"`
		MusicSize int `json:"musicSize"`
		MVSize    int `json:"mvSize"`
	} `json:"artist"`
	// Blacklist              bool `json:"blacklist"`
	// PreferShow             int  `json:"preferShow"`
	// ShowPriMsg             bool `json:"showPriMsg"`
	// SecondaryExpertIdentiy []struct {
	// 	ExpertIdentiyID    int    `json:"expertIdentiyId"`
	// 	ExpertIdentiyName  string `json:"expertIdentiyName"`
	// 	ExpertIdentiyCount int    `json:"expertIdentiyCount"`
	// } `json:"secondaryExpertIdentiy"`
}

// ArtistDetailResponse represents the response from /artist/detail
type ArtistDetailResponse struct {
	BaseResponse
	Data ArtistDetail `json:"data"`
}

// Album represents an album from Netease API
type Album struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	Size            int      `json:"size"`
	PicID           int64    `json:"picId"`
	BlurPicURL      string   `json:"blurPicUrl"`
	CompanyID       int      `json:"companyId"`
	Pic             int64    `json:"pic"`
	PicURL          string   `json:"picUrl"`
	PublishTime     int64    `json:"publishTime"`
	Description     string   `json:"description"`
	Tags            string   `json:"tags"`
	Company         string   `json:"company"`
	BriefDesc       string   `json:"briefDesc"`
	Artist          Artist   `json:"artist"`
	Songs           *[]Song  `json:"songs"`
	Alias           []string `json:"alias"`
	Status          int      `json:"status"`
	CopyrightID     int      `json:"copyrightId"`
	CommentThreadID string   `json:"commentThreadId"`
	Artists         []Artist `json:"artists"`
	Paid            bool     `json:"paid"`
	OnSale          bool     `json:"onSale"`
	PicIDStr        string   `json:"picId_str"`
	Alg             string   `json:"alg"`
	Mark            int64    `json:"mark"`
	ContainedSong   string   `json:"containedSong"`
	SubType         string   `json:"subType,omitempty"`
}

// AlbumDetail represents detailed album information
type AlbumDetail struct {
	Album
	Songs []Song `json:"songs"`
	Info  struct {
		CommentThread struct {
			ID           string `json:"id"`
			ResourceInfo struct {
				ID      int64  `json:"id"`
				UserID  int64  `json:"userId"`
				Name    string `json:"name"`
				ImgURL  string `json:"imgUrl"`
				Creator struct {
					UserID   int64  `json:"userId"`
					Nickname string `json:"nickname"`
				} `json:"creator"`
			} `json:"resourceInfo"`
			ResourceType int `json:"resourceType"`
		} `json:"commentThread"`
		LatestLikedUsers []User `json:"latestLikedUsers"`
		Liked            bool   `json:"liked"`
		Comments         []struct {
			User      User `json:"user"`
			BeReplied []struct {
				User    User   `json:"user"`
				Content string `json:"content"`
			} `json:"beReplied"`
			CommentID     int64  `json:"commentId"`
			Content       string `json:"content"`
			Time          int64  `json:"time"`
			LikedCount    int    `json:"likedCount"`
			Liked         bool   `json:"liked"`
			ExpressionURL string `json:"expressionUrl"`
		} `json:"comments"`
		CommentCount int `json:"commentCount"`
	} `json:"info"`
}

// Song represents a song from Netease API
type Song struct {
	ID                   int64    `json:"id"`
	Name                 string   `json:"name"`
	Artists              []Artist `json:"artists"`
	Album                Album    `json:"album"`
	Duration             int      `json:"duration"`
	CopyrightID          int      `json:"copyrightId"`
	Status               int      `json:"status"`
	Alias                []string `json:"alias"`
	Rtype                int      `json:"rtype"`
	Ftype                int      `json:"ftype"`
	MVId                 int64    `json:"mvid"`
	Fee                  int      `json:"fee"`
	RURL                 string   `json:"rUrl"`
	Mark                 int64    `json:"mark"`
	OriginCoverType      int      `json:"originCoverType"`
	OriginSongSimpleData struct {
		SongID    int64    `json:"songId"`
		Name      string   `json:"name"`
		Artists   []Artist `json:"artists"`
		AlbumMeta struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"albumMeta"`
	} `json:"originSongSimpleData,omitempty"`
	TagPicList        []interface{} `json:"tagPicList"`
	ResourceState     bool          `json:"resourceState"`
	Version           int           `json:"version"`
	SongJumpInfo      interface{}   `json:"songJumpInfo"`
	EntertainmentTags interface{}   `json:"entertainmentTags"`
	AwardTags         interface{}   `json:"awardTags"`
	Single            int           `json:"single"`
	NoCopyrightRcmd   interface{}   `json:"noCopyrightRcmd"`
	Mst               int           `json:"mst"`
	Cp                int           `json:"cp"`
	PublishTime       int64         `json:"publishTime"`
	Tns               []string      `json:"tns"`
}

// User represents a user from Netease API
type User struct {
	UserID     int64             `json:"userId"`
	Nickname   string            `json:"nickname"`
	AvatarURL  string            `json:"avatarUrl"`
	AuthStatus int               `json:"authStatus"`
	ExpertTags []string          `json:"expertTags"`
	Experts    map[string]string `json:"experts"`
	VipType    int               `json:"vipType"`
	RemarkName string            `json:"remarkName"`
	Mutual     bool              `json:"mutual"`
	Followed   bool              `json:"followed"`
	VipRights  struct {
		Associator struct {
			VipCode int  `json:"vipCode"`
			Rights  bool `json:"rights"`
		} `json:"associator"`
		MusicPackage struct {
			VipCode int  `json:"vipCode"`
			Rights  bool `json:"rights"`
		} `json:"musicPackage"`
		RedVipAnnualCount int `json:"redVipAnnualCount"`
		RedVipLevel       int `json:"redVipLevel"`
	} `json:"vipRights"`
}

// Search response structures
type SearchResponse struct {
	BaseResponse
	Result SearchResult `json:"result"`
}

type SearchResult struct {
	HasMore          bool       `json:"hasMore,omitempty"`
	Songs            []Song     `json:"songs,omitempty"`
	Albums           []Album    `json:"albums,omitempty"`
	Artists          []Artist   `json:"artists,omitempty"`
	Playlists        []Playlist `json:"playlists,omitempty"`
	SongCount        int        `json:"songCount,omitempty"`
	AlbumCount       int        `json:"albumCount,omitempty"`
	ArtistCount      int        `json:"artistCount,omitempty"`
	PlaylistCount    int        `json:"playlistCount,omitempty"`
	HlWords          []string   `json:"hlWords,omitempty"`
	SearchQcReminder *string    `json:"searchQcReminder,omitempty"`
}

// Playlist represents a playlist from Netease API
type Playlist struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	CoverImgURL  string   `json:"coverImgUrl"`
	Creator      User     `json:"creator"`
	Subscribed   bool     `json:"subscribed"`
	TrackCount   int      `json:"trackCount"`
	UserID       int64    `json:"userId"`
	PlayCount    int64    `json:"playCount"`
	BookCount    int      `json:"bookCount"`
	SpecialType  int      `json:"specialType"`
	OfficialType int      `json:"officialType"`
	CopyWriter   string   `json:"copywriter"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
}

// Album detail response
type AlbumDetailResponse struct {
	BaseResponse
	Songs []Song `json:"songs"`
	Album Album  `json:"album"`
}

// Similar artists response
type SimilarArtistsResponse struct {
	BaseResponse
	Artists []Artist `json:"artists"`
}

// Artist top songs response
type ArtistTopSongsResponse struct {
	BaseResponse
	Songs []Song `json:"songs"`
	More  bool   `json:"more"`
	Total int    `json:"total"`
}
