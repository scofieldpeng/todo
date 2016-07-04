DROP TABLE IF EXISTS `category`;

CREATE TABLE `category` (
  `categoryid` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '分类id',
  `userid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `category_name` varchar(10) NOT NULL DEFAULT '' COMMENT '分类名',
  `order` int(10) NOT NULL DEFAULT '0' COMMENT '排序',
  PRIMARY KEY (`categoryid`),
  KEY `userid` (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务类别';

DROP TABLE IF EXISTS `comment`;

CREATE TABLE `comment` (
  `commentid` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `todoid` int(11) NOT NULL DEFAULT '0' COMMENT 'todoid',
  `userid` int(11) NOT NULL DEFAULT '0' COMMENT 'userid',
  `parent_commentid` int(11) NOT NULL DEFAULT '0' COMMENT '父级commentid',
  `time` int(10) NOT NULL DEFAULT '0' COMMENT '评论时间',
  `content` text NOT NULL COMMENT '评论内容',
  PRIMARY KEY (`commentid`),
  KEY `todoid` (`todoid`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COMMENT='todo评论表';

DROP TABLE IF EXISTS `todo`;

CREATE TABLE `todo` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `todo_name` varchar(100) NOT NULL DEFAULT '' COMMENT '任务名称',
  `userid` int(10) unsigned NOT NULL COMMENT '用户id',
  `categoryid` int(10) unsigned NOT NULL COMMENT '分类id',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `start_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '开始时间戳',
  `end_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '结束时间戳',
  `status` int(10) NOT NULL DEFAULT '0' COMMENT '状态',
  `star` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '重要程度,从1-4,1一般,2重要,3紧急,4重要且紧急',
  `score` int(10) unsigned DEFAULT '0' COMMENT '积分,用户自设,默认为0',
  `remark` text NOT NULL COMMENT '备注',
  PRIMARY KEY (`id`),
  KEY `user_todo` (`userid`,`categoryid`,`status`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COMMENT='todo';

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `userid` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',
  `email` varchar(100) NOT NULL DEFAULT '' COMMENT '用户邮箱',
  `password` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` varchar(8) NOT NULL DEFAULT '' COMMENT 'salt',
  `create_time` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `last_login` int(10) NOT NULL DEFAULT '0' COMMENT '上次登录时间',
  `unfinish_num` int(10) NOT NULL DEFAULT '0' COMMENT '待完成任务数量',
  `score` int(10) unsigned DEFAULT '0' COMMENT '用户积分',
  PRIMARY KEY (`userid`),
  KEY `user_name` (`user_name`),
  KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';