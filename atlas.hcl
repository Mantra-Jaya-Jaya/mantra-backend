data "external_schema" "gorm" {
  program = [
    "go", "run", "-mod=mod",
    "ariga.io/atlas-provider-gorm", "load",
    "--path", "./models", "--dialect", "postgres",
  ]
}

variable "db_url" {
  type    = string
  default = "postgres://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@${getenv("DB_HOST")}:${getenv("DB_PORT")}/${getenv("DB_NAME")}?sslmode=disable"
}

variable "dev_url" {
  type    = string
  default = "postgres://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@${getenv("DB_HOST")}:${getenv("DB_PORT")}/mantra_dev?sslmode=disable"
}

env "local" {
  url  = var.db_url
  from = var.db_url
  to   = data.external_schema.gorm.url
  dev  = var.dev_url
}