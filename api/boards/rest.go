package boards

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
)

func Bootstrap(rg *gin.RouterGroup) {
	a := api{s: &service{}}
	rg.GET("/orgs/:orgSlug/boards", a.getInfos)
	rg.GET("/orgs/:orgSlug/boards/:boardSlug", a.getBoard)
	rg.GET("/orgs/:orgSlug/boards/:boardSlug/ws", handleWebSocket)
	rg.POST("/orgs/:orgSlug/boards/:boardSlug", a.handleUpdateBoard)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = sync.Map{}

type api struct {
	s Boarder
}

func (a *api) getInfos(c *gin.Context) {
	orgSlug := c.Param("orgSlug")
	infos, err := a.s.getInfos(orgSlug)
	if err != nil {
		respondWithJSON(c, http.StatusBadRequest, nil, err)
		return
	}
	respondWithJSON(c, http.StatusOK, infos, nil)
}

func (a *api) getBoard(c *gin.Context) {
	orgSlug := c.Param("orgSlug")
	boardSlug := c.Param("boardSlug")
	board, err := a.s.getBoard(orgSlug, boardSlug)
	if err != nil {
		respondWithJSON(c, http.StatusBadRequest, nil, err)
		return
	}
	respondWithJSON(c, http.StatusOK, board, nil)
}

func (a *api) handleUpdateBoard(c *gin.Context) {
	orgSlug := c.Param("orgSlug")
	boardSlug := c.Param("boardSlug")
	var body json.RawMessage
	if err := c.ShouldBindJSON(&body); err != nil {
		respondWithJSON(c, http.StatusBadRequest, nil, err)
		return
	}
	board, err := a.s.storeBoard(orgSlug, boardSlug, body)
	if err != nil {
		respondWithJSON(c, http.StatusBadRequest, nil, err)
		return
	}
	respondWithJSON(c, http.StatusOK, board, nil)
}

func respondWithJSON(c *gin.Context, status int, o interface{}, err error) {
	c.Writer.WriteHeader(status)
	var errs []string
	if err != nil {
		errs = []string{err.Error()}
	}
	body := map[string]interface{}{"errors": errs, "data": o}
	b, _ := json.Marshal(body)
	_,_ = c.Writer.Write(b)
}

func handleWebSocket(c *gin.Context) {
	clientID := c.Query("clientId")
	orgSlug := c.Param("orgSlug")
	boardSlug := c.Param("boardSlug")
	boardID := fmt.Sprintf("%s.%s", orgSlug, boardSlug)
	clientID = fmt.Sprintf("%s.%s", boardID, clientID)

	log.Printf("Initial request from client '%s', upgrading connection...", clientID)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to websockets. Error: %+v", err)
		return
	}

	clients.Store(clientID, conn)
	log.Printf("Connection successfully upgraded. Starting to read messages.")
	for {
		mType, msg, err := conn.ReadMessage()
		if err != nil {
			clients.Delete(clientID)
			log.Printf("Failed to read message from connection. Removing client '%s'. Error: %+v", clientID, err)
			break
		}
		log.Printf("Message received from client '%s': %s", clientID, msg)
		go func() {
			notifyAllClients(boardID, clientID, mType, msg)
		}()
	}
}

func notifyAllClients(boardID string, origin string, messageType int, msg []byte) {
	clients.Range(func(k, v interface{}) bool {
		key := fmt.Sprintf("%s", k)
		// do not notify the origin or any other board
		if key == origin || !strings.HasPrefix(key, boardID){
			return true
		}
		conn := v.(*websocket.Conn)
		_ = conn.WriteMessage(messageType, msg)
		return true
	})
}

type event struct {
	Source string `json:"source"`
	Content interface{} `json:"interface"`
	Timestamp int64 `json:"timestamp"`
}