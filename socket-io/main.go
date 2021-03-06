package socketIo2

import (
	"backend-qrcode/middleware"
	"log"
	"strconv"
	"sync"

	"github.com/dgrijalva/jwt-go"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

var once sync.Once

// SocketIO ...
type SocketIO struct {
	Server *gosocketio.Server
}

// VerifyJWTReturn ...
type VerifyJWTReturn struct {
	UserID string
	RoleID string
}

// VerifyJWT ...
func VerifyJWT(tokenString string) (*VerifyJWTReturn, *error) {
	claims, err := middleware.VerifyToken(tokenString)
	if err != nil {
		return nil, &err
	}
	UserID := strconv.FormatFloat(claims.(jwt.MapClaims)["user_id"].(float64), 'g', 1, 64)
	RoleID := strconv.FormatFloat(claims.(jwt.MapClaims)["role_id"].(float64), 'g', 1, 64)

	return &VerifyJWTReturn{UserID, RoleID}, nil
}

// Channel
type Channel struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

// Run ...
func (s *SocketIO) Run() {

	s.Server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client consnected")
		c.Emit("onReconnect", "asd")

	})

	s.Server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("New Disconect")

	})

	s.Server.On(gosocketio.OnError, func(c *gosocketio.Channel) {
		log.Println("New Error")

	})

	s.Server.On("/join", func(c *gosocketio.Channel, channel Channel) string {

		c.Join(channel.Name)
		return "joined to " + "room"
	})
}

var instance *SocketIO

// GetSocketIO ....
func GetSocketIO() *SocketIO {

	once.Do(func() {

		server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

		instance = &SocketIO{Server: server}
	})

	return instance
}
