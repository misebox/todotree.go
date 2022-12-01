
CREATE DATABASE IF NOT EXISTS todo_test character set utf8mb4 collate utf8mb4_general_ci;
CREATE USER IF NOT EXISTS 'todo_test'@'%' IDENTIFIED BY 'todo_test';
GRANT ALL ON `todo_test`.* TO 'todo_test'@'%';
FLUSH PRIVILEGES;