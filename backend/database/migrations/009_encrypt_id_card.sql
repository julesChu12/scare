-- 009_encrypt_id_card.sql
-- 身份证号加密存储：扩容 id_card 列 + 新增 id_card_masked 列

ALTER TABLE `users` ADD COLUMN `id_card_masked` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '身份证号脱敏值' AFTER `id_card_hmac`;
ALTER TABLE `users` MODIFY COLUMN `id_card` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '身份证号（AES-256-GCM 密文）';
