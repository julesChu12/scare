INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id
FROM `roles` r
JOIN `permissions` p ON p.code = 'service:task:pool'
WHERE r.code = 'staff'
  AND NOT EXISTS (
    SELECT 1
    FROM `role_permissions` rp
    WHERE rp.role_id = r.id
      AND rp.permission_id = p.id
  );

INSERT INTO `users` (
  `phone`,
  `password_hash`,
  `name`,
  `email`,
  `gender`,
  `birth_date`,
  `id_card`,
  `id_card_hmac`,
  `id_card_masked`,
  `station_id`,
  `status`,
  `created_at`,
  `updated_at`
)
SELECT
  '13800000012',
  '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK',
  '李小勇',
  'lixy@example.com',
  'male',
  '1978-04-20',
  '110105197804201234',
  '',
  '1101**********1234',
  NULL,
  'active',
  NOW(3),
  NOW(3)
WHERE NOT EXISTS (
  SELECT 1
  FROM `users`
  WHERE `phone` = '13800000012'
    AND `deleted_at` IS NULL
);

INSERT INTO `user_identities` (
  `user_id`,
  `identity_type`,
  `is_primary`,
  `station_id`,
  `status`,
  `granted_at`,
  `created_at`,
  `updated_at`
)
SELECT
  u.id,
  'family',
  TRUE,
  NULL,
  'active',
  NOW(),
  NOW(),
  NOW()
FROM `users` u
WHERE u.phone = '13800000012'
  AND u.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1
    FROM `user_identities` ui
    WHERE ui.user_id = u.id
      AND ui.identity_type = 'family'
      AND ui.deleted_at IS NULL
  );

INSERT INTO `customer_profiles` (
  `user_id`,
  `id_card`,
  `address`,
  `customer_type`,
  `gender`,
  `birth_date`,
  `health_status`,
  `emergency_contact`,
  `created_at`,
  `updated_at`
)
SELECT
  u.id,
  '110105197804201234',
  '北京市昌平区霍营街道龙锦苑东一区3号楼',
  'family',
  'male',
  '1978-04-20',
  NULL,
  NULL,
  NOW(3),
  NOW(3)
FROM `users` u
WHERE u.phone = '13800000012'
  AND u.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1
    FROM `customer_profiles` cp
    WHERE cp.user_id = u.id
      AND cp.deleted_at IS NULL
  );
