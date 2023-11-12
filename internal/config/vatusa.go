package config

import "os"
import _ "github.com/joho/godotenv/autoload"

var VATUSA_API2_URL = os.Getenv("VATUSA_API2_URL")
var VATUSA_API2_KEY = os.Getenv("VATUSA_API2_KEY")
var SERVICE_ACCOUNT_JSON = os.Getenv("SERVICE_ACCOUNT_JSON")
var GOOGLE_CUSTOMER_ID = os.Getenv("GOOGLE_CUSTOMER_ID")

var MAIL_HOST = os.Getenv("MAIL_HOST")
var MAIL_PORT = os.Getenv("MAIL_PORT")
var MAIL_USERNAME = os.Getenv("MAIL_USERNAME")
var MAIL_PASSWORD = os.Getenv("MAIL_PASSWORD")
