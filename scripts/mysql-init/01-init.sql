-- MySQL Initialization Script
-- This script runs automatically when MySQL container starts for the first time

-- Create additional databases
CREATE DATABASE IF NOT EXISTS `lumine_dev` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS `lumine_test` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create users with proper permissions
CREATE USER IF NOT EXISTS 'dev'@'%' IDENTIFIED BY 'dev';
GRANT ALL PRIVILEGES ON `lumine_%`.* TO 'dev'@'%';

-- Flush privileges
FLUSH PRIVILEGES;

-- Log initialization
SELECT 'MySQL initialization complete!' AS message;
