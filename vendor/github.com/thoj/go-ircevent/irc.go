// Copyright 2009 Thomas Jager <mail@jager.no>  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
This package provides an event based IRC client library. It allows to
register callbacks for the events you need to handle. Its features
include handling standard CTCP, reconnecting on errors and detecting
stones servers.
Details of the IRC protocol can be found in the following RFCs:
https://tools.ietf.org/html/rfc1459
https://tools.ietf.org/html/rfc2810
https://tools.ietf.org/html/rfc2811
https://tools.ietf.org/html/rfc2812
https://tools.ietf.org/html/rfc2813
The details of the client-to-client protocol (CTCP) can be found here: http://www.irchelp.org/irchelp/rfc/ctcpspec.html
*/

package irc

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	VERSION = "go-ircevent v2.1"
)

var ErrDisconnected = errors.New("Disconnect Called")

// Read data from a connection. To be used as a goroutine.
func (irc *Connection) readLoop() {
	defer irc.Done()
	br := bufio.NewReaderSize(irc.socket, 512)

	errChan := irc.ErrorChan()

	for {
		select {
		case <-irc.end:
			return
		default:
			// Set a read deadline based on the combined timeout and ping frequency
			// We should ALWAYS have received a response from the server within the timeout
			// after our own pings
			if irc.socket != nil {
				irc.socket.SetReadDeadline(time.Now().Add(irc.Timeout + irc.PingFreq))
			}

			msg, err := br.ReadString('\n')

			// We got past our blocking read, so bin timeout
			if irc.socket != nil {
				var zero time.Time
				irc.socket.SetReadDeadline(zero)
			}

			if err != nil {
				errChan <- err
				return
			}

			if irc.Debug {
				irc.Log.Printf("<-- %s\n", strings.TrimSpace(msg))
			}

			irc.lastMessageMutex.Lock()
			irc.lastMessage = time.Now()
			irc.lastMessageMutex.Unlock()
			event, err := parseToEvent(msg)
			event.Connection = irc
			if err == nil {
				/* XXX: len(args) == 0: args should be empty */
				irc.RunCallbacks(event)
			}
		}
	}
}

// Unescape tag values as defined in the IRCv3.2 message tags spec
// http://ircv3.net/specs/core/message-tags-3.2.html
func unescapeTagValue(value string) string {
	value = strings.Replace(value, "\\:", ";", -1)
	value = strings.Replace(value, "\\s", " ", -1)
	value = strings.Replace(value, "\\\\", "\\", -1)
	value = strings.Replace(value, "\\r", "\r", -1)
	value = strings.Replace(value, "\\n", "\n", -1)
	return value
}

//Parse raw irc messages
func parseToEvent(msg string) (*Event, error) {
	msg = strings.TrimSuffix(msg, "\n") //Remove \r\n
	msg = strings.TrimSuffix(msg, "\r")
	event := &Event{Raw: msg}
	if len(msg) < 5 {
		return nil, errors.New("Malformed msg from server")
	}

	if msg[0] == '@' {
		// IRCv3 Message Tags
		if i := strings.Index(msg, " "); i > -1 {
			event.Tags = make(map[string]string)
			tags := strings.Split(msg[1:i], ";")
			for _, data := range tags {
				parts := strings.SplitN(data, "=", 2)
				if len(parts) == 1 {
					event.Tags[parts[0]] = ""
				} else {
					event.Tags[parts[0]] = unescapeTagValue(parts[1])
				}
			}
			msg = msg[i+1 : len(msg)]
		} else {
			return nil, errors.New("Malformed msg from server")
		}
	}

	if msg[0] == ':' {
		if i := strings.Index(msg, " "); i > -1 {
			event.Source = msg[1:i]
			msg = msg[i+1 : len(msg)]

		} else {
			return nil, errors.New("Malformed msg from server")
		}

		if i, j := strings.Index(event.Source, "!"), strings.Index(event.Source, "@"); i > -1 && j > -1 && i < j {
			event.Nick = event.Source[0:i]
			event.User = event.Source[i+1 : j]
			event.Host = event.Source[j+1 : len(event.Source)]
		}
	}

	split := strings.SplitN(msg, " :", 2)
	args := strings.Split(split[0], " ")
	event.Code = strings.ToUpper(args[0])
	event.Arguments = args[1:]
	if len(split) > 1 {
		event.Arguments = append(event.Arguments, split[1])
	}
	return event, nil

}

