# GenVideo 系统架构

## 概述

AI 视频生成系统，支持视频搬运和脚本转视频两大功能。

## 系统架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (React)                         │
│                     Port: 3003                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Backend (Go)                            │
│                     Port: 3004                                  │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  │
│  │ Handler │  │ Service │  │   Repo  │  │  Model  │  │  Config │  │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘  │
└─────────────────────────────────────────────────────────────────┘
         │                                           │
         ▼                                           ▼
┌─────────────────┐                    ┌─────────────────────────┐
│      MySQL      │                    │          Redis          │
│   (Docker)      │                    │        (Docker)         │
└─────────────────┘                    └─────────────────────────┘
         │                                           │
         ▼                                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                    External Services                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │
│  │  Minimax │  │    n8n   │  │  ComfyUI │  │   Ali OSS        │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## 核心模块

### Backend (Go)

| 模块 | 路径 | 说明 |
|------|------|------|
| handler | `internal/handler/` | HTTP 处理器，处理 API 请求 |
| service | `internal/service/` | 业务逻辑层，核心功能实现 |
| repository | `internal/repository/` | 数据访问层，MySQL/Redis 操作 |
| model | `internal/model/` | 数据模型定义 |
| config | `internal/config/` | 配置管理 |
| middleware | `internal/middleware/` | 中间件（CORS、日志、限流） |

### Frontend (React)

待实现

## API 端点

### 视频搬运

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/video/download` | 下载视频 |
| POST | `/api/v1/video/recreate` | 二次创作 |
| POST | `/api/v1/video/publish` | 发布视频 |

### 脚本转视频

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/script/generate` | 生成剧本 |
| POST | `/api/v1/script/storyboard` | 生成分镜 |
| POST | `/api/v1/script/frames` | 生成首尾帧 |
| POST | `/api/v1/script/video` | 生成视频 |

### 配置

| 方法 | 路径 | 说明 |
|------|------|------|
| GET/POST | `/api/v1/config/video-providers` | 视频平台配置 |
| GET/POST | `/api/v1/config/ai-providers` | AI 服务配置 |
| GET/POST | `/api/v1/config/oss` | OSS 配置 |

## 视频生成流程

### 方式一：外部 AI API

```
用户配置 API Key → 后端调用 Minimax API → 获取结果
```

### 方式二：自建 n8n + ComfyUI + Ollama

```
用户配置 → 后端调用 n8n Webhook → n8n 触发 ComfyUI → ComfyUI 调用 Ollama → 返回结果
```

## 数据流

1. **视频下载**: 用户请求 → Handler → Service → 调用 yt-dlp/平台工具 → 上传 OSS
2. **视频创作**: 用户请求 → Handler → Service → AI 处理 → 生成视频 → 上传 OSS
3. **脚本转视频**: 用户输入 → Handler → Service → AI 生成剧本/分镜/帧 → ComfyUI 生成视频 → 上传 OSS
4. **发布**: 用户请求 → Handler → Service → n8n 模拟人工操作 → 发布到目标平台

## 定时任务

- **每周日 23:00**: 自动更新风控规则
- **每 12 小时**: 更新平台风控规则

## 目录结构

```
genvideo/
├── cmd/                    # 入口
│   └── server/            # 服务入口
├── internal/              # 内部包
│   ├── config/           # 配置
│   ├── handler/          # HTTP 处理器
│   ├── middleware/       # 中间件
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问
│   └── service/         # 业务逻辑
├── pkg/                  # 公共包
│   ├── errors/          # 错误定义
│   ├── response/        # 响应格式
│   └── utils/           # 工具函数
├── docs/                 # 文档
│   ├── schema.sql       # 数据库 schema
│   └── architecture.md  # 架构文档
└── go.mod               # Go 模块
```
