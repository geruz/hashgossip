package storage

import (
	"bytes"
	"log"
	"sync"

	"github.com/DemonVex/hashgossip/models"
)

type PeerStorage interface {
	List() []models.Peer
	Add(models.Peer)
	Merge([]models.Peer)
	IsIn(models.Peer) bool
	IsEmpty() bool
}

type peerStorage struct {
	list  []models.Peer
	//RW mutex?
	mutex *sync.Mutex
}

func NewPeerStorage() PeerStorage {
	return &peerStorage{mutex: &sync.Mutex{}}
}

func (p *peerStorage) List() []models.Peer {
	//lock? может быть ситуация когда gossiper получает список и в это время список меняется
	// не стоит ли возвращать копию?
	// Падения скорее всего не будет но build --race  может детектить
	return p.list
}

func (p *peerStorage) Add(peer models.Peer) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.unsafeAdd(peer)
}

func (p *peerStorage) unsafeAdd(peer models.Peer) {
	if !p.IsIn(peer) {
		p.list = append(p.list, peer)
		log.Printf("New peer %v", peer)
	}
}

// lock? публичный но не безопасный (сделать приватным?)
func (p *peerStorage) IsIn(peer models.Peer) bool {
	for _, v := range p.list {
		if v.Port == peer.Port && bytes.Equal(v.IP, peer.IP) {
			return true
		}
	}
	return false
}

func (p *peerStorage) Merge(list []models.Peer) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, v := range list {
		p.unsafeAdd(v)
	}
}

func (p *peerStorage) IsEmpty() bool {
	return len(p.list) == 0
}
