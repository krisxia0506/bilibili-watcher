package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	baseURL = "https://api.bilibili.com/x/player/wbi/v2"
	// TODO: Dynamically load SESSDATA instead of hardcoding
	sessData = "SESSDATA=a486b214%2C1761995605%2Ce8064%2A51CjAjbbi4oaDoEXCN7yjPThSija81Url7d8duiZqLF-IVReywvw-pC5bbiw4O0IFLqZwSVlZUZ3lDMktiU2NSYTMwYTVqRnoxMmoxRG5jVEVQaVlXb0tzdm1zX1k3VlpKcGQ2aFk2SzMyMVJ5SlNRc1g4MmtNckdkTFp0RjZJM0V0WnpyMnRUMUNBIIEC"
)

// VideoProgressData 定义 Bilibili 视频进度 API 响应中 data 字段的结构体
type VideoProgressData struct {
	Aid                     int64         `json:"aid"`
	Bvid                    string        `json:"bvid"`
	AllowBp                 bool          `json:"allow_bp"`
	NoShare                 bool          `json:"no_share"`
	Cid                     int64         `json:"cid"`
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
	LastPlayTime            int           `json:"last_play_time"` // 观看进度，单位毫秒
	LastPlayCid             int64         `json:"last_play_cid"`
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

// IpInfo IP 相关信息
type IpInfo struct {
	Ip       string `json:"ip"`
	ZoneIp   string `json:"zone_ip"`
	ZoneId   int    `json:"zone_id"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
}

// LevelInfo 等级信息
type LevelInfo struct {
	CurrentLevel int   `json:"current_level"`
	CurrentMin   int   `json:"current_min"`
	CurrentExp   int   `json:"current_exp"`
	NextExp      int   `json:"next_exp"`
	LevelUp      int64 `json:"level_up"`
}

// VipInfo VIP 相关信息
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

// VipLabel VIP 标签信息
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

// AvatarIcon 头像图标信息
type AvatarIcon struct {
	IconResource interface{} `json:"icon_resource"` // 结构未知，暂用 interface{}
}

// Options 播放器选项
type Options struct {
	Is360      bool `json:"is_360"`
	WithoutVip bool `json:"without_vip"`
}

// OnlineSwitch 在线开关配置
type OnlineSwitch struct {
	EnableGrayDashPlayback string `json:"enable_gray_dash_playback"`
	NewBroadcast           string `json:"new_broadcast"`
	RealtimeDm             string `json:"realtime_dm"`
	SubtitleSubmitSwitch   string `json:"subtitle_submit_switch"`
}

// Fawkes 配置信息
type Fawkes struct {
	ConfigVersion int `json:"config_version"`
	FfVersion     int `json:"ff_version"`
}

// ShowSwitch 显示开关配置
type ShowSwitch struct {
	LongProgress bool `json:"long_progress"`
}

// ElecHighLevel 高能等级信息
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

// VideoProgressResponse 定义 Bilibili 视频进度 API 的响应结构体
// 注意：这里需要根据实际 API 返回的 JSON 结构进行定义
type VideoProgressResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	TTL     int               `json:"ttl"`
	Data    VideoProgressData `json:"data"` // 更新为具体的 Data 结构体
}

// Client Bilibili API 客户端结构体
type Client struct {
	httpClient *http.Client
	baseURL    string
	sessData   string
}

// NewClient 创建一个新的 Bilibili API 客户端实例
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
		sessData:   sessData,
	}
}

// GetVideoProgress 调用 Bilibili API 获取指定视频的观看进度
// aid: 视频的 AV 号 (不带 'av' 前缀)
// cid: 视频的分 P ID
func (c *Client) GetVideoProgress(aid, cid string) (*VideoProgressResponse, error) {
	// 构建请求 URL
	requestURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	q := requestURL.Query()
	q.Set("aid", aid)
	q.Set("cid", cid)
	requestURL.RawQuery = q.Encode()

	// 创建 HTTP GET 请求
	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cookie", c.sessData) // 使用硬编码的 SESSDATA

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析 JSON 响应
	var progressResponse VideoProgressResponse
	if err := json.Unmarshal(body, &progressResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json response: %w, body: %s", err, string(body))
	}

	// 检查 Bilibili API 返回的业务状态码
	if progressResponse.Code != 0 {
		return &progressResponse, fmt.Errorf("bilibili api error: code=%d, message=%s", progressResponse.Code, progressResponse.Message)
	}

	return &progressResponse, nil
}