// Loop to write to a connection. To be used as a goroutine.
func (irc *Connection) writeLoop() {
	defer irc.Done()
	errChan := irc.ErrorChan()
	for {
		select {
		case <-irc.end:
			return
		case b, ok := <-irc.pwrite:
			if !ok || b == "" || irc.socket == nil {
				return
			}

			if irc.Debug {
				irc.Log.Printf("--> %s\n", strings.TrimSpace(b))
			}

			// Set a write deadline based on the time out
			irc.socket.SetWriteDeadline(time.Now().Add(irc.Timeout))

			_, err := irc.socket.Write([]byte(b))

			// Past blocking write, bin timeout
			var zero time.Time
			irc.socket.SetWriteDeadline(zero)

			if err != nil {
				errChan <- err
				return
			}
		}
	}
}

// Pings the server if we have not received any messages for 5 minutes
// to keep the connection alive. To be used as a goroutine.
func (irc *Connection) pingLoop() {
	defer irc.Done()
	ticker := time.NewTicker(1 * time.Minute) // Tick every minute for monitoring
	ticker2 := time.NewTicker(irc.PingFreq)   // Tick at the ping frequency.
	for {
		select {
		case <-ticker.C:
			//Ping if we haven't received anything from the server within the keep alive period
			irc.lastMessageMutex.Lock()
			if time.Since(irc.lastMessage) >= irc.KeepAlive {
				irc.SendRawf("PING %d", time.Now().UnixNano())
			}
			irc.lastMessageMutex.Unlock()
		case <-ticker2.C:
			//Ping at the ping frequency
			irc.SendRawf("PING %d", time.Now().UnixNano())
			//Try to recapture nickname if it's not as configured.
			irc.Lock()
			if irc.nick != irc.nickcurrent {
				irc.nickcurrent = irc.nick
				irc.SendRawf("NICK %s", irc.nick)
			}
			irc.Unlock()
		case <-irc.end:
			ticker.Stop()
			ticker2.Stop()
			return
		}
	}
}

func (irc *Connection) isQuitting() bool {
	irc.Lock()
	defer irc.Unlock()
	return irc.quit
}

// Main loop to control the connection.
func (irc *Connection) Loop() {
	errChan := irc.ErrorChan()
	for !irc.isQuitting() {
		err := <-errChan
		close(irc.end)
		irc.Wait()
		for !irc.isQuitting() {
			irc.Log.Printf("Error, disconnected: %s\n", err)
			if err = irc.Reconnect(); err != nil {
				irc.Log.Printf("Error while reconnecting: %s\n", err)
				time.Sleep(60 * time.Second)
			} else {
				errChan = irc.ErrorChan()
				break
			}
		}
	}
}

// Quit the current connection and disconnect from the server
// RFC 1459 details: https://tools.ietf.org/html/rfc1459#section-4.1.6
func (irc *Connection) Quit() {
	quit := "QUIT"

	if irc.QuitMessage != "" {
		quit = fmt.Sprintf("QUIT :%s", irc.QuitMessage)
	}

	irc.SendRaw(quit)
	irc.Lock()
	irc.stopped = true
	irc.quit = true
	irc.Unlock()
}

// Use the connection to join a given channel.
// RFC 1459 details: https://tools.ietf.org/html/rfc1459#section-4.2.1
func (irc *Connection) Join(channel string) {
	irc.pwrite <- fmt.Sprintf("JOIN %s\r\n", channel)
}

// Leave a given channel.
// RFC 1459 details: https://tools.ietf.org/html/rfc1459#section-4.2.2
func (irc *Connection) Part(channel string) {
	irc.pwrite <- fmt.Sprintf("PART %s\r\n", channel)
}

// Send a notification to a nickname. This is similar to Privmsg but must not receive replies.
// RFC 1459 details: https://tools.ietf.org/html/rfc1459#section-4.4.2
func (irc *Connection) Notice(target, message string) {
	irc.pwrite <- fmt.Sprintf("NOTICE %s :%s\r\n", target, message)
}

// Send a formated notification to a nickname.
// RFC 1459 details: https://tools.ietf.org/html/rfc1459#section-4.4.2
func (irc *Connection) Noticef(target, format string, a ...interface{}) {
	irc.Notice(target, fmt.Sprintf(format, a...))
}

// Send (action) message to a target (channel or nickname).
// No clear RFC on this one...
func (irc *Connection) Action(target, message string) {
	irc.pwrite <- fmt.Sprintf("PRIVMSG %s :\001ACTION %s\001\r\n", target, message)
}

// Send formatted (action) message to a target (channel or nickname).
func (irc *Connection) Actionf(target, format string, a ...interface{}) {
	irc.Action(target, fmt.Sprintf(format, a...))
}

// Send (private) message to a target (channel or nickname).
// RFC 1459 details: https://tools.ietf.org/html/rfc1459#section-4.4.1
func (irc *Connection) Privmsg(target, message string) {
	irc.pwrite <- fmt.Sprintf("PRIVMSG %s :%s\r\n", target, message)
}

