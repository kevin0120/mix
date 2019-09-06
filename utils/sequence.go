package utils

import "sync"

const (
	DEFAULT_SEQUENCE = 1
	MAX_SEQUENCE     = 99999
)

func CreateSequence() *Sequence {
	return &Sequence{
		sequence: DEFAULT_SEQUENCE,
		mtxSeq:   sync.Mutex{},
	}
}

type Sequence struct {
	sequence uint32
	mtxSeq   sync.Mutex
}

func (s *Sequence) GetSequence() uint32 {

	s.mtxSeq.Lock()
	defer s.mtxSeq.Unlock()

	seq := s.sequence
	if s.sequence == MAX_SEQUENCE {
		s.sequence = MAX_SEQUENCE
	} else {
		s.sequence++
	}

	return seq
}
