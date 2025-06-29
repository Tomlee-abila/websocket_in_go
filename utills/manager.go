package utills

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin: checkOrigin,
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct{
	clients ClientList
	sync.RWMutex

	handlers map[string]EventHandler
}

func NewManager() *Manager{
	m := &Manager{
		clients: make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}



func (m *Manager) setupEventHandlers(){
	m.handlers[EventSendMessage] = SendMessage
}

func SendMessage (event Event, c *Client) error{
	fmt.Println(event)
	return nil
}

func (m *Manager) routEvent(event Event, c *Client) error{
	// Checks if the event is part of the handlers
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil{
			return err
		}
		return nil
	}else{
		return errors.New("there is no such event type")
	}
}

func (m *Manager) ServeWs( w http.ResponseWriter, r *http.Request){
	log.Println("new connection")

	//upgrade regular http connection into websocket
	conn, err := websocketUpgrader.Upgrade(w, r, nil)

	if err != nil{
		log.Println(err)
		return
	}

	client := NewClient(conn, m);
	m.addClient(client)

	//Start client processes
	go client.readMessages()
	go client.writeMessages();
}

func (m *Manager) addClient(client *Client){
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true;
}

func (m *Manager) removeClint(client *Client){
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok{
		client.connection.Close()
		delete(m.clients, client)
	}
}

func checkOrigin (r *http.Request)bool{
	origin := r.Header.Get("Origin")

	switch origin{
	case "http://localhost:8080":
		return true
	default:
		return false
	}
}