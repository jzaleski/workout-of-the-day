CREATE TABLE IF NOT EXISTS workout (
  date timestamp without time zone PRIMARY KEY,
  goal character varying NOT NULL,
  description text NOT NULL,
  sms_to character varying NULL,
  mail_to character varying NULL,
  completed int NOT NULL DEFAULT 0,
  voting_enabled boolean NOT NULL DEFAULT true
);
