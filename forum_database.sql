SET XACT_ABORT ON

BEGIN TRANSACTION QUICKDBD

CREATE TABLE [Users] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Username] VARCHAR(30)  NOT NULL ,
    [Password] VARCHAR(20)  NOT NULL ,
    [Email_address] VARCHAR(30)  NOT NULL ,
    [Avatar] UUID  NULL ,
    [Banner] UUID  NULL ,
    [Desc] VARCHAR(255)  NULL ,
    [Type] VARCHAR(10)  NOT NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_Users_Created] DEFAULT (getutcdate()),
    [Membership] JSON  NOT NULL ,
    [Followers] JSON  NOT NULL ,
    [Following] JSON  NOT NULL ,
    [Bookmarks_user] JSON  NOT NULL ,
    [Bookmarks_channel] JSON  NOT NULL ,
    [Bookmarks_post] JSON  NOT NULL ,
    [Bookmarks_comment] JSON  NOT NULL ,
    [Posts] JSON  NOT NULL ,
    [Comments] JSON  NOT NULL ,
    [Flagged_users] JSON  NOT NULL ,
    [Flagged_posts] JSON  NOT NULL ,
    [Flagged_comments] JSON  NOT NULL ,
    [Mod_of] JSON  NOT NULL ,
    [Reactions] JSON  NOT NULL ,
    [Is_flagged] JSON  NOT NULL ,
    [Is_flagged_approved] JSON  NOT NULL ,
    CONSTRAINT [PK_Users] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

CREATE TABLE [Channels] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Name] VARCHAR(30)  NOT NULL ,
    [Avatar] UUID  NULL ,
    [Desc] VARCHAR(555)  NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_Channels_Created] DEFAULT (getutcdate()),
    [Rules] VARCHAR  NULL ,
    [Privacy] BOOL  NOT NULL ,
    [Members] JSON  NOT NULL ,
    [Mods] JSON  NOT NULL ,
    [Posts] JSON  NOT NULL ,
    CONSTRAINT [PK_Channels] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

CREATE TABLE [Posts] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Title] VARCHAR(80)  NOT NULL ,
    [Content] VARCHAR  NOT NULL ,
    [Images] JSON  NOT NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_Posts_Created] DEFAULT (getutcdate()),
    [Commentable] BOOL  NOT NULL ,
    [Author] INT  NOT NULL ,
    [Reactions] JSON  NOT NULL ,
    [Comments] JSON  NOT NULL ,
    [Channel] INT  NOT NULL ,
    [Is_flagged] JSON  NOT NULL ,
    [Is_flagged_approved] JSON  NOT NULL ,
    CONSTRAINT [PK_Posts] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

CREATE TABLE [Comments] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Content] VARCHAR  NOT NULL ,
    [Images] JSON  NOT NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_Comments_Created] DEFAULT (getutcdate()),
    [Author] INT  NOT NULL ,
    [Reactions] JSON  NOT NULL ,
    [Replies] JSON  NOT NULL ,
    [Channel] INT  NOT NULL ,
    [Is_flagged] JSON  NOT NULL ,
    [Is_flagged_approved] JSON  NOT NULL ,
    [Commented_post] INT  NOT NULL ,
    [Commented_comment] INT  NOT NULL ,
    CONSTRAINT [PK_Comments] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

CREATE TABLE [Reactions] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Liked] BOOL  NOT NULL ,
    [Disliked] BOOL  NOT NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_Reactions_Created] DEFAULT (getutcdate()),
    [ParentID] JSON  NOT NULL ,
    [Author] INT  NOT NULL ,
    [Channel] INT  NOT NULL ,
    [Reacted_post] INT  NOT NULL ,
    [Reacted_comment] INT  NOT NULL ,
    [Reacted_post_author] INT  NOT NULL ,
    [Reacted_comment_author] INT  NOT NULL ,
    CONSTRAINT [PK_Reactions] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

CREATE TABLE [Flags] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Type] VARCHAR  NOT NULL ,
    [Content] VARCHAR(160)  NOT NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_Flags_Created] DEFAULT (getutcdate()),
    [Approved] BOOL  NOT NULL ,
    [Author] INT  NOT NULL ,
    [Flagged_user] INT  NOT NULL ,
    [Flagged_post] INT  NOT NULL ,
    [Flagged_comment] INT  NOT NULL ,
    [Flagged_post_author] INT  NOT NULL ,
    [Flagged_comment_author] INT  NOT NULL ,
    CONSTRAINT [PK_Flags] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Membership] FOREIGN KEY([Membership])
