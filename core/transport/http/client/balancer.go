package client

import (
	"errors"
	"strconv"
)

type Balancer interface {
	Select() (*Address, error)
	Update([]*Address)
}

type Peer struct {
	weight          int
	effectiveWeight int
	currentWeight   int
	Addr            *Address
}

type Wrr struct {
	peers []*Peer
	Next  int
}

func (w *Wrr) Select() (*Address, error) {
	total := 0
	var best *Peer
	for _, peer := range w.peers {
		peer.currentWeight += peer.effectiveWeight
		total += peer.effectiveWeight

		if peer.effectiveWeight < peer.weight {
			peer.effectiveWeight += 1
		}

		if best == nil || best.currentWeight < peer.currentWeight {
			best = peer
		}
	}

	if best != nil {
		best.currentWeight -= total
	}

	if best == nil {
		return nil, errors.New("service unavailable: produced zero addresses")
	}

	return best.Addr, nil
}

func (w *Wrr) Update(addresses []*Address) {
	peers := make([]*Peer, 0, len(addresses))
	for _, v := range addresses {
		weight, _ := strconv.Atoi(v.Metadata["weight"])
		if weight == 0 {
			weight = 1
		}
		peer := &Peer{
			weight:          weight,
			effectiveWeight: weight,
			currentWeight:   0,
			Addr:            v,
		}
		peers = append(peers, peer)
	}
	w.peers = peers
}

type RR struct {
	Nodes []*Address
	Next  int
}

func (r *RR) Select() (*Address, error) {
	if len(r.Nodes) == 0 {
		return nil, errors.New("service unavailable: produced zero addresses")
	}
	addr := r.Nodes[r.Next]
	r.Next = (r.Next + 1) % len(r.Nodes)

	return addr, nil
}

func (r *RR) Update(addresses []*Address) {
	r.Nodes = addresses
}
