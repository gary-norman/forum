-- Exported from QuickDBD: https://www.quickdatabasediagrams.com/
-- NOTE! If you have used non-SQL datatypes in your design, you will have to change these here.

-- Modify this code to update the DB schema diagram.
-- To reset the sample schema, replace everything with
-- two dots ('..' - without quotes).

SET XACT_ABORT ON

BEGIN TRANSACTION QUICKDBD

CREATE TABLE [user] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Username] VARCHAR(30)  NOT NULL ,
    [Password] VARCHAR(20)  NOT NULL ,
    [Email_address] VARCHAR(30)  NOT NULL ,
    [Avatar] UUID  NULL ,
    [Banner] UUID  NULL ,
    [Desc] VARCHAR(255)  NULL ,
    [Type] VARCHAR(10)  NOT NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_user_Created] DEFAULT (getutcdate()),
    [Membership] JSON  NOT NULL ,
    [Followers] JSON  NOT NULL ,
    [Following] JSON  NOT NULL ,
    [Bookmarks_user] JSON  NOT NULL ,
    [Bookmarks_channel] JSON  NOT NULL ,
    [Bookmarks_post] JSON  NOT NULL ,
    [Bookmarks_comment] JSON  NOT NULL ,
    [Posts] JSON  NOT NULL ,
    [Comments] JSON  NOT NULL ,
    [Flagged_items] JSON  NOT NULL ,
    [Mod_of] JSON  NOT NULL ,
    [Reactions] JSON  NOT NULL ,
    [Is_flagged] JSON  NOT NULL ,
    [Is_flagged_approved] JSON  NOT NULL ,
    CONSTRAINT [PK_user] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

CREATE TABLE [channel] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Name] VARCHAR(30)  NOT NULL ,
    [Avatar] UUID  NULL ,
    [Desc] VARCHAR(555)  NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_channel_Created] DEFAULT (getutcdate()),
    [Rules] VARCHAR  NULL ,
    [Privacy] BOOL  NOT NULL ,
    [Members] JSON  NOT NULL ,
    [Mods] JSON  NOT NULL ,
    [Posts] JSON  NOT NULL ,
    CONSTRAINT [PK_channel] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

CREATE TABLE [post] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Title] VARCHAR(80)  NOT NULL ,
    [Content] VARCHAR  NOT NULL ,
    [Images] JSON  NOT NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_post_Created] DEFAULT (getutcdate()),
    [Commentable] BOOL  NOT NULL ,
    [Author] VARCHAR  NOT NULL ,
    [Reactions] JSON  NOT NULL ,
    [Comments] JSON  NOT NULL ,
    [Channel] INT  NOT NULL ,
    [Is_flagged] JSON  NOT NULL ,
    [Is_flagged_approved] JSON  NOT NULL ,
    CONSTRAINT [PK_post] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

CREATE TABLE [comment] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Content] VARCHAR  NOT NULL ,
    [Images] JSON  NOT NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_comment_Created] DEFAULT (getutcdate()),
    [Author] INT  NOT NULL ,
    [Reactions] JSON  NOT NULL ,
    [Replies] JSON  NOT NULL ,
    [Channel] INT  NOT NULL ,
    [Is_flagged] JSON  NOT NULL ,
    [Is_flagged_approved] JSON  NOT NULL ,
    -- figure out how to add multiple
    [Parent_post] INT  NOT NULL ,
    [Parent_comment] INT  NOT NULL ,
    CONSTRAINT [PK_comment] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

CREATE TABLE [reaction] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Liked] BOOL  NOT NULL ,
    [Disliked] BOOL  NOT NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_reaction_Created] DEFAULT (getutcdate()),
    [Author] INT  NOT NULL ,
    [Channel] INT  NOT NULL ,
    -- figure out how to add multiple
    [Parent_post] JSON  NOT NULL ,
    [Parent_comment] JSON  NOT NULL ,
    -- figure out how to add multiple
    [Parent_post_author] INT  NOT NULL ,
    [Parent_comment_author] INT  NOT NULL ,
    CONSTRAINT [PK_reaction] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

