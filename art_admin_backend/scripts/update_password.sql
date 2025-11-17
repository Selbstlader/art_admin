-- 更新用户密码为 123456
USE gin_admin;

UPDATE sys_user SET password = '$2a$10$hW6cCywvV3KqFpep3Qd.buD/04YsVcb6OsL4Bmh3DyvVNVn..AN8S' WHERE user_name = 'admin';
UPDATE sys_user SET password = '$2a$10$hW6cCywvV3KqFpep3Qd.buD/04YsVcb6OsL4Bmh3DyvVNVn..AN8S' WHERE user_name = 'user';
UPDATE sys_user SET password = '$2a$10$hW6cCywvV3KqFpep3Qd.buD/04YsVcb6OsL4Bmh3DyvVNVn..AN8S' WHERE user_name = 'guest';

SELECT user_name, '密码已更新为: 123456' as message FROM sys_user WHERE user_name IN ('admin', 'user', 'guest');

