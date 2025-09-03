CREATE TABLE  IF NOT EXISTS `ats_domain`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT,
  `domain_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `create_time` bigint(19) NULL DEFAULT NULL,
  PRIMARY KEY (`kid`) USING BTREE,
  UNIQUE INDEX `uk_domain_id`(`domain_id`) USING BTREE,
  INDEX `domain_id`(`domain_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;


CREATE TABLE IF NOT EXISTS `ats_supeve`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `seid` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '事件的ID',
  `etime` bigint(19) NOT NULL,
  `create_time` bigint(19) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `seid`(`seid`) USING BTREE,
  INDEX `domain_id`(`domain_id`) USING BTREE,
  CONSTRAINT `ats_supeve_ibfk_1` FOREIGN KEY (`domain_id`) REFERENCES `ats_domain` (`domain_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 10000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;


CREATE TABLE IF NOT EXISTS `ats_extras`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `seid` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `exid` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `reqdata` json NULL,
  `uagent` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `source_ip` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `method` char(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `requrl` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`, `seid`) USING BTREE,
  UNIQUE INDEX `exid_unique`(`exid`) USING BTREE,
  INDEX `exid`(`exid`) USING BTREE,
  INDEX `seid`(`seid`) USING BTREE,
  INDEX `domain_id`(`domain_id`) USING BTREE,
  CONSTRAINT `ats_extras_ibfk_1` FOREIGN KEY (`seid`) REFERENCES `ats_supeve` (`seid`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `ats_extras_ibfk_2` FOREIGN KEY (`domain_id`) REFERENCES `ats_domain` (`domain_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 10000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;


CREATE TABLE IF NOT EXISTS `ats_auditlog`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `seid` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `eid` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '事件ID',
  `user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '操作用户',
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `service` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '服务名称',
  `resource_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '事件名称',
  `rating` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '事件级别',
  `message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '消息',
  `extras` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '额外附加内容',
  `etime` bigint(20) NULL DEFAULT NULL COMMENT '事件发生时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `eid_unique`(`eid`) USING BTREE,
  INDEX `extras`(`extras`) USING BTREE,
  INDEX `seid`(`seid`) USING BTREE,
  INDEX `domain_id`(`domain_id`) USING BTREE,
  CONSTRAINT `ats_auditlog_ibfk_1` FOREIGN KEY (`extras`) REFERENCES `ats_extras` (`exid`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `ats_auditlog_ibfk_2` FOREIGN KEY (`seid`) REFERENCES `ats_supeve` (`seid`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `ats_auditlog_ibfk_3` FOREIGN KEY (`domain_id`) REFERENCES `ats_domain` (`domain_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 10000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;


CREATE TABLE IF NOT EXISTS `ats_logname`  (
  `id` int(19) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `enus` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `zhcn` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `create_time` bigint(19) NOT NULL,
  `update_time` bigint(19) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;
