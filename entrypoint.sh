#!/bin/sh
set -e

# Path to the database
DB_PATH="/var/lib/db-codex/forum_database.db"

# Insert the user into the Users table
sqlite3 "$DB_PATH" <<EOF
INSERT INTO Users (Username, EmailAddress, Avatar, Banner, Description, Usertype, Created, IsFlagged, SessionToken, CsrfToken, HashedPassword)
VALUES (
    'TheCodexDonkey',
    'donkey@codex.com',
    'K1xcM1YYVT6xkHoEUICS4g==.png',
    '',
    "I'm such a friendly donkey, and I'm here to show you around the wonderful world of Codex. I've already added you to my channel, where you'll find updates, information and any cool stuff I feel like sharing. Welcome to Codex!",
    0,
    DateTime('now'),
    0,
    '',
    '',
    '$2a$14$qK2P4N12utI8c4dPS6AMaueafDJygKtdVHVLgNVq2wJM5MW5xjdVm'
);
EOF

# Insert the channel into the Channels table
sqlite3 "$DB_PATH" <<EOF
INSERT INTO Channels (Name, Avatar, Description, Created, Privacy, Banner, OwnerID, IsMuted, IsFlagged)
VALUES (
    'WelcomeToCodex',
    '-ou7VWQ7pK1JHiOvK7lJiw==.png',
    "Welcome to Codex! This channel will guide you through the forum, and give you updates, insights, and generally keep you up to date with everything that's going on.",
    DateTime('now'),
    0,
    '',
    1,
    0,
    0
);
EOF

# Insert the post into the posts table
sqlite3 "$DB_PATH" <<EOF
INSERT INTO Posts (Title, Content, Images, Created, IsCommentable, Author, AuthorID, AuthorAvatar, IsFlagged)
VALUES (
  'Welcome to Codex!',
  "We are so glad you could join us! Since you're here, you probably already know this, but just in case, this is a place for us all to share everything we know about coding - our hopes, fears, plans, anxieties... Everything that makes up this world in which we find ourselves. Hopefully, Codex will make it less daunting, and a world we enjoy living in. So, welcome! We're so glad you're here.",
  'noimage',
  DateTime('now'),
  0,
  'TheCodexDonkey',
  1,
  'K12cM1YYVT6xkHoEUICS4g==.png',
  0
);
EOF

# Attach the post to the channel
sqlite3 "$DB_PATH" <<EOF
INSERT INTO PostChannels (PostID, ChannelID, Created)
VALUES (
    1,
    1,
    DateTime('now')
);
EOF

# Continue with the original command
exec "$@"
