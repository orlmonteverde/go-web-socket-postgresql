package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/olahol/melody.v1"

	"github.com/labstack/echo"
	"github.com/lib/pq"
)

func main() {
	e := echo.New()
	m := melody.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Web")
	})

	e.GET("/ws", func(c echo.Context) error {
		go connect(m)
		m.HandleRequest(c.Response(), c.Request())
		return nil
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	fmt.Println(e.Start(":8080"))
}

//connect to postgres and  activate listener for database events
func connect(m *melody.Melody) {
	connString := "dbname=name user=user password=pwd"

	_, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(connString, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")
	if err != nil {
		panic(err)
	}

	for {
		notificationListener(listener, m)
	}
}

//listener wait for db events
func notificationListener(l *pq.Listener, m *melody.Melody) {
	for {
		select {
		case n := <-l.Notify:
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
			if err != nil {
				fmt.Println("Error processing JSON: ", err)
				return
			}
			m.Broadcast(prettyJSON.Bytes())
			return
		}
	}
}
