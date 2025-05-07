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
WHERE ch.Name IN (SELECT name FROM Channels)
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

.width 20 10

-- View number of channels per user
SELECT
  u.Username,
  COUNT(m.ChannelID) AS Channels
FROM Memberships m
JOIN Users u ON u.ID = m.UserID
GROUP BY m.UserID
ORDER BY Channels DESC;

.width 35 10
-- view number of posts per channel
SELECT
  ch.Name AS ChannelName,
  COUNT(pc.PostID) AS Posts
FROM PostChannels pc
JOIN Channels ch ON ch.ID = pc.ChannelID
GROUP BY ch.ID
ORDER BY Posts DESC;

.width 14 8

-- view number of posts per user
SELECT
  u.Username,
  COUNT(p.ID) AS Posts
FROM Posts p
JOIN Users u ON u.ID = p.AuthorID
GROUP BY p.AuthorID
ORDER BY Posts DESC;

-- view stats
SELECT * FROM Stats;

-- clean up
DROP TABLE IF EXISTS Stats;
