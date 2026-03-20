# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AI video generation system with agent teams that auto-evolve. Two main capabilities:

1. **Video Repurposing**: Download popular videos from platforms (likes/views/favorites), recreate with AI, publish to other platforms
2. **Script-to-Video**: Generate original videos from keywords, documents, or novels through: script → storyboard → first/last frame images → video → auto-publish

## Tech Stack

- **Frontend**: React + Vite
- **Backend**: Go + Gin framework
- **Database**: MySQL 8 (Docker)
- **Cache**: Redis 7 (Docker)
- **AI**: Minimax model (fallback: wait 5 hours if token exhausted)

## Architecture

- **n8n**: Workflow automation for publishing
- **ComfyUI**: Image/video generation via AI workflows
- **Ollama**: Local LLM for self-hosted option
- **Video Generation Options** (user-selectable):
  1. External AI API (user provides API key)
  2. Self-hosted n8n + ComfyUI + Ollama (user configures)

## Code Repository

- Push to: https://github.com/kubegogo/genvideo.git
- All documentation: `docs/` directory

## Development Environment

- Development: local Docker environment
- Testing: Docker container
- MySQL/Redis: local Docker self-hosted
- Frontend port: 3003
- Backend port: 3004

## Agent Team Behavior

- Enable with: `$env:CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS = "1"`
- Auto-evolve based on requirements
- Search GitHub for best solutions (including kubegogo repos)
- Create repos in kubegogo organization
- Use minimax model; wait 5 hours for token reset if exhausted
- Parallel execution when possible

## Core Flows

### Video Repurposing
```
Platform → Search Popular Videos → Download → AI Recreate → Publish
```

### Script-to-Video
```
Input (keywords/doc/novel)
    ↓
Generate Script (Minimax/Ollama)
    ↓
Generate Storyboard
    ↓
Generate First/Last Frame Images (ComfyUI)
    ↓
Generate Video (ComfyUI)
    ↓
Publish (n8n simulates human operation)
```

## Risk Control

- All automation must mimic human behavior (slow typing 3s-60s intervals, screen sliding, clicking)
- Update platform risk control rules every 12 hours via web search
- Log all failures and account exceptions to files
- Auto-update weekly (Sundays 23:00)

## Video Data

- All video files must be synced to Alibaba Cloud OSS using ossutil
- OSS configuration is provided by user
