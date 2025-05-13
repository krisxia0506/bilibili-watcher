package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/krisxia0506/bilibili-watcher/internal/application"
	"github.com/krisxia0506/bilibili-watcher/internal/interfaces/api/rest/dto"
	"github.com/krisxia0506/bilibili-watcher/pkg/response"
)

// VideoAnalyticsHandler 处理视频分析相关的 API 请求。
type VideoAnalyticsHandler struct {
	appService application.VideoAnalyticsService
}

// NewVideoAnalyticsHandler 创建 VideoAnalyticsHandler 实例。
func NewVideoAnalyticsHandler(appService application.VideoAnalyticsService) *VideoAnalyticsHandler {
	return &VideoAnalyticsHandler{appService: appService}
}

// RegisterRoutes 在 Gin 路由组上注册视频分析相关的路由。
func (h *VideoAnalyticsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/video/watch-segments", h.GetWatchedSegments)
}

// GetWatchedSegments 处理获取观看分段的请求。
// @Summary 获取指定时间范围和间隔的视频观看时长分段
// @Description 根据提供的AID或BVID、开始/结束时间和时间间隔，计算每个时间段内的观看时长。
// @Tags VideoAnalytics
// @Accept json
// @Produce json
// @Param request body dto.GetWatchedSegmentsRequest true "查询参数"
// @Success 200 {object} response.APIResponse{data=dto.GetWatchedSegmentsResponse} "成功响应"
// @Failure 400 {object} response.APIResponse "请求参数错误"
// @Failure 500 {object} response.APIResponse "服务器内部错误"
// @Router /api/v1/video/watch-segments [post]
func (h *VideoAnalyticsHandler) GetWatchedSegments(c *gin.Context) {
	var req dto.GetWatchedSegmentsRequest
	// 绑定并验证请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	// 验证 AID 或 BVID 至少提供一个
	if req.AID == "" && req.BVID == "" {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "Either aid or bvid must be provided")
		return
	}

	// 解析时间字符串
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, fmt.Sprintf("Invalid start_time format: %v", err))
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, fmt.Sprintf("Invalid end_time format: %v", err))
		return
	}

	// 解析时间间隔字符串
	var interval time.Duration
	intervalStr := req.Interval

	if strings.HasSuffix(intervalStr, "d") {
		// 如果单位是 'd' (天)
		daysStr := strings.TrimSuffix(intervalStr, "d")
		days, parseErr := strconv.Atoi(daysStr)
		if parseErr != nil {
			response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, fmt.Sprintf("Invalid interval format (days): %v", parseErr))
			return
		}
		// 将天转换为小时，再用 time.ParseDuration 解析
		intervalStr = fmt.Sprintf("%dh", days*24)
	}

	interval, err = time.ParseDuration(intervalStr)
	if err != nil {
		// 如果转换后或原始格式仍然无效 (例如，gin binding oneof 漏掉了某些校验)
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, fmt.Sprintf("Invalid interval format: %v for input '%s' (parsed as '%s')", err, req.Interval, intervalStr))
		return
	}

	// 调用应用服务
	analyticsResult, err := h.appService.GetWatchedSegments(c.Request.Context(), req.AID, req.BVID, startTime, endTime, interval)
	if err != nil {
		// 根据应用层返回的错误类型决定 HTTP 状态码和业务码
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, fmt.Sprintf("Failed to calculate watched segments: %v", err))
		return
	}

	// 映射结果到响应 DTO
	respData := dto.GetWatchedSegmentsResponse{
		Segments:                make([]dto.WatchedSegment, 0, len(analyticsResult.Segments)),
		TotalWatchedDurationSec: int64(analyticsResult.TotalWatchedDuration.Seconds()),
	}
	for _, seg := range analyticsResult.Segments {
		respData.Segments = append(respData.Segments, dto.WatchedSegment{
			SegmentStartTime:   seg.SegmentStartTime,
			SegmentEndTime:     seg.SegmentEndTime,
			WatchedDurationSec: int64(seg.WatchedDuration.Seconds()), // 转换为秒
		})
	}

	response.Success(c, respData)
}
