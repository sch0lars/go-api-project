/*
 * This script creates the music database and album table
 * for the go-api-project.
 */

-- Grant the docker user all privileges on the music database
GRANT ALL PRIVILEGES ON DATABASE music TO docker;
-- Connect to the music database
\c music;
-- Create the albums table
CREATE TABLE albums (
id bigserial primary key,
title varchar(255) not null,
artist varchar(255) not null,
genre varchar(255),
release_date date);
