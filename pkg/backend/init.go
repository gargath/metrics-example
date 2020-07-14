package backend

import (
	"database/sql"
	"time"

	sqldbstats "github.com/krpn/go-sql-db-stats"
	// import SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	dbDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "metrics_example_backend_query_duration_seconds",
		Help: "Duration of database queries.",
	}, []string{"resource"})
)

func init() {
	db, err := sql.Open("sqlite3", "./anaximander.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users ( id TEXT PRIMARY KEY, name TEXT NOT NULL, dob TEXT NOT NULL, address, TEXT );")
	if err != nil {
		panic(err)
	}

	_ = sqldbstats.StartCollectPrometheusMetrics(db, 30*time.Second, "entity_db")

}
