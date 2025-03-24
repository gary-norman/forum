
CREATE TABLE Posts (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Title TEXT NOT NULL,
  Content TEXT NOT NULL,
  Images TEXT,
  Created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  IsCommentable BOOLEAN NOT NULL,
  Author TEXT NOT NULL,
  AuthorID INTEGER NOT NULL,
  AuthorAvatar TEXT,
  IsFlagged BOOLEAN,
  FOREIGN KEY (Author) REFERENCES Users(Username),
  FOREIGN KEY (AuthorID) REFERENCES Users(ID),
  FOREIGN KEY (AuthorAvatar) REFERENCES Users(Avatar),
  FOREIGN KEY (ChannelName) REFERENCES Channels(Name),
  FOREIGN KEY (ChannelID) REFERENCES Channels(ID)
);

INSERT INTO Posts_New (ID, Title, Content, Images, Created, IsCommentable, Author, AuthorID, AuthorAvatar, IsFlagged)
SELECT ID, Title, Content, Images, Created, IsCommentable, Author, AuthorID, AuthorAvatar, IsFlagged FROM Posts;

DROP TABLE Posts;

ALTER TABLE Posts_New RENAME TO Posts;
