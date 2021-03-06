package handler

import (
	"backend-qrcode/db"
	"backend-qrcode/model"

	customHTTP "backend-qrcode/http"
	"encoding/json"
	"net/http"
)

// Login ...
func Login(w http.ResponseWriter, r *http.Request) {

	var params model.LoginParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		customHTTP.NewErrorResponse(w, http.StatusUnauthorized, "Error: "+err.Error())
		return
	}

	var user model.User

	if db.DB.First(&user, &model.User{
		Username: params.Username,
	}).RecordNotFound() {
		customHTTP.NewErrorResponse(w, http.StatusNotFound, "Error: NotFound")
		return
	}

	if user.CheckPassword(params.Password) {
		if token, err := user.GenerateJWT(); err != nil {
			customHTTP.NewErrorResponse(w, http.StatusUnauthorized, "Error: Password Wrong")
		} else {
			result := &model.LoginReturn{
				RoleID: user.RoleID,
				Token:  token.Token,
			}
			json.NewEncoder(w).Encode(&result)
		}
		return
	}

	customHTTP.NewErrorResponse(w, http.StatusUnauthorized, "Error: Password Wrong")

}
