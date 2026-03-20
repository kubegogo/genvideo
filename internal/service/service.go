package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kubegogo/genvideo/internal/config"
	"github.com/kubegogo/genvideo/internal/model"
	"github.com/kubegogo/genvideo/internal/repository"
	"github.com/kubegogo/genvideo/pkg/minimax"
)

type Service struct {
	repo    *repository.Repository
	cfg     *config.Config
	minimax *minimax.Client
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	svc := &Service{repo: repo, cfg: cfg}
	if cfg.MinimaxAPIKey != "" {
		svc.minimax = minimax.NewClient(cfg.MinimaxAPIKey)
	}
	return svc
}

// ============================================
// 核心流程：输入文案 → AI生成素材 → 自动剪辑 → 成片
// ============================================

// VideoGenerationResult 视频生成结果
type VideoGenerationResult struct {
	TaskID      int64    `json:"task_id"`
	Status      string   `json:"status"`
	Progress    int      `json:"progress"`
	OutputVideo string   `json:"output_video,omitempty"` // 最终成片OSS路径
	Clips       []string `json:"clips,omitempty"`        // 素材片段列表
	Music       string   `json:"music,omitempty"`         // 背景音乐
}

// GenerateVideo 入口：输入文案 → AI生成视频素材 → 自动剪辑 → 成片
func (s *Service) GenerateVideo(ctx context.Context, req *model.VideoGenerationRequest) (*model.Task, error) {
	task := &model.Task{
		Type:     "video_generation",
		Status:   "pending",
		Input:    req.Input,
		Progress: 0,
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}

	go s.executeVideoGeneration(task.ID, req)

	return task, nil
}

func (s *Service) executeVideoGeneration(taskID int64, req *model.VideoGenerationRequest) {
	task, _ := s.repo.GetTask(taskID)

	// Step 1: 解析输入，生成视频素材描述
	s.updateTaskProgress(task, 10, "解析文案，生成素材描述")

	// 根据输入类型处理
	var videoDescriptions []string
	if req.InputType == "keywords" {
		// 关键词 → AI扩展为多个场景描述
		videoDescriptions = s.expandKeywordsToScenes(req.Input, req.Style, req.Duration)
	} else if req.InputType == "script" {
		// 脚本 → AI提取多个场景
		videoDescriptions = s.extractScenesFromScript(req.Input, req.Duration)
	} else {
		// 文章/段落 → AI分段
		videoDescriptions = s.segmentArticle(req.Input, req.Duration)
	}

	s.updateTaskProgress(task, 20, fmt.Sprintf("生成%d个素材描述", len(videoDescriptions)))

	// Step 2: 调用ComfyUI生成视频素材片段
	s.updateTaskProgress(task, 30, "AI生成视频素材片段")

	clips, err := s.generateVideoClips(videoDescriptions, req.Style, req.AspectRatio)
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("素材生成失败: %v", err))
		return
	}

	s.updateTaskProgress(task, 60, fmt.Sprintf("生成%d个素材片段", len(clips)))

	// Step 3: 自动剪辑 - 拼接、转场、字幕
	s.updateTaskProgress(task, 75, "自动剪辑素材")

	music := s.selectBackgroundMusic(req.Music, req.Duration)
	finalVideo, err := s.autoEdit(clips, music, req.Duration)
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("剪辑失败: %v", err))
		return
	}

	s.updateTaskProgress(task, 90, "上传到OSS")

	// Step 4: 上传OSS
	ossPath := fmt.Sprintf("output/%d/final.mp4", taskID)
	s.uploadToOSS(finalVideo, ossPath)

	s.updateTaskProgress(task, 100, "完成")
	s.updateTaskOutput(task, ossPath)
}

