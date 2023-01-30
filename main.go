/* Package declaration */

package main


/* Import(s) */

import (
  "context"
  "fmt"
  "html/template"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"
  "github.com/gin-gonic/gin"
  "github.com/jackc/pgx/v4"
)


/* Struct(s) */

type Workout struct {
  Date string
  Goal template.HTML
  Description template.HTML
  SmsTo string
  MailTo string
  MarkedCompleted bool
  Completed int
  VotingEnabled bool
  QuestionsEnabled bool
  Historical bool
}


/* Constant(s) */

const ANY_IPV4_ADDRESS = "0.0.0.0"
const BIND_ADDRESS_TEMPLATE = "%s:%s"
const COMMA = ","
const CURRENT = "current"
const DATABASE_URL_KEY = "DATABASE_URL"
const DEBUG_KEY = "DEBUG"
const DEFAULT_DATE = "1970-01-01"
const DEFAULT_INTERFACE = "localhost"
const DEFAULT_PORT = "5001"
const EMPTY = ""
const ENV_OR_PANIC_MESSAGE_TEMPLATE = `Key: "%s" was not found in the environment`
const FALSE = "false"
const HISTORICAL = "all"
const HISTORICAL_WORKOUTS_HTML_TEMPLATE = "historical-workouts.html"
const INTERFACE_KEY = "INTERFACE"
const LOCALHOST = "localhost"
const PORT_KEY = "PORT"
const SERVER_PUBLIC_ADDRESS_KEY = "SERVER_PUBLIC_ADDRESS"
const SESSION_COOKIE = "_wod"
const SINGLE_WORKOUT_HTML_TEMPLATE = "single-workout.html"
const TRUE = "true"
const TRUSTED_PROXIES_KEY = "TRUSTED_PROXIES"
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
  return envOrDefault(INTERFACE_KEY, DEFAULT_INTERFACE)
}

func bindPort() string {
  return envOrDefault(PORT_KEY, DEFAULT_PORT)
}

func cookieDomain() string {
  if gin.Mode() == gin.ReleaseMode {
    return envOrPanic(SERVER_PUBLIC_ADDRESS_KEY)
  }
  return LOCALHOST
}

func cookieExists(ginContext *gin.Context) bool {
  _, cookieError := ginContext.Cookie(cookieName())
  return cookieError == nil
}

func cookieValue(ginContext *gin.Context) int64 {
  if (!cookieExists(ginContext)) {
    return 0
  }

  cookieValue, cookieError := ginContext.Cookie(cookieName())
  if (cookieError != nil) {
    panic(cookieError)
  }

  parseIntValue, parseIntError := strconv.ParseInt(cookieValue, 10, 64)
  if (parseIntError != nil) {
    panic(parseIntError)
  }

  return parseIntValue
}

func cookieName() string {
  return time.Now().Format(WORKOUT_DATE_FORMAT)
}

func databaseConnection() *pgx.Conn {
  databaseConnection, databaseConnectionError := pgx.Connect(context.Background(), databaseUrl())
  if databaseConnectionError != nil {
    panic(databaseConnectionError)
  }
  return databaseConnection
}

func databaseUrl() string {
  return envOrPanic(DATABASE_URL_KEY)
}

func debug() bool {
  return envOrDefault(DEBUG_KEY, FALSE) == TRUE
}

func envOrDefault(key string, defaultValue string) string {
  result, found := os.LookupEnv(key)
  if found {
    return result
  }
  return defaultValue
}

func envOrPanic(key string) string {
  result, found := os.LookupEnv(key)
  if !found {
    panic(fmt.Sprintf(ENV_OR_PANIC_MESSAGE_TEMPLATE, key))
  }
  return result
}

func getHistoricalWorkouts(ginContext *gin.Context) []Workout {
  databaseConnection := databaseConnection()

  var rows, queryError = databaseConnection.Query(
    context.Background(),
    `
    SELECT
      date,
      goal,
      description,
      COALESCE(sms_to, '') AS sms_to,
      COALESCE(mail_to, '') AS mail_to,
      false AS marked_completed,
      completed,
      voting_enabled,
      true AS historical
    FROM workout
    WHERE date <> $1::DATE
    ORDER BY date DESC
    `,
    DEFAULT_DATE,
  )

  defer databaseConnection.Close(context.Background())

  if queryError != nil {
    panic(queryError)
  }

  var workoutCompleted int
  var workoutDate time.Time
  var workoutGoal, workoutDescription, workoutSmsTo, workoutMailTo string
  var workoutMarkedCompleted, workoutVotingEnabled, workoutHistorical bool
  var workout Workout

  var workouts []Workout
  for rows.Next() {
    rows.Scan(
      &workoutDate,
      &workoutGoal,
      &workoutDescription,
      &workoutSmsTo,
      &workoutMailTo,
      &workoutMarkedCompleted,
      &workoutCompleted,
      &workoutVotingEnabled,
      &workoutHistorical,
    )

    workout = Workout{
      Date: workoutDate.Format(WORKOUT_DATE_FORMAT),
      Goal: template.HTML(strings.TrimSpace(workoutGoal)),
      Description: template.HTML(strings.TrimSpace(workoutDescription)),
      SmsTo: strings.TrimSpace(workoutSmsTo),
      MailTo: strings.TrimSpace(workoutMailTo),
      MarkedCompleted: workoutMarkedCompleted,
      Completed: workoutCompleted,
      VotingEnabled: workoutVotingEnabled,
      QuestionsEnabled: len(workoutMailTo) > 0 || len(workoutSmsTo) > 0,
      Historical: workoutHistorical,
    }

    workouts = append(workouts, workout)
  }

  return workouts
}

