# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AI video generation system with agent teams that auto-evolve. Two main capabilities:

1. **Video Repurposing**: Download popular videos from platforms (likes/views/favorites), recreate with AI, publish to other platforms
2. **Script-to-Video**: Generate original videos from keywords, documents, or novels through: script → storyboard → first/last frame images → video → auto-publish

## Tech Stack

- **Frontend**: React
- **Backend**: Go
- **Database**: MySQL 8 (local Docker)
- **Cache**: Redis 7 (local Docker)
- **AI**: Minimax model (fallback: wait 5 hours if token exhausted)

## Architecture

- **n8n**: Workflow automation for video generation pipeline
- **ComfyUI**: Image/video generation via AI workflows
- **Ollama**: Local LLM for self-hosted option
- **Video Generation Options** (user-selectable):
  1. External AI API (user provides API key)
  2. Self-hosted n8n + ComfyUI + Ollama (user configures)

## Code Repository

- Push to: https://github.com/kubegogo/genvideo.git
- All documentation: `docs/` directory
- Include: product workflow, architecture, dependencies, database schema

## Development Environment

- Development: local machine
- Testing: Docker container
- MySQL/Redis: local Docker self-hosted

## Agent Team Behavior

- Enable with: `$env:CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS = "1"`
- Auto-evolve based on requirements
- Search GitHub for best solutions (including kubegogo repos)
- Create repos in kubegogo organization
- Use minimax model; wait 5 hours for token reset if exhausted