CREATE TABLE [flag] (
    [ID] INT IDENTITY(1,1) NOT NULL ,
    [Type] VARCHAR  NOT NULL ,
    [Content] VARCHAR(160)  NOT NULL ,
    [Created] DATETIME  NOT NULL CONSTRAINT [DF_flag_Created] DEFAULT (getutcdate()),
    [Approved] BOOL  NOT NULL ,
    [Author] INT  NOT NULL ,
    -- figure out how to add multiple
    [Parent_post] JSON  NOT NULL ,
    [Parent_comment] JSON  NOT NULL ,
    -- figure out how to add multiple
    [Parent_post_author] INT  NOT NULL ,
    [Parent_comment_author] INT  NOT NULL ,
    CONSTRAINT [PK_flag] PRIMARY KEY CLUSTERED (
        [ID] ASC
    )
)

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Membership] FOREIGN KEY([Membership])
REFERENCES [channel] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Membership]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Followers] FOREIGN KEY([Followers])
REFERENCES [user] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Followers]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Following] FOREIGN KEY([Following])
REFERENCES [user] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Following]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Bookmarks_user] FOREIGN KEY([Bookmarks_user])
REFERENCES [user] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Bookmarks_user]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Bookmarks_channel] FOREIGN KEY([Bookmarks_channel])
REFERENCES [channel] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Bookmarks_channel]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Bookmarks_post] FOREIGN KEY([Bookmarks_post])
REFERENCES [post] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Bookmarks_post]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Bookmarks_comment] FOREIGN KEY([Bookmarks_comment])
REFERENCES [comment] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Bookmarks_comment]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Posts] FOREIGN KEY([Posts])
REFERENCES [post] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Posts]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Comments] FOREIGN KEY([Comments])
REFERENCES [comment] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Comments]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Flagged_items_Is_flagged_approved] FOREIGN KEY([Flagged_items], [Is_flagged_approved])
REFERENCES [flag] ([ID], [Approved])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Flagged_items_Is_flagged_approved]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Mod_of] FOREIGN KEY([Mod_of])
REFERENCES [channel] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Mod_of]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Reactions] FOREIGN KEY([Reactions])
REFERENCES [reaction] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Reactions]

ALTER TABLE [user] WITH CHECK ADD CONSTRAINT [FK_user_Is_flagged] FOREIGN KEY([Is_flagged])
REFERENCES [flag] ([ID])

ALTER TABLE [user] CHECK CONSTRAINT [FK_user_Is_flagged]

ALTER TABLE [channel] WITH CHECK ADD CONSTRAINT [FK_channel_Members] FOREIGN KEY([Members])
REFERENCES [user] ([ID])

ALTER TABLE [channel] CHECK CONSTRAINT [FK_channel_Members]

ALTER TABLE [channel] WITH CHECK ADD CONSTRAINT [FK_channel_Mods] FOREIGN KEY([Mods])
REFERENCES [user] ([ID])

ALTER TABLE [channel] CHECK CONSTRAINT [FK_channel_Mods]

ALTER TABLE [channel] WITH CHECK ADD CONSTRAINT [FK_channel_Posts] FOREIGN KEY([Posts])
REFERENCES [post] ([ID])

ALTER TABLE [channel] CHECK CONSTRAINT [FK_channel_Posts]

ALTER TABLE [post] WITH CHECK ADD CONSTRAINT [FK_post_Author] FOREIGN KEY([Author])
REFERENCES [user] ([ID])

ALTER TABLE [post] CHECK CONSTRAINT [FK_post_Author]

ALTER TABLE [post] WITH CHECK ADD CONSTRAINT [FK_post_Reactions] FOREIGN KEY([Reactions])
REFERENCES [reaction] ([ID])

ALTER TABLE [post] CHECK CONSTRAINT [FK_post_Reactions]

ALTER TABLE [post] WITH CHECK ADD CONSTRAINT [FK_post_Comments] FOREIGN KEY([Comments])
REFERENCES [comment] ([ID])

ALTER TABLE [post] CHECK CONSTRAINT [FK_post_Comments]

ALTER TABLE [post] WITH CHECK ADD CONSTRAINT [FK_post_Channel] FOREIGN KEY([Channel])
REFERENCES [channel] ([ID])

ALTER TABLE [post] CHECK CONSTRAINT [FK_post_Channel]

ALTER TABLE [post] WITH CHECK ADD CONSTRAINT [FK_post_Is_flagged_Is_flagged_approved] FOREIGN KEY([Is_flagged], [Is_flagged_approved])
REFERENCES [flag] ([ID], [Approved])

ALTER TABLE [post] CHECK CONSTRAINT [FK_post_Is_flagged_Is_flagged_approved]

ALTER TABLE [comment] WITH CHECK ADD CONSTRAINT [FK_comment_Author] FOREIGN KEY([Author])
REFERENCES [user] ([ID])

