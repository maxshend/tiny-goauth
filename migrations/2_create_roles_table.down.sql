DROP INDEX IF EXISTS index_user_roles_on_user_id;
DROP INDEX IF EXISTS index_user_roles_on_role_id_and_user_id;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
