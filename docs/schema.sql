-- GenVideo 数据库Schema
-- MySQL 8.0

CREATE DATABASE IF NOT EXISTS genvideo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE genvideo;

-- 任务表
CREATE TABLE tasks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(50) NOT NULL COMMENT 'repurposing, script_to_video',
    status VARCHAR(50) NOT NULL DEFAULT 'pending' COMMENT 'pending, processing, completed, failed',
    input TEXT NOT NULL COMMENT '输入内容',
    output TEXT COMMENT '输出路径',
    error TEXT COMMENT '错误信息',
    progress INT DEFAULT 0 COMMENT '0-100',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- AI服务商配置表
CREATE TABLE ai_providers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(50) NOT NULL COMMENT 'minimax, self_hosted',
    api_key VARCHAR(255) COMMENT 'API密钥',
    base_url VARCHAR(255) COMMENT '自建服务URL',
    is_active BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- OSS配置表
CREATE TABLE oss_config (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    endpoint VARCHAR(255) NOT NULL,
    access_key VARCHAR(255) COMMENT 'Access Key',
    secret_key VARCHAR(255) COMMENT 'Secret Key',
    bucket VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 视频平台配置表
CREATE TABLE video_providers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) NOT NULL UNIQUE COMMENT 'douyin, kuaishou, bilibili, xiaohongshu',
    cookie TEXT COMMENT '登录Cookie',
    status VARCHAR(20) DEFAULT 'active' COMMENT 'active, inactive',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 发布平台表
CREATE TABLE publish_platforms (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) NOT NULL UNIQUE COMMENT 'youtube, tiktok, instagram',
    config TEXT COMMENT '平台配置',
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
