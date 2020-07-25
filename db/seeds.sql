INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed
) VALUES (
  1,
  '1970-01-01',
  '<u>Goal:</u>&nbsp;<i>General Fitness</i>',
  '<b>Exercise(s):</b><br/><ul><li>150 Air Squats</li><li>100 Sit-Ups</li><li>50 Push-Ups</li><li>Run 3 Miles</li></ul>',
  '+1-617-455-7595',
  'JonathanZaleski@gmail.com',
  0
) ON CONFLICT DO NOTHING;
