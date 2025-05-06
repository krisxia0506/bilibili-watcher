package model

// VideoPage 代表视频的一个分P及其信息。
// 这是领域层进行观看时长计算所需的基础结构。
type VideoPage struct {
	Cid      int64  // 分P的唯一标识符
	Duration int64  // 分P的持续时间（单位：秒）
	Part     string // 分P的标题
	Page     int    // 分P的序号 (从1开始)
}
