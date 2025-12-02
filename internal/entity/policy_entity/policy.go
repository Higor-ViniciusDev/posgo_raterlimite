package policy_entity

import (
	"os"
	"strconv"
	"time"
)

type Policy struct {
	RequestPerSecond int64
	WindowPerSecond  int64
	TTL              int64
	Fonte            string
	start_at         int64 //Pensei que por ser um valor um importante deixar ele privado para que não seja alterado depois da criação
}

const (
	FONTE_IP     = "IP"
	FONTE_TOLKEN = "TOLKEN"
)

func NewPolicyTolken() *Policy {
	request := os.Getenv("REQUEST_PER_SECOND_TOLKEN")
	widowsRequest := os.Getenv("REQUEST_PER_WINDOW")
	tll := os.Getenv("TOLKEN_EXPIRATION")

	start := time.Now().Unix()

	requestNumber, _ := strconv.Atoi(request)
	widowsRequestNumber, _ := strconv.Atoi(widowsRequest)
	tllNumber, _ := strconv.Atoi(tll)

	return &Policy{
		Fonte:            "TOLKEN",
		WindowPerSecond:  int64(widowsRequestNumber),
		RequestPerSecond: int64(requestNumber),
		TTL:              int64(tllNumber),
		start_at:         start,
	}
}

func NewPolicyIP() *Policy {
	request := os.Getenv("REQUEST_PER_SECOND_IP")
	widowsRequest := os.Getenv("REQUEST_PER_WINDOW")
	tll := os.Getenv("TOLKEN_EXPIRATION")

	start := time.Now().Unix()

	requestNumber, _ := strconv.Atoi(request)
	widowsRequestNumber, _ := strconv.Atoi(widowsRequest)
	tllNumber, _ := strconv.Atoi(tll)

	return &Policy{
		Fonte:            "IP",
		WindowPerSecond:  int64(widowsRequestNumber),
		RequestPerSecond: int64(requestNumber),
		TTL:              int64(tllNumber),
		start_at:         start,
	}
}

func (p *Policy) GetTimeStartad() int64 {
	return p.start_at
}

// SetStartAt define o tempo inicial (usado ao ler do armazenamento)
func (p *Policy) SetStartAt(t int64) {
	p.start_at = t
}
