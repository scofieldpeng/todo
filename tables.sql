DROP TABLE IF EXISTS `category`;

CREATE TABLE `category` (
  `categoryid` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '分类id',
  `userid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `category_name` varchar(10) NOT NULL DEFAULT '' COMMENT '分类名',
  `order` int(10) NOT NULL DEFAULT '0' COMMENT '排序',
  PRIMARY KEY (`categoryid`),
  KEY `userid` (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务类别';

DROP TABLE IF EXISTS `todo`;

CREATE TABLE `todo` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `todo_name` varchar(100) NOT NULL DEFAULT '' COMMENT '任务名称',
  `userid` int(10) unsigned NOT NULL COMMENT '用户id',
  `categoryid` int(10) unsigned NOT NULL COMMENT '分类id',
  `start_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '开始时间戳',
  `end_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '结束时间戳',
  `status` int(10) NOT NULL DEFAULT '0' COMMENT '状态',
  `remark` text NOT NULL COMMENT '备注',
  PRIMARY KEY (`id`),
  KEY `user_todo` (`userid`,`categoryid`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='todo';

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `userid` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` varchar(8) NOT NULL DEFAULT '' COMMENT 'salt',
  `create_time` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `last_login` int(10) NOT NULL DEFAULT '0' COMMENT '上次登录时间',
  `unfinish_num` int(10) NOT NULL DEFAULT '0' COMMENT '待完成任务数量',
  PRIMARY KEY (`userid`),
  KEY `user_name` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';