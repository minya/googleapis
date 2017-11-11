package googleapis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/minya/goutils/web"
	"io/ioutil"
	"net/http"
)

func SignInWithEmailAndPassword(email string, password string, apiKey string) (VerifyPasswordResponse, error) {
	url := fmt.Sprintf("https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=%v", apiKey)
	client := http.Client{
		Transport: web.DefaultTransport(1000),
	}
	var reqStruct LoginAndPasswordRequest
	var responseStruct VerifyPasswordResponse
	reqStruct.Email = email
	reqStruct.Password = password
	reqStruct.ReturnSecureToken = true
	reqBytes, err := json.Marshal(reqStruct)
	if err != nil {
		return responseStruct, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		return responseStruct, err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	response, err := client.Do(req)

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return responseStruct, err
	}

	if response.StatusCode >= 400 {
		var errStruct ErrorResponse
		if err = json.Unmarshal(responseBytes, &errStruct); err != nil {
			return responseStruct, err
		}
		return responseStruct, fmt.Errorf("%v", errStruct.Error.Message)
	}
	err = json.Unmarshal(responseBytes, &responseStruct)
	if err != nil {
		return responseStruct, err
	}
	return responseStruct, nil
}

type LoginAndPasswordRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type VerifyPasswordResponse struct {
	Kind         string `json:"kind"`
	LocalId      string `json:"localId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	IdToken      string `json:"idToken"`
	Registered   bool   `json:"registered"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
