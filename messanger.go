package messanger

import (
	"github.com/gorilla/websocket"
)

// Subcriber :
type Subcriber struct {
	id  string
	con *websocket.Conn
}

// Subcription :
type Subcription struct {
	Topic     string
	Subcriber *Subcriber
}

// Messanger :
type Messanger struct {
	Subcribers   []Subcriber
	Subcriptions []Subcription
}

// Sub : empty topic if yo just subribe broadcast message
func (m *Messanger) Sub(subcriber *Subcriber, topik string) *Messanger {
	// cek apakah topik dan user sudah ada?
	isExist := false
	for _, sub := range m.Subcriptions {
		if sub.Subcriber.id == subcriber.id {
			if topik == sub.Topic {
				return m
			}
			isExist = true
		}
	}
	if !isExist {
		m.Subcribers = append(m.Subcribers, *subcriber)
	}
	if topik != "" {
		m.Subcriptions = append(m.Subcriptions, Subcription{Subcriber: subcriber, Topic: topik})
	}
	return m
}

// Unsub : empty string topik is unsub for all topic
func (m *Messanger) Unsub(subcriber *Subcriber, topik string) *Messanger {
	for i, c := range m.Subcriptions {
		if subcriber.id == c.Subcriber.id && (topik == c.Topic || topik == "") {
			m.Subcriptions = append(m.Subcriptions[:i], m.Subcriptions[i+1:]...)
		}
	}
	return m
}

// RemoveSubcriber : always remove before you close
func (m *Messanger) RemoveSubcriber(subcriber *Subcriber) *Messanger {
	m.Unsub(subcriber, "")
	for i := range m.Subcribers {
		m.Subcribers = append(m.Subcribers[:i], m.Subcribers[i+1:]...)
	}
	return m
}

// SendMessage : to broadcast, pass empty string to topic
func (m *Messanger) SendMessage(message, topic string) {
	if topic != "" {
		for _, subriber := range m.Subcriptions {
			if topic == subriber.Topic {
				subriber.Subcriber.con.WriteJSON(`{"topic": ` + topic + `, "message": "` + message + `"}`)
			}
		}
	}
	for _, subriber := range m.Subcribers {
		subriber.con.WriteJSON(`{"topic": ` + topic + `, "message": "` + message + `"}`)
	}
}
