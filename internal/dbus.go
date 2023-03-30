package internal

import (
	"github.com/godbus/dbus/v5"
)

type DbusConnection struct {
	Inner *dbus.Conn
}

func DbusConnect() (*DbusConnection, error) {
	connection := new(DbusConnection)

	inner, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, err
	}

	connection.Inner = inner

	return connection, nil
}

func (c *DbusConnection) Object(dest string, name string) dbus.BusObject {
	object := c.Inner.Object(dest, dbus.ObjectPath(name))
	return object
}

func (c *DbusConnection) GetActiveSession() (dbus.BusObject, error) {
	userObj := c.Object("org.freedesktop.login1", "/org/freedesktop/login1/user/self")

	rawSessions, err := userObj.GetProperty("org.freedesktop.login1.User.Sessions")
	if err != nil {
		return nil, err
	}

	var sessions []struct {
		Index string
		Name  string
	}

	err = rawSessions.Store(&sessions)
	if err != nil {
		return nil, err
	}

	var activeSessions []dbus.BusObject

	for _, s := range sessions {
		sessionObj := c.Object("org.freedesktop.login1", s.Name)

		rawState, err := sessionObj.GetProperty("org.freedesktop.login1.Session.State")
		if err != nil {
			return nil, err
		}

		var state string
		if err := rawState.Store(&state); err != nil {
			return nil, err
		}

		if state != "active" {
			continue
		}

		activeSessions = append(activeSessions, sessionObj)
	}

	return activeSessions[0], nil
}
