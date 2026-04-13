ALTER TABLE `service_requests`
  ADD COLUMN `service_location_lat` DECIMAL(10,7) DEFAULT NULL AFTER `submit_location_lng`,
  ADD COLUMN `service_location_lng` DECIMAL(10,7) DEFAULT NULL AFTER `service_location_lat`,
  ADD COLUMN `source_station_id` BIGINT DEFAULT NULL AFTER `urgency`,
  ADD COLUMN `dispatch_basis` VARCHAR(50) DEFAULT NULL AFTER `station_id`,
  ADD COLUMN `needs_manual_verify` TINYINT(1) DEFAULT 0 AFTER `dispatch_basis`,
  ADD KEY `idx_service_requests_source_station_id` (`source_station_id`);
