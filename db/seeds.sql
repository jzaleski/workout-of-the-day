INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed,
  voting_enabled
) VALUES (
  1,
  '1970-01-01',
  '<u>Goal:</u>&nbsp;<i>General Fitness</i>',
  '<b>Exercise(s):</b><br/><ul><li>150 Air Squats</li><li>100 Sit-Ups</li><li>50 Push-Ups</li><li>Run 3 Miles</li></ul>',
  '+1-617-455-7595',
  'JonathanZaleski@gmail.com',
  0,
  true
) ON CONFLICT DO NOTHING;

INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed,
  voting_enabled
) VALUES (
  2,
  '2020-07-25',
  '<u>Goal:</u>&nbsp;<i>Cardio &amp; Strength Training</i>',
  '<b>Exercise(s):</b><br/><ul><li>100 Dumbbell Hammer Curls</li><li>100 Dumbbell Bicep Curls</li><li>50 Dumbbell Tricep Extensions</li><li>5 Uphill Sprints (~100m)</li></ul><b>Cooldown:</b><br/><ul><li>Run 1 Mile</li></ul>',
  '+1-617-455-7595',
  'JonathanZaleski@gmail.com',
  0,
  true
) ON CONFLICT DO NOTHING;

INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed,
  voting_enabled
) VALUES (
  3,
  '2020-07-26',
  '<u>Goal:</u>&nbsp;<i>Aerobic &amp; Strength Training</i>',
  '<b>Warmup:</b><br/><ul><li>Run 1 Mile</li></ul><b>Exercise(s):</b><br/><ul><li><span>50&nbsp;<a href="https://www.youtube.com/watch?v=Xjo_fY9Hl9w" style="color: #00f;" target="_blank">Dumbbell Goblet Squats</a></span></li><li><span>50&nbsp;<a href="https://www.youtube.com/watch?v=D7KaRcUTQeE" style="color: #00f;" target="_blank">Dumbbell Lunges</a></span></li><li><span>50&nbsp;<a href="https://www.youtube.com/watch?v=LktGPg-AkvY" style="color: #00f;" target="_blank">Dumbbell Bent-Over Rows</a></span></li><li><span>25&nbsp;<a href="https://www.youtube.com/watch?v=ir5PsbniVSc" style="color: #00f;" target="_blank">Dumbbell Skull Crushers</a></span></li><li><span>25&nbsp;<a href="https://www.youtube.com/watch?v=RT_MTXaLKxU" style="color: #00f;" target="_blank">Dumbbell Squat Cleans</a></span></li></ul><b>Cooldown:</b><br/><ul><li>Run 1 Mile</li></ul>',
  '+1-617-455-7595',
  'JonathanZaleski@gmail.com',
  0,
  true
) ON CONFLICT DO NOTHING;

INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed,
  voting_enabled
) VALUES (
  4,
  '2020-07-27',
  '<u>Goal:</u>&nbsp;<i>Rest Day</i>',
  '<b>Exercise(s):</b><br/><br/>No Workout Today - Enjoy the day off!<br/><br/>',
  null,
  null,
  0,
  false
) ON CONFLICT DO NOTHING;

INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed,
  voting_enabled
) VALUES (
  5,
  '2020-07-30',
  '<u>Goal:</u>&nbsp;<i>General Fitness</i>',
  '<b>Exercise(s):</b><br/><ul><li>100 Air Squats</li><li>50 Sit-Ups</li><li>25 Push-Ups</li><li>Run 3 Miles</li></ul>',
  '+1-617-455-7595',
  'JonathanZaleski@gmail.com',
  0,
  true
) ON CONFLICT DO NOTHING;

INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed,
  voting_enabled
) VALUES (
  6,
  '2020-07-31',
  '<u>Goal:</u>&nbsp;<i>Aerobic &amp; Strength Training</i>',
  '<b>Warmup:</b><br/><ul><li>Run 1 Mile</li></ul><b>Exercise(s):</b><br/><ul><li><span>50&nbsp;<a href="https://www.youtube.com/watch?v=Xjo_fY9Hl9w" style="color: #00f;" target="_blank">Dumbbell Goblet Squats</a></span></li><li><span>50&nbsp;<a href="https://www.youtube.com/watch?v=D7KaRcUTQeE" style="color: #00f;" target="_blank">Dumbbell Lunges</a></span></li><li><span>50&nbsp;<a href="https://www.youtube.com/watch?v=LktGPg-AkvY" style="color: #00f;" target="_blank">Dumbbell Bent-Over Rows</a></span></li><li><span>25&nbsp;<a href="https://www.youtube.com/watch?v=ir5PsbniVSc" style="color: #00f;" target="_blank">Dumbbell Skull Crushers</a></span></li><li>25 Dumbbell Bench Presses</li><li><span>25&nbsp;<a href="https://www.youtube.com/watch?v=RT_MTXaLKxU" style="color: #00f;" target="_blank">Dumbbell Squat Cleans</a></span></li></ul><b>Cooldown:</b><br/><ul><li>Run 1 Mile</li></ul>',
  '+1-617-455-7595',
  'JonathanZaleski@gmail.com',
  0,
  true
) ON CONFLICT DO NOTHING;

INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed,
  voting_enabled
) VALUES (
  7,
  '2020-08-01',
  '<u>Goal:</u>&nbsp;<i>Aerobic Training / General Fitness</i>',
  '<b>Exercise(s):</b><br/><ul><li>Run 1 Mile</li><li>50 Sit-Ups</li><li>45 Jumping Jacks</li><li>20 Lunges (10 on each side)</li><li>35 Air Squats</li><li>Run in place for 60 seconds</li><li>Plank for 60 seconds</li><li style="color: green; font-weight: 500;">Rest for 60 seconds</li><li>50 Sit-Ups</li><li>45 Jumping Jacks</li><li>20 Lunges (10 on each side)</li><li>35 Air Squats</li><li>40 Knee Highs (20 on each side)</li><li style="color: green; font-weight: 500;">Rest for 60 seconds</li><li>45 Jumping Jacks</li><li>20 Lunges (10 on each side)</li><li>40 Donkey Kicks (20 on each side)</li><li>20 Lunges (10 on each side)</li></ul>',
  '+1-617-455-7595',
  'JonathanZaleski@gmail.com',
  0,
  true
) ON CONFLICT DO NOTHING;

INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed,
  voting_enabled
) VALUES (
  8,
  '2020-08-02',
  '<u>Goal:</u>&nbsp;General Fitness</i>',
  '<b>Exercise(s):</b><br/><ul><li>3 Wall Sits (1 minute each)</li><li>30 Air Squats</li><li>30 Sit-Ups</li><li>30 Push-Ups</li><li>3 Uphill Sprints (~100m)</li><li>Run 3 Miles</li></ul>',
  '+1-617-455-7595',
  'JonathanZaleski@gmail.com',
  0,
  true
) ON CONFLICT DO NOTHING;

INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed,
  voting_enabled
) VALUES (
  9,
  '2020-08-03',
  '<u>Goal:</u>&nbsp;<i>Rest Day</i>',
  '<b>Exercise(s):</b><br/><br/>No Workout Today - Enjoy the day off!<br/><br/>',
  null,
  null,
  0,
  false
) ON CONFLICT DO NOTHING;

INSERT INTO workout (
  id,
  date,
  goal,
  description,
  sms_to,
  mail_to,
  completed,
  voting_enabled
) VALUES (
  10,
  '2020-08-04',
  '<center><i>2 song Tuesday...</i></center>',
  '<b>Exercise(s):</b><br><ul><li><a href="https://www.youtube.com/watch?v=6A2V9Bu80J4" style="color: #00f;" target="_blank">Flower, Moby</a> - Every time it says "Sally up" you stand, every time it says "Sally down" you squat and hold until it says "Sally up" again</li><li style="color: green; font-weight: 500;">Rest for 3 minutes</li><li><a href="https://www.youtube.com/watch?v=v2AC41dglnM" style="color: #00f;" target="_blank">Thunderstruck, AC/DC</a> - Every time you hear "thunder", drop down and do a burpee. And not the wimpy kind, make sure it has a push-up at the bottom and a full jump and clap above your head at the top!</li></ul>',
  '+1-617-455-7595',
  'JonathanZaleski@gmail.com',
  0,
  true
) ON CONFLICT DO NOTHING;
