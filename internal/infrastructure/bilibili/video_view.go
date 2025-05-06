package bilibili

import (
	"context"
	"fmt"
	"net/url"

	"github.com/krisxia0506/bilibili-watcher/internal/application"
)

// --- 响应结构体定义 ---

// VideoViewPage 视频分P信息。
type VideoViewPage struct {
	Cid        int64              `json:"cid"`
	Page       int                `json:"page"`
	From       string             `json:"from"`
	Part       string             `json:"part"`
	Duration   int64              `json:"duration"`
	Vid        string             `json:"vid"`
	Weblink    string             `json:"weblink"`
	Dimension  VideoViewDimension `json:"dimension"`
	FirstFrame string             `json:"first_frame,omitempty"` // 可能不存在
	Ctime      int64              `json:"ctime,omitempty"`       // 可能不存在
}

// VideoViewRights 视频权限信息。
type VideoViewRights struct {
	Bp            int `json:"bp"`
	Elec          int `json:"elec"`
	Download      int `json:"download"`
	Movie         int `json:"movie"`
	Pay           int `json:"pay"`
	Hd5           int `json:"hd5"`
	NoReprint     int `json:"no_reprint"`
	Autoplay      int `json:"autoplay"`
	UgcPay        int `json:"ugc_pay"`
	IsCooperation int `json:"is_cooperation"`
	UgcPayPreview int `json:"ugc_pay_preview"`
	NoBackground  int `json:"no_background"`
	CleanMode     int `json:"clean_mode"`
	IsSteinGate   int `json:"is_stein_gate"`
	Is360         int `json:"is_360"`
	NoShare       int `json:"no_share"`
	ArcPay        int `json:"arc_pay"`
	FreeWatch     int `json:"free_watch"`
}

// VideoViewOwner UP主信息。
type VideoViewOwner struct {
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
	Face string `json:"face"`
}

// VideoViewStat 视频统计信息。
type VideoViewStat struct {
	Aid        int64  `json:"aid"`
	View       int    `json:"view"`
	Danmaku    int    `json:"danmaku"`
	Reply      int    `json:"reply"`
	Favorite   int    `json:"favorite"`
	Coin       int    `json:"coin"`
	Share      int    `json:"share"`
	NowRank    int    `json:"now_rank"`
	HisRank    int    `json:"his_rank"`
	Like       int    `json:"like"`
	Dislike    int    `json:"dislike"`
	Evaluation string `json:"evaluation"`
	Vt         int    `json:"vt"`
}

// VideoViewDimension 视频尺寸信息。
type VideoViewDimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Rotate int `json:"rotate"`
}

// VideoViewDescV2 新版视频简介条目。
type VideoViewDescV2 struct {
	RawText string `json:"raw_text"`
	Type    int    `json:"type"`
	BizID   int    `json:"biz_id"`
}

// VideoViewArgueInfo 争议信息。
type VideoViewArgueInfo struct {
	ArgueMsg  string `json:"argue_msg"`
	ArgueType int    `json:"argue_type"`
	ArgueLink string `json:"argue_link"`
}

// VideoViewSubtitle 字幕信息。
type VideoViewSubtitle struct {
	AllowSubmit bool          `json:"allow_submit"`
	List        []interface{} `json:"list"` // 字幕列表结构未知，暂用 interface{}
}

// VideoViewUserGarb 用户装扮信息。
type VideoViewUserGarb struct {
	UrlImageAniCut string `json:"url_image_ani_cut"`
}

// VideoViewHonorReply 荣誉信息。
type VideoViewHonorReply struct {
	// 结构未知
}

