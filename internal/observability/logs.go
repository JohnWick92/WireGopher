package observability

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"time"
)

var syslogConn net.Conn
var loc = time.Local
var hostname string

func Init(host string, port string) error {
	var err error

	hostname, err = os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	address := net.JoinHostPort(host, port)

	syslogConn, err = net.Dial("udp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to syslog at %s: %w", address, err)
	}

	return nil
}

func formatSyslogRFC5424(priority int, appName, message string) string {
	timestamp := time.Now().In(loc)

	return fmt.Sprintf(
		"<%d>1 %s %s %s %s - - %s",
		priority,
		timestamp.Format("2006-01-02T15:04:05-07:00"),
		hostname,
		appName,
		strconv.Itoa(os.Getpid()),
		message,
	)
}

func LogInfo(msg string) {
	timestamp := time.Now().In(loc)
	processedMsg := fmt.Sprintf(
		"[INFO] %s : [TIMESTAMP] %s",
		msg,
		timestamp.Format("2006-01-02 15:04:05"),
	)

	rfc5424 := formatSyslogRFC5424(14, "WireGopher", processedMsg)

	if _, err := syslogConn.Write([]byte(rfc5424)); err != nil {
		slog.Error("Failed to send syslog", "err", err)
	}
}

func LogAuth() {}

func LogLogout() {}
