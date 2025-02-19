-- Create database
CREATE DATABASE IF NOT EXISTS online_tutoring_platform;
USE online_tutoring_platform;

-- User table
CREATE TABLE users (
    -- gorm.Model fields
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- User authentication fields
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    
    -- User information
    role ENUM('parent', 'kid', 'tutor', 'admin') NOT NULL DEFAULT 'kid',
    phone VARCHAR(20) NULL,
    full_name VARCHAR(255) NULL,
    
    -- Account status and verification
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    status ENUM('active', 'inactive', 'suspended', 'banned') NOT NULL DEFAULT 'inactive',
    email_verification_token VARCHAR(100) NULL,
    
    -- Security and tracking
    last_login_at BIGINT NULL,
    login_attempts INT UNSIGNED NOT NULL DEFAULT 0,
    account_locked BOOLEAN NOT NULL DEFAULT FALSE,
    password_changed_at BIGINT NULL,
    
    -- Password reset
    password_reset_token VARCHAR(100) NULL,
    password_reset_expires BIGINT NULL,
    
    -- Constraints and indexes
    CONSTRAINT users_username_unique UNIQUE (username),
    CONSTRAINT users_email_unique UNIQUE (email),
    
    -- Validation constraints
    CONSTRAINT chk_username CHECK (LENGTH(username) >= 3 AND LENGTH(username) <= 50),
    CONSTRAINT chk_password CHECK (LENGTH(password) >= 8),
    CONSTRAINT chk_full_name CHECK (full_name IS NULL OR LENGTH(full_name) >= 2),
    CONSTRAINT chk_phone CHECK (phone IS NULL OR LENGTH(phone) <= 20),
    
    -- Indexes for performance
    INDEX idx_users_status (status),
    INDEX idx_users_role (role),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
