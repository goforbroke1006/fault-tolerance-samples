CREATE DATABASE test_db;
CREATE USER 'dbwebapp'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON test_db.* TO 'dbwebapp'@'%';