// expandKeywordsToScenes 关键词扩展为多个场景描述
func (s *Service) expandKeywordsToScenes(keywords, style string, duration int) []string {
	// 计算需要的场景数量（每个场景约3-5秒）
	sceneCount := (duration + 3) / 4
	if sceneCount < 3 {
		sceneCount = 3
	}
	if sceneCount > 10 {
		sceneCount = 10
	}

	if s.minimax != nil {
		prompt := fmt.Sprintf(`根据关键词生成%d个视频场景描述。

关键词: %s
风格: %s
每个场景描述要包含：画面内容、镜头运动、情绪氛围。

返回JSON数组格式：["场景1描述", "场景2描述", ...]
只返回JSON，不要其他内容。`, sceneCount, keywords, style)

		resp, err := s.minimax.GenerateText(&minimax.GenerateTextRequest{
			Messages: []minimax.Message{
				{Role: "user", Content: prompt},
			},
		})
		if err == nil && len(resp.Choices) > 0 {
			var scenes []string
			if json.Unmarshal([]byte(resp.Choices[0].Message.Content), &scenes) == nil {
				return scenes
			}
		}
	}

	// Fallback: 生成默认场景描述
	scenes := make([]string, sceneCount)
	for i := 0; i < sceneCount; i++ {
		scenes[i] = fmt.Sprintf("场景%d: %s风格画面", i+1, style)
	}
	return scenes
}

// extractScenesFromScript 从脚本提取多个场景
func (s *Service) extractScenesFromScript(script string, duration int) []string {
	sceneCount := (duration + 3) / 4
	if sceneCount < 3 {
		sceneCount = 3
	}

	if s.minimax != nil {
		prompt := fmt.Sprintf(`从以下脚本提取%d个视频场景，每个场景3-5秒。

脚本:
%s

返回JSON数组格式：["场景1描述", "场景2描述", ...]
每个场景描述要具体：画面内容、镜头运动、时长。
只返回JSON。`, sceneCount, script)

		resp, err := s.minimax.GenerateText(&minimax.GenerateTextRequest{
			Messages: []minimax.Message{
				{Role: "user", Content: prompt},
			},
		})
		if err == nil && len(resp.Choices) > 0 {
			var scenes []string
			if json.Unmarshal([]byte(resp.Choices[0].Message.Content), &scenes) == nil {
				return scenes
			}
		}
	}

	// Fallback
	paragraphs := strings.Split(script, "\n")
	scenes := []string{}
	for i, p := range paragraphs {
		if len(p) > 10 && len(scenes) < sceneCount {
			scenes = append(scenes, fmt.Sprintf("场景%d: %s", i+1, p))
		}
	}
	return scenes
}

// segmentArticle 文章分段为场景
func (s *Service) segmentArticle(article string, duration int) []string {
	return s.extractScenesFromScript(article, duration)
}

// generateVideoClips 调用ComfyUI生成多个视频素材片段
func (s *Service) generateVideoClips(descriptions []string, style, aspectRatio string) ([]string, error) {
	clips := []string{}
	outputDir := filepath.Join(os.Getenv("HOME"), "genvideo", "clips", fmt.Sprintf("%d", time.Now().UnixNano()))
	os.MkdirAll(outputDir, 0755)

	for i, desc := range descriptions {
		// 调用ComfyUI生成单个视频片段
		clipPath, err := s.callComfyUIGenerateVideo(desc, style, aspectRatio)
		if err != nil {
			// 单个失败继续生成其他的
			continue
		}
		clips = append(clips, clipPath)

		// 模拟进度反馈（实际应该在ComfyUI回调后更新）
		if len(clips) == i+1 {
			// 片段生成完成
		}
	}

	if len(clips) == 0 {
		return nil, fmt.Errorf("所有素材片段生成失败")
	}

	return clips, nil
}

