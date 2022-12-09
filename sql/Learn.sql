USE [kequ5060]
GO

/****** Object:  Table [dbo].[guild_learn]    Script Date: 2022/6/6 10:56:45 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[guild_learn](
	[ID] [bigint] IDENTITY(1,1) NOT NULL,
	[ask] [nvarchar](200) NOT NULL,
	[guild_id] [varchar](32) NOT NULL,
	[channel_id] [varchar](32) NOT NULL,
	[admin_id] [varchar](32) NULL,
	[answer] [nvarchar](2000) NULL,
	[gmt_modified] [datetime2](7) NULL,
	[pass] [bit] NULL
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO

