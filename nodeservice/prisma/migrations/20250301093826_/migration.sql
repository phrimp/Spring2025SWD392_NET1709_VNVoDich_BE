-- CreateTable
CREATE TABLE `User` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `deleted_at` DATETIME(3) NULL,
    `username` VARCHAR(50) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `role` ENUM('Parent', 'Children', 'Tutor', 'Admin') NOT NULL DEFAULT 'Parent',
    `phone` VARCHAR(20) NULL,
    `full_name` VARCHAR(255) NULL,
    `picture` VARCHAR(191) NULL,
    `is_verified` BOOLEAN NOT NULL DEFAULT false,
    `status` ENUM('Active', 'Suspended', 'Banned') NOT NULL DEFAULT 'Active',
    `last_login_at` BIGINT NULL,
    `account_locked` BOOLEAN NOT NULL DEFAULT false,
    `password_changed_at` BIGINT NULL,

    UNIQUE INDEX `User_username_key`(`username`),
    UNIQUE INDEX `User_email_key`(`email`),
    INDEX `User_status_idx`(`status`),
    INDEX `User_role_idx`(`role`),
    INDEX `User_deleted_at_idx`(`deleted_at`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Tutor` (
    `id` INTEGER NOT NULL,
    `bio` VARCHAR(191) NOT NULL,
    `qualifications` VARCHAR(191) NOT NULL,
    `teaching_style` VARCHAR(191) NOT NULL,
    `is_available` BOOLEAN NOT NULL DEFAULT false,
    `demo_video_url` VARCHAR(191) NULL,
    `image` VARCHAR(191) NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Parent` (
    `id` INTEGER NOT NULL,
    `preferred_language` VARCHAR(191) NOT NULL,
    `notifications_enable` BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Children` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `age` INTEGER NOT NULL,
    `grade_level` VARCHAR(191) NOT NULL,
    `learning_goals` VARCHAR(191) NOT NULL,
    `full_name` VARCHAR(191) NOT NULL,
    `password` VARCHAR(191) NOT NULL,
    `parent_id` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `TutorSpecialty` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `subject` VARCHAR(191) NOT NULL,
    `level` VARCHAR(191) NOT NULL,
    `certification` VARCHAR(191) NOT NULL,
    `years_experience` INTEGER NOT NULL,
    `tutor_id` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Course` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(191) NOT NULL,
    `description` VARCHAR(191) NULL,
    `price` DOUBLE NOT NULL,
    `total_lessons` INTEGER NOT NULL,
    `grade` INTEGER NOT NULL,
    `subject` VARCHAR(191) NOT NULL,
    `status` ENUM('Draft', 'Published') NOT NULL DEFAULT 'Draft',
    `image` VARCHAR(191) NULL,
    `tutor_id` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Lesson` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(191) NOT NULL,
    `description` VARCHAR(191) NULL,
    `learning_objectives` VARCHAR(191) NULL,
    `materials_needed` VARCHAR(191) NULL,
    `course_id` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Availability` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `tutor_id` INTEGER NOT NULL,
    `timeGap` INTEGER NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    UNIQUE INDEX `Availability_tutor_id_key`(`tutor_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `DayAvailability` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `availability_id` INTEGER NOT NULL,
    `day` ENUM('MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY', 'SUNDAY') NOT NULL,
    `startTime` DATETIME(3) NOT NULL,
    `endTime` DATETIME(3) NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `CourseSubscription` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `status` VARCHAR(191) NOT NULL,
    `sessions_remaining` INTEGER NULL,
    `course_id` INTEGER NOT NULL,
    `children_id` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `TeachingSession` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `google_meet_id` VARCHAR(191) NULL,
    `startTime` DATETIME(3) NOT NULL,
    `endTime` DATETIME(3) NOT NULL,
    `status` VARCHAR(191) NOT NULL,
    `topics_covered` VARCHAR(191) NULL,
    `homework_assigned` VARCHAR(191) NULL,
    `subscription_id` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `SessionFeedback` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `rating` INTEGER NOT NULL,
    `comment` VARCHAR(191) NULL,
    `teaching_quality` VARCHAR(191) NULL,
    `session_id` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `CourseSubscriptionSchedule` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `subscription_id` INTEGER NOT NULL,
    `day` ENUM('MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY', 'SUNDAY') NOT NULL,
    `startTime` DATETIME(3) NOT NULL,
    `endTime` DATETIME(3) NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `TutorReview` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `rating` INTEGER NOT NULL,
    `review_content` VARCHAR(191) NULL,
    `createAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `tutor_id` INTEGER NOT NULL,
    `parent_id` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `CourseReview` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `rating` INTEGER NOT NULL,
    `review_content` VARCHAR(191) NULL,
    `createAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `course_id` INTEGER NOT NULL,
    `parent_id` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- AddForeignKey
ALTER TABLE `Tutor` ADD CONSTRAINT `Tutor_id_fkey` FOREIGN KEY (`id`) REFERENCES `User`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `Parent` ADD CONSTRAINT `Parent_id_fkey` FOREIGN KEY (`id`) REFERENCES `User`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `Children` ADD CONSTRAINT `Children_parent_id_fkey` FOREIGN KEY (`parent_id`) REFERENCES `Parent`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `TutorSpecialty` ADD CONSTRAINT `TutorSpecialty_tutor_id_fkey` FOREIGN KEY (`tutor_id`) REFERENCES `Tutor`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `Course` ADD CONSTRAINT `Course_tutor_id_fkey` FOREIGN KEY (`tutor_id`) REFERENCES `Tutor`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `Lesson` ADD CONSTRAINT `Lesson_course_id_fkey` FOREIGN KEY (`course_id`) REFERENCES `Course`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `Availability` ADD CONSTRAINT `Availability_tutor_id_fkey` FOREIGN KEY (`tutor_id`) REFERENCES `Tutor`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `DayAvailability` ADD CONSTRAINT `DayAvailability_availability_id_fkey` FOREIGN KEY (`availability_id`) REFERENCES `Availability`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `CourseSubscription` ADD CONSTRAINT `CourseSubscription_course_id_fkey` FOREIGN KEY (`course_id`) REFERENCES `Course`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `CourseSubscription` ADD CONSTRAINT `CourseSubscription_children_id_fkey` FOREIGN KEY (`children_id`) REFERENCES `Children`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `TeachingSession` ADD CONSTRAINT `TeachingSession_subscription_id_fkey` FOREIGN KEY (`subscription_id`) REFERENCES `CourseSubscription`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `SessionFeedback` ADD CONSTRAINT `SessionFeedback_session_id_fkey` FOREIGN KEY (`session_id`) REFERENCES `TeachingSession`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `CourseSubscriptionSchedule` ADD CONSTRAINT `CourseSubscriptionSchedule_subscription_id_fkey` FOREIGN KEY (`subscription_id`) REFERENCES `CourseSubscription`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `TutorReview` ADD CONSTRAINT `TutorReview_tutor_id_fkey` FOREIGN KEY (`tutor_id`) REFERENCES `Tutor`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `TutorReview` ADD CONSTRAINT `TutorReview_parent_id_fkey` FOREIGN KEY (`parent_id`) REFERENCES `Parent`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `CourseReview` ADD CONSTRAINT `CourseReview_course_id_fkey` FOREIGN KEY (`course_id`) REFERENCES `Course`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `CourseReview` ADD CONSTRAINT `CourseReview_parent_id_fkey` FOREIGN KEY (`parent_id`) REFERENCES `Parent`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
