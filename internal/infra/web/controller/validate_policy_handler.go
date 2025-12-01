package controller

import (
	"net/http"
)

type PolicyController struct {
}

func NewPolicyController() *PolicyController {
	return &PolicyController{}
}

func (tc *PolicyController) ValidatePolicy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}
