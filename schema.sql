DROP TABLE IF EXISTS `resource`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `resource` (
  `uuid` varchar(8) COLLATE latin1_bin NOT NULL,
  `reference_type` tinyint(4) NOT NULL,
  `reference_id` varchar(16) COLLATE latin1_bin DEFAULT NULL,
  `team_uuid` varchar(8) COLLATE latin1_bin DEFAULT NULL,
  `project_uuid` varchar(16) COLLATE latin1_bin NOT NULL,
  `owner_uuid` varchar(8) COLLATE latin1_bin DEFAULT NULL,
  `type` tinyint(4) NOT NULL DEFAULT '-1',
  `source` tinyint(4) NOT NULL DEFAULT '0',
  `ext_id` varchar(100) COLLATE latin1_bin DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 NOT NULL,
  `create_time` bigint(20) NOT NULL DEFAULT '0',
  `modify_time` bigint(20) NOT NULL DEFAULT '0',
  `status` tinyint(4) NOT NULL DEFAULT '-1',
  `description` varchar(64) CHARACTER SET utf8mb4 NOT NULL,
  `modifier` varchar(8) COLLATE latin1_bin DEFAULT NULL COMMENT '修改者',
  `callback_url` varchar(100) NOT NULL COMMENT '上传回调地址',
  `callback_body` varchar(1000) NOT NULL COMMENT '上传回调内容',
  PRIMARY KEY (`uuid`),
  KEY `index_ext_id` (`ext_id`) USING BTREE,
  KEY `idx_team_uuid_reference_id` (`team_uuid`,`reference_id`),
  KEY `idx_team_uuid_project_uuid_create_time` (`team_uuid`,`project_uuid`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin;