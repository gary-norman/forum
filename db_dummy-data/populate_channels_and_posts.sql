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

WITH IndexedChannelNames AS (
  SELECT name, ROW_NUMBER() OVER () AS rn
  FROM ChannelNames
),
IndexedTempUsers AS (
  SELECT UserID, Username, Avatar, ROW_NUMBER() OVER () AS rn
  FROM TempUsersForChannels
)
INSERT INTO Channels (OwnerID, Name, Avatar, Banner, Description, Created, Privacy, IsFlagged, IsMuted)
SELECT
  tu.UserID,
  cn.name,
  'noimage',
  'default.png',
  'Community focused on ' || cn.name || '. Join discussions, share ideas, and collaborate with fellow tech enthusiasts.',
  datetime('now', printf("-%d days", (25 - cn.rn) * 10)),
  (cn.rn % 2),
  ((cn.rn + 1) % 2),
  ((cn.rn + 2) % 2)
FROM IndexedChannelNames cn
JOIN IndexedTempUsers tu ON tu.rn = ((cn.rn - 1) % 5) + 1;

-- Create a temporary table to hold the IDs and names of the newly created channels
-- This is crucial for dynamically linking posts to these specific channels.
CREATE TEMP TABLE TempNewChannelIDs AS
SELECT
    ch.ID AS channel_id,
    ch.Name AS channel_name,
  ch.OwnerID AS channel_owner,
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
SELECT channel_id, channel_name, channel_owner, ROW_NUMBER() OVER (ORDER BY channel_id ASC) as rn
FROM TempNewChannelIDs;
DROP TABLE TempNewChannelIDs;
ALTER TABLE TempNewChannelIDsOrdered RENAME TO TempNewChannelIDs;

-- ******* Memberships ********

-- Step 4: Insert owner memberships
INSERT INTO Memberships (UserID, ChannelID, Created)
SELECT
  channel_owner,
  channel_id,
  datetime('now')
FROM TempNewChannelIDs;

-- Step 5: Add ONE non-owner member per channel (randomly chosen)
WITH PossibleMembers AS (
  SELECT
    tnc.channel_id,
    tu.UserID,
    tnc.channel_owner,
    ROW_NUMBER() OVER (PARTITION BY tnc.channel_id ORDER BY RANDOM()) AS rn,
    tnc.channel_id AS ChannelID,
    tnc.channel_name,
    ch.Created
  FROM TempNewChannelIDs tnc
  JOIN TempUsersForChannels tu ON tu.UserID != tnc.channel_owner
  JOIN Channels ch ON ch.ID = tnc.channel_id
)
INSERT INTO Memberships (UserID, ChannelID, Created)
SELECT
  UserID,
  ChannelID,
  datetime(Created, '+' || (ABS(RANDOM()) % 24) || ' hours')
FROM PossibleMembers
WHERE rn = 1;