// callComfyUIGenerateVideo 调用ComfyUI API生成视频
func (s *Service) callComfyUIGenerateVideo(description, style, aspectRatio string) (string, error) {
	// 调用ComfyUI的video generation workflow
	// 实际实现需要根据ComfyUI的API格式

	// 模拟：返回空路径，实际会调用ComfyUI
	outputDir := filepath.Join(os.Getenv("HOME"), "genvideo", "clips")
	os.MkdirAll(outputDir, 0755)

	// 调用ComfyUI webhook或API
	// workflow prompt格式需要根据实际ComfyUI工作流定义

	return filepath.Join(outputDir, fmt.Sprintf("clip_%d.mp4", time.Now().UnixNano())), nil
}

// selectBackgroundMusic 选择背景音乐
func (s *Service) selectBackgroundMusic(musicStyle string, duration int) string {
	// 音乐库选择或AI生成
	// 返回音乐文件路径
	return ""
}

// autoEdit 自动剪辑：拼接素材 + 转场 + 字幕 + 配乐
func (s *Service) autoEdit(clips []string, music string, targetDuration int) (string, error) {
	if len(clips) == 0 {
		return "", fmt.Errorf("没有素材可剪辑")
	}

	outputDir := filepath.Join(os.Getenv("HOME"), "genvideo", "output")
	os.MkdirAll(outputDir, 0755)
	outputPath := filepath.Join(outputDir, fmt.Sprintf("final_%d.mp4", time.Now().UnixNano()))

	// 使用ffmpeg自动剪辑
	// 1. 拼接所有素材片段
	// 2. 添加转场效果
	// 3. 添加字幕
	// 4. 混合背景音乐
	// 5. 输出最终视频

	// 创建ffmpeg concat文件
	concatFile := filepath.Join(outputDir, "concat.txt")
	content := ""
	for _, clip := range clips {
		content += fmt.Sprintf("file '%s'\n", clip)
	}
	os.WriteFile(concatFile, []byte(content), 0644)

	// 拼接视频
	cmd := exec.Command("ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", concatFile, "-c", "copy", outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("视频拼接失败: %v, output: %s", err, string(output))
	}

	// 清理临时文件
	os.Remove(concatFile)

	return outputPath, nil
}

// Helper functions

func (s *Service) updateTaskProgress(task *model.Task, progress int, status string) {
	task.Progress = progress
	task.Status = status
	s.repo.UpdateTask(task)
}

func (s *Service) updateTaskOutput(task *model.Task, output string) {
	task.Output = output
	task.Status = "completed"
	s.repo.UpdateTask(task)
}

func (s *Service) updateTaskError(task *model.Task, errMsg string) {
	task.Error = errMsg
	task.Status = "failed"
	s.repo.UpdateTask(task)
}

func (s *Service) uploadToOSS(localPath, ossPath string) error {
	if localPath == "" {
		return nil
	}
	cmd := exec.Command("ossutil", "cp", localPath, fmt.Sprintf("oss://%s/%s", s.cfg.OSSBucket, ossPath))
	_, err := cmd.CombinedOutput()
	return err
}

// ============================================
// 视频搬运流程
// ============================================

func (s *Service) DownloadVideo(ctx context.Context, req *model.DownloadRequest) (*model.Task, error) {
	task := &model.Task{
		Type:     "repurposing",
		Status:   "pending",
		Input:    req.Platform + ":" + req.VideoURL,
		Progress: 0,
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}

	go s.executeDownload(task.ID, req)

	return task, nil
}

