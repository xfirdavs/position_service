ALTER TABLE IF EXISTS "position_attributes" DROP CONSTRAINT IF EXISTS "position_attributes_position_id_fkey";
ALTER TABLE IF EXISTS "position_attributes" DROP CONSTRAINT IF EXISTS "position_attributes_attribute_id_fkey";
DROP TABLE IF EXISTS "position_attributes";