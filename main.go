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

const ALL = "all"
const ALL_WORKOUTS_HTML_TEMPLATE = "workout-list.html.tmpl"
const ANY_IPV4_ADDRESS = "0.0.0.0"
const BIND_ADDRESS_TEMPLATE = "%s:%s"
const COMMA = ","
const CURRENT = "current"
const DATABASE_URL_KEY = "DATABASE_URL"
const DEBUG_KEY = "DEBUG"
const DEFAULT_INTERFACE = "localhost"
const DEFAULT_PORT = "5001"
const EMPTY = ""
const ENV_OR_PANIC_MESSAGE_TEMPLATE = `Key: "%s" was not found in the environment`;
const FALSE = "false"
const INTERFACE_KEY = "INTERFACE"
const LOCALHOST = "localhost"
const PORT_KEY = "PORT"
const SERVER_PUBLIC_ADDRESS_KEY = "SERVER_PUBLIC_ADDRESS"
const SESSION_COOKIE = "_wod"
const TRUE = "true"
const TRUSTED_PROXIES_KEY = "TRUSTED_PROXIES"
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
    return 0;
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

func getAllWorkouts(ginContext *gin.Context) []Workout {
  databaseConnection := databaseConnection()


  var rows, queryError = databaseConnection.Query(
    context.Background(),
    `
      SELECT
        id,
        date,
        goal,
        description,
        COALESCE(sms_to, '') AS sms_to,
        COALESCE(mail_to, '') AS mail_to,
        false AS marked_completed,
        completed,
        false AS voting_enabled
      FROM workout
      WHERE id > 1
      ORDER BY date DESC
    `,
  )

  defer databaseConnection.Close(context.Background())

  if queryError != nil {
    panic(queryError)
  }

  var workoutId, workoutCompleted int
  var workoutDate time.Time
  var workoutGoal, workoutDescription, workoutSmsTo, workoutMailTo string
  var workoutMarkedCompleted, workoutVotingEnabled bool
  var workout Workout

  var workouts []Workout
  for rows.Next() {
    rows.Scan(
      &workoutId,
      &workoutDate,
      &workoutGoal,
      &workoutDescription,
      &workoutSmsTo,
      &workoutMailTo,
      &workoutMarkedCompleted,
      &workoutCompleted,
      &workoutVotingEnabled,
    )

    workout = Workout{
      Id: workoutId,
      Date: workoutDate.Format(WORKOUT_DATE_FORMAT),
      Goal: template.HTML(strings.TrimSpace(workoutGoal)),
      Description: template.HTML(strings.TrimSpace(workoutDescription)),
      SmsTo: strings.TrimSpace(workoutSmsTo),
      MailTo: strings.TrimSpace(workoutMailTo),
      MarkedCompleted: workoutMarkedCompleted,
      Completed: workoutCompleted,
      VotingEnabled: workoutVotingEnabled,
      QuestionsEnabled: len(workoutMailTo) > 0 || len(workoutSmsTo) > 0,
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
    dateFilterAndDisplayDate = fmt.Sprintf("'%s'", workoutDateParam)
  }

  var workoutId, workoutCompleted int
  var workoutDate time.Time
  var workoutGoal, workoutDescription, workoutSmsTo, workoutMailTo string
  var workoutVotingEnabled bool

  var queryRowError = databaseConnection.QueryRow(
    context.Background(),
    `
    SELECT
      id,
      GREATEST(date, ` + dateFilterAndDisplayDate + `),
      goal,
      description,
      COALESCE(sms_to, '') AS sms_to,
      COALESCE(mail_to, '') AS mail_to,
      completed,
      voting_enabled
    FROM workout
    WHERE date::DATE = ` + dateFilterAndDisplayDate + `::DATE OR id = 1
    ORDER BY date DESC, id DESC
    LIMIT 1
    `,
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

  defer databaseConnection.Close(context.Background())

  if queryRowError != nil {
    panic(queryRowError)
  }

  return Workout{
    Id: workoutId,
    Date: workoutDate.Format(WORKOUT_DATE_FORMAT),
    Goal: template.HTML(strings.TrimSpace(workoutGoal)),
    Description: template.HTML(strings.TrimSpace(workoutDescription)),
    SmsTo: strings.TrimSpace(workoutSmsTo),
    MailTo: strings.TrimSpace(workoutMailTo),
    MarkedCompleted: cookieExists(ginContext),
    Completed: workoutCompleted,
    VotingEnabled: (workoutDateParam == EMPTY || workoutDateParam == CURRENT) && workoutVotingEnabled,
    QuestionsEnabled: len(workoutMailTo) > 0 || len(workoutSmsTo) > 0,
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
  cookieValue := cookieValue(ginContext);

  workoutId := ginContext.Param("workoutId")

  databaseConnection := databaseConnection()

  _, execError := databaseConnection.Exec(
    context.Background(),
    `
    UPDATE workout
    SET completed = completed + 1
    WHERE id = $1
    `,
    workoutId,
  )

  if execError != nil {
    panic(execError)
  }

  ginContext.SetCookie(
    cookieName(),
    fmt.Sprintf("%d", cookieValue + 1),
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

  if workoutDateParam == ALL {
    ginContext.HTML(
      http.StatusOK,
      ALL_WORKOUTS_HTML_TEMPLATE,
      getAllWorkouts(ginContext),
    )
  } else {
    ginContext.HTML(
      http.StatusOK,
      WORKOUT_HTML_TEMPLATE,
      getWorkout(ginContext, workoutDateParam),
    )
  }
}


/* Application entry-point */

func main() {
  gin.SetMode(ginMode())
  router := gin.New()
  router.SetTrustedProxies(trustedProxies())
  router.LoadHTMLGlob("templates/*.tmpl")
  router.Use(gin.Logger(), gin.Recovery())
  router.GET("/", currentWorkoutHandler)
  router.GET("/workout/:workoutDate", workoutMetaHandler)
  router.POST("/workout/:workoutId/completed", workoutCompletedHandler)
  router.StaticFile("/favicon.ico", "assets/favicon.ico")
  router.Static("/assets", "assets")
  router.Run(bindAddress())
}
