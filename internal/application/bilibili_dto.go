package application

// VideoProgressDTO 应用层关心的视频进度数据
type VideoProgressDTO struct {
	AID          int64  `json:"aid"`
	BVID         string `json:"bvid"`
	LastPlayTime int64  `json:"last_play_time"` // 观看进度，单位毫秒
	LastPlayCid  int64  `json:"last_play_cid"`  // 上次播放的视频分 P ID
	// 可以根据需要从 bilibili.VideoProgressData 添加更多字段
}

// VideoViewPageDTO 应用层关心的分P信息
type VideoViewPageDTO struct {
	Cid      int64  `json:"cid"`
	Part     string `json:"part"`     // 分P标题
	Duration int64  `json:"duration"` // 分P时长（秒）
	Page     int    `json:"page"`     // 分P序号
}

// VideoViewDTO 应用层关心的视频视图数据
type VideoViewDTO struct {
	Bvid      string             `json:"bvid"`
	Aid       int64              `json:"aid"`
	Title     string             `json:"title"`
	Desc      string             `json:"desc"`
	Pubdate   int64              `json:"pubdate"`  // 发布时间戳
	Duration  int64              `json:"duration"` // 总时长(秒)
	OwnerName string             `json:"owner_name"`
	Pages     []VideoViewPageDTO `json:"pages"` // 新增：分P信息列表
	// 可以根据需要从 bilibili.VideoViewData 添加更多字段
}
