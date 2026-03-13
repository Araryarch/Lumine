-- PostgreSQL Initialization Script
-- This script runs automatically when PostgreSQL container starts for the first time

-- Create additional databases
CREATE DATABASE lumine_dev;
CREATE DATABASE lumine_test;

-- Create users
CREATE USER dev WITH PASSWORD 'dev';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE lumine TO dev;
GRANT ALL PRIVILEGES ON DATABASE lumine_dev TO dev;
GRANT ALL PRIVILEGES ON DATABASE lumine_test TO dev;

-- Enable extensions
\c lumine
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

\c lumine_dev
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

\c lumine_test
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
