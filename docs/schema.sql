-- GenVideo Database Schema
-- MySQL 8.0

CREATE DATABASE IF NOT EXISTS genvideo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE genvideo;

-- Tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(50) NOT NULL COMMENT 'repurposing or script_to_video',
    status VARCHAR(50) NOT NULL DEFAULT 'pending' COMMENT 'pending, processing, completed, failed',
    input TEXT NOT NULL COMMENT 'input data or reference',
    output TEXT COMMENT 'output result or path',
    error TEXT COMMENT 'error message if failed',
    progress INT DEFAULT 0 COMMENT '0-100',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Video providers table
CREATE TABLE IF NOT EXISTS video_providers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) NOT NULL UNIQUE COMMENT 'douyin, kuaishou, bilibili, xiaohongshu',
    cookie TEXT COMMENT 'login cookie for platform',
    status VARCHAR(20) DEFAULT 'active' COMMENT 'active, inactive',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_platform (platform)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- AI providers table
CREATE TABLE IF NOT EXISTS ai_providers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(50) NOT NULL COMMENT 'minimax, self_hosted',
    api_key VARCHAR(255) COMMENT 'API key (encrypted)',
    base_url VARCHAR(255) COMMENT 'Base URL for self-hosted',
    is_active BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- OSS configuration table
CREATE TABLE IF NOT EXISTS oss_config (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    endpoint VARCHAR(255) NOT NULL,
    access_key VARCHAR(255) COMMENT 'Access key (encrypted)',
    secret_key VARCHAR(255) COMMENT 'Secret key (encrypted)',
    bucket VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Publish platforms table
CREATE TABLE IF NOT EXISTS publish_platforms (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) NOT NULL UNIQUE COMMENT 'youtube, tiktok, instagram, etc.',
    config TEXT COMMENT 'Platform-specific configuration',
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_platform (platform)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Risk control rules table (auto-updated weekly)
CREATE TABLE IF NOT EXISTS risk_control_rules (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(50) NOT NULL,
    rule_type VARCHAR(50) NOT NULL COMMENT 'typing_speed, click_interval, etc.',
    rule_value TEXT NOT NULL COMMENT 'JSON value',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_platform (platform),
    UNIQUE KEY uk_platform_rule (platform, rule_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Logs table
CREATE TABLE IF NOT EXISTS logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    level VARCHAR(20) NOT NULL DEFAULT 'info',
    message TEXT NOT NULL,
    context TEXT COMMENT 'JSON context',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_level (level),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
