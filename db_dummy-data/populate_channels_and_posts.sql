BEGIN TRANSACTION;

-- Temporary table to cache 5 users (for channel owners)
-- Ensures that if Users table is small, we don't run out of distinct users.
-- If Users table has < 5 users, this will take all of them.
CREATE TEMP TABLE TempUsersForChannels AS
SELECT ID AS UserID, Username, Avatar FROM Users ORDER BY RANDOM() LIMIT 5;

-- Descriptive channel names
CREATE TEMP TABLE ChannelNames(name TEXT PRIMARY KEY); -- Added PRIMARY KEY for potential safer joins
INSERT INTO ChannelNames (name) VALUES
('Go Concurrency Deep Dive'),
('CSS Advanced Techniques'),
('C# for Beginners'),
('Java Memory Management'),
('JavaScript Event Loop Explained'),
('Rust Ownership Patterns'),
('Python Data Science Hub'),
('Node.js Performance Tuning'),
('React State Management'),
('Vue Composition API'),
('Machine Learning in Practice'),
('Docker and Containerization'),
('Kubernetes for Devs'),
('PostgreSQL Performance Tips'),
('SQLite for Mobile Apps'),
('Microservices Architecture'),
('API Security Fundamentals'),
('OAuth2 and Identity'),
('Web Accessibility Matters'),
('TypeScript Type Systems'),
('DevOps CI/CD Strategies'),
('Clean Code Principles'),
('Design Patterns in Go'),
('Functional Programming in JS'),
('SvelteKit for Web Apps');

-- Create 25 channels with static names
-- Ensure OwnerID comes from the limited set in TempUsersForChannels
INSERT INTO Channels (OwnerID, Name, Avatar, Banner, Description, Created, Privacy, IsFlagged, IsMuted)
SELECT
  (SELECT UserID FROM TempUsersForChannels ORDER BY RANDOM() LIMIT 1), -- Select a random owner from the cached users
  cn.name,
  'noimage',
  'default.png',
  'Community focused on ' || cn.name || '. Join discussions, share ideas, and collaborate with fellow tech enthusiasts.',
  -- Generate a staggered creation date for channels within the last year, capped by current date
  MIN(datetime('now', '-' || (25 - ROWID) * 10 || ' days'), CURRENT_TIMESTAMP),
  (ROWID % 2), -- Alternating privacy
  (ROWID % 2), -- Alternating IsFlagged
  (ROWID % 2) -- Alternating IsMuted
FROM ChannelNames cn;

-- Create a temporary table to hold the IDs and names of the newly created channels
-- This is crucial for dynamically linking posts to these specific channels.
CREATE TEMP TABLE TempNewChannelIDs AS
SELECT
    ch.ID AS channel_id,
    ch.Name AS channel_name,
    -- Assign a row number from 1 to 25, ordered by the actual ID of the new channels.
    -- This allows us to cycle through them predictably.
    ROW_NUMBER() OVER (ORDER BY ch.ID ASC) AS rn
FROM Channels ch
WHERE ch.Name IN (SELECT name FROM ChannelNames); -- Filter to ensure we only get the 25 channels just created

-- Insert memberships: Each channel gets the owner + up to 2 other random users from TempUsersForChannels
-- This part might need adjustment if TempUsersForChannels has very few users (e.g., 1 or 2)
-- For simplicity, let's ensure each user in TempUsersForChannels becomes a member of a few channels.
INSERT INTO Memberships (UserID, ChannelID, Created)
SELECT
    tu.UserID,
    tnc.channel_id,
    MIN(datetime(c.Created, '+' || (ABS(RANDOM()) % 24) || ' hours'), CURRENT_TIMESTAMP)
