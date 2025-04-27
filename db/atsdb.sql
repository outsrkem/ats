/*
 Navicat Premium Data Transfer

 Source Server         : 10.10.10.16
 Source Server Type    : MySQL
 Source Server Version : 50744
 Source Host           : 10.10.10.16:3306
 Source Schema         : atsdb

 Target Server Type    : MySQL
 Target Server Version : 50744
 File Encoding         : 65001

 Date: 28/04/2025 10:04:57
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for auditlog
-- ----------------------------
DROP TABLE IF EXISTS `auditlog`;
CREATE TABLE `auditlog`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
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
  CONSTRAINT `auditlog_ibfk_1` FOREIGN KEY (`extras`) REFERENCES `extras` (`exid`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `auditlog_ibfk_2` FOREIGN KEY (`seid`) REFERENCES `supeve` (`seid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1001 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for extras
-- ----------------------------
DROP TABLE IF EXISTS `extras`;
CREATE TABLE `extras`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `seid` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `exid` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `reqdata` varchar(3000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `uagent` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `source_ip` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `method` char(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `requrl` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`, `seid`) USING BTREE,
  UNIQUE INDEX `exid_unique`(`exid`) USING BTREE,
  INDEX `exid`(`exid`) USING BTREE,
  INDEX `seid`(`seid`) USING BTREE,
  CONSTRAINT `extras_ibfk_1` FOREIGN KEY (`seid`) REFERENCES `supeve` (`seid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1001 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for logname
-- ----------------------------
DROP TABLE IF EXISTS `logname`;
CREATE TABLE `logname`  (
  `id` int(19) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `enus` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `zhcn` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `create_time` bigint(19) NOT NULL,
  `update_time` bigint(19) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1001 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for supeve
-- ----------------------------
DROP TABLE IF EXISTS `supeve`;
CREATE TABLE `supeve`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `seid` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '事件的ID',
  `etime` bigint(19) NOT NULL,
  `create_time` bigint(19) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `seid`(`seid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1001 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
