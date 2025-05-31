package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
	"github.com/rudychandra/lagi/config"
	"github.com/rudychandra/lagi/model"
)

func SendNotification(c *gin.Context) {
	ctx := context.Background()

	client, err := config.FirebaseApp.Messaging(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Firebase init failed",
			"result":  "error",
			"data":    err.Error(),
		})
		return
	}

	var notification model.FCMMessage

	if err := c.ShouldBind(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Notification sent failed",
			"result":  "error",
			"data":    "Format data salah",
		})
		return
	}

	fmt.Println("Token Devices:", notification.Message.Tokens)
	tokens := notification.Message.Tokens // sekarang tokens adalah []string

	var count = 0
	var failedTokens []string

	for _, token := range tokens {
		msg := &messaging.Message{
			Token: token,
			Notification: &messaging.Notification{
				Title: notification.Message.Notification.Title,
				Body:  notification.Message.Notification.Body,
			},
			Data: map[string]string{
				"screen":      notification.Message.Data.Screen,
				"jadwal_id":   notification.Message.Data.JadwalID,
				"waktu_mulai": notification.Message.Data.WaktuMulai,
			},
		}

		_, err := client.Send(ctx, msg)
		if err != nil {
			failedTokens = append(failedTokens, token)
			log.Printf("Failed to send to %s: %v\n", token, err)
		} else {
			count++
		}
	}

	if count == len(tokens) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Notification sent successfully",
			"result":  "success",
			"data":    notification,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Some notifications failed",
			"result":  "partial_success",
			"data": gin.H{
				"sent":   count,
				"failed": failedTokens,
			},
		})
	}
}