// VideoViewData /x/web-interface/view API 响应中的 data 字段结构体。
type VideoViewData struct {
	Bvid                    string              `json:"bvid"`
	Aid                     int64               `json:"aid"`
	Videos                  int                 `json:"videos"` // 分P数量
	Tid                     int                 `json:"tid"`
	TidV2                   int                 `json:"tid_v2"`
	Tname                   string              `json:"tname"`
	TnameV2                 string              `json:"tname_v2"`
	Copyright               int                 `json:"copyright"`
	Pic                     string              `json:"pic"`
	Title                   string              `json:"title"`
	Pubdate                 int64               `json:"pubdate"` // 发布时间戳
	Ctime                   int64               `json:"ctime"`   // 投稿时间戳
	Desc                    string              `json:"desc"`
	DescV2                  []VideoViewDescV2   `json:"desc_v2"`
	State                   int                 `json:"state"`
	Duration                int64               `json:"duration"` // 总时长(秒)
	MissionID               int64               `json:"mission_id,omitempty"`
	Rights                  VideoViewRights     `json:"rights"`
	Owner                   VideoViewOwner      `json:"owner"`
	Stat                    VideoViewStat       `json:"stat"`
	ArgueInfo               VideoViewArgueInfo  `json:"argue_info"`
	Dynamic                 string              `json:"dynamic"`
	Cid                     int64               `json:"cid"` // 当前访问的分P的cid
	Dimension               VideoViewDimension  `json:"dimension"`
	Premiere                interface{}         `json:"premiere"` // 首映信息，结构未知
	TeenageMode             int                 `json:"teenage_mode"`
	IsChargeableSeason      bool                `json:"is_chargeable_season"`
	IsStory                 bool                `json:"is_story"`
	IsUpowerExclusive       bool                `json:"is_upower_exclusive"`
	IsUpowerPlay            bool                `json:"is_upower_play"`
	IsUpowerPreview         bool                `json:"is_upower_preview"`
	EnableVt                int                 `json:"enable_vt"`
	VtDisplay               string              `json:"vt_display"`
	IsUpowerExclusiveWithQa bool                `json:"is_upower_exclusive_with_qa"`
	NoCache                 bool                `json:"no_cache"`
	Pages                   []VideoViewPage     `json:"pages"` // 分P列表
	Subtitle                VideoViewSubtitle   `json:"subtitle"`
	IsSeasonDisplay         bool                `json:"is_season_display"`
	UserGarb                VideoViewUserGarb   `json:"user_garb"`
	HonorReply              VideoViewHonorReply `json:"honor_reply"`
	LikeIcon                string              `json:"like_icon"`
	NeedJumpBv              bool                `json:"need_jump_bv"`
	DisableShowUpInfo       bool                `json:"disable_show_up_info"`
	IsStoryPlay             int                 `json:"is_story_play"`
	IsViewSelf              bool                `json:"is_view_self"`
}

// VideoViewResponse /x/web-interface/view API 的响应结构体。
type VideoViewResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	TTL     int           `json:"ttl"`
	Data    VideoViewData `json:"data,omitempty"` // 成功时才有 data
}

// --- API 方法 ---

// GetVideoView 调用 Bilibili API 获取视频的详细信息，并返回应用层 DTO。
// 实现 application.BilibiliClient 接口的一部分。
// aid: 视频的 AV 号 (可选)
// bvid: 视频的 BV 号 (可选)
// aid 和 bvid 必须提供一个。
func (c *Client) GetVideoView(ctx context.Context, aid, bvid string) (*application.VideoViewDTO, error) {
	const path = "/x/web-interface/view"

	if aid == "" && bvid == "" {
		return nil, fmt.Errorf("either aid or bvid must be provided for GetVideoView")
	}

	params := url.Values{}
	if aid != "" {
		params.Set("aid", aid)
	} else {
		params.Set("bvid", bvid)
	}

	var resp VideoViewResponse
	err := c.Get(path, params, &resp)
	if err != nil {
		// 底层 Get 方法已处理 HTTP 和解码错误
		// 尝试提取 Bilibili API 的特定错误信息
		if apiErrResp, ok := err.(interface{ GetResponse() *VideoViewResponse }); ok && apiErrResp.GetResponse() != nil {
			return nil, fmt.Errorf("bilibili api error: code=%d, message=%s, underlying error: %w",
				apiErrResp.GetResponse().Code, apiErrResp.GetResponse().Message, err)
		}
		return nil, err
	}

	// 检查 Bilibili API 返回的业务状态码
	if resp.Code != 0 {
		// 返回包含错误信息的响应体
		return nil, fmt.Errorf("bilibili api error: code=%d, message=%s", resp.Code, resp.Message)
	}

	// 映射到应用层 DTO
	// 检查 Data 是否有效
	if resp.Data.Aid == 0 && resp.Data.Bvid == "" {
		// 如果关键标识符缺失，认为数据无效
		return nil, fmt.Errorf("received invalid data from GetVideoView: missing aid and bvid")
	}

	// 映射 Pages
	pagesDTO := make([]application.VideoViewPageDTO, 0, len(resp.Data.Pages))
	for _, p := range resp.Data.Pages {
		pagesDTO = append(pagesDTO, application.VideoViewPageDTO{
			Cid:      p.Cid,
			Part:     p.Part,
			Duration: p.Duration,
			Page:     p.Page,
		})
	}

	dto := &application.VideoViewDTO{
		Bvid:      resp.Data.Bvid,
		Aid:       resp.Data.Aid,
		Title:     resp.Data.Title,
		Desc:      resp.Data.Desc, // TODO: Consider using resp.Data.DescV2 for richer description
		Pubdate:   resp.Data.Pubdate,
		Duration:  resp.Data.Duration,
		OwnerName: resp.Data.Owner.Name,
		Pages:     pagesDTO, // 填充 Pages DTO
	}

	return dto, nil
}
