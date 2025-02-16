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

-- Parent table
CREATE TABLE parents (
    parent_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    preferred_language VARCHAR(50) DEFAULT 'English',
    notifications_enabled BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id)
);

-- Tutor table
CREATE TABLE tutors (
    tutor_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    bio TEXT,
    qualifications TEXT NOT NULL,
    hourly_rate DECIMAL(10,2) NOT NULL,
    teaching_style TEXT,
    is_available BOOLEAN DEFAULT TRUE,
    demo_video_url VARCHAR(255),
    meet_test_passed BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    INDEX idx_is_available (is_available)
);

-- Student table
CREATE TABLE students (
    student_id INT AUTO_INCREMENT PRIMARY KEY,
    parent_id INT NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    age INT NOT NULL,
    grade_level VARCHAR(50) NOT NULL,
    learning_goals TEXT,
    preferred_language VARCHAR(50) DEFAULT 'English',
    FOREIGN KEY (parent_id) REFERENCES parents(parent_id) ON DELETE CASCADE,
    INDEX idx_parent_id (parent_id)
);

-- Booking table
CREATE TABLE bookings (
    booking_id INT AUTO_INCREMENT PRIMARY KEY,
    parent_id INT NOT NULL,
    tutor_id INT NOT NULL,
    student_id INT NOT NULL,
    booking_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    session_start DATETIME NOT NULL,
    session_end DATETIME NOT NULL,
    status ENUM('pending', 'confirmed', 'cancelled', 'completed') NOT NULL,
    is_recurring BOOLEAN DEFAULT FALSE,
    recurrence_pattern VARCHAR(50),
    FOREIGN KEY (parent_id) REFERENCES parents(parent_id),
    FOREIGN KEY (tutor_id) REFERENCES tutors(tutor_id),
    FOREIGN KEY (student_id) REFERENCES students(student_id),
    INDEX idx_session_start (session_start),
    INDEX idx_status (status)
);

-- Teaching Session table (updated with google_meet_id)
CREATE TABLE teaching_sessions (
    session_id INT AUTO_INCREMENT PRIMARY KEY,
    course_id INT NOT NULL,
    student_id INT NOT NULL,
    google_meet_id VARCHAR(255) NOT NULL,
    actual_start DATETIME,
    actual_end DATETIME,
    status ENUM('scheduled', 'in_progress', 'completed', 'cancelled') NOT NULL,
    topics_covered TEXT,
    homework_assigned TEXT,
    recording_enabled BOOLEAN DEFAULT FALSE,
    session_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (course_id) REFERENCES courses(course_id),
    FOREIGN KEY (student_id) REFERENCES students(student_id),
    INDEX idx_status (status),
    INDEX idx_google_meet_id (google_meet_id)
);

-- Session Recording table
CREATE TABLE session_recordings (
    recording_id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INT NOT NULL,
    recording_url VARCHAR(255) NOT NULL,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    duration INT NOT NULL,
    storage_path VARCHAR(255) NOT NULL,
    expiry_date DATETIME NOT NULL,
    FOREIGN KEY (session_id) REFERENCES teaching_sessions(session_id),
    INDEX idx_expiry_date (expiry_date)
);

-- Availability table
CREATE TABLE availabilities (
    availability_id INT AUTO_INCREMENT PRIMARY KEY,
    tutor_id INT NOT NULL,
    day_of_week ENUM('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday') NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    timezone VARCHAR(50) NOT NULL,
    is_recurring BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (tutor_id) REFERENCES tutors(tutor_id),
    INDEX idx_tutor_availability (tutor_id, day_of_week)
);

-- Payment table
CREATE TABLE payments (
    payment_id INT AUTO_INCREMENT PRIMARY KEY,
    subscription_id INT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payment_status ENUM('pending', 'completed', 'failed', 'refunded') NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    transaction_id VARCHAR(255) UNIQUE,
    FOREIGN KEY (subscription_id) REFERENCES course_subscriptions(subscription_id),
    INDEX idx_payment_status (payment_status)
);

