package util

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// USAGE List down all health check for services
func HealthCheck(w http.ResponseWriter, r *http.Request, ac *AppContext, startTimestamp time.Time, m map[string]interface{}) {
	status := true
	var statusError string

	postgresPingSuccess := true
	if err := Ping(ac.DB.DB); err != nil {
		postgresPingSuccess = false

		status = false
		statusError = fmt.Sprintf("postgres ping failed: %s", err.Error())
	}

	uptime := time.Since(startTimestamp)
	for key, value := range map[string]interface{}{
		"postgres_ping_success": postgresPingSuccess,
		"uptime_seconds":        int64(uptime.Seconds()),
		"uptime":                uptime.String(),
	} {
		m[key] = value
	}

	statusCode := http.StatusOK
	if !status {
		statusCode = http.StatusInternalServerError
		m["error"] = statusError
	}

	WriteJSONWithStatus(w, m, statusCode)
}

// USAGE: Check DB Connections
func Ping(conn *sql.DB) error {
	if conn == nil {
		return errors.New("no db connection")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return conn.PingContext(ctx)
}