// Send formated string to specified target (channel or nickname).
func (irc *Connection) Privmsgf(target, format string, a ...interface{}) {
	irc.Privmsg(target, fmt.Sprintf(format, a...))
}

// Kick <user> from <channel> with <msg>. For no message, pass empty string ("")
func (irc *Connection) Kick(user, channel, msg string) {
	var cmd bytes.Buffer
	cmd.WriteString(fmt.Sprintf("KICK %s %s", channel, user))
	if msg != "" {
		cmd.WriteString(fmt.Sprintf(" :%s", msg))
	}
	cmd.WriteString("\r\n")
	irc.pwrite <- cmd.String()
}

// Kick all <users> from <channel> with <msg>. For no message, pass
// empty string ("")
func (irc *Connection) MultiKick(users []string, channel string, msg string) {
	var cmd bytes.Buffer
	cmd.WriteString(fmt.Sprintf("KICK %s %s", channel, strings.Join(users, ",")))
	if msg != "" {
		cmd.WriteString(fmt.Sprintf(" :%s", msg))
	}
	cmd.WriteString("\r\n")
	irc.pwrite <- cmd.String()
}

// Send raw string.
func (irc *Connection) SendRaw(message string) {
	irc.pwrite <- message + "\r\n"
}

// Send raw formated string.
func (irc *Connection) SendRawf(format string, a ...interface{}) {
	irc.SendRaw(fmt.Sprintf(format, a...))
}

// Set (new) nickname.
// RFC 1459 details: https://tools.ietf.org/html/rfc1459#section-4.1.2
func (irc *Connection) Nick(n string) {
	irc.nick = n
	irc.SendRawf("NICK %s", n)
}

// Determine nick currently used with the connection.
func (irc *Connection) GetNick() string {
	return irc.nickcurrent
}

// Query information about a particular nickname.
// RFC 1459: https://tools.ietf.org/html/rfc1459#section-4.5.2
func (irc *Connection) Whois(nick string) {
	irc.SendRawf("WHOIS %s", nick)
}

// Query information about a given nickname in the server.
// RFC 1459 details: https://tools.ietf.org/html/rfc1459#section-4.5.1
func (irc *Connection) Who(nick string) {
	irc.SendRawf("WHO %s", nick)
}

// Set different modes for a target (channel or nickname).
// RFC 1459 details: https://tools.ietf.org/html/rfc1459#section-4.2.3
func (irc *Connection) Mode(target string, modestring ...string) {
	if len(modestring) > 0 {
		mode := strings.Join(modestring, " ")
		irc.SendRawf("MODE %s %s", target, mode)
		return
	}
	irc.SendRawf("MODE %s", target)
}

func (irc *Connection) ErrorChan() chan error {
	return irc.Error
}

// Returns true if the connection is connected to an IRC server.
func (irc *Connection) Connected() bool {
	return !irc.stopped
}

// A disconnect sends all buffered messages (if possible),
// stops all goroutines and then closes the socket.
func (irc *Connection) Disconnect() {
	irc.Lock()
	defer irc.Unlock()

	if irc.end != nil {
		close(irc.end)
	}

	irc.end = nil

	if irc.pwrite != nil {
		close(irc.pwrite)
	}

	irc.Wait()
	if irc.socket != nil {
		irc.socket.Close()
	}
	irc.ErrorChan() <- ErrDisconnected
}

// Reconnect to a server using the current connection.
func (irc *Connection) Reconnect() error {
	irc.end = make(chan struct{})
	return irc.Connect(irc.Server)
}