-- ******** Posts *********

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
  CASE
  WHEN tnc.channel_name LIKE '%Go%' THEN 'Best Practices in ' || tnc.channel_name
  WHEN tnc.channel_name LIKE '%CSS%' THEN 'Modern Layouts in ' || tnc.channel_name
  WHEN tnc.channel_name LIKE '%C#%' THEN 'Getting Started with ' || tnc.channel_name
  WHEN tnc.channel_name LIKE '%Java Memory%' THEN 'Tuning Garbage Collection in Java'
  WHEN tnc.channel_name LIKE '%JavaScript%' THEN 'Why ' || tnc.channel_name || ' Still Matters'
  WHEN tnc.channel_name LIKE '%Rust%' THEN 'Understanding Ownership in Rust'
  WHEN tnc.channel_name LIKE '%Python%' THEN 'Data Analysis Tips with Python'
  WHEN tnc.channel_name LIKE '%Node.js%' THEN 'Speeding Up ' || tnc.channel_name
  WHEN tnc.channel_name LIKE '%React%' THEN 'Managing State in Large React Apps'
  WHEN tnc.channel_name LIKE '%Vue%' THEN 'Composition API: A Game Changer'
  WHEN tnc.channel_name LIKE '%Machine Learning%' THEN 'Real-World ML Workflows'
  WHEN tnc.channel_name LIKE '%Docker%' THEN 'Lightweight Container Patterns'
  WHEN tnc.channel_name LIKE '%Kubernetes%' THEN 'K8s for Daily DevOps Tasks'
  WHEN tnc.channel_name LIKE '%PostgreSQL%' THEN 'Optimizing Indexes in PostgreSQL'
  WHEN tnc.channel_name LIKE '%SQLite%' THEN 'Mobile-Friendly SQL with SQLite'
  WHEN tnc.channel_name LIKE '%Microservices%' THEN 'Service Boundaries & Data Ownership'
  WHEN tnc.channel_name LIKE '%API Security%' THEN 'Avoiding Common API Security Pitfalls'
  WHEN tnc.channel_name LIKE '%OAuth2%' THEN 'Token Lifecycles & Refresh Logic'
  WHEN tnc.channel_name LIKE '%Accessibility%' THEN 'Building Inclusive Interfaces'
  WHEN tnc.channel_name LIKE '%TypeScript%' THEN 'Advanced TypeScript Tricks'
  WHEN tnc.channel_name LIKE '%CI/CD%' THEN 'Faster Pipelines with CI/CD'
  WHEN tnc.channel_name LIKE '%Clean Code%' THEN 'Refactoring for Readability'
  WHEN tnc.channel_name LIKE '%Design Patterns%' THEN 'Applying Patterns in Go'
  WHEN tnc.channel_name LIKE '%Functional%' THEN 'Pure Functions in JavaScript'
  WHEN tnc.channel_name LIKE '%SvelteKit%' THEN 'Why SvelteKit Feels Instant'
  ELSE 'Topic Discussion: ' || tnc.channel_name
END AS Title,

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

-- Create a temp view that assigns row numbers to posts for deterministic mapping
CREATE TEMP VIEW IF NOT EXISTS IndexedPosts AS
SELECT ID AS PostID, ROW_NUMBER() OVER (ORDER BY ID ASC) AS rn FROM Posts;

-- Now map each post to one of the 25 channels using modulo logic
WITH MatchedPosts AS (
  SELECT
    ip.PostID,
    tnc.channel_id,
    MIN(p.Created, CURRENT_TIMESTAMP) AS Created
  FROM IndexedPosts ip
  JOIN Posts p ON p.ID = ip.PostID
  JOIN TempNewChannelIDs tnc ON tnc.rn = ((ip.rn - 1) % 25) + 1
  LEFT JOIN PostChannels pc ON pc.PostID = ip.PostID
  WHERE pc.PostID IS NULL
)

INSERT OR IGNORE INTO PostChannels (PostID, ChannelID, Created)
SELECT * FROM MatchedPosts;

.mode box
.headers on
.width 35 20 10

SELECT
  ch.Name AS ChannelName,
  u.Username AS MemberUsername,
  CASE
    WHEN u.ID = ch.OwnerID THEN 'Owner'
    ELSE 'Member'
  END AS Role
FROM Channels ch
JOIN Memberships m ON m.ChannelID = ch.ID
JOIN Users u ON u.ID = m.UserID
WHERE ch.Name IN (SELECT name FROM ChannelNames)
ORDER BY ch.ID, Role DESC;


CREATE TABLE IF NOT EXISTS Stats (
  Name TEXT PRIMARY KEY,
  Value INTEGER NOT NULL
);

INSERT INTO Stats (Name, Value)
VALUES
  ('Channels',      (SELECT COUNT(*) FROM Channels)),
  ('Memberships',   (SELECT COUNT(*) FROM Memberships)),
  ('Posts',         (SELECT COUNT(*) FROM Posts)),
  ('PostChannels',  (SELECT COUNT(*) FROM PostChannels));

.mode box
.headers on
.width 20 20

-- View number of channels per user
SELECT
  u.Username,
  COUNT(m.ChannelID) AS ChannelMemberships
FROM Memberships m
JOIN Users u ON u.ID = m.UserID
GROUP BY m.UserID
ORDER BY ChannelMemberships DESC;

.width 14 8

-- view stats
SELECT * FROM Stats;

-- Drop temporary tables and views used for seeding
DROP TABLE IF EXISTS TempUsersForChannels;
DROP TABLE IF EXISTS TempUsersForPosts;
DROP TABLE IF EXISTS TempNewChannelIDs;
DROP TABLE IF EXISTS ChannelNames;
DROP TABLE IF EXISTS Stats;
DROP VIEW IF EXISTS IndexedPosts;

COMMIT;
