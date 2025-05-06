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
  -- MODIFIED: Use printf for robust date string formatting and cn.ROWID for explicitness.
  -- This calculates a date in the past, staggered based on cn.ROWID.
  -- (25 - cn.ROWID) * 10 results in offsets from 0 (for ROWID=25) to 240 (for ROWID=1).
  datetime('now', printf("-%d days", (25 - cn.ROWID) * 10)),
  (cn.ROWID % 2), -- Alternating privacy, using cn.ROWID for clarity
  ((cn.ROWID + 1) % 2), -- Alternating IsFlagged, using cn.ROWID for clarity
  ((cn.ROWID + 2) % 2) -- Alternating IsMuted, using cn.ROWID for clarity
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
-- Ensure we are selecting from the channels just created by matching their names.
-- This assumes channel names are unique for this batch.
WHERE ch.Name IN (SELECT name FROM ChannelNames)
-- Additionally, to be more robust, filter by creation time if possible,
-- though this might be tricky if other channels could be created concurrently.
-- For this script, Name matching should be sufficient.
ORDER BY ch.ID DESC LIMIT 25; -- Get the last 25 inserted matching names, then re-order by ID ASC for ROW_NUMBER

-- Re-create TempNewChannelIDs with correct ordering for rn after ensuring we got the right 25 channels
CREATE TEMP TABLE TempNewChannelIDsOrdered AS
SELECT channel_id, channel_name, ROW_NUMBER() OVER (ORDER BY channel_id ASC) as rn
FROM TempNewChannelIDs;
DROP TABLE TempNewChannelIDs;
ALTER TABLE TempNewChannelIDsOrdered RENAME TO TempNewChannelIDs;


-- Insert memberships: Each channel gets the owner + up to 2 other random users from TempUsersForChannels
INSERT INTO Memberships (UserID, ChannelID, Created)
SELECT
    tu.UserID,
    tnc.channel_id,
    -- Get the actual creation time of the channel for membership creation time
    MIN(datetime(ch.Created, '+' || (ABS(RANDOM()) % 24) || ' hours'), CURRENT_TIMESTAMP)
FROM TempNewChannelIDs tnc
JOIN Channels ch ON tnc.channel_id = ch.ID -- Join to get actual channel creation time
CROSS JOIN TempUsersForChannels tu
WHERE NOT EXISTS (
    SELECT 1 FROM Channels ch_owner
    WHERE ch_owner.ID = tnc.channel_id AND ch_owner.OwnerID = tu.UserID
)
GROUP BY tu.UserID, tnc.channel_id
ORDER BY RANDOM()
LIMIT 50; -- Create up to 50 additional memberships, ensuring variety

-- Create a temporary table to cache 5 users for use in post generation.
CREATE TEMP TABLE TempUsersForPosts AS
SELECT ID AS UserID, Username, Avatar FROM Users ORDER BY RANDOM() LIMIT 5;

-- Prepare indexed users by assigning a row number to each user in the temporary table.
WITH IndexedUsers AS (
  SELECT UserID, Username, Avatar, ROW_NUMBER() OVER () AS rn
  FROM TempUsersForPosts
),
-- Generate a sequence of numbers from 1 to 1000 using a recursive CTE.
Counter(x) AS (
  SELECT 1
  UNION ALL
  SELECT x + 1 FROM Counter WHERE x < 1000 -- Generate 1000 numbers (ensure DB supports recursion depth)
)
-- Insert posts, linking them to the newly created channels dynamically
INSERT INTO Posts (Title, Content, Images, Created, IsCommentable, IsFlagged, Author, AuthorID, AuthorAvatar)
SELECT
  'Discussion: [' || tnc.channel_name || '] PostNum:' || c.x AS Title,
  'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque vel sem eget justo consequat convallis. Integer porta purus at egestas tincidunt.

   Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Curabitur id nunc id nulla dapibus fermentum.',
  'noimage',
  MIN(datetime('2025-01-01', '+' || (c.x % 365) || ' days'), CURRENT_TIMESTAMP),
  (c.x % 2),
  ((c.x + 1) % 2),
  iu.Username,
  iu.UserID,
  iu.Avatar
FROM Counter c
JOIN TempNewChannelIDs tnc ON tnc.rn = (((c.x - 1) % 25) + 1)
JOIN IndexedUsers iu ON iu.rn = (((c.x - 1) % 5) + 1);

-- Link each post to its respective channel in PostChannels
INSERT INTO PostChannels (PostID, ChannelID, Created)
SELECT
  p.ID AS PostID,
  tnc.channel_id AS ChannelID,
  MIN(p.Created, CURRENT_TIMESTAMP) AS Created
FROM Posts p
JOIN TempNewChannelIDs tnc ON tnc.rn = (
  (
    ( CAST( SUBSTR(p.Title, INSTR(p.Title, 'PostNum:') + LENGTH('PostNum:')) AS INTEGER) - 1) % 25
  ) + 1
)
WHERE p.Title LIKE 'Discussion: [%] PostNum:%'
AND p.ID NOT IN (SELECT PostID FROM PostChannels WHERE PostID = p.ID);

COMMIT;