FROM TempNewChannelIDs tnc
JOIN Channels c ON tnc.channel_id = c.ID -- To get channel's creation time
CROSS JOIN TempUsersForChannels tu        -- Each user from the temp table
WHERE NOT EXISTS ( -- Ensure user is not already the owner (implicitly handled if OwnerID is also in TempUsersForChannels)
    SELECT 1 FROM Channels ch_owner
    WHERE ch_owner.ID = tnc.channel_id AND ch_owner.OwnerID = tu.UserID
)
GROUP BY tu.UserID, tnc.channel_id -- Avoid duplicate memberships if RANDOM picks same user/channel
ORDER BY RANDOM()
LIMIT 50; -- Create up to 50 additional memberships

-- Create a temporary table to cache 5 users for use in post generation.
-- This is the second TempUsers table as in the original script, intended for post authors.
CREATE TEMP TABLE TempUsersForPosts AS
SELECT ID AS UserID, Username, Avatar FROM Users ORDER BY RANDOM() LIMIT 5;

-- Prepare indexed users by assigning a row number to each user in the temporary table.
-- This allows efficient mapping of users to posts in subsequent queries.
WITH IndexedUsers AS (
  SELECT UserID, Username, Avatar, ROW_NUMBER() OVER () AS rn
  FROM TempUsersForPosts -- Use the correct temp table for post authors
),
-- Generate a sequence of numbers from 1 to 1000 using a recursive CTE.
-- This sequence is used to create multiple posts for testing purposes.
Counter(x) AS (
  SELECT 1
  UNION ALL
  SELECT x + 1 FROM Counter WHERE x < 1000 -- Generate 1000 numbers
)
-- Insert posts, linking them to the newly created channels dynamically
INSERT INTO Posts (Title, Content, Images, Created, IsCommentable, IsFlagged, Author, AuthorID, AuthorAvatar)
SELECT
  -- Modified title to make 'x' parsing safer: "Discussion: [Channel Name] PostNum:X"
  'Discussion: [' || tnc.channel_name || '] PostNum:' || c.x AS Title,
  'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque vel sem eget justo consequat convallis. Integer porta purus at egestas tincidunt.

   Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Curabitur id nunc id nulla dapibus fermentum.',
  'noimage',
  -- This calculates a pseudo-random date within the year 2025, capped at the current date.
  MIN(datetime('2025-01-01', '+' || (c.x % 365) || ' days'), CURRENT_TIMESTAMP),
  (c.x % 2), -- Alternating IsCommentable
  0, -- IsFlagged set to 0
  iu.Username,
  iu.UserID,
  iu.Avatar
FROM Counter c
-- Join with TempNewChannelIDs to pick one of the 25 newly created channels
-- The 'rn' in TempNewChannelIDs goes from 1 to 25.
JOIN TempNewChannelIDs tnc ON tnc.rn = (((c.x - 1) % 25) + 1)
-- Join with IndexedUsers to pick an author for the post
JOIN IndexedUsers iu ON iu.rn = (((c.x - 1) % 5) + 1);

-- Link each post to its respective channel in PostChannels
-- This uses the Title of the post to extract 'x' and determine the correct channel.
INSERT INTO PostChannels (PostID, ChannelID, Created)
SELECT
  p.ID AS PostID,
  tnc.channel_id AS ChannelID,
  MIN(p.Created, CURRENT_TIMESTAMP) AS Created
FROM Posts p
-- Parse 'x' from the post title. Example Title: "Discussion: [Some Channel] PostNum:123"
-- 1. Find 'PostNum:'
-- 2. Get the substring after 'PostNum:'
-- 3. Convert to INTEGER
JOIN TempNewChannelIDs tnc ON tnc.rn = (
  (
    ( CAST( SUBSTR(p.Title, INSTR(p.Title, 'PostNum:') + LENGTH('PostNum:')) AS INTEGER) - 1) % 25
  ) + 1
)
-- Process only posts created in this batch, identifiable by their title format.
-- This WHERE clause is important if the Posts table might contain other posts.
WHERE p.Title LIKE 'Discussion: [%] PostNum:%'
AND p.ID NOT IN (SELECT PostID FROM PostChannels WHERE PostID = p.ID); -- Avoid re-inserting if script is run multiple times on same posts

COMMIT;
