package pqtimeouts

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

func init() {
	sql.Register("pq-timeouts", timeoutDriver{dialOpen: pq.DialOpen})
}

type timeoutDriver struct {
	dialOpen func(pq.Dialer, string) (driver.Conn, error) // Allow this to be stubbed for testing
}

func (t timeoutDriver) Open(connection string) (_ driver.Conn, err error) {
	// Look for read_timeout and write_timeout in the connection string and extract the values.
	// read_timeout and write_timeout need to be removed from the connection string before calling pq as well.
	var newConnectionSettings []string
	var readTimeout time.Duration
	var writeTimeout time.Duration

	// If the connection is specified as a URL, use the parsing function in lib/pq to turn it into options.
	if strings.HasPrefix(connection, "postgres://") || strings.HasPrefix(connection, "postgresql://") {
		connection, err = pq.ParseURL(connection)
		if err != nil {
			return nil, err
		}
	}
	fmt.Printf("%s\n",connection)
	for _, setting := range strings.Fields(connection) {
		s := strings.Split(setting, "=")
		fmt.Printf("Strings split: %s\n",s)
		if s[0] == "read_timeout" {
			fmt.Printf("Value: %s\n",s[1])
			val, err := strconv.Atoi(strings.Trim(s[1],"'\""))
			fmt.Printf("%s %s\n",val,err)
			if err != nil {
				return nil, fmt.Errorf("Error interpreting value for read_timeout")
			}
			readTimeout = time.Duration(val) * time.Millisecond // timeout is in milliseconds
		} else if s[0] == "write_timeout" {
			val, err := strconv.Atoi(strings.Trim(s[1],"'\""))
			if err != nil {
				return nil, fmt.Errorf("Error interpreting value for write_timeout")
			}
			writeTimeout = time.Duration(val) * time.Millisecond // timeout is in milliseconds
		} else {
			newConnectionSettings = append(newConnectionSettings, setting)
		}
	}

	newConnectionStr := strings.Join(newConnectionSettings, " ")

	return t.dialOpen(
		timeoutDialer{
			netDial:        net.Dial,
			netDialTimeout: net.DialTimeout,
			readTimeout:    readTimeout,
			writeTimeout:   writeTimeout},
		newConnectionStr)
}
