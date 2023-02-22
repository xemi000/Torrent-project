package torr

import (
	"encoding/binary"
)

/*

This code defines a struct called "Announce" which contains various fields
for storing information related to a BitTorrent announce request and response,
as well as methods for serializing and deserializing the struct to/from a byte slice.
The struct has fields for storing information such as the connection ID, info hash, peer ID,
downloaded and uploaded amounts, event, IP, key, and number of peers wanted. The newAnnounce function is used
to create a new announce struct with the provided info hash, peer ID, and left value, along with a randomly generated
connection ID. The serializeAnnounce method converts the struct into a byte slice, while the readAnnounce method takes
in a byte slice and converts it back into an Announce struct. This code is likely used
 in a BitTorrent client to send announce requests to a tracker and to parse the responses received from the tracker.
*/

const announceLen = 98

type Announce struct {
	Action        uint32 // request & response
	TransactionID []byte // request & response

	ConnectionID []byte   // request
	InfoHash     [20]byte // request
	PeerID       [20]byte // request
	Downloaded   uint64   // request
	Left         uint64   // request
	Uploaded     uint64   // request
	Event        uint32   // request
	IP           uint32   // request
	Key          []byte   // request
	NumWant      int      // request
	Port         uint16   // request

	Interval uint32 // response
	Leechers uint32 // response
	Seeders  uint32 // response
	Peers    []byte // response
}

func newAnnounce(infoHash, peerID [20]byte, left int, connectionID []byte) *Announce {
	return &Announce{
		ConnectionID:  connectionID,
		Action:        1,
		TransactionID: generateRandomID(4),
		InfoHash:      infoHash,
		PeerID:        peerID,
		Downloaded:    0,
		Left:          uint64(left),
		Uploaded:      0,
		Event:         0,
		IP:            0,
		Key:           generateRandomID(4),
		NumWant:       -1,
		Port:          0,
	}
}

func (a *Announce) serializeAnnounce() []byte {
	buf := make([]byte, announceLen)
	copy(buf[:8], a.ConnectionID[:])
	binary.BigEndian.PutUint32(buf[8:12], a.Action)
	copy(buf[12:16], a.TransactionID[:])
	copy(buf[16:36], a.InfoHash[:])
	copy(buf[36:56], a.PeerID[:])
	binary.BigEndian.PutUint64(buf[56:64], a.Downloaded)
	binary.BigEndian.PutUint64(buf[64:72], a.Left)
	binary.BigEndian.PutUint64(buf[72:80], a.Uploaded)
	binary.BigEndian.PutUint32(buf[80:84], a.Event)
	binary.BigEndian.PutUint32(buf[84:88], a.IP)
	copy(buf[88:92], a.Key[:])
	binary.BigEndian.PutUint32(buf[92:96], uint32(a.NumWant))
	binary.BigEndian.PutUint16(buf[96:98], a.Port)
	return buf
}

func readAnnounce(buf []byte) *Announce {
	announceReq := make([]byte, 20)
	copy(announceReq, buf[:20])

	actionBuf := make([]byte, 4)
	transactionIDBuf := make([]byte, 4)
	intervalBuf := make([]byte, 4)
	leechersBuf := make([]byte, 4)
	seedersBuf := make([]byte, 4)

	copy(actionBuf, announceReq[0:4])
	copy(transactionIDBuf[:], announceReq[4:8])
	copy(intervalBuf, announceReq[8:12])
	copy(leechersBuf, announceReq[12:16])
	copy(seedersBuf, announceReq[16:20])

	peersBuf := make([]byte, len(buf)-20)
	copy(peersBuf, buf[20:])

	ar := Announce{
		Action:        binary.BigEndian.Uint32(actionBuf),
		TransactionID: transactionIDBuf[:],
		Interval:      binary.BigEndian.Uint32(intervalBuf),
		Leechers:      binary.BigEndian.Uint32(leechersBuf),
		Seeders:       binary.BigEndian.Uint32(seedersBuf),
		Peers:         peersBuf,
	}
	return &ar
}
