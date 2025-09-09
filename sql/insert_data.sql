/*
 * This script inserts the sample album data into the albums table
 */
-- Copy the sample data from the CSV file into the albums table
COPY albums (title, artist, genre, release_date) FROM '/tmp/data.csv' WITH (FORMAT CSV, HEADER, DELIMITER ',');
