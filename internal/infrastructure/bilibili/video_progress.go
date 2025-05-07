package bilibili

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/krisxia0506/bilibili-watcher/internal/application"
)

// VideoProgressData 定义 Bilibili 视频进度 API 响应中 data 字段的结构体。
type VideoProgressData struct {
	Aid                     int64         `json:"aid"`
	Bvid                    string        `json:"bvid"`
	AllowBp                 bool          `json:"allow_bp"`
	NoShare                 bool          `json:"no_share"`
	Cid                     int64         `json:"cid"` // 注意：API 返回的是当前视频页面的 CID，而不是上次播放的 CID
	MaxLimit                int           `json:"max_limit"`
	PageNo                  int           `json:"page_no"`
	HasNext                 bool          `json:"has_next"`
	IpInfo                  IpInfo        `json:"ip_info"`
	LoginMid                int           `json:"login_mid"`
	LoginMidHash            string        `json:"login_mid_hash"`
	IsOwner                 bool          `json:"is_owner"`
	Name                    string        `json:"name"`
	Permission              string        `json:"permission"`
	LevelInfo               LevelInfo     `json:"level_info"`
	Vip                     VipInfo       `json:"vip"`
	AnswerStatus            int           `json:"answer_status"`
	BlockTime               int           `json:"block_time"`
	Role                    string        `json:"role"`
	LastPlayTime            int64         `json:"last_play_time"` // 观看进度，单位毫秒
	LastPlayCid             int64         `json:"last_play_cid"`  // 上次播放的视频分 P ID
	NowTime                 int64         `json:"now_time"`
	OnlineCount             int           `json:"online_count"`
	NeedLoginSubtitle       bool          `json:"need_login_subtitle"`
	ViewPoints              []any         `json:"view_points"` // 根据实际情况可能需要具体类型
	PreviewToast            string        `json:"preview_toast"`
	Options                 Options       `json:"options"`
	GuideAttention          []any         `json:"guide_attention"` // 根据实际情况可能需要具体类型
	JumpCard                []any         `json:"jump_card"`       // 根据实际情况可能需要具体类型
	OperationCard           []any         `json:"operation_card"`  // 根据实际情况可能需要具体类型
	OnlineSwitch            OnlineSwitch  `json:"online_switch"`
	Fawkes                  Fawkes        `json:"fawkes"`
	ShowSwitch              ShowSwitch    `json:"show_switch"`
	BgmInfo                 interface{}   `json:"bgm_info"` // 可能为 null，使用 interface{}
	ToastBlock              bool          `json:"toast_block"`
	IsUpowerExclusive       bool          `json:"is_upower_exclusive"`
	IsUpowerPlay            bool          `json:"is_upower_play"`
	IsUgcPayPreview         bool          `json:"is_ugc_pay_preview"`
	ElecHighLevel           ElecHighLevel `json:"elec_high_level"`
	DisableShowUpInfo       bool          `json:"disable_show_up_info"`
	IsUpowerExclusiveWithQa bool          `json:"is_upower_exclusive_with_qa"`
}

// IpInfo IP 相关信息。
type IpInfo struct {
	Ip       string `json:"ip"`
	ZoneIp   string `json:"zone_ip"`
	ZoneId   int    `json:"zone_id"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
}

// LevelInfo 等级信息。
type LevelInfo struct {
	CurrentLevel int   `json:"current_level"`
	CurrentMin   int   `json:"current_min"`
	CurrentExp   int   `json:"current_exp"`
	NextExp      int   `json:"next_exp"`
	LevelUp      int64 `json:"level_up"`
}

// VipInfo VIP 相关信息。
type VipInfo struct {
	Type               int        `json:"type"`
	Status             int        `json:"status"`
	DueDate            int64      `json:"due_date"`
	VipPayType         int        `json:"vip_pay_type"`
	ThemeType          int        `json:"theme_type"`
	Label              VipLabel   `json:"label"`
	AvatarSubscript    int        `json:"avatar_subscript"`
	NicknameColor      string     `json:"nickname_color"`
	Role               int        `json:"role"`
	AvatarSubscriptUrl string     `json:"avatar_subscript_url"`
	TvVipStatus        int        `json:"tv_vip_status"`
	TvVipPayType       int        `json:"tv_vip_pay_type"`
	TvDueDate          int64      `json:"tv_due_date"`
	AvatarIcon         AvatarIcon `json:"avatar_icon"`
}

// VipLabel VIP 标签信息。
type VipLabel struct {
	Path                  string `json:"path"`
	Text                  string `json:"text"`
	LabelTheme            string `json:"label_theme"`
	TextColor             string `json:"text_color"`
	BgStyle               int    `json:"bg_style"`
	BgColor               string `json:"bg_color"`
	BorderColor           string `json:"border_color"`
	UseImgLabel           bool   `json:"use_img_label"`
	ImgLabelUriHans       string `json:"img_label_uri_hans"`
	ImgLabelUriHant       string `json:"img_label_uri_hant"`
	ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"`
	ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"`
}

// AvatarIcon 头像图标信息。
type AvatarIcon struct {
	IconResource interface{} `json:"icon_resource"` // 结构未知，暂用 interface{}
}

// Options 播放器选项。
type Options struct {
	Is360      bool `json:"is_360"`
	WithoutVip bool `json:"without_vip"`
}

// OnlineSwitch 在线开关配置。
type OnlineSwitch struct {
	EnableGrayDashPlayback string `json:"enable_gray_dash_playback"`
	NewBroadcast           string `json:"new_broadcast"`
	RealtimeDm             string `json:"realtime_dm"`
	SubtitleSubmitSwitch   string `json:"subtitle_submit_switch"`
}

// Fawkes 配置信息。
type Fawkes struct {
	ConfigVersion int `json:"config_version"`
	FfVersion     int `json:"ff_version"`
}

// ShowSwitch 显示开关配置。
type ShowSwitch struct {
	LongProgress bool `json:"long_progress"`
}

// ElecHighLevel 高能等级信息。
type ElecHighLevel struct {
	PrivilegeType int    `json:"privilege_type"`
	Title         string `json:"title"`
	SubTitle      string `json:"sub_title"`
	ShowButton    bool   `json:"show_button"`
	ButtonText    string `json:"button_text"`
	JumpUrl       string `json:"jump_url"`
	Intro         string `json:"intro"`
	New           bool   `json:"new"`
	QuestionText  string `json:"question_text"`
	QaTitle       string `json:"qa_title"`
}

// -----------------------------------------

// VideoProgressResponse 定义 Bilibili 视频进度 API 的响应结构体。
type VideoProgressResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	TTL     int               `json:"ttl"`
	Data    VideoProgressData `json:"data,omitempty"`
}

// GetVideoProgress 调用 Bilibili API 获取指定视频的观看进度，并返回应用层 DTO。
// 实现 application.BilibiliClient 接口的一部分。
// aidStr (视频稿件 avid) 和 bvidStr (视频稿件 bvid) 必须提供一个。
// cidStr (视频分P的 ID) 必须提供。
func (c *Client) GetVideoProgress(ctx context.Context, aidStr, bvidStr, cidStr string) (*application.VideoProgressDTO, error) {
	const path = "/x/player/wbi/v2"

	var finalAidStr string

	if aidStr == "" && bvidStr == "" {
		return nil, fmt.Errorf("GetVideoProgress requires either aid or bvid")
	}
	if cidStr == "" {
		return nil, fmt.Errorf("GetVideoProgress requires cid")
	}

	if aidStr != "" {
		finalAidStr = aidStr
	} else {
		// aid 为空，bvid 提供了，需要通过 bvid 获取 aid
		// 注意：这里会发生一次额外的 API 调用来获取视频视图信息
		videoView, err := c.GetVideoView(ctx, "", bvidStr) // GetVideoView 接受 ctx
		if err != nil {
			return nil, fmt.Errorf("failed to get video info for bvid %s to resolve aid: %w", bvidStr, err)
		}
		if videoView == nil || videoView.Aid == 0 {
			return nil, fmt.Errorf("could not resolve aid from bvid %s: no view data or aid is zero", bvidStr)
		}
		finalAidStr = strconv.FormatInt(videoView.Aid, 10)
	}

	params := url.Values{}
	params.Set("aid", finalAidStr)
	params.Set("cid", cidStr)

	var resp VideoProgressResponse
	err := c.Get(ctx, path, params, &resp) // Pass ctx to c.Get
	if err != nil {
		if apiErrResp, ok := err.(interface{ GetResponse() *VideoProgressResponse }); ok && apiErrResp.GetResponse() != nil {
			return nil, fmt.Errorf("bilibili api error: code=%d, message=%s, underlying error: %w",
				apiErrResp.GetResponse().Code, apiErrResp.GetResponse().Message, err)
		}
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("bilibili api error: code=%d, message=%s", resp.Code, resp.Message)
	}

	if resp.Data.LastPlayCid == 0 && resp.Data.LastPlayTime == 0 {
		return nil, nil // No progress data, not an error
	}

	dto := &application.VideoProgressDTO{
		AID:          resp.Data.Aid,
		BVID:         resp.Data.Bvid,
		LastPlayTime: resp.Data.LastPlayTime,
		LastPlayCid:  resp.Data.LastPlayCid,
	}

	return dto, nil
}
