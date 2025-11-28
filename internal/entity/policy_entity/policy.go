package policy_entity

import (
	"os"
	"strconv"
)

type Policy struct {
	RequestPerSecond int64
	WindowPerSecond  int64
	TTL              int64
	Fonte            string
}

const (
	FONTE_IP     = "IP"
	FONTE_TOLKEN = "TOLKEN"
)

func NewPolicyTolken() *Policy {
	request := os.Getenv("REQUEST_PER_SECOND_TOLKEN")
	widowsRequest := os.Getenv("REQUEST_PER_WINDOW")
	tll := os.Getenv("TOLKEN_EXPIRATION")

	requestNumber, _ := strconv.Atoi(request)
	widowsRequestNumber, _ := strconv.Atoi(widowsRequest)
	tllNumber, _ := strconv.Atoi(tll)

	return &Policy{
		Fonte:            "TOLKEN",
		WindowPerSecond:  int64(widowsRequestNumber),
		RequestPerSecond: int64(requestNumber),
		TTL:              int64(tllNumber),
	}
}

func NewPolicyIP() *Policy {
	request := os.Getenv("REQUEST_PER_SECOND_IP")
	widowsRequest := os.Getenv("REQUEST_PER_WINDOW")
	tll := os.Getenv("TOLKEN_EXPIRATION")

	requestNumber, _ := strconv.Atoi(request)
	widowsRequestNumber, _ := strconv.Atoi(widowsRequest)
	tllNumber, _ := strconv.Atoi(tll)

	return &Policy{
		Fonte:            "IP",
		WindowPerSecond:  int64(widowsRequestNumber),
		RequestPerSecond: int64(requestNumber),
		TTL:              int64(tllNumber),
	}
}
