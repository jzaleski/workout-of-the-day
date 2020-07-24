INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to
) VALUES (
  1,
  '2020-07-23',
  'GOAL',
  'DESCRIPTION',
  'SMS_TO',
  'MAIL_TO'
) ON CONFLICT DO NOTHING;
