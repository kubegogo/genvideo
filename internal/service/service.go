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
)

type Service struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

// Video Repurposing Functions

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

	// Start async download
	go s.executeDownload(task.ID, req)

	return task, nil
}

func (s *Service) executeDownload(taskID int64, req *model.DownloadRequest) {
	task, err := s.repo.GetTask(taskID)
	if err != nil {
		return
	}

	s.updateTaskProgress(task, 10, "starting download")

	// Simulate download based on platform
	outputDir := filepath.Join(os.Getenv("HOME"), "genvideo", "downloads")
	os.MkdirAll(outputDir, 0755)

	var cmd *exec.Cmd
	switch req.Platform {
	case "douyin":
		// Use yt-dlp or similar tool for Douyin
		outputPath := filepath.Join(outputDir, fmt.Sprintf("video_%d.mp4", taskID))
		cmd = exec.Command("yt-dlp", "-o", outputPath, req.VideoURL)
	case "kuaishou":
		outputPath := filepath.Join(outputDir, fmt.Sprintf("video_%d.mp4", taskID))
		cmd = exec.Command("ks", "download", "-o", outputPath, req.VideoURL)
	default:
		s.updateTaskError(task, "unsupported platform: "+req.Platform)
		return
	}

	s.updateTaskProgress(task, 30, "downloading")

	output, err := cmd.CombinedOutput()
	if err != nil {
		s.updateTaskError(task, fmt.Sprintf("download failed: %v, output: %s", err, string(output)))
		return
	}

	s.updateTaskProgress(task, 70, "uploading to OSS")

	// Upload to OSS
	ossPath := fmt.Sprintf("downloads/%d/video.mp4", taskID)
	if err := s.uploadToOSS(outputDir, ossPath); err != nil {
		s.updateTaskError(task, fmt.Sprintf("OSS upload failed: %v", err))
		return
	}

	s.updateTaskProgress(task, 100, "completed")
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

	s.updateTaskProgress(task, 10, "analyzing original video")

	// Call AI to recreate video based on style
	s.updateTaskProgress(task, 50, "AI processing")

	// Generate new video via n8n or ComfyUI
	outputDir := filepath.Join(os.Getenv("HOME"), "genvideo", "output")
	os.MkdirAll(outputDir, 0755)
	outputPath := filepath.Join(outputDir, fmt.Sprintf("recreated_%d.mp4", taskID))

	s.updateTaskProgress(task, 80, "generating video")

	// Call ComfyUI workflow
	workflowResult := s.callComfyUI(req.Style, req.OriginalVideo)
	if workflowResult == "" {
		s.updateTaskError(task, "ComfyUI workflow failed")
		return
	}

	// Copy result to output path (simulated)
	// In real implementation, this would be the actual generated video path

	s.updateTaskProgress(task, 90, "uploading to OSS")

	ossPath := fmt.Sprintf("output/%d/video.mp4", taskID)
	s.uploadToOSS(outputPath, ossPath)

	s.updateTaskProgress(task, 100, "completed")
	s.updateTaskOutput(task, ossPath)
}

// Script-to-Video Functions

func (s *Service) GenerateScript(ctx context.Context, req *model.ScriptRequest) (*model.Task, error) {
	task := &model.Task{
		Type:     "script_to_video",
		Status:   "pending",
		Input:    string(req.InputType) + ":" + req.Input,
		Progress: 0,
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}

	go s.executeScriptGeneration(task.ID, req)

	return task, nil
}

func (s *Service) executeScriptGeneration(taskID int64, req *model.ScriptRequest) {
	task, _ := s.repo.GetTask(taskID)

	s.updateTaskProgress(task, 10, "parsing input")

	// Generate script using AI
	s.updateTaskProgress(task, 50, "generating script")

	var script string
	if s.cfg.AIProvider == "minimax" {
		script = s.callMinimaxScript(req.Input, req.Style, req.Duration)
	} else {
		script = s.callOllamaScript(req.Input, req.Style, req.Duration)
	}

	s.updateTaskProgress(task, 100, "completed")
	s.updateTaskOutput(task, script)
}

func (s *Service) GenerateStoryboard(ctx context.Context, req *model.StoryboardRequest) (*model.Task, error) {
	task := &model.Task{
		Type:     "script_to_video",
		Status:   "pending",
		Input:    "storyboard:" + req.SceneCount,
		Progress: 0,
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}

	go s.executeStoryboardGeneration(task.ID, req)

	return task, nil
}

