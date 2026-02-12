CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS units (
id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
unit_guid    text NOT NULL UNIQUE,
created_at   timestamptz NOT NULL DEFAULT now(),
updated_at   timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tsv_files (
id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
file_name     text NOT NULL UNIQUE,
file_sha256   text NOT NULL,
status        text NOT NULL,
error_message text NULL,
created_at    timestamptz NOT NULL DEFAULT now(),
updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tsv_files_status ON tsv_files(status);

CREATE TABLE IF NOT EXISTS tsv_records (
id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
file_id      uuid NOT NULL REFERENCES tsv_files(id) ON DELETE CASCADE,
unit_guid    text NOT NULL,
line_no      int NOT NULL,
msg_id       text NULL,
payload      jsonb NOT NULL,
created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tsv_records_unit_guid ON tsv_records(unit_guid);
CREATE INDEX IF NOT EXISTS idx_tsv_records_file_id ON tsv_records(file_id);

CREATE TABLE IF NOT EXISTS parse_errors (
id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
file_id     uuid NULL REFERENCES tsv_files(id) ON DELETE SET NULL,
file_name   text NOT NULL,
unit_guid   text NULL,
line_no     int NULL,
error       text NOT NULL,
raw_line    text NULL,
created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_parse_errors_file_name ON parse_errors(file_name);

CREATE TABLE IF NOT EXISTS reports (
id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
unit_guid   text NOT NULL,
file_path   text NOT NULL,
created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reports_unit_guid ON reports(unit_guid);
