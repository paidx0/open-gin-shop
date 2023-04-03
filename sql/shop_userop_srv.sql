-- MySQL dump 10.13  Distrib 8.0.32, for macos12.6 (x86_64)
--
-- Host: localhost    Database: shop_userop_srv
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
-- Table structure for table `address`
--

DROP TABLE IF EXISTS `address`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `address` (
  `id` int NOT NULL AUTO_INCREMENT,
  `add_time` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  `update_time` datetime NOT NULL,
  `user` int NOT NULL,
  `province` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `city` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `district` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `address` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `signer_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `signer_mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `address`
--

LOCK TABLES `address` WRITE;
/*!40000 ALTER TABLE `address` DISABLE KEYS */;
INSERT INTO `address` VALUES (2,'2020-11-19 20:55:34',0,'2020-11-19 20:55:34',1,'陕西省','宜昌市','华中','Helen Lewis','George Allen','18686868686'),(3,'2020-11-19 23:48:50',0,'2020-11-19 23:48:50',1,'天津市','天津市','和平区','北京市','bobby','18782222220'),(4,'2020-12-11 23:44:06',0,'2020-12-11 23:44:06',1,'北京市','北京市','西城区','北京市','bobby','18787878787'),(5,'2020-12-15 18:11:27',0,'2020-12-15 18:11:27',1,'河北省','唐山市','路南区','bobby','bobby','18789898987');
/*!40000 ALTER TABLE `address` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `leavingmessages`
--

DROP TABLE IF EXISTS `leavingmessages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `leavingmessages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `add_time` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  `update_time` datetime NOT NULL,
  `user` int NOT NULL,
  `message_type` int NOT NULL,
  `subject` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `file` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `leavingmessages`
--

LOCK TABLES `leavingmessages` WRITE;
/*!40000 ALTER TABLE `leavingmessages` DISABLE KEYS */;
INSERT INTO `leavingmessages` VALUES (1,'2020-11-19 20:56:58',0,'2020-11-19 20:56:58',1,1,'occaecat aute voluptate dolor','ad','rlogin://nsmyiuyc.tt/dxrieny'),(4,'2020-12-15 18:17:34',0,'2020-12-15 18:17:34',1,5,'口罩 ','继续一批kn95口罩','http://mxshop-files.oss-cn-hangzhou.aliyuncs.com/mxshop-images/csp2019真题.zip');
/*!40000 ALTER TABLE `leavingmessages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `userfav`
--

DROP TABLE IF EXISTS `userfav`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `userfav` (
  `add_time` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  `update_time` datetime NOT NULL,
  `user` int NOT NULL,
  `goods` int NOT NULL,
  PRIMARY KEY (`user`,`goods`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `userfav`
--

LOCK TABLES `userfav` WRITE;
/*!40000 ALTER TABLE `userfav` DISABLE KEYS */;
INSERT INTO `userfav` VALUES ('2020-12-14 22:03:32',0,'2020-12-14 22:03:32',1,423),('2020-12-14 23:12:46',0,'2020-12-14 23:12:46',1,430);
/*!40000 ALTER TABLE `userfav` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'shop_userop_srv'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-04-04  5:10:05
