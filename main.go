/* Package declaration */

package main


/* Import(s) */

import (
  "fmt"
  "html/template"
  "net/http"
  "os"
  "strings"
  "time"
  "github.com/gin-gonic/gin"
)


/* Struct(s) */

type Workout struct {
  Date string
  Goal template.HTML
  Description template.HTML
  SmsTo string
  MailTo string
}

/* Constant(s) */

const ANY_IPV4_ADDRESS = "0.0.0.0"
const BIND_ADDRESS_TEMPLATE = "%s:%s"
const DEFAULT_PORT = "5000"
const DEFAULT_WORKOUT_DESCRIPTION = "<b>Warm-up:</b><br/><ul><li>Run 1 Mile</li></ul><b>Exercise(s):</b><br/><ul><li>99 Air Squats</li><li>66 Sit-Ups</li><li>33 Push-Ups</li></ul>"
const DEFAULT_WORKOUT_GOAL = "<u>Goal:</u>&nbsp;<i>General Fitness</i>"
const DEFAULT_WORKOUT_MAIL_TO = "JonathanZaleski@gmail.com"
const DEFAULT_WORKOUT_SMS_TO = "+1-617-455-7595"
const INDEX_HTML_TEMPLATE = "index.html.tmpl"
const LOCALHOST = "127.0.0.1"
const PORT_KEY = "PORT"
const WORKOUT_DATE_FORMAT = "2006-01-02"
const WORKOUT_DESCRIPTION_KEY = "WORKOUT_DESCRIPTION"
const WORKOUT_GOAL_KEY = "WORKOUT_GOAL"
const WORKOUT_MAIL_TO_KEY = "WORKOUT_MAIL_TO"
const WORKOUT_SMS_TO_KEY = "WORKOUT_SMS_TO"


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
  return Workout{
    Date: strings.TrimSpace(workoutDate()),
    Goal: template.HTML(strings.TrimSpace(workoutGoal())),
    Description: template.HTML(strings.TrimSpace(workoutDescription())),
    SmsTo: strings.TrimSpace(workoutSmsTo()),
    MailTo: strings.TrimSpace(workoutMailTo()),
  }
}

func envOrDefault(key string, defaultValue string) string {
  result, found := os.LookupEnv(key)
  if found {
    return result
  }
  return defaultValue
}

func workoutDate() string {
  return time.Now().Format(WORKOUT_DATE_FORMAT)
}

func workoutDescription() string {
  return envOrDefault(WORKOUT_DESCRIPTION_KEY, DEFAULT_WORKOUT_DESCRIPTION)
}

func workoutGoal() string {
  return envOrDefault(WORKOUT_GOAL_KEY, DEFAULT_WORKOUT_GOAL)
}

func workoutMailTo() string {
  return envOrDefault(WORKOUT_MAIL_TO_KEY, DEFAULT_WORKOUT_MAIL_TO)
}

func workoutSmsTo() string {
  return envOrDefault(WORKOUT_SMS_TO_KEY, DEFAULT_WORKOUT_SMS_TO)
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