ALTER TABLE [comment] CHECK CONSTRAINT [FK_comment_Author]

ALTER TABLE [comment] WITH CHECK ADD CONSTRAINT [FK_comment_Reactions] FOREIGN KEY([Reactions])
REFERENCES [reaction] ([ID])

ALTER TABLE [comment] CHECK CONSTRAINT [FK_comment_Reactions]

ALTER TABLE [comment] WITH CHECK ADD CONSTRAINT [FK_comment_Replies] FOREIGN KEY([Replies])
REFERENCES [comment] ([ID])

ALTER TABLE [comment] CHECK CONSTRAINT [FK_comment_Replies]

ALTER TABLE [comment] WITH CHECK ADD CONSTRAINT [FK_comment_Channel] FOREIGN KEY([Channel])
REFERENCES [channel] ([ID])

ALTER TABLE [comment] CHECK CONSTRAINT [FK_comment_Channel]

ALTER TABLE [comment] WITH CHECK ADD CONSTRAINT [FK_comment_Is_flagged_Is_flagged_approved] FOREIGN KEY([Is_flagged], [Is_flagged_approved])
REFERENCES [flag] ([ID], [Approved])

ALTER TABLE [comment] CHECK CONSTRAINT [FK_comment_Is_flagged_Is_flagged_approved]

ALTER TABLE [comment] WITH CHECK ADD CONSTRAINT [FK_comment_Parent_post] FOREIGN KEY([Parent_post])
REFERENCES [post] ([ID])

ALTER TABLE [comment] CHECK CONSTRAINT [FK_comment_Parent_post]

ALTER TABLE [comment] WITH CHECK ADD CONSTRAINT [FK_comment_Parent_comment] FOREIGN KEY([Parent_comment])
REFERENCES [comment] ([ID])

ALTER TABLE [comment] CHECK CONSTRAINT [FK_comment_Parent_comment]

ALTER TABLE [reaction] WITH CHECK ADD CONSTRAINT [FK_reaction_Author] FOREIGN KEY([Author])
REFERENCES [user] ([ID])

ALTER TABLE [reaction] CHECK CONSTRAINT [FK_reaction_Author]

ALTER TABLE [reaction] WITH CHECK ADD CONSTRAINT [FK_reaction_Channel] FOREIGN KEY([Channel])
REFERENCES [channel] ([ID])

ALTER TABLE [reaction] CHECK CONSTRAINT [FK_reaction_Channel]

ALTER TABLE [reaction] WITH CHECK ADD CONSTRAINT [FK_reaction_Parent_post_Parent_post_author] FOREIGN KEY([Parent_post], [Parent_post_author])
REFERENCES [post] ([ID], [Author])

ALTER TABLE [reaction] CHECK CONSTRAINT [FK_reaction_Parent_post_Parent_post_author]

ALTER TABLE [reaction] WITH CHECK ADD CONSTRAINT [FK_reaction_Parent_comment_Parent_comment_author] FOREIGN KEY([Parent_comment], [Parent_comment_author])
REFERENCES [comment] ([ID], [Author])

ALTER TABLE [reaction] CHECK CONSTRAINT [FK_reaction_Parent_comment_Parent_comment_author]

ALTER TABLE [flag] WITH CHECK ADD CONSTRAINT [FK_flag_Author] FOREIGN KEY([Author])
REFERENCES [user] ([ID])

ALTER TABLE [flag] CHECK CONSTRAINT [FK_flag_Author]

ALTER TABLE [flag] WITH CHECK ADD CONSTRAINT [FK_flag_Parent_post_Parent_post_author] FOREIGN KEY([Parent_post], [Parent_post_author])
REFERENCES [post] ([ID], [Author])

ALTER TABLE [flag] CHECK CONSTRAINT [FK_flag_Parent_post_Parent_post_author]

ALTER TABLE [flag] WITH CHECK ADD CONSTRAINT [FK_flag_Parent_comment_Parent_comment_author] FOREIGN KEY([Parent_comment], [Parent_comment_author])
REFERENCES [comment] ([ID], [Author])

ALTER TABLE [flag] CHECK CONSTRAINT [FK_flag_Parent_comment_Parent_comment_author]

CREATE INDEX [idx_user_Username]
ON [user] ([Username])

CREATE INDEX [idx_channel_Name]
ON [channel] ([Name])

COMMIT TRANSACTION QUICKDBD