package config

import "os"
import _ "github.com/joho/godotenv/autoload"

var VATUSA_API2_URL = os.Getenv("VATUSA_API2_URL")
