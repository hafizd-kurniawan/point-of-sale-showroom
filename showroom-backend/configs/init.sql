-- Initialize database
\echo 'Creating database and extensions...'

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\echo 'Database initialization completed.'