-- Session Feedback table
CREATE TABLE session_feedbacks (
    feedback_id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INT NOT NULL,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comments TEXT,
    technical_quality TEXT,
    teaching_quality TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES teaching_sessions(session_id),
    INDEX idx_rating (rating)
);

-- Tutor Specialty table
CREATE TABLE tutor_specialties (
    specialty_id INT AUTO_INCREMENT PRIMARY KEY,
    tutor_id INT NOT NULL,
    subject VARCHAR(100) NOT NULL,
    level VARCHAR(50) NOT NULL,
    certification TEXT,
    years_experience INT NOT NULL,
    FOREIGN KEY (tutor_id) REFERENCES tutors(tutor_id),
    INDEX idx_subject (subject)
);

-- System Metrics table
CREATE TABLE system_metrics (
    metric_id INT AUTO_INCREMENT PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    metric_type VARCHAR(50) NOT NULL,
    value FLOAT NOT NULL,
    description TEXT,
    INDEX idx_timestamp (timestamp),
    INDEX idx_metric_type (metric_type)
);

-- Report table (New)
CREATE TABLE reports (
    report_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    report_type ENUM('session_summary', 'tutor_performance', 'student_progress', 'financial', 'system_status') NOT NULL,
    report_period_start DATE NOT NULL,
    report_period_end DATE NOT NULL,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    report_data JSON,
    status ENUM('processing', 'completed', 'failed') NOT NULL DEFAULT 'processing',
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    INDEX idx_user_reports (user_id, report_type),
    INDEX idx_generated_at (generated_at)
);

-- Course table
CREATE TABLE courses (
    course_id INT AUTO_INCREMENT PRIMARY KEY,
    tutor_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    total_sessions INT NOT NULL,
    level VARCHAR(50) NOT NULL,
    subject VARCHAR(100) NOT NULL,
    status ENUM('draft', 'published', 'archived') NOT NULL DEFAULT 'draft',
    syllabus TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tutor_id) REFERENCES tutors(tutor_id),
    INDEX idx_subject_level (subject, level),
    INDEX idx_status (status)
);

-- Course Subscription table
CREATE TABLE course_subscriptions (
    subscription_id INT AUTO_INCREMENT PRIMARY KEY,
    course_id INT NOT NULL,
    parent_id INT NOT NULL,
    student_id INT NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    status ENUM('active', 'completed', 'cancelled', 'expired') NOT NULL DEFAULT 'active',
    sessions_remaining INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (course_id) REFERENCES courses(course_id),
    FOREIGN KEY (parent_id) REFERENCES parents(parent_id),
    FOREIGN KEY (student_id) REFERENCES students(student_id),
    INDEX idx_status (status)
);


-- Triggers
DELIMITER //

-- Trigger to update sessions_remaining when a session is completed
CREATE TRIGGER after_session_completion_update_subscription
AFTER UPDATE ON teaching_sessions
FOR EACH ROW
BEGIN
    IF NEW.status = 'completed' AND OLD.status != 'completed' THEN
        UPDATE course_subscriptions cs
        JOIN courses c ON cs.course_id = c.course_id
        SET cs.sessions_remaining = cs.sessions_remaining - 1
        WHERE c.course_id = NEW.course_id;
    END IF;
END//

-- Trigger to set subscription status to completed when sessions are exhausted
CREATE TRIGGER after_subscription_update
AFTER UPDATE ON course_subscriptions
FOR EACH ROW
BEGIN
    IF NEW.sessions_remaining = 0 AND OLD.sessions_remaining > 0 THEN
        UPDATE course_subscriptions
        SET status = 'completed'
        WHERE subscription_id = NEW.subscription_id;
    END IF;
END//

-- Trigger to log system metric when new booking is created
CREATE TRIGGER after_booking_creation
AFTER INSERT ON bookings
FOR EACH ROW
BEGIN
    INSERT INTO system_metrics (metric_type, value, description)
    VALUES ('new_booking', 1, CONCAT('New booking created for tutor_id: ', NEW.tutor_id));
END//

-- More logging trigger in the future

DELIMITER ;
