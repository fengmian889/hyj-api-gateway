SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for call_log
-- ----------------------------
DROP TABLE IF EXISTS `call_log`;
CREATE TABLE `call_log` (
  `id` bigint(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_time` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `answered_time` timestamp NULL DEFAULT NULL COMMENT '接听时间',
  `call_time` datetime NOT NULL COMMENT '呼叫时间',
  `hangup_time` timestamp NULL DEFAULT NULL COMMENT '挂断时间',
  `duration` int(11) unsigned DEFAULT NULL COMMENT '通话时长',
  `dani` varchar(64) NOT NULL DEFAULT '' COMMENT 'dani',
  `sip_code` varchar(128) NOT NULL DEFAULT '' COMMENT 'sip错误码',
  `org_id` bigint(10) unsigned NOT NULL DEFAULT '0' COMMENT '公司id',
  `team_id` bigint(10) unsigned NOT NULL DEFAULT '0' COMMENT '团队id',
  `agent_org_id` bigint(10) unsigned DEFAULT NULL COMMENT '坐席id',
  `call_uuid` varchar(45) NOT NULL COMMENT '通话记录唯一键',
  `callee` varchar(64) NOT NULL DEFAULT '' COMMENT '被叫号码',
  `caller` varchar(64) NOT NULL DEFAULT '' COMMENT '主叫号码',
  `call_duration` bigint(20) NOT NULL DEFAULT '0' COMMENT '通话时长(毫秒值)'
) AUTO_INCREMENT = 1000074576 DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '通话记录'