func (s *Service) executeStoryboardGeneration(taskID int64, req *model.StoryboardRequest) {
	task, _ := s.repo.GetTask(taskID)

	s.updateTaskProgress(task, 30, "generating storyboard")

	// Call AI to generate storyboard
	var storyboard string
	if s.cfg.AIProvider == "minimax" {
		storyboard = s.callMinimaxStoryboard(req.Script, req.SceneCount)
	} else {
		storyboard = s.callOllamaStoryboard(req.Script, req.SceneCount)
	}

	s.updateTaskProgress(task, 100, "completed")
	s.updateTaskOutput(task, storyboard)
}

func (s *Service) GenerateFrames(ctx context.Context, req *model.FrameRequest) (*model.Task, error) {
	task := &model.Task{
		Type:     "script_to_video",
		Status:   "pending",
		Input:    "frames",
		Progress: 0,
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}

	go s.executeFrameGeneration(task.ID, req)

	return task, nil
}

func (s *Service) executeFrameGeneration(taskID int64, req *model.FrameRequest) {
	task, _ := s.repo.GetTask(taskID)

	s.updateTaskProgress(task, 30, "generating first/last frames")

	// Call ComfyUI to generate images
	frames := s.callComfyUIFrames(req.Storyboard, req.Style)

	s.updateTaskProgress(task, 100, "completed")
	s.updateTaskOutput(task, strings.Join(frames, ","))
}

func (s *Service) GenerateVideo(ctx context.Context, req *model.VideoGenerationRequest) (*model.Task, error) {
	task := &model.Task{
		Type:     "script_to_video",
		Status:   "pending",
		Input:    "video generation",
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

	s.updateTaskProgress(task, 20, "preparing video generation")

	// Call ComfyUI video workflow
	s.updateTaskProgress(task, 50, "generating video via ComfyUI")

	outputPath := s.callComfyUIVideo(req.Storyboard, req.Frames, req.Duration)

	s.updateTaskProgress(task, 80, "uploading to OSS")

	ossPath := fmt.Sprintf("output/%d/final.mp4", taskID)
	s.uploadToOSS(outputPath, ossPath)

	s.updateTaskProgress(task, 100, "completed")
	s.updateTaskOutput(task, ossPath)
}

// Publish Functions

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
		s.updateTaskProgress(task, progress, "publishing to "+platform)

		// Simulate human-like publishing with delays
		time.Sleep(3 * time.Second)

		// Call n8n workflow to publish
		if err := s.callN8nPublish(platform, req.VideoPath, req.Caption, req.Tags); err != nil {
			s.updateTaskError(task, fmt.Sprintf("publish failed to %s: %v", platform, err))
			return
		}
	}

	s.updateTaskProgress(task, 100, "completed")
	s.updateTaskOutput(task, "published to "+strings.Join(req.Platforms, ","))
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
	// Use ossutil or oss SDK to upload
	cmd := exec.Command("ossutil", "cp", localPath, fmt.Sprintf("oss://%s/%s", s.cfg.OSSBucket, ossPath))
	_, err := cmd.CombinedOutput()
	return err
}

func (s *Service) callMinimaxScript(input, style string, duration int) string {
	// Call Minimax API for script generation
	// In real implementation, this would call the actual API
	return fmt.Sprintf("Generated script for %s style, %d seconds duration", style, duration)
}

func (s *Service) callOllamaScript(input, style string, duration int) string {
	// Call Ollama for script generation
	return fmt.Sprintf("Ollama generated script for %s style, %d seconds duration", style, duration)
}

func (s *Service) callMinimaxStoryboard(script string, sceneCount int) string {
	return fmt.Sprintf("Storyboard with %d scenes", sceneCount)
}

func (s *Service) callOllamaStoryboard(script string, sceneCount int) string {
	return fmt.Sprintf("Ollama storyboard with %d scenes", sceneCount)
}

func (s *Service) callComfyUI(style, input string) string {
	// Call ComfyUI API
	return ""
}

func (s *Service) callComfyUIFrames(storyboard, style string) []string {
	// Call ComfyUI for image generation
	return []string{"frame1.png", "frame2.png"}
}

func (s *Service) callComfyUIVideo(storyboard string, frames []string, duration int) string {
	// Call ComfyUI video workflow
	return filepath.Join(os.Getenv("HOME"), "genvideo", "output", "video.mp4")
}

func (s *Service) callN8nPublish(platform, videoPath, caption string, tags []string) error {
	// Call n8n webhook to trigger publish workflow
	return nil
}
