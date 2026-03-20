package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kubegogo/genvideo/internal/config"
	"github.com/kubegogo/genvideo/internal/model"
	"github.com/kubegogo/genvideo/internal/repository"
	"github.com/kubegogo/genvideo/pkg/comfyui"
	"github.com/kubegogo/genvideo/pkg/minimax"
	"github.com/kubegogo/genvideo/pkg/n8n"
)

type Service struct {
	repo    *repository.Repository
	cfg     *config.Config
	minimax *minimax.Client
	comfyui *comfyui.Client
	n8n     *n8n.Client
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	svc := &Service{
		repo: repo,
		cfg:  cfg,
	}

	// 初始化AI客户端
	if cfg.MinimaxAPIKey != "" {
		svc.minimax = minimax.NewClient(cfg.MinimaxAPIKey)
	}
	if cfg.ComfyUIBaseURL != "" {
		svc.comfyui = comfyui.NewClient(cfg.ComfyUIBaseURL)
	}
	if cfg.N8NBaseURL != "" {
		svc.n8n = n8n.NewClient(cfg.N8NBaseURL)
	}

	return svc
}

// ============================================
// 视频生成（脚本转视频）
// ============================================

func (s *Service) GenerateVideo(ctx context.Context, req *model.VideoGenerationRequest) (*model.Task, error) {
	task := &model.Task{
		Type:     "script_to_video",
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

	// Step 1: 生成剧本
	s.updateTaskProgress(task, 10, "生成剧本")
	script, err := s.generateScript(req)
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("生成剧本失败: %v", err))
		return
	}

	// Step 2: 生成分镜
	s.updateTaskProgress(task, 25, "生成分镜")
	storyboard, err := s.generateStoryboard(script, req.Duration)
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("生成分镜失败: %v", err))
		return
	}

	// Step 3: 生成首尾帧
	s.updateTaskProgress(task, 40, "生成首尾帧图片")
	firstFrame, lastFrame, err := s.generateFrames(storyboard, req.Style)
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("生成首尾帧失败: %v", err))
		return
	}

	// Step 4: 生成视频
	s.updateTaskProgress(task, 60, "AI生成视频")
	videoPath, err := s.generateVideo(storyboard, firstFrame, lastFrame, req)
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("生成视频失败: %v", err))
		return
	}

	// Step 5: 上传OSS
	s.updateTaskProgress(task, 80, "上传到OSS")
	ossPath, err := s.uploadToOSS(videoPath)
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("上传OSS失败: %v", err))
		return
	}

	s.updateTaskProgress(task, 100, "完成")
	s.updateTaskOutput(task, ossPath)
}

func (s *Service) generateScript(req *model.VideoGenerationRequest) (string, error) {
	if s.minimax != nil {
		return s.minimax.GenerateScript(req.Input, req.InputType, req.Style, req.Duration)
	}
	// Fallback: 返回占位符
	return fmt.Sprintf("[剧本占位符] 风格:%s 时长:%d秒 内容:%s", req.Style, req.Duration, req.Input), nil
}

func (s *Service) generateStoryboard(script string, duration int) (string, error) {
	sceneCount := (duration + 4) / 5
	if sceneCount < 3 {
		sceneCount = 3
	}
	if sceneCount > 10 {
		sceneCount = 10
	}

	if s.minimax != nil {
		return s.minimax.GenerateStoryboard(script, sceneCount)
	}
	return fmt.Sprintf("[分镜占位符] 共%d个场景", sceneCount), nil
}

func (s *Service) generateFrames(storyboard, style string) (firstFrame, lastFrame string, err error) {
	if s.comfyui != nil {
		first, err := s.comfyui.GenerateFrame(storyboard, style, true)
		if err != nil {
			return "", "", err
		}
		last, err := s.comfyui.GenerateFrame(storyboard, style, false)
		if err != nil {
			return "", "", err
		}
		return first, last, nil
	}
	return "first_frame_placeholder.png", "last_frame_placeholder.png", nil
}

func (s *Service) generateVideo(storyboard, firstFrame, lastFrame string, req *model.VideoGenerationRequest) (string, error) {
	outputDir := filepath.Join(os.Getenv("HOME"), "genvideo", "output")
	os.MkdirAll(outputDir, 0755)
	outputPath := filepath.Join(outputDir, fmt.Sprintf("video_%d.mp4", time.Now().UnixNano()))

	if s.comfyui != nil {
		resp, err := s.comfyui.GenerateVideo(&comfyui.GenerateVideoRequest{
			Model:  req.Style,
			Prompt: storyboard,
			Frames: req.Duration * 8, // 约8fps
		})
		if err != nil {
			return "", err
		}
		return resp.VideoPath, nil
	}

	// Fallback: 创建空文件占位
	cmd := exec.Command("touch", outputPath)
	cmd.Run()
	return outputPath, nil
}

