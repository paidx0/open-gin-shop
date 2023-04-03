-- MySQL dump 10.13  Distrib 8.0.32, for macos12.6 (x86_64)
--
-- Host: localhost    Database: shop_order_srv
-- ------------------------------------------------------
-- Server version	8.0.32

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `ordergoods`
--

DROP TABLE IF EXISTS `ordergoods`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ordergoods` (
  `id` int NOT NULL AUTO_INCREMENT,
  `add_time` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  `update_time` datetime NOT NULL,
  `order` int NOT NULL,
  `goods` int NOT NULL,
  `goods_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `goods_image` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `goods_price` decimal(10,5) NOT NULL,
  `nums` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=75 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ordergoods`
--

LOCK TABLES `ordergoods` WRITE;
/*!40000 ALTER TABLE `ordergoods` DISABLE KEYS */;
INSERT INTO `ordergoods` VALUES (1,'2020-11-18 17:08:45',0,'2020-11-18 17:08:45',1,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,5),(2,'2020-11-18 17:08:45',0,'2020-11-18 17:08:45',1,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,3),(3,'2020-11-18 17:10:22',0,'2020-11-18 17:10:22',2,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,5),(4,'2020-11-18 17:10:22',0,'2020-11-18 17:10:22',2,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,3),(5,'2020-11-18 17:10:52',0,'2020-11-18 17:10:52',3,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,5),(6,'2020-11-18 17:10:52',0,'2020-11-18 17:10:52',3,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,3),(7,'2020-11-18 17:14:11',0,'2020-11-18 17:14:11',4,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,5),(8,'2020-11-18 17:14:11',0,'2020-11-18 17:14:11',4,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,3),(9,'2020-11-19 13:43:52',0,'2020-11-19 13:43:52',5,423,'越南进口红心火龙果 4个装 红肉中果 单','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/b39672c6abebe124b982250642cb9a0f',27.90000,5),(10,'2020-11-19 13:45:26',0,'2020-11-19 13:45:26',6,423,'越南进口红心火龙果 4个装 红肉中果 单','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/b39672c6abebe124b982250642cb9a0f',27.90000,5),(11,'2020-11-19 13:47:49',0,'2020-11-19 13:47:49',7,423,'越南进口红心火龙果 4个装 红肉中果 单','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/b39672c6abebe124b982250642cb9a0f',27.90000,5),(12,'2020-11-19 15:52:34',0,'2020-11-19 15:52:34',8,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,5),(13,'2020-11-19 15:52:34',0,'2020-11-19 15:52:34',8,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,11),(14,'2020-11-19 15:52:34',0,'2020-11-19 15:52:34',8,423,'越南进口红心火龙果 4个装 红肉中果 单','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/b39672c6abebe124b982250642cb9a0f',27.90000,5),(15,'2020-11-19 22:40:57',0,'2020-11-19 22:40:57',9,430,'Zespri佳沛 新西兰阳光金奇异果 6','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/f02d5ff01123557697f897eb2172c4e4',69.90000,1),(16,'2020-11-26 16:47:20',0,'2020-11-26 16:47:20',10,430,'Zespri佳沛 新西兰阳光金奇异果 6','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/f02d5ff01123557697f897eb2172c4e4',69.90000,1),(17,'2020-11-26 20:43:30',0,'2020-11-26 20:43:30',11,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,5),(18,'2020-11-26 20:43:30',0,'2020-11-26 20:43:30',11,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,11),(19,'2020-11-26 20:43:30',0,'2020-11-26 20:43:30',11,423,'越南进口红心火龙果 4个装 红肉中果 单','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/b39672c6abebe124b982250642cb9a0f',27.90000,5),(20,'2020-11-26 20:47:40',0,'2020-11-26 20:47:40',12,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,5),(21,'2020-11-26 20:47:40',0,'2020-11-26 20:47:40',12,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,11),(22,'2020-11-26 20:47:40',0,'2020-11-26 20:47:40',12,423,'越南进口红心火龙果 4个装 红肉中果 单','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/b39672c6abebe124b982250642cb9a0f',27.90000,5),(23,'2020-11-27 02:59:22',0,'2020-11-27 02:59:22',13,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,5),(24,'2020-11-27 02:59:22',0,'2020-11-27 02:59:22',13,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,11),(25,'2020-11-27 03:39:45',0,'2020-11-27 03:39:45',14,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(26,'2020-11-27 03:39:45',0,'2020-11-27 03:39:45',14,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,20),(27,'2020-11-27 03:41:58',0,'2020-11-27 03:41:58',15,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(28,'2020-11-27 03:41:58',0,'2020-11-27 03:41:58',15,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,20),(31,'2020-11-27 05:39:14',0,'2020-11-27 05:39:14',17,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(32,'2020-11-27 05:39:14',0,'2020-11-27 05:39:14',17,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,20),(33,'2020-11-27 05:42:39',0,'2020-11-27 05:42:39',18,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(34,'2020-11-27 05:42:39',0,'2020-11-27 05:42:39',18,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,20),(35,'2020-11-27 05:46:21',0,'2020-11-27 05:46:21',19,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(36,'2020-11-27 05:46:21',0,'2020-11-27 05:46:21',19,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,20),(37,'2020-11-27 05:50:01',0,'2020-11-27 05:50:01',20,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(38,'2020-11-27 05:50:01',0,'2020-11-27 05:50:01',20,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,20),(39,'2020-11-27 05:51:36',0,'2020-11-27 05:51:36',21,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(40,'2020-11-27 05:51:36',0,'2020-11-27 05:51:36',21,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,20),(41,'2020-11-27 05:53:22',0,'2020-11-27 05:53:22',22,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(42,'2020-11-27 05:53:22',0,'2020-11-27 05:53:22',22,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,20),(43,'2020-11-27 05:55:44',0,'2020-11-27 05:55:44',23,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(44,'2020-11-27 05:55:44',0,'2020-11-27 05:55:44',23,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,20),(45,'2020-11-27 05:57:22',0,'2020-11-27 05:57:22',24,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(46,'2020-11-27 05:59:35',0,'2020-11-27 05:59:35',25,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(47,'2020-11-27 06:00:10',0,'2020-11-27 06:00:10',26,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(48,'2020-11-27 06:01:18',0,'2020-11-27 06:01:18',27,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(49,'2020-11-27 06:03:23',0,'2020-11-27 06:03:23',28,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(50,'2020-11-27 06:06:48',0,'2020-11-27 06:06:48',29,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(51,'2020-11-27 06:08:58',0,'2020-11-27 06:08:58',30,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(52,'2020-11-27 07:00:17',0,'2020-11-27 07:00:17',31,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(53,'2020-11-27 07:00:17',0,'2020-11-27 07:00:17',31,422,'西州蜜瓜25号哈密瓜 2粒装 单果1.2','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73',36.90000,20),(54,'2020-11-27 14:48:58',0,'2020-11-27 14:48:58',32,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(55,'2020-11-27 14:53:37',0,'2020-11-27 14:53:37',33,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(56,'2020-11-27 14:56:55',0,'2020-11-27 14:56:55',34,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(57,'2020-11-27 15:00:20',0,'2020-11-27 15:00:20',35,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(58,'2020-11-27 15:03:04',0,'2020-11-27 15:03:04',36,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(59,'2020-11-27 15:07:07',0,'2020-11-27 15:07:07',37,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(60,'2020-11-29 09:24:51',0,'2020-11-29 09:24:51',38,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(61,'2020-11-29 09:32:20',0,'2020-11-29 09:32:20',39,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(62,'2020-11-29 09:39:35',0,'2020-11-29 09:39:35',40,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(63,'2020-11-29 09:43:50',0,'2020-11-29 09:43:50',41,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(64,'2020-11-29 09:48:35',0,'2020-11-29 09:48:35',42,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(65,'2020-12-01 19:53:53',0,'2020-12-01 19:53:53',43,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(66,'2020-12-01 19:55:29',0,'2020-12-01 19:55:29',44,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(67,'2020-12-01 20:37:13',0,'2020-12-01 20:37:13',45,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(68,'2020-12-01 20:53:53',0,'2020-12-01 20:53:53',46,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(69,'2020-12-01 21:19:36',0,'2020-12-01 21:19:36',47,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(70,'2020-12-01 21:20:37',0,'2020-12-01 21:20:37',48,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(71,'2020-12-01 21:51:47',0,'2020-12-01 21:51:47',49,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(72,'2020-12-01 21:57:23',0,'2020-12-01 21:57:23',50,421,'烟台红富士苹果12个 净重2.6kg以上','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1',44.90000,10),(73,'2020-12-13 01:47:19',0,'2020-12-13 01:47:19',51,443,'国产水蜜桃 新鲜桃子 精选特级果3kg装','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/6cb86dd7a499aeea3898afbe06ee129f',49.90000,1),(74,'2020-12-13 01:47:19',0,'2020-12-13 01:47:19',51,562,'【已通过核酸检测】极地湾 冷冻多春鱼 1','https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/f02565ef3883975ba51df6440a71e312',59.90000,1);
/*!40000 ALTER TABLE `ordergoods` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `orderinfo`
--

DROP TABLE IF EXISTS `orderinfo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orderinfo` (
  `id` int NOT NULL AUTO_INCREMENT,
  `add_time` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  `update_time` datetime NOT NULL,
  `user` int NOT NULL,
  `order_sn` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `pay_type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `status` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `trade_no` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `order_mount` float NOT NULL,
  `pay_time` datetime DEFAULT NULL,
  `address` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `signer_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `singer_mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `post` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `orderinfo_order_sn` (`order_sn`) USING BTREE,
  UNIQUE KEY `orderinfo_trade_no` (`trade_no`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orderinfo`
--

LOCK TABLES `orderinfo` WRITE;
/*!40000 ALTER TABLE `orderinfo` DISABLE KEYS */;
INSERT INTO `orderinfo` VALUES (1,'2020-11-18 17:08:45',0,'2020-11-18 17:08:45',1,'20201118170844170','alipay','paying',NULL,335.2,NULL,'北京市','bobby','18787878787','请尽快发货'),(2,'2020-11-18 17:10:22',0,'2020-11-18 17:10:22',1,'20201118171022198','alipay','paying',NULL,335.2,NULL,'北京市','bobby','18787878787','请尽快发货'),(3,'2020-11-18 17:10:52',0,'2020-11-18 17:10:52',1,'20201118171051179','alipay','paying',NULL,335.2,NULL,'北京市','bobby','18787878787','请尽快发货'),(4,'2020-11-18 17:14:11',0,'2020-11-18 17:14:11',1,'20201118171410162','alipay','paying',NULL,335.2,NULL,'北京市','bobby','18787878787','请尽快发货'),(7,'2020-11-19 13:47:50',0,'2020-11-19 13:47:50',1,'20201119134749110','alipay','paying',NULL,139.5,NULL,'其它区','Susan Thompson','18788989898','请尽快发货'),(8,'2020-11-19 15:52:34',0,'2020-11-19 15:52:34',1,'20201119155234198','alipay','paying',NULL,769.9,NULL,'武清区','Lisa Williams','18788989898','请尽快发货'),(9,'2020-11-19 22:40:57',0,'2020-11-19 22:40:57',1,'20201119224057124','alipay','paying',NULL,69.9,NULL,'陕西省宜昌市华中Helen Lewis','George Allen','18686868686','第一个订单'),(10,'2020-11-26 16:47:29',0,'2020-11-26 16:47:29',1,'20201126164709163','wechat','paying',NULL,69.9,NULL,'陕西省宜昌市华中Helen Lewis','George Allen','18686868686','hhhh'),(31,'2020-11-27 07:00:26',0,'2020-11-27 07:01:26',1,'20201127070017114','alipay','TRADE_CLOSED',NULL,1187,NULL,'承德县','Linda Hernandez','18788989898','请尽快发货'),(32,'2020-11-27 14:49:00',0,'2020-11-27 14:53:21',1,'20201127144846190','wechat','TRADE_CLOSED',NULL,449,NULL,'长白朝鲜族自治县','Kimberly Jones','18788989898','请尽快发货'),(33,'2020-11-27 14:53:45',0,'2020-11-27 14:54:48',1,'20201127145325137','wechat','TRADE_CLOSED',NULL,449,NULL,'长白朝鲜族自治县','Kimberly Jones','18788989898','请尽快发货'),(34,'2020-11-27 14:56:57',0,'2020-11-27 14:57:58',1,'20201127145648156','wechat','TRADE_CLOSED',NULL,449,NULL,'长白朝鲜族自治县','Kimberly Jones','18788989898','请尽快发货'),(35,'2020-11-27 15:00:21',0,'2020-11-27 15:01:21',1,'20201127150015115','wechat','TRADE_CLOSED',NULL,449,NULL,'长白朝鲜族自治县','Kimberly Jones','18788989898','请尽快发货'),(36,'2020-11-27 15:03:12',0,'2020-11-27 15:04:12',1,'20201127150220121','wechat','TRADE_CLOSED',NULL,449,NULL,'长白朝鲜族自治县','Kimberly Jones','18788989898','请尽快发货'),(37,'2020-11-27 15:07:09',0,'2020-11-27 15:08:10',1,'20201127150657145','wechat','TRADE_CLOSED',NULL,449,NULL,'长白朝鲜族自治县','Kimberly Jones','18788989898','请尽快发货'),(38,'2020-11-29 09:24:51',0,'2020-11-29 09:25:51',1,'20201129092450166','alipay','TRADE_CLOSED',NULL,449,NULL,'博山区','Frank Lee','18788989898','请尽快发货'),(39,'2020-11-29 09:32:20',0,'2020-11-29 09:33:20',1,'20201129093219132','alipay','TRADE_CLOSED',NULL,449,NULL,'博山区','Frank Lee','18788989898','请尽快发货'),(40,'2020-11-29 09:39:35',0,'2020-11-29 09:40:35',1,'20201129093935115','alipay','TRADE_CLOSED',NULL,449,NULL,'博山区','Frank Lee','18788989898','请尽快发货'),(41,'2020-11-29 09:43:50',0,'2020-11-29 09:44:50',1,'20201129094350157','alipay','TRADE_CLOSED',NULL,449,NULL,'博山区','Frank Lee','18788989898','请尽快发货'),(42,'2020-11-29 09:48:35',0,'2020-11-29 09:49:35',1,'20201129094834113','alipay','TRADE_CLOSED',NULL,449,NULL,'博山区','Frank Lee','18788989898','请尽快发货'),(43,'2020-12-01 19:53:54',0,'2020-12-01 19:54:54',1,'20201201195352175','alipay','TRADE_CLOSED',NULL,449,NULL,'-','Kimberly Robinson','18788989898','请尽快发货'),(44,'2020-12-01 19:55:29',0,'2020-12-01 19:56:29',1,'20201201195529184','alipay','TRADE_CLOSED',NULL,449,NULL,'-','Kimberly Robinson','18788989898','请尽快发货'),(45,'2020-12-01 20:37:13',0,'2020-12-01 20:38:13',1,'20201201203712117','alipay','TRADE_CLOSED',NULL,449,NULL,'怀仁县','Michael Moore','18788989898','请尽快发货'),(46,'2020-12-01 20:53:53',0,'2020-12-01 20:54:53',1,'20201201205352195','alipay','TRADE_CLOSED',NULL,449,NULL,'怀仁县','Michael Moore','18788989898','请尽快发货'),(47,'2020-12-01 21:19:36',0,'2020-12-01 21:20:36',1,'20201201211936129','alipay','TRADE_CLOSED',NULL,449,NULL,'怀仁县','Michael Moore','18788989898','请尽快发货'),(48,'2020-12-01 21:20:37',0,'2020-12-01 21:21:37',1,'20201201212036147','alipay','TRADE_CLOSED',NULL,449,NULL,'怀仁县','Michael Moore','18788989898','请尽快发货'),(49,'2020-12-01 21:51:47',0,'2020-12-01 21:52:48',1,'20201201215147159','alipay','TRADE_CLOSED',NULL,449,NULL,'怀仁县','Michael Moore','18788989898','请尽快发货'),(50,'2020-12-01 21:57:23',0,'2020-12-01 21:58:23',1,'20201201215723121','alipay','TRADE_CLOSED',NULL,449,NULL,'怀仁县','Michael Moore','18788989898','请尽快发货'),(51,'2020-12-13 01:47:19',0,'2020-12-13 01:48:19',1,'20201213014717160','alipay','TRADE_CLOSED',NULL,109.8,NULL,'陕西省宜昌市华中Helen Lewis','George Allen','18686868686','尽快发货');
/*!40000 ALTER TABLE `orderinfo` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `shoppingcart`
--

DROP TABLE IF EXISTS `shoppingcart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `shoppingcart` (
  `id` int NOT NULL AUTO_INCREMENT,
  `add_time` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  `update_time` datetime NOT NULL,
  `user` int NOT NULL,
  `goods` int NOT NULL,
  `nums` int NOT NULL,
  `checked` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `shoppingcart`
--

LOCK TABLES `shoppingcart` WRITE;
/*!40000 ALTER TABLE `shoppingcart` DISABLE KEYS */;
INSERT INTO `shoppingcart` VALUES (1,'2020-11-18 17:01:51',1,'2020-11-18 22:10:29',1,421,10,1),(2,'2020-11-18 17:04:40',1,'2020-11-18 22:13:10',1,422,20,1),(3,'2020-11-18 20:37:57',1,'2020-11-18 20:37:57',1,423,5,1),(4,'2020-11-19 22:26:33',1,'2020-11-19 22:40:48',1,430,1,1),(5,'2020-11-19 23:47:35',1,'2020-11-19 23:47:35',1,430,1,1),(6,'2020-12-12 01:05:20',1,'2020-12-12 01:05:20',1,443,1,1),(7,'2020-12-13 01:45:04',1,'2020-12-13 01:45:04',1,562,1,1);
/*!40000 ALTER TABLE `shoppingcart` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'shop_order_srv'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-04-04  5:09:24
