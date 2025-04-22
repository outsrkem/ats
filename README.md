## ats

### Audit Traces Server

### 组件构成

- MySQL
- go
- Docker

### 前提条件

- UIAS服务
- 提前安装好docker
- 提前安装好mysql

### 初始化数据库

```sql
-- 创建数据库
CREATE DATABASE /*!32312 IF NOT EXISTS */ atsdb /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

-- 创建用户
CREATE USER 'ats'@'%' identified by '123456##ats';

-- 授权数据库权限给用户
GRANT ALL ON atsdb.* TO 'ats'@'%';

-- 修改密码
SET PASSWORD FOR 'ats'@'%' = '12345678';
```

### 初始化表和数据

> 导入 `db/atsdb.sql` 文件到`atsdb`库

```
mysql -u ats -p atsdb < db/atsdb.sql
```

### 构建镜像

```
./.cid/build.sh
```
