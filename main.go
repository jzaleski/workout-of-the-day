/* Package declaration */

package main


/* Import(s) */

import (
  "database/sql"
  "fmt"
  "html/template"
  "net/http"
  "os"
  "strings"
  "time"
  "github.com/gin-gonic/gin"
  _ "github.com/lib/pq"
)


/* Struct(s) */

type Workout struct {
  Id int
  Date string
  Goal template.HTML
  Description template.HTML
  SmsTo string
  MailTo string
}

/* Constant(s) */

const ANY_IPV4_ADDRESS = "0.0.0.0"
const BIND_ADDRESS_TEMPLATE = "%s:%s"
const DATABASE_URL_KEY = "DATABASE_URL"
const DEFAULT_DATABASE_URL = "postgres://localhost:5432/workout_of_the_day"
const DEFAULT_PORT = "5000"
const DEFAULT_WORKOUT_DESCRIPTION = "<b>Exercise(s):</b><br/><ul><li>150 Air Squats</li><li>100 Sit-Ups</li><li>50 Push-Ups</li><li>Run 3 Miles</li></ul>"
const DEFAULT_WORKOUT_DESCRIPTION_KEY = "DEFAULT_WORKOUT_DESCRIPTION"
const DEFAULT_WORKOUT_GOAL = "<u>Goal:</u>&nbsp;<i>General Fitness</i>"
const DEFAULT_WORKOUT_GOAL_KEY = "DEFAULT_WORKOUT_GOAL"
const DEFAULT_WORKOUT_MAIL_TO = "JonathanZaleski@gmail.com"
const DEFAULT_WORKOUT_MAIL_TO_KEY = "DEFAULT_WORKOUT_MAIL_TO"
const DEFAULT_WORKOUT_SMS_TO = "+1-617-455-7595"
const DEFAULT_WORKOUT_SMS_TO_KEY = "DEFAULT_WORKOUT_SMS_TO"
const INDEX_HTML_TEMPLATE = "index.html.tmpl"
const LOCALHOST = "127.0.0.1"
const PORT_KEY = "PORT"
const SSLMODE_SUFFIX = "?sslmode=disable"
const WORKOUT_DATE_FORMAT = "2006-01-02"


/* Helper(s) */

func bindAddress() string {
  return fmt.Sprintf(
    BIND_ADDRESS_TEMPLATE,
    bindInterface(),
    bindPort(),
  )
}

func bindInterface() string {
  if gin.Mode() == gin.ReleaseMode {
    return ANY_IPV4_ADDRESS
  }
  return LOCALHOST
}

func bindPort() string {
  return envOrDefault(PORT_KEY, DEFAULT_PORT)
}

func currentWorkout() Workout {
  databaseConnection := databaseConnection()

  var workoutId int
  var workoutDate time.Time
  var workoutGoal, workoutDescription, workoutSmsTo, workoutMailTo string

  var queryRowError = databaseConnection.QueryRow(`
    SELECT
      id,
      date,
      goal,
      description,
      sms_to,
      mail_to
    FROM workout
    WHERE date::DATE = NOW()::DATE
    LIMIT 1`,
  ).Scan(
    &workoutId,
    &workoutDate,
    &workoutGoal,
    &workoutDescription,
    &workoutSmsTo,
    &workoutMailTo,
  )

  if queryRowError != nil {
    workoutId = defaultWorkoutId()
    workoutDate = defaultWorkoutDate()
    workoutGoal = defaultWorkoutGoal()
    workoutDescription = defaultWorkoutDescription()
    workoutSmsTo = defaultWorkoutSmsTo()
    workoutMailTo = defaultWorkoutMailTo()
  }

  defer databaseConnection.Close()

  return Workout{
    Id: workoutId,
    Date: workoutDate.Format(WORKOUT_DATE_FORMAT),
    Goal: template.HTML(strings.TrimSpace(workoutGoal)),
    Description: template.HTML(strings.TrimSpace(workoutDescription)),
    SmsTo: strings.TrimSpace(workoutSmsTo),
    MailTo: strings.TrimSpace(workoutMailTo),
  }
}

func databaseConnection() *sql.DB {
  databaseConnection, databaseConnectionError := sql.Open("postgres", databaseUrl())
  if databaseConnectionError != nil {
    panic(databaseConnectionError)
  }
  return databaseConnection
}

func databaseUrl() string {
  var databaseUrl = envOrDefault(DATABASE_URL_KEY, DEFAULT_DATABASE_URL)
  if !strings.HasSuffix(databaseUrl, SSLMODE_SUFFIX) {
    return databaseUrl + SSLMODE_SUFFIX
  }
  return databaseUrl
}

func defaultWorkoutDate() time.Time {
  return time.Now()
}

func defaultWorkoutDescription() string {
  return envOrDefault(DEFAULT_WORKOUT_DESCRIPTION_KEY, DEFAULT_WORKOUT_DESCRIPTION)
}

func defaultWorkoutGoal() string {
  return envOrDefault(DEFAULT_WORKOUT_GOAL_KEY, DEFAULT_WORKOUT_GOAL)
}

func defaultWorkoutId() int {
  return -1
}

func defaultWorkoutMailTo() string {
  return envOrDefault(DEFAULT_WORKOUT_MAIL_TO_KEY, DEFAULT_WORKOUT_MAIL_TO)
}

func defaultWorkoutSmsTo() string {
  return envOrDefault(DEFAULT_WORKOUT_SMS_TO_KEY, DEFAULT_WORKOUT_SMS_TO)
}

func envOrDefault(key string, defaultValue string) string {
  result, found := os.LookupEnv(key)
  if found {
    return result
  }
  return defaultValue
}

/* Handler(s) */

func indexHandler(context *gin.Context) {
  context.HTML(
    http.StatusOK,
    INDEX_HTML_TEMPLATE,
    currentWorkout(),
  )
}


/* Application entry-point */

func main() {
  router := gin.New()
  router.LoadHTMLGlob("templates/*.tmpl")
  router.Use(gin.Logger(), gin.Recovery())
  router.GET("/", indexHandler)
  router.StaticFile("/favicon.ico", "assets/favicon.ico")
  router.Static("/assets", "assets")
  router.Run(bindAddress())
}
