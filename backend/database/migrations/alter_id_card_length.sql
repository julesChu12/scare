-- 扩大 id_card 列长度以支持加密存储
ALTER TABLE users MODIFY COLUMN id_card VARCHAR(64) DEFAULT NULL;
ALTER TABLE customer_profiles MODIFY COLUMN id_card VARCHAR(64) DEFAULT NULL COMMENT '身份证号';
