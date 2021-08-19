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

const ANY_IPV4_ADDRESS = "0.0.0.0"
const BIND_ADDRESS_TEMPLATE = "%s:%s"
const DATABASE_URL_KEY = "DATABASE_URL"
const DEFAULT_DATABASE_URL = "postgres://postgres:postgres@localhost:5432/workout_of_the_day?sslmode=disable"
const DEFAULT_PORT = "5000"
const ENV_OR_PANIC_MESSAGE_TEMPLATE = `Key: "%s" was not found in the environment`;
const LOCALHOST = "localhost"
const PORT_KEY = "PORT"
const SERVER_PUBLIC_ADDRESS_KEY = "SERVER_PUBLIC_ADDRESS"
const SESSION_COOKIE = "_wod"
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

func getWorkout(ginContext *gin.Context) Workout {
  databaseConnection := databaseConnection()

  var dateFilterAndDisplayDate string
  workoutDateParam := ginContext.Param("workoutDate")
  if workoutDateParam == "" || workoutDateParam == "current" {
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
      COALESCE(sms_to, ''),
      COALESCE(mail_to, ''),
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

  if queryRowError != nil {
    panic(queryRowError)
  }

  defer databaseConnection.Close(context.Background())

  return Workout{
    Id: workoutId,
    Date: workoutDate.Format(WORKOUT_DATE_FORMAT),
    Goal: template.HTML(strings.TrimSpace(workoutGoal)),
    Description: template.HTML(strings.TrimSpace(workoutDescription)),
    SmsTo: strings.TrimSpace(workoutSmsTo),
    MailTo: strings.TrimSpace(workoutMailTo),
    MarkedCompleted: cookieExists(ginContext),
    Completed: workoutCompleted,
    VotingEnabled: (workoutDateParam == "" || workoutDateParam == "current") && workoutVotingEnabled,
    QuestionsEnabled: len(workoutMailTo) > 0 || len(workoutSmsTo) > 0,
  }
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

func workoutForDateHandler(ginContext *gin.Context) {
  ginContext.HTML(
    http.StatusOK,
    WORKOUT_HTML_TEMPLATE,
    getWorkout(ginContext),
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
