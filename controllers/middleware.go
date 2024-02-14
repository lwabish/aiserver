package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	authPrivateKeyPath = "certs/private.pem"
	authOpenIDName     = "openid"
	authTicketName     = "ticket"
)

var (
	MidCtl = &MiddlewareController{}
)

type MiddlewareController struct {
	L            *logrus.Logger
	TicketExpire bool
}

type MiddlewareControllerCfg struct {
	L            *logrus.Logger
	TicketExpire bool
}

func (m *MiddlewareController) Setup(cfg *MiddlewareControllerCfg) {
	m.L = cfg.L
	m.TicketExpire = cfg.TicketExpire
}

func (m *MiddlewareController) Authenticate(c *gin.Context) {
	openid := c.PostForm(authOpenIDName)
	ticket := c.PostForm(authTicketName)

	if openid == "" || ticket == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		c.Abort()
		return
	}

	// Validate ticket
	decryptedTicket, err := decryptTicket(ticket)
	if err != nil || !strings.HasPrefix(decryptedTicket, openid) {
		m.L.Errorf("decryptTicket error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication"})
		c.Abort()
		return
	}

	ticketTimestamp, err := strconv.ParseInt(decryptedTicket[len(openid):], 10, 64)
	if err != nil || (m.TicketExpire && time.Now().Unix()-ticketTimestamp > 60) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication expired"})
		c.Abort()
		return
	}

	c.Next()
}

func decryptTicket(ticket string) (string, error) {
	privateKey, err := os.ReadFile(authPrivateKeyPath)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	decodedTicket, err := base64.StdEncoding.DecodeString(ticket)
	if err != nil {
		return "", err
	}

	rng := rand.Reader
	plaintext, err := rsa.DecryptPKCS1v15(rng, key.(*rsa.PrivateKey), decodedTicket)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
