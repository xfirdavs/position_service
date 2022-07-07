CREATE TABLE IF NOT EXISTS position_attributes(
    id uuid primary key,
   attribute_id uuid not null,
    position_id uuid not null,
    value varchar(255) not null

);

ALTER TABLE "position_attributes" ADD  CONSTRAINT "position_attributes_attribute_id_fkey" FOREIGN KEY ("attribute_id") REFERENCES "attribute" ("id");
ALTER TABLE "position_attributes" ADD  CONSTRAINT "position_attributes_position_id_fkey" FOREIGN KEY ("position_id") REFERENCES "position" ("id");