REFERENCES [Channels] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Membership]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Followers] FOREIGN KEY([Followers])
REFERENCES [Users] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Followers]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Following] FOREIGN KEY([Following])
REFERENCES [Users] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Following]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Bookmarks_user] FOREIGN KEY([Bookmarks_user])
REFERENCES [Users] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Bookmarks_user]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Bookmarks_channel] FOREIGN KEY([Bookmarks_channel])
REFERENCES [Channels] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Bookmarks_channel]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Bookmarks_post] FOREIGN KEY([Bookmarks_post])
REFERENCES [Posts] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Bookmarks_post]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Bookmarks_comment] FOREIGN KEY([Bookmarks_comment])
REFERENCES [Comments] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Bookmarks_comment]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Posts] FOREIGN KEY([Posts])
REFERENCES [Posts] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Posts]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Comments] FOREIGN KEY([Comments])
REFERENCES [Comments] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Comments]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Flagged_users_Is_flagged_approved] FOREIGN KEY([Flagged_users], [Is_flagged_approved])
REFERENCES [Flags] ([ID], [Approved])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Flagged_users_Is_flagged_approved]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Flagged_posts] FOREIGN KEY([Flagged_posts])
REFERENCES [Flags] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Flagged_posts]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Flagged_comments] FOREIGN KEY([Flagged_comments])
REFERENCES [Flags] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Flagged_comments]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Mod_of] FOREIGN KEY([Mod_of])
REFERENCES [Channels] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Mod_of]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Reactions] FOREIGN KEY([Reactions])
REFERENCES [Reactions] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Reactions]

ALTER TABLE [Users] WITH CHECK ADD CONSTRAINT [FK_Users_Is_flagged] FOREIGN KEY([Is_flagged])
REFERENCES [Flags] ([ID])

ALTER TABLE [Users] CHECK CONSTRAINT [FK_Users_Is_flagged]

ALTER TABLE [Channels] WITH CHECK ADD CONSTRAINT [FK_Channels_Members] FOREIGN KEY([Members])
REFERENCES [Users] ([ID])

ALTER TABLE [Channels] CHECK CONSTRAINT [FK_Channels_Members]

ALTER TABLE [Channels] WITH CHECK ADD CONSTRAINT [FK_Channels_Mods] FOREIGN KEY([Mods])
REFERENCES [Users] ([ID])

ALTER TABLE [Channels] CHECK CONSTRAINT [FK_Channels_Mods]

ALTER TABLE [Channels] WITH CHECK ADD CONSTRAINT [FK_Channels_Posts] FOREIGN KEY([Posts])
REFERENCES [Posts] ([ID])

ALTER TABLE [Channels] CHECK CONSTRAINT [FK_Channels_Posts]

ALTER TABLE [Posts] WITH CHECK ADD CONSTRAINT [FK_Posts_Author] FOREIGN KEY([Author])
REFERENCES [Users] ([ID])

ALTER TABLE [Posts] CHECK CONSTRAINT [FK_Posts_Author]

ALTER TABLE [Posts] WITH CHECK ADD CONSTRAINT [FK_Posts_Reactions] FOREIGN KEY([Reactions])
REFERENCES [Reactions] ([ID])

ALTER TABLE [Posts] CHECK CONSTRAINT [FK_Posts_Reactions]

ALTER TABLE [Posts] WITH CHECK ADD CONSTRAINT [FK_Posts_Comments] FOREIGN KEY([Comments])
REFERENCES [Comments] ([ID])

ALTER TABLE [Posts] CHECK CONSTRAINT [FK_Posts_Comments]

ALTER TABLE [Posts] WITH CHECK ADD CONSTRAINT [FK_Posts_Channel] FOREIGN KEY([Channel])
REFERENCES [Channels] ([ID])

ALTER TABLE [Posts] CHECK CONSTRAINT [FK_Posts_Channel]

ALTER TABLE [Posts] WITH CHECK ADD CONSTRAINT [FK_Posts_Is_flagged_Is_flagged_approved] FOREIGN KEY([Is_flagged], [Is_flagged_approved])
REFERENCES [Flags] ([ID], [Approved])

ALTER TABLE [Posts] CHECK CONSTRAINT [FK_Posts_Is_flagged_Is_flagged_approved]

ALTER TABLE [Comments] WITH CHECK ADD CONSTRAINT [FK_Comments_Author] FOREIGN KEY([Author])
REFERENCES [Users] ([ID])

ALTER TABLE [Comments] CHECK CONSTRAINT [FK_Comments_Author]

ALTER TABLE [Comments] WITH CHECK ADD CONSTRAINT [FK_Comments_Reactions] FOREIGN KEY([Reactions])
REFERENCES [Reactions] ([ID])