func getWorkout(ginContext *gin.Context, workoutDateParam string) Workout {
  databaseConnection := databaseConnection()

  var dateFilterAndDisplayDate string
  if workoutDateParam == EMPTY || workoutDateParam == CURRENT {
    dateFilterAndDisplayDate = "NOW()"
  } else {
    dateFilterAndDisplayDate = workoutDateParam
  }

  var workoutCompleted int
  var workoutDate time.Time
  var workoutGoal, workoutDescription, workoutSmsTo, workoutMailTo string
  var workoutVotingEnabled, workoutHistorical bool

  var queryRowError = databaseConnection.QueryRow(
    context.Background(),
    `
    SELECT
      GREATEST(date, $1::DATE),
      goal,
      description,
      COALESCE(sms_to, '') AS sms_to,
      COALESCE(mail_to, '') AS mail_to,
      completed,
      voting_enabled,
      $1::DATE < NOW()::DATE AS historical
    FROM workout
    WHERE date::DATE = $1::DATE OR date = $2::DATE
    ORDER BY date DESC
    LIMIT 1
    `,
    dateFilterAndDisplayDate,
    DEFAULT_DATE,
  ).Scan(
    &workoutDate,
    &workoutGoal,
    &workoutDescription,
    &workoutSmsTo,
    &workoutMailTo,
    &workoutCompleted,
    &workoutVotingEnabled,
    &workoutHistorical,
  )

  defer databaseConnection.Close(context.Background())

  if queryRowError != nil {
    panic(queryRowError)
  }

  return Workout{
    Date: workoutDate.Format(WORKOUT_DATE_FORMAT),
    Goal: template.HTML(strings.TrimSpace(workoutGoal)),
    Description: template.HTML(strings.TrimSpace(workoutDescription)),
    SmsTo: strings.TrimSpace(workoutSmsTo),
    MailTo: strings.TrimSpace(workoutMailTo),
    MarkedCompleted: cookieExists(ginContext) && workoutCompleted > 0,
    Completed: workoutCompleted,
    VotingEnabled: workoutVotingEnabled,
    QuestionsEnabled: len(workoutMailTo) > 0 || len(workoutSmsTo) > 0,
    Historical: workoutHistorical,
  }
}

func ginMode() string {
  if debug() {
    return gin.DebugMode
  }
  return gin.ReleaseMode
}

func trustedProxies() []string {
  trustedProxies := envOrDefault(TRUSTED_PROXIES_KEY, EMPTY)
  if trustedProxies == EMPTY {
    return nil
  }
  return strings.Split(trustedProxies, COMMA)
}


/* Handler(s) */

func currentWorkoutHandler(ginContext *gin.Context) {
  ginContext.Redirect(http.StatusFound, "/workout/current")
}

func workoutCompletedHandler(ginContext *gin.Context) {
  workoutDate := ginContext.Param("workoutDate")

  databaseConnection := databaseConnection()

  _, insertError := databaseConnection.Exec(
    context.Background(),
    `
    INSERT INTO workout (
      SELECT
        NOW()::DATE,
        goal,
        description,
        sms_to,
        mail_to,
        0,
        voting_enabled
      FROM workout
      WHERE date = $1::DATE
    ) ON CONFLICT DO NOTHING
    `,
    DEFAULT_DATE,
  )

  if insertError != nil {
    panic(insertError)
  }

  _, updateError := databaseConnection.Exec(
    context.Background(),
    `
    UPDATE workout
    SET completed = completed + 1
    WHERE date = $1
    `,
    workoutDate,
  )

  if updateError != nil {
    panic(updateError)
  }

  ginContext.SetCookie(
    cookieName(),
    fmt.Sprintf("%d", 1),
    86400,
    "/",
    cookieDomain(),
    false,
    true,
  )

  defer databaseConnection.Close(context.Background())

  ginContext.Redirect(http.StatusFound, "/workout/current")
}

func workoutMetaHandler(ginContext *gin.Context) {
  workoutDateParam := ginContext.Param("workoutDate")

  if workoutDateParam == HISTORICAL {
    ginContext.HTML(
      http.StatusOK,
      HISTORICAL_WORKOUTS_HTML_TEMPLATE,
      getHistoricalWorkouts(ginContext),
    )
  } else {
    ginContext.HTML(
      http.StatusOK,
      SINGLE_WORKOUT_HTML_TEMPLATE,
      getWorkout(ginContext, workoutDateParam),
    )
  }
}


/* Application entry-point */

func main() {
  gin.SetMode(ginMode())
  router := gin.New()
  router.SetTrustedProxies(trustedProxies())
  router.LoadHTMLGlob("templates/*.html")
  router.Use(gin.Logger(), gin.Recovery())
  router.GET("/", currentWorkoutHandler)
  router.GET("/workout/:workoutDate", workoutMetaHandler)
  router.POST("/workout/:workoutDate/completed", workoutCompletedHandler)
  router.StaticFile("/favicon.ico", "assets/favicon.ico")
  router.Static("/assets", "assets")
  router.Run(bindAddress())
}
