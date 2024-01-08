package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"vitopass.com/vitogosocket/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type EvenetTypeHandler func(conn *websocket.Conn, data map[string]interface{})

var eventHandlers = map[string]EvenetTypeHandler{
	"scoreUpdate": handleScoreUpdate,
}

func main() {
	server := gin.Default()

	server.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			// Handle error
			return
		}
		fmt.Println("Client connected")

		cookie, err := c.Cookie("token") // Replace "token" with your cookie name
		if err != nil {
			// Handle error or no cookie found
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You must be logged in"})
			return
			fmt.Println("please login")
		}

		extractedTokenClaims, err1 := utils.VerifyToken(cookie)
		// Here you can perform token verification
		if err1 != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unable to authenticate"})
			return
		}

		fmt.Println(extractedTokenClaims)
		defer conn.Close()

		for {
			// Read message from the client
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message: ", err)
				break
			}
			// Print the message to the console.
			fmt.Println("Received:", msg)

			// Here, you can parse and handle the incoming message
			// For example, assuming the message is in JSON format
			var data map[string]interface{}
			err = json.Unmarshal(msg, &data)
			if err != nil {
				log.Println("Error parsing JSON:", err)
				continue
			}

			eventType, exists := data["event"].(string)

			if !exists {
				log.Println("Invalid or missing event type")
				continue
			}

			handler, found := eventHandlers[eventType]
			if !found {
				log.Println("Handler not found for event type:", eventType)
				continue
			}

			// Call the respective handler for this event type
			handler(conn, data)
		}
	})

	server.Run(":2026")
}

func handleScoreUpdate(conn *websocket.Conn, data map[string]interface{}) {

	fmt.Println(data)
	response := []byte("Received score update")
	err := conn.WriteMessage(websocket.TextMessage, response)
	if err != nil {
		log.Println("Error writing response:", err)
	}
	// Emit an event back to the client
	eventData := map[string]interface{}{
		"message": "Updated the total leaderboard",
		// Other data fields for the 'updatedleaderboard' event, if needed
	}

	event := map[string]interface{}{
		"event": "updatedleaderboard",
		"data":  eventData,
	}
	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Println("Error emitting event:", err)
	}
}