ALTER TABLE [Comments] CHECK CONSTRAINT [FK_Comments_Reactions]

ALTER TABLE [Comments] WITH CHECK ADD CONSTRAINT [FK_Comments_Replies] FOREIGN KEY([Replies])
REFERENCES [Comments] ([ID])

ALTER TABLE [Comments] CHECK CONSTRAINT [FK_Comments_Replies]

ALTER TABLE [Comments] WITH CHECK ADD CONSTRAINT [FK_Comments_Channel] FOREIGN KEY([Channel])
REFERENCES [Channels] ([ID])

ALTER TABLE [Comments] CHECK CONSTRAINT [FK_Comments_Channel]

ALTER TABLE [Comments] WITH CHECK ADD CONSTRAINT [FK_Comments_Is_flagged_Is_flagged_approved] FOREIGN KEY([Is_flagged], [Is_flagged_approved])
REFERENCES [Flags] ([ID], [Approved])

ALTER TABLE [Comments] CHECK CONSTRAINT [FK_Comments_Is_flagged_Is_flagged_approved]

ALTER TABLE [Comments] WITH CHECK ADD CONSTRAINT [FK_Comments_Commented_post] FOREIGN KEY([Commented_post])
REFERENCES [Posts] ([ID])

ALTER TABLE [Comments] CHECK CONSTRAINT [FK_Comments_Commented_post]

ALTER TABLE [Comments] WITH CHECK ADD CONSTRAINT [FK_Comments_Commented_comment] FOREIGN KEY([Commented_comment])
REFERENCES [Comments] ([ID])

ALTER TABLE [Comments] CHECK CONSTRAINT [FK_Comments_Commented_comment]

ALTER TABLE [Reactions] WITH CHECK ADD CONSTRAINT [FK_Reactions_Author] FOREIGN KEY([Author])
REFERENCES [Users] ([ID])

ALTER TABLE [Reactions] CHECK CONSTRAINT [FK_Reactions_Author]

ALTER TABLE [Reactions] WITH CHECK ADD CONSTRAINT [FK_Reactions_Channel] FOREIGN KEY([Channel])
REFERENCES [Channels] ([ID])

ALTER TABLE [Reactions] CHECK CONSTRAINT [FK_Reactions_Channel]

ALTER TABLE [Reactions] WITH CHECK ADD CONSTRAINT [FK_Reactions_Reacted_post_Reacted_post_author] FOREIGN KEY([Reacted_post], [Reacted_post_author])
REFERENCES [Posts] ([ID], [Author])

ALTER TABLE [Reactions] CHECK CONSTRAINT [FK_Reactions_Reacted_post_Reacted_post_author]

ALTER TABLE [Reactions] WITH CHECK ADD CONSTRAINT [FK_Reactions_Reacted_comment_Reacted_comment_author] FOREIGN KEY([Reacted_comment], [Reacted_comment_author])
REFERENCES [Comments] ([ID], [Author])

ALTER TABLE [Reactions] CHECK CONSTRAINT [FK_Reactions_Reacted_comment_Reacted_comment_author]

ALTER TABLE [Flags] WITH CHECK ADD CONSTRAINT [FK_Flags_Author] FOREIGN KEY([Author])
REFERENCES [Users] ([ID])

ALTER TABLE [Flags] CHECK CONSTRAINT [FK_Flags_Author]

ALTER TABLE [Flags] WITH CHECK ADD CONSTRAINT [FK_Flags_Flagged_user] FOREIGN KEY([Flagged_user])
REFERENCES [Users] ([ID])

ALTER TABLE [Flags] CHECK CONSTRAINT [FK_Flags_Flagged_user]

ALTER TABLE [Flags] WITH CHECK ADD CONSTRAINT [FK_Flags_Flagged_post_Flagged_post_author] FOREIGN KEY([Flagged_post], [Flagged_post_author])
REFERENCES [Posts] ([ID], [Author])

ALTER TABLE [Flags] CHECK CONSTRAINT [FK_Flags_Flagged_post_Flagged_post_author]

ALTER TABLE [Flags] WITH CHECK ADD CONSTRAINT [FK_Flags_Flagged_comment_Flagged_comment_author] FOREIGN KEY([Flagged_comment], [Flagged_comment_author])
REFERENCES [Comments] ([ID], [Author])

ALTER TABLE [Flags] CHECK CONSTRAINT [FK_Flags_Flagged_comment_Flagged_comment_author]

CREATE INDEX [idx_Users_Username]
ON [Users] ([Username])

CREATE INDEX [idx_Channels_Name]
ON [Channels] ([Name])

COMMIT TRANSACTION QUICKDBD
