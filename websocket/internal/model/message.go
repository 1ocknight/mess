package model

import wsdto "github.com/1ocknight/mess/shared/dto/ws"

type Message struct {
	SubjectID string
	WSMessage *wsdto.WSMessage
}
