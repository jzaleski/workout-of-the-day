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

INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed
) VALUES (
  2,
  '2020-07-25',
  '<u>Goal:</u>&nbsp;<i>Cardio &amp; Strength Training</i>',
  '<b>Exercise(s):</b><br/><ul><li>100 Dumbbell Hammer Curls</li><li>100 Dumbbell Bicep Curls</li><li>50 Dumbbell Tricep Extensions</li><li>5 Uphill Sprints (~100m)</li></ul><b>Cooldown:</b><br/><ul><li>Run 1 Mile</li></ul>',
  '+1-617-455-7595',
  'JonathanZaleski@gmail.com',
  0
) ON CONFLICT DO NOTHING;
