package controller

// func (h *UsuarioHandler) PegaJWT(w http.ResponseWriter, r *http.Request) {
// 	var JWTDto dto.GetJWT

// 	err := json.NewDecoder(r.Body).Decode(&JWTDto)

// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		erros := Error{Mensagem: err.Error()}
// 		json.NewEncoder(w).Encode(erros)
// 		return
// 	}

// 	u, err := h.UserDB.ProcuraPorEmail(JWTDto.Email)
// 	fmt.Println(u)

// 	if err != nil {
// 		w.WriteHeader(http.StatusNotFound)
// 		erros := Error{Mensagem: err.Error()}
// 		json.NewEncoder(w).Encode(erros)
// 		return
// 	}

// 	if !u.ValidarSenha(JWTDto.Senha) {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		erros := Error{Mensagem: "Senha incorreta"}
// 		json.NewEncoder(w).Encode(erros)
// 		return
// 	}

// 	_, tolkenString, _ := h.Jwt.Encode(map[string]interface{}{
// 		"sub": u.ID.String(),
// 		"exp": time.Now().Add(time.Second * time.Duration(h.JwtTempo)).Unix(),
// 	})

// 	tolken := dto.GetJWTOutput{AccessTolken: tolkenString}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(tolken)
// }
