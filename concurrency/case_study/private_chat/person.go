package privatechat

import (
	"log"
	"time"

	"github.com/ductran999/shared-pkg/scrypto/caesar"
)

type person struct {
	name string // Name of the person

	sendChan     chan string // Channel used by the person to send messages
	receivedChan chan string // Channel used by the person to receive messages
	secretChan   chan int    // Channel through which the person receives the nonce (secret value) from a generator

	turn     int      // Indicates the person's turn in a communication or protocol sequence
	messages []string // Stores messages sent by the person
}

func NewPerson(
	name string,
	messages []string,
	sendChan chan string,
	receivedChan chan string,
	turn int,
) *person {
	return &person{
		name:         name,
		messages:     messages,
		sendChan:     sendChan,
		receivedChan: receivedChan,
		secretChan:   make(chan int),
		turn:         turn,
	}
}

func (p *person) ReceiveSecret(secret int) {
	p.secretChan <- secret
}

func (p *person) Chat() {
	sent := p.turn
	msgIndex := 0

	for s := range p.secretChan {
		// My turn to SEND
		if sent%2 == 0 {
			if msgIndex >= len(p.messages) {
				log.Printf("%s: all messages sent - closing chat", p.name)
				close(p.sendChan) // Notify other party
				return
			}

			cipher := caesar.CaesarEncrypt(p.messages[msgIndex], s)
			p.sendChan <- cipher
			msgIndex++
			time.Sleep(time.Second)
		} else {
			// My turn to RECEIVE
			cipher, ok := <-p.receivedChan
			if !ok {
				log.Printf("%s: receive channel closed - ending chat", p.name)
				return
			}

			log.Printf("%s got cipher: %q", p.name, cipher)

			// Decrypt the received message
			plain := caesar.CaesarDecrypt(cipher, s)
			log.Printf("%s received message: %q", p.name, plain)
			log.Println("---------------------------------------------------------")
		}

		sent++
	}
}
