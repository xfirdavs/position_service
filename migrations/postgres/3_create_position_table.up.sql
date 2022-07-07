CREATE TABLE IF NOT EXISTS position(
    id uuid primary key,
    name varchar(255) not null,
    profession_id uuid not null,
    company_id uuid not null

);

ALTER TABLE "position" ADD  CONSTRAINT "position_profession_id_fkey" FOREIGN KEY ("profession_id") REFERENCES "profession" ("id");



