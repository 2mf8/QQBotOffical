USE [kequ5060]
GO

/****** Object:  Table [dbo].[guild_price]    Script Date: 2022/6/6 10:54:26 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[guild_price](
	[ID] [bigint] IDENTITY(1,1) NOT NULL,
	[guild_id] [varchar](32) NOT NULL,
	[channel_id] [varchar](32) NOT NULL,
	[brand] [nvarchar](50) NULL,
	[item] [nvarchar](100) NOT NULL,
	[price] [varchar](100) NULL,
	[shipping] [nvarchar](100) NULL,
	[updater] [varchar](32) NOT NULL,
	[gmt_modified] [datetime2](7) NULL
) ON [PRIMARY]
GO

