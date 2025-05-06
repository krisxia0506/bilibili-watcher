package bilibili

import (
	"context"
	"fmt"
	"log"
	"net/url"

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
// aid: 视频的 AV 号 (不带 'av' 前缀)
// cid: 视频的分 P ID (当前页面的 CID)
func (c *Client) GetVideoProgress(ctx context.Context, aid, cid string) (*application.VideoProgressDTO, error) {
	log.Printf("Getting video progress for AID: %s, CID: %s", aid, cid)
	const path = "/x/player/wbi/v2" // TODO: Handle WBI signing if necessary

	params := url.Values{}
	params.Set("aid", aid)
	params.Set("cid", cid)

	var resp VideoProgressResponse
	// 注意：这里的 c.Get 方法还没有 context 参数，暂时忽略
	err := c.Get(path, params, &resp)
	if err != nil {
		// 底层 Get 方法已处理 HTTP 和解码错误，这里直接返回
		// 尝试从错误中提取 Bilibili API 的特定错误信息（如果底层 Get 返回了带有 resp 的错误）
		if apiErrResp, ok := err.(interface{ GetResponse() *VideoProgressResponse }); ok && apiErrResp.GetResponse() != nil {
			return nil, fmt.Errorf("bilibili api error: code=%d, message=%s, underlying error: %w",
				apiErrResp.GetResponse().Code, apiErrResp.GetResponse().Message, err)
		}
		return nil, err // 返回通用错误
	}

	// 检查 Bilibili API 返回的业务状态码
	if resp.Code != 0 {
		// 返回包含错误信息的响应体，让调用者能看到具体错误
		return nil, fmt.Errorf("bilibili api error: code=%d, message=%s", resp.Code, resp.Message)
	}

	// 映射到应用层 DTO
	// 检查 Data 字段是否存在且有效
	if resp.Data.LastPlayCid == 0 && resp.Data.LastPlayTime == 0 {
		// 如果关键进度信息都为 0 或无效，可能表示没有进度记录
		// 这里返回 nil DTO 和 nil error，由应用层决定如何处理
		return nil, nil
	}

	dto := &application.VideoProgressDTO{
		AID:          resp.Data.Aid,
		BVID:         resp.Data.Bvid,
		LastPlayTime: resp.Data.LastPlayTime,
		LastPlayCid:  resp.Data.LastPlayCid,
	}

	return dto, nil
}
