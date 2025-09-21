-- Database initialization script for Yakka Backend
-- This script runs when the PostgreSQL container starts for the first time

-- Create development database if it doesn't exist
SELECT 'CREATE DATABASE yakka_dev'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'yakka_dev')\gexec

-- Create production database if it doesn't exist
SELECT 'CREATE DATABASE yakka_prod'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'yakka_prod')\gexec

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE yakka_dev TO postgres;
GRANT ALL PRIVILEGES ON DATABASE yakka_prod TO postgres;

-- Create extensions and custom types for development
\c yakka_dev;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create custom types for user status and roles
DO $$ BEGIN
    CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'pending');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('user', 'admin', 'moderator');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create extensions and custom types for production
\c yakka_prod;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create custom types for user status and roles
DO $$ BEGIN
    CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'pending');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('user', 'admin', 'moderator');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
