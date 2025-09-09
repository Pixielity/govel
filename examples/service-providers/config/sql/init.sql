-- Service Provider Example Database Initialization
-- This script creates sample tables and data for demonstrating the PostgreSQL service provider

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- User profiles table
CREATE TABLE IF NOT EXISTS user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Posts table
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    content TEXT,
    published BOOLEAN DEFAULT false,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Comments table
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    approved BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Login logs table
CREATE TABLE IF NOT EXISTS login_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN DEFAULT true,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Application settings table
CREATE TABLE IF NOT EXISTS app_settings (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(active);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_published ON posts(published);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_login_logs_user_id ON login_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_login_logs_timestamp ON login_logs(timestamp);

-- Trigger function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers to tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_posts_updated_at BEFORE UPDATE ON posts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_comments_updated_at BEFORE UPDATE ON comments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_app_settings_updated_at BEFORE UPDATE ON app_settings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data
INSERT INTO users (name, email, password_hash, active) VALUES 
    ('John Doe', 'john@example.com', '$2a$10$example_hash_here', true),
    ('Jane Smith', 'jane@example.com', '$2a$10$example_hash_here', true),
    ('Bob Johnson', 'bob@example.com', '$2a$10$example_hash_here', false),
    ('Alice Brown', 'alice@example.com', '$2a$10$example_hash_here', true)
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_profiles (user_id, first_name, last_name, city, country) VALUES
    (1, 'John', 'Doe', 'New York', 'USA'),
    (2, 'Jane', 'Smith', 'London', 'UK'),
    (3, 'Bob', 'Johnson', 'Toronto', 'Canada'),
    (4, 'Alice', 'Brown', 'Sydney', 'Australia');

INSERT INTO posts (user_id, title, content, published, published_at) VALUES
    (1, 'Welcome to Service Providers', 'This is an example post demonstrating the PostgreSQL service provider.', true, CURRENT_TIMESTAMP - INTERVAL '1 day'),
    (1, 'Database Connections', 'Managing database connections efficiently with connection pooling.', true, CURRENT_TIMESTAMP - INTERVAL '2 hours'),
    (2, 'Redis Caching', 'Using Redis for high-performance caching in web applications.', true, CURRENT_TIMESTAMP - INTERVAL '3 hours'),
    (4, 'Draft Post', 'This is a draft post that is not yet published.', false, NULL);

INSERT INTO comments (post_id, user_id, content, approved) VALUES
    (1, 2, 'Great introduction! Looking forward to more posts.', true),
    (1, 4, 'Very informative, thanks for sharing.', true),
    (2, 3, 'Connection pooling is indeed crucial for performance.', true),
    (3, 1, 'Redis has been a game-changer for our application performance.', true);

INSERT INTO app_settings (key, value, description) VALUES
    ('app.name', 'Service Provider Example', 'Application name'),
    ('app.version', '1.0.0', 'Application version'),
    ('feature.caching', 'true', 'Enable caching features'),
    ('feature.logging', 'true', 'Enable detailed logging'),
    ('max_upload_size', '10485760', 'Maximum file upload size in bytes (10MB)')
ON CONFLICT (key) DO NOTHING;

-- Create a view for user statistics
CREATE OR REPLACE VIEW user_stats AS
SELECT 
    u.id,
    u.name,
    u.email,
    u.active,
    COUNT(DISTINCT p.id) as post_count,
    COUNT(DISTINCT c.id) as comment_count,
    MAX(ll.timestamp) as last_login,
    u.created_at as joined_date
FROM users u
LEFT JOIN posts p ON u.id = p.user_id
LEFT JOIN comments c ON u.id = c.user_id
LEFT JOIN login_logs ll ON u.id = ll.user_id AND ll.success = true
GROUP BY u.id, u.name, u.email, u.active, u.created_at
ORDER BY u.id;

-- Create a function for health checks
CREATE OR REPLACE FUNCTION health_check()
RETURNS JSON AS $$
DECLARE
    result JSON;
    user_count INTEGER;
    post_count INTEGER;
    active_users INTEGER;
    db_size TEXT;
BEGIN
    SELECT COUNT(*) INTO user_count FROM users;
    SELECT COUNT(*) INTO post_count FROM posts WHERE published = true;
    SELECT COUNT(*) INTO active_users FROM users WHERE active = true;
    SELECT pg_size_pretty(pg_database_size(current_database())) INTO db_size;
    
    result := json_build_object(
        'status', 'healthy',
        'timestamp', CURRENT_TIMESTAMP,
        'database_size', db_size,
        'total_users', user_count,
        'active_users', active_users,
        'published_posts', post_count
    );
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Grant permissions (adjust as needed for your security requirements)
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO postgres;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO postgres;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO postgres;

-- Final success message
DO $$
BEGIN
    RAISE NOTICE 'Database initialization completed successfully!';
    RAISE NOTICE 'Sample data has been inserted for demonstration purposes.';
    RAISE NOTICE 'You can now test the PostgreSQL service provider with this database.';
END $$;
