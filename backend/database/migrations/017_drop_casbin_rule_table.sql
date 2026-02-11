-- 017_drop_casbin_rule_table.sql
-- 删除已废弃的 casbin_rule 表（权限系统已迁移到 permissions 表）

DROP TABLE IF EXISTS `casbin_rule`;
