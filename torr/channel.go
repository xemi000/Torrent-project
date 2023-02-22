package torr

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type Channel struct {
	Conn     net.Conn
	Choked   bool
	Bitfield Bitfield
	peer     Peer
	infoHash [20]byte
	peerID   [20]byte
}

func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) error {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	request := newHandshake(infoHash, peerID)
	_, err := conn.Write(request.serializeHandshake())
	if err != nil {
		return err
	}

	result, err := readHandshake(conn)
	if err != nil {
		return err
	}

	if !bytes.Equal(result.InfoHash[:], infoHash[:]) {
		err := fmt.Errorf("expected infohash %x but got %x", infoHash, request.InfoHash)
		return err
	}

	return nil
}

func receiveBitfield(conn net.Conn) (Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	msg, err := readMessage(conn)
	if err != nil {
		return nil, err
	}

	if msg == nil {
		err := fmt.Errorf("expected bitfield but got keep alive")
		return nil, err
	}

	if msg.ID != bitfield {
		err := fmt.Errorf("expected bitfield but got message ID %d", msg.ID)
		return nil, err
	}

	return msg.Payload, nil
}

func (t *Torrent) newChannel(peer Peer, peerID, infoHash [20]byte) (*Channel, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 5*time.Second)
	if err != nil {
		return nil, err
	}

	err = completeHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := receiveBitfield(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	t.activePeers++

	return &Channel{
		Conn:     conn,
		Choked:   true,
		Bitfield: bf,
		peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}

func (ch *Channel) read() (*Message, error) {
	msg, err := readMessage(ch.Conn)
	return msg, err
}

func (ch *Channel) sendRequest(index, begin, length int) error {
	req := createRequestMessage(index, begin, length)
	_, err := ch.Conn.Write(req.serializeMessage())
	return err
}

func (ch *Channel) sendInterested() error {
	msg := Message{ID: interested}
	_, err := ch.Conn.Write(msg.serializeMessage())
	return err
}

func (ch *Channel) sendNotInterested() error {
	msg := Message{ID: notInterested}
	_, err := ch.Conn.Write(msg.serializeMessage())
	return err
}

func (ch *Channel) sendUnchoke() error {
	msg := Message{ID: unchoke}
	_, err := ch.Conn.Write(msg.serializeMessage())
	return err
}

func (ch *Channel) sendHave(index int) error {
	msg := createHaveMessage(index)
	_, err := ch.Conn.Write(msg.serializeMessage())
	return err
}