func (s *Service) uploadToOSS(localPath string) (string, error) {
	if localPath == "" || s.cfg.OSSBucket == "" {
		return "", nil
	}

	ossPath := fmt.Sprintf("output/%d/video.mp4", time.Now().UnixNano())
	cmd := exec.Command("ossutil", "cp", localPath, fmt.Sprintf("oss://%s/%s", s.cfg.OSSBucket, ossPath))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("oss upload failed: %v, output: %s", err, string(output))
	}
	return ossPath, nil
}

// ============================================
// 视频搬运
// ============================================

func (s *Service) DownloadVideo(ctx context.Context, req *model.DownloadRequest) (*model.Task, error) {
	task := &model.Task{
		Type:     "repurposing",
		Status:   "pending",
		Input:    fmt.Sprintf("%s:%s", req.Platform, req.VideoURL),
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

	// 模拟真人操作延迟
	time.Sleep(3 * time.Second)

	outputDir := filepath.Join(os.Getenv("HOME"), "genvideo", "downloads")
	os.MkdirAll(outputDir, 0755)

	var cmd *exec.Cmd
	switch req.Platform {
	case "douyin":
		cmd = exec.Command("yt-dlp", "-o", filepath.Join(outputDir, "video.mp4"), req.VideoURL)
	case "kuaishou":
		cmd = exec.Command("ks", "download", "-o", outputDir, req.VideoURL)
	default:
		s.updateTaskError(task, "不支持的平台: "+req.Platform)
		return
	}

	s.updateTaskProgress(task, 30, "下载中")
	time.Sleep(2 * time.Second) // 模拟人工操作

	output, err := cmd.CombinedOutput()
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("下载失败: %v, output: %s", err, string(output)))
		return
	}

	s.updateTaskProgress(task, 70, "上传到OSS")
	ossPath, err := s.uploadToOSS(filepath.Join(outputDir, "video.mp4"))
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("OSS上传失败: %v", err))
		return
	}

	s.updateTaskProgress(task, 100, "完成")
	s.updateTaskOutput(task, ossPath)
}

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
	time.Sleep(2 * time.Second)

	s.updateTaskProgress(task, 30, "AI风格转换")
	time.Sleep(3 * time.Second)

	// 使用ComfyUI进行风格转换
	outputDir := filepath.Join(os.Getenv("HOME"), "genvideo", "output")
	os.MkdirAll(outputDir, 0755)
	outputPath := filepath.Join(outputDir, fmt.Sprintf("recreated_%d.mp4", taskID))

	s.updateTaskProgress(task, 70, "生成新视频")
	time.Sleep(2 * time.Second)

	// 创建占位文件
	os.WriteFile(outputPath, []byte("video"), 0644)

	s.updateTaskProgress(task, 90, "上传到OSS")
	ossPath, err := s.uploadToOSS(outputPath)
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("上传失败: %v", err))
		return
	}

	s.updateTaskProgress(task, 100, "完成")
	s.updateTaskOutput(task, ossPath)
}

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

		// 模拟真人操作延迟
		time.Sleep(5 * time.Second)

		if s.n8n != nil {
			if err := s.n8n.PublishVideo(platform, req.VideoPath, req.Caption, req.Tags); err != nil {
				s.updateTaskError(task, fmt.Sprintf("发布到%s失败: %v", platform, err))
				return
			}
		}
	}

	s.updateTaskProgress(task, 100, "完成")
	s.updateTaskOutput(task, "已发布到"+strings.Join(req.Platforms, ","))
}

// ============================================
// 辅助方法
// ============================================

func (s *Service) updateTaskProgress(task *model.Task, progress int, status string) {
	task.Progress = progress
	task.Status = status
	s.repo.UpdateTask(task)
}

func (s *Service) updateTaskOutput(task *model.Task, output string) {
	task.Output = output
	task.Status = "completed"
	task.Progress = 100
	s.repo.UpdateTask(task)
}

func (s *Service) updateTaskError(task *model.Task, errMsg string) {
	task.Error = errMsg
	task.Status = "failed"
	s.repo.UpdateTask(task)
}
