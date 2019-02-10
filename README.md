
# MySQLを利用するには
```
go get github.com/go-sql-driver/mysql
```

# テーブル定義
```sql
CREATE TABLE `tasks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `message` text COLLATE utf8mb4_unicode_ci,
  `status` int(1) DEFAULT '0',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

# WebAPI仕様

|method|path|require request body|
|:--|:--|:--|
|GET|/task|none|
|POST|/task|{"message": "add new task"}|
|GET|/task/{task_id}|none|
|POST|/task/{task_id}|{"id": 1, "message": "update task message", "status": 1}|
|DELETE|/task/{task_id}|none|
