CREATE TABLE IF NOT EXISTS workout (
  id SERIAL PRIMARY KEY,
  date timestamp without time zone NOT NULL,
  goal character varying NOT NULL,
  description text NOT NULL,
  sms_to character varying NULL,
  mail_to character varying NULL,
  completed int NOT NULL DEFAULT 0,
  voting_enabled boolean NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_workout_date_desc ON workout (date DESC);