// Connect to a given server using the current connection configuration.
// This function also takes care of identification if a password is provided.
// RFC 1459 details: https://tools.ietf.org/html/rfc1459#section-4.1
func (irc *Connection) Connect(server string) error {
	irc.Server = server
	// mark Server as stopped since there can be an error during connect
	irc.stopped = true

	// make sure everything is ready for connection
	if len(irc.Server) == 0 {
		return errors.New("empty 'server'")
	}
	if strings.Count(irc.Server, ":") != 1 {
		return errors.New("wrong number of ':' in address")
	}
	if strings.Index(irc.Server, ":") == 0 {
		return errors.New("hostname is missing")
	}
	if strings.Index(irc.Server, ":") == len(irc.Server)-1 {
		return errors.New("port missing")
	}
	// check for valid range
	ports := strings.Split(irc.Server, ":")[1]
	port, err := strconv.Atoi(ports)
	if err != nil {
		return errors.New("extracting port failed")
	}
	if !((port >= 0) && (port <= 65535)) {
		return errors.New("port number outside valid range")
	}
	if irc.Log == nil {
		return errors.New("'Log' points to nil")
	}
	if len(irc.nick) == 0 {
		return errors.New("empty 'nick'")
	}
	if len(irc.user) == 0 {
		return errors.New("empty 'user'")
	}

	if irc.UseTLS {
		dialer := &net.Dialer{Timeout: irc.Timeout}
		irc.socket, err = tls.DialWithDialer(dialer, "tcp", irc.Server, irc.TLSConfig)
	} else {
		irc.socket, err = net.DialTimeout("tcp", irc.Server, irc.Timeout)
	}
	if err != nil {
		return err
	}

	irc.stopped = false
	irc.Log.Printf("Connected to %s (%s)\n", irc.Server, irc.socket.RemoteAddr())

	irc.pwrite = make(chan string, 10)
	irc.Error = make(chan error, 2)
	irc.Add(3)
	go irc.readLoop()
	go irc.writeLoop()
	go irc.pingLoop()

	if len(irc.WebIRC) > 0 {
		irc.pwrite <- fmt.Sprintf("WEBIRC %s\r\n", irc.WebIRC)
	}

	if len(irc.Password) > 0 {
		irc.pwrite <- fmt.Sprintf("PASS %s\r\n", irc.Password)
	}

	err = irc.negotiateCaps()
	if err != nil {
		return err
	}

	realname := irc.user
	if irc.RealName != "" {
		realname = irc.RealName
	}

	irc.pwrite <- fmt.Sprintf("NICK %s\r\n", irc.nick)
	irc.pwrite <- fmt.Sprintf("USER %s 0.0.0.0 0.0.0.0 :%s\r\n", irc.user, realname)
	return nil
}

// Negotiate IRCv3 capabilities
func (irc *Connection) negotiateCaps() error {
	saslResChan := make(chan *SASLResult)
	if irc.UseSASL {
		irc.RequestCaps = append(irc.RequestCaps, "sasl")
		irc.setupSASLCallbacks(saslResChan)
	}

	if len(irc.RequestCaps) == 0 {
		return nil
	}

	cap_chan := make(chan bool, len(irc.RequestCaps))
	irc.AddCallback("CAP", func(e *Event) {
		if len(e.Arguments) != 3 {
			return
		}
		command := e.Arguments[1]

		if command == "LS" {
			missing_caps := len(irc.RequestCaps)
			for _, cap_name := range strings.Split(e.Arguments[2], " ") {
				for _, req_cap := range irc.RequestCaps {
					if cap_name == req_cap {
						irc.pwrite <- fmt.Sprintf("CAP REQ :%s\r\n", cap_name)
						missing_caps--
					}
				}
			}

			for i := 0; i < missing_caps; i++ {
				cap_chan <- true
			}
		} else if command == "ACK" || command == "NAK" {
			for _, cap_name := range strings.Split(strings.TrimSpace(e.Arguments[2]), " ") {
				if cap_name == "" {
					continue
				}

				if command == "ACK" {
					irc.AcknowledgedCaps = append(irc.AcknowledgedCaps, cap_name)
				}
				cap_chan <- true
			}
		}
	})

	irc.pwrite <- "CAP LS\r\n"

	if irc.UseSASL {
		select {
		case res := <-saslResChan:
			if res.Failed {
				close(saslResChan)
				return res.Err
			}
		case <-time.After(time.Second * 15):
			close(saslResChan)
			return errors.New("SASL setup timed out. This shouldn't happen.")
		}
	}

	// Wait for all capabilities to be ACKed or NAKed before ending negotiation
	for i := 0; i < len(irc.RequestCaps); i++ {
		<-cap_chan
	}
	irc.pwrite <- fmt.Sprintf("CAP END\r\n")

	realname := irc.user
	if irc.RealName != "" {
		realname = irc.RealName
	}

	irc.pwrite <- fmt.Sprintf("NICK %s\r\n", irc.nick)
	irc.pwrite <- fmt.Sprintf("USER %s 0.0.0.0 0.0.0.0 :%s\r\n", irc.user, realname)
	return nil
}

// Create a connection with the (publicly visible) nickname and username.
// The nickname is later used to address the user. Returns nil if nick
// or user are empty.
func IRC(nick, user string) *Connection {
	// catch invalid values
	if len(nick) == 0 {
		return nil
	}
	if len(user) == 0 {
		return nil
	}

	irc := &Connection{
		nick:        nick,
		nickcurrent: nick,
		user:        user,
		Log:         log.New(os.Stdout, "", log.LstdFlags),
		end:         make(chan struct{}),
		Version:     VERSION,
		KeepAlive:   4 * time.Minute,
		Timeout:     1 * time.Minute,
		PingFreq:    15 * time.Minute,
		SASLMech:    "PLAIN",
		QuitMessage: "",
	}
	irc.setupCallbacks()
	return irc
}
