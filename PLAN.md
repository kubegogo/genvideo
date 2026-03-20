# GenVideo 开发计划

## 项目概述

AI 视频生成系统，支持：
- **功能一**：视频搬运（下载 → AI二次创作 → 发布）
- **功能二**：脚本转视频（剧本 → 分镜 → 首尾帧 → 视频 → 发布）

## 技术栈

| 组件 | 技术 | 端口 |
|------|------|------|
| 前端 | React + Vite | 3003 |
| 后端 | Go + Gin | 3004 |
| 数据库 | MySQL 8 | 3306 |
| 缓存 | Redis 7 | 6379 |
| AI | Minimax API / Ollama | - |

## 开发阶段

### Phase 1: 项目初始化
- [ ] 初始化 Git 仓库
- [ ] 创建项目目录结构
- [ ] 编写 CLAUDE.md

### Phase 2: 后端开发
- [ ] Go 模块初始化
- [ ] 配置管理 (config)
- [ ] 数据模型 (model)
- [ ] 数据库层 (repository) - MySQL + Redis
- [ ] 业务逻辑层 (service)
- [ ] HTTP 处理器 (handler)
- [ ] 路由和中间件

### Phase 3: 前端开发
- [ ] React + Vite 初始化
- [ ] 页面组件
- [ ] API 调用
- [ ] 样式

### Phase 4: 文档
- [ ] 系统架构文档
- [ ] 数据库字典
- [ ] API 文档
- [ ] 部署文档

### Phase 5: Docker 环境
- [ ] docker-compose.yml
- [ ] Dockerfile.backend
- [ ] Dockerfile.frontend
- [ ] Makefile

## 核心流程

### 视频搬运流程
```
平台选择 → 搜索热门视频 → 下载 → AI风格转换 → 发布
```

### 脚本转视频流程
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

## 文件结构

```
genvideo/
├── cmd/server/           # 入口
├── internal/
│   ├── config/          # 配置
│   ├── handler/         # HTTP
│   ├── middleware/      # 中间件
│   ├── model/          # 数据模型
│   ├── repository/      # 数据访问
│   └── service/        # 业务逻辑
├── pkg/
│   ├── minimax/        # AI客户端
│   ├── n8n/            # n8n客户端
│   ├── comfyui/         # ComfyUI客户端
│   └── oss/            # OSS客户端
├── frontend/           # React前端
├── docs/               # 文档
├── docker-compose.yml
├── Dockerfile.backend
└── Dockerfile.frontend
```

## 部署

- 开发: `docker-compose up -d --build`
- 测试: Docker 容器内运行
- 端口: 前端3003, 后端3004
