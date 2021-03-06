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
  MarkedCompleted bool
  Completed int
  VotingEnabled bool
  QuestionsEnabled bool
}

/* Constant(s) */

const ANY_IPV4_ADDRESS = "0.0.0.0"
const BIND_ADDRESS_TEMPLATE = "%s:%s"
const DATABASE_URL_KEY = "DATABASE_URL"
const DEFAULT_DATABASE_URL = "postgres://localhost:5432/workout_of_the_day"
const DEFAULT_PORT = "5000"
const LOCALHOST = "localhost"
const PORT_KEY = "PORT"
const PRODUCTION_DOMAIN = "wod.jzaleski.com"
const SESSION_COOKIE = "_wod"
const SSLMODE_SUFFIX = "?sslmode=disable"
const WORKOUT_DATE_FORMAT = "2006-01-02"
const WORKOUT_HTML_TEMPLATE = "workout.html.tmpl"


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

func cookieExists(context *gin.Context) bool {
  _, cookieError := context.Cookie(cookieName())
  return cookieError == nil
}

func cookieName() string {
  return time.Now().Format(WORKOUT_DATE_FORMAT)
}

func getWorkout(context *gin.Context) Workout {
  databaseConnection := databaseConnection()

  var dateFilterAndDisplayDate string
  workoutDateParam := context.Param("workoutDate")
  if workoutDateParam == "" || workoutDateParam == "current" {
    dateFilterAndDisplayDate = "NOW()"
  } else {
    dateFilterAndDisplayDate = fmt.Sprintf("'%s'", workoutDateParam)
  }

  var workoutId, workoutCompleted int
  var workoutDate time.Time
  var workoutGoal, workoutDescription, workoutSmsTo, workoutMailTo string
  var workoutVotingEnabled bool

  var queryRowError = databaseConnection.QueryRow(`
    SELECT
      id,
      GREATEST(date, ` + dateFilterAndDisplayDate + `),
      goal,
      description,
      COALESCE(sms_to, ''),
      COALESCE(mail_to, ''),
      completed,
      voting_enabled
    FROM workout
    WHERE date::DATE = ` + dateFilterAndDisplayDate + `::DATE OR id = 1
    ORDER BY date DESC, id DESC
    LIMIT 1`,
  ).Scan(
    &workoutId,
    &workoutDate,
    &workoutGoal,
    &workoutDescription,
    &workoutSmsTo,
    &workoutMailTo,
    &workoutCompleted,
    &workoutVotingEnabled,
  )

  if queryRowError != nil {
    panic(queryRowError)
  }

  defer databaseConnection.Close()

  return Workout{
    Id: workoutId,
    Date: workoutDate.Format(WORKOUT_DATE_FORMAT),
    Goal: template.HTML(strings.TrimSpace(workoutGoal)),
    Description: template.HTML(strings.TrimSpace(workoutDescription)),
    SmsTo: strings.TrimSpace(workoutSmsTo),
    MailTo: strings.TrimSpace(workoutMailTo),
    MarkedCompleted: cookieExists(context),
    Completed: workoutCompleted,
    VotingEnabled: (workoutDateParam == "" || workoutDateParam == "current") && workoutVotingEnabled,
    QuestionsEnabled: len(workoutMailTo) > 0 || len(workoutSmsTo) > 0,
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

func envOrDefault(key string, defaultValue string) string {
  result, found := os.LookupEnv(key)
  if found {
    return result
  }
  return defaultValue
}

/* Handler(s) */

func currentWorkoutHandler(context *gin.Context) {
  context.Redirect(http.StatusFound, "/workout/current")
}

func workoutCompletedHandler(context *gin.Context) {
  if !cookieExists(context) {
    workoutId := context.Param("workoutId")

    databaseConnection := databaseConnection()

    _, execError := databaseConnection.Exec(`
      UPDATE workout
      SET completed = completed + 1
      WHERE id = $1`,
      workoutId,
    )

    if execError != nil {
      panic(execError)
    }

    context.SetCookie(
      cookieName(),
      "completed",
      86400,
      "/",
      ".",
      false,
      true,
    )

    defer databaseConnection.Close()
  }

  context.Redirect(http.StatusFound, "/workout/current")
}

func workoutForDateHandler(context *gin.Context) {
  context.HTML(
    http.StatusOK,
    WORKOUT_HTML_TEMPLATE,
    getWorkout(context),
  )
}


/* Application entry-point */

func main() {
  router := gin.New()
  router.LoadHTMLGlob("templates/*.tmpl")
  router.Use(gin.Logger(), gin.Recovery())
  router.GET("/", currentWorkoutHandler)
  router.GET("/workout/:workoutDate", workoutForDateHandler)
  router.POST("/workout/:workoutId/completed", workoutCompletedHandler)
  router.StaticFile("/favicon.ico", "assets/favicon.ico")
  router.Static("/assets", "assets")
  router.Run(bindAddress())
}
