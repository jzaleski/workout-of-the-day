CREATE TABLE IF NOT EXISTS workout (
  id SERIAL PRIMARY KEY,
  date timestamp without time zone NOT NULL,
  goal character varying NOT NULL,
  description text NOT NULL,
  sms_to character varying NULL,
  mail_to character varying NULL
);
