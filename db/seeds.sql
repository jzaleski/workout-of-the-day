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

INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed
) VALUES (
  3,
  '2020-07-26',
  '<u>Goal:</u>&nbsp;<i>Aerobic &amp; Strength Training</i>',
  '<b>Warmup:</b><br/><ul><li>Run 1 Mile</li></ul><b>Exercise(s):</b><br/><ul><li><span>50&nbsp;<a href="https://www.youtube.com/watch?v=Xjo_fY9Hl9w" style="color: #00f;" target="_blank">Dumbbell Goblet Squats</a></span></li><li><span>50&nbsp;<a href="https://www.youtube.com/watch?v=D7KaRcUTQeE" style="color: #00f;" target="_blank">Dumbbell Lunges</a></span></li><li><span>50&nbsp;<a href="https://www.youtube.com/watch?v=LktGPg-AkvY" style="color: #00f;" target="_blank">Dumbbell Bent-Over Rows</a></span></li><li><span>25&nbsp;<a href="https://www.youtube.com/watch?v=ir5PsbniVSc" style="color: #00f;" target="_blank">Dumbbell Skull Crushers</a></span></li><li><span>25&nbsp;<a href="https://www.youtube.com/watch?v=RT_MTXaLKxU" style="color: #00f;" target="_blank">Dumbbell Squat Cleans</a></span></li></ul><b>Cooldown:</b><br/><ul><li>Run 1 Mile</li></ul>',
  '+1-617-455-7595',
  'JonathanZaleski@gmail.com',
  0
) ON CONFLICT DO NOTHING;
