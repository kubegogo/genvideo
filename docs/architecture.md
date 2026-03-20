# GenVideo 系统架构

## 概述

AI 视频生成系统，支持视频搬运和脚本转视频两大功能。

## 核心流程

### 脚本转视频（剪映模式）
```
输入(关键词/文档/小说)
    ↓
生成剧本 (Minimax/Ollama)
    ↓
生成分镜剧本
    ↓
生成首尾帧图片 (ComfyUI)
    ↓
生成视频 (ComfyUI)
    ↓
发布 (n8n 模拟人工操作)
```

### 视频搬运
```
平台选择 → 搜索热门视频 → 下载 → AI二次创作 → 发布
```

## 技术架构

```
┌──────────────┐     ┌──────────────┐
│   Frontend   │────▶│   Backend    │
│   (React)    │◀────│    (Go)      │
│   :3003      │     │    :3004     │
└──────────────┘     └──────┬───────┘
                            │
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
         ┌────────┐  ┌────────┐  ┌────────┐
         │  MySQL │  │  Redis │  │   AI   │
         │  :3306 │  │  :6379 │  │Services│
         └────────┘  └────────┘  └────────┘
```

## 目录结构

```
genvideo/
├── cmd/server/           # 服务入口
├── internal/
│   ├── config/          # 配置管理
│   ├── handler/         # HTTP处理器
│   ├── middleware/      # 中间件
│   ├── model/          # 数据模型
│   ├── repository/     # 数据访问层
│   └── service/        # 业务逻辑
├── pkg/
│   ├── minimax/        # Minimax API客户端
│   ├── comfyui/        # ComfyUI客户端
│   ├── n8n/            # n8n客户端
│   └── response/       # 响应工具
├── frontend/            # React前端
├── docs/               # 文档
└── docker-compose.yml  # Docker编排
```

## 启动方式

```bash
# 构建并启动
docker-compose up -d --build

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

## 端口

| 服务 | 端口 |
|------|------|
| Frontend | http://localhost:3003 |
| Backend | http://localhost:3004 |
| MySQL | localhost:3306 |
| Redis | localhost:6379 |