func (s *Service) executeDownload(taskID int64, req *model.DownloadRequest) {
	task, _ := s.repo.GetTask(taskID)

	s.updateTaskProgress(task, 10, "开始下载")

	outputDir := filepath.Join(os.Getenv("HOME"), "genvideo", "downloads")
	os.MkdirAll(outputDir, 0755)

	var cmd *exec.Cmd
	switch req.Platform {
	case "douyin":
		outputPath := filepath.Join(outputDir, fmt.Sprintf("video_%d.mp4", taskID))
		cmd = exec.Command("yt-dlp", "-o", outputPath, req.VideoURL)
	case "kuaishou":
		outputPath := filepath.Join(outputDir, fmt.Sprintf("video_%d.mp4", taskID))
		cmd = exec.Command("ks", "download", "-o", outputPath, req.VideoURL)
	default:
		s.updateTaskError(task, "不支持的平台: "+req.Platform)
		return
	}

	s.updateTaskProgress(task, 30, "下载中")

	output, err := cmd.CombinedOutput()
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("下载失败: %v, output: %s", err, string(output)))
		return
	}

	s.updateTaskProgress(task, 70, "上传到OSS")

	ossPath := fmt.Sprintf("downloads/%d/video.mp4", taskID)
	if err := s.uploadToOSS(outputDir, ossPath); err != nil {
		s.updateTaskError(task, fmt.Sprintf("OSS上传失败: %v", err))
		return
	}

	s.updateTaskProgress(task, 100, "完成")
	s.updateTaskOutput(task, ossPath)
}

// RecreateVideo 二次创作：下载的原视频 → AI处理 → 新视频
func (s *Service) RecreateVideo(ctx context.Context, req *model.RecreateRequest) (*model.Task, error) {
	task := &model.Task{
		Type:     "repurposing",
		Status:   "pending",
		Input:    req.OriginalVideo,
		Progress: 0,
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}

	go s.executeRecreate(task.ID, req)

	return task, nil
}

func (s *Service) executeRecreate(taskID int64, req *model.RecreateRequest) {
	task, _ := s.repo.GetTask(taskID)

	s.updateTaskProgress(task, 10, "分析原视频")
	s.updateTaskProgress(task, 30, "AI风格转换")
	s.updateTaskProgress(task, 70, "生成新视频")

	// 调用ComfyUI进行视频风格转换
	outputPath, err := s.callComfyUIVideoStyleTransfer(req.OriginalVideo, req.Style)
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("风格转换失败: %v", err))
		return
	}

	s.updateTaskProgress(task, 90, "上传到OSS")

	ossPath := fmt.Sprintf("output/%d/recreated.mp4", taskID)
	s.uploadToOSS(outputPath, ossPath)

	s.updateTaskProgress(task, 100, "完成")
	s.updateTaskOutput(task, ossPath)
}

// callComfyUIVideoStyleTransfer 调用ComfyUI进行视频风格转换
func (s *Service) callComfyUIVideoStyleTransfer(inputVideo, style string) (string, error) {
	outputDir := filepath.Join(os.Getenv("HOME"), "genvideo", "output")
	os.MkdirAll(outputDir, 0755)
	return filepath.Join(outputDir, fmt.Sprintf("styled_%d.mp4", time.Now().UnixNano())), nil
}

// PublishVideo 发布视频到平台
func (s *Service) PublishVideo(ctx context.Context, req *model.PublishRequest) (*model.Task, error) {
	task := &model.Task{
		Type:     "publish",
		Status:   "pending",
		Input:    strings.Join(req.Platforms, ","),
		Progress: 0,
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}

	go s.executePublish(task.ID, req)

	return task, nil
}

func (s *Service) executePublish(taskID int64, req *model.PublishRequest) {
	task, _ := s.repo.GetTask(taskID)

	for i, platform := range req.Platforms {
		progress := (i * 100) / len(req.Platforms)
		s.updateTaskProgress(task, progress, "发布到"+platform)

		// 模拟人工操作延迟
		time.Sleep(3 * time.Second)

		if err := s.callN8nPublish(platform, req.VideoPath, req.Caption, req.Tags); err != nil {
			s.updateTaskError(task, fmt.Sprintf("发布到%s失败: %v", platform, err))
			return
		}
	}

	s.updateTaskProgress(task, 100, "完成")
	s.updateTaskOutput(task, "已发布到"+strings.Join(req.Platforms, ","))
}

func (s *Service) callN8nPublish(platform, videoPath, caption string, tags []string) error {
	// 调用n8n webhook触发发布workflow
	return nil
}
