package postgres

import (
	"github.com/sirupsen/logrus"

	"github.com/dselans/welcome-svc/backends/db"
)

type Postgres struct {
	DB      *db.Storage
	options *Options
	log     *logrus.Entry

	// Models will go here
}

type Options struct {
	Host          string
	Name          string
	User          string
	Pass          string
	Port          int
	SSLOptions    *SSLOptions
	RunMigrations bool
	MigrationsDir string
}

type SSLOptions struct{}

func New(db *db.Storage) *Postgres {

	return &Postgres{
		DB:  db,
		log: logrus.WithField("pkg", "postgres"),
	}

}
