package controllers

import (
	"encoding/base64"     
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
    "encoding/pem"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/phapsu86/tamlinh/api/auth"
	"github.com/phapsu86/tamlinh/api/models"
	"github.com/phapsu86/tamlinh/api/responses"
	"github.com/phapsu86/tamlinh/api/utils/formaterror"
	"github.com/phapsu86/tamlinh/api/utils/formatresult"
)

type ParamRSA struct { 

	//Status   int  `json:"status"` 
	Data string  `json:"data"` 
	Otp string 	`json:"otp"`
}

func encodeBase64(b []byte) string {                                                                                                                                                                        
    return base64.StdEncoding.EncodeToString(b)                                                                                                                                                             
} 

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
    privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
    return privkey, &privkey.PublicKey
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
    privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
    privkey_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "RSA PRIVATE KEY",
                    Bytes: privkey_bytes,
            },
    )
    return string(privkey_pem)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode([]byte(privPEM))
    if block == nil {
            return nil, errors.New("failed to parse PEM block containing the key")
    }

    priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
            return nil, err
    }

    return priv, nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
    pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
    if err != nil {
            return "", err
    }
    pubkey_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "RSA PUBLIC KEY",
                    Bytes: pubkey_bytes,
            },
    )

    return string(pubkey_pem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(pubPEM))
    if block == nil {
            return nil, errors.New("failed to parse PEM block containing the key")
    }

    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
            return nil, err
    }

    switch pub := pub.(type) {
    case *rsa.PublicKey:
            return pub, nil
    default:
            break // fall through
    }
    return nil, errors.New("Key type is not RSA")
}



func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
 // Export the keys to pem string
 priv_pem := ExportRsaPrivateKeyAsPemStr(privateKey)
 pub_pem, _ := ExportRsaPublicKeyAsPemStr(&privateKey.PublicKey) 
	fmt.Printf("private:%s", priv_pem)
	fmt.Printf("public:%s", pub_pem)
	// The public key is a part of the *rsa.PrivateKey struct
	///publicKey := privateKey.PublicKey

	priv_parsed, _ := ParseRsaPrivateKeyFromPemStr(priv_pem)
    pub_parsed, _ := ParseRsaPublicKeyFromPemStr(pub_pem)	
	fmt.Printf("publickey:%s", priv_parsed)

	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pub_parsed,
		[]byte(`{
		"nickname":"phrapsfu2323",
        "email": "hungndf3grvfnd@geeeemail.com",
        "mobile":"09496519143",
        "address":"hanoiphog",
        "password": "chaoemcogainamdhong",
        "token_devices":"AAAAJDDDHHJXNRdddXXDKDKDGMGSMSSO",
        "otp":"442545"
        }`),
		nil)
	if err != nil {
		panic(err)
	}
//	rs:= formatresult.ReturnGlobalArray(pub_pem +"[]"+priv_pem)
//rs:= formatresult.ReturnGlobalArray(encryptedBytes)
	//responses.JSON(w, http.StatusUnprocessableEntity, rs)
	fmt.Println("encrypted bytes: ", encryptedBytes)

	decryptedBytes, err := priv_parsed.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}

	fmt.Println("decrypted message: ", string(decryptedBytes))
	//=========================================

	resultFail := ResultFail{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//formattedError := formaterror.FormatError(err.Error())
		resultFail.Status = http.StatusUnprocessableEntity
		resultFail.Msg = err.Error()
		responses.JSON(w, http.StatusUnprocessableEntity, resultFail)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		resultFail.Status = http.StatusUnprocessableEntity
		resultFail.Msg = err.Error()
		responses.JSON(w, http.StatusUnprocessableEntity, resultFail)
		return
	}
	user.Prepare()
	err = user.Validate("create")
	if err != nil {
		resultFail.Status = http.StatusUnprocessableEntity
		resultFail.Msg = err.Error()
		responses.JSON(w, http.StatusOK, resultFail)
		return
	}
	// otpModel := models.Otp{}
	// otpData, err := otpModel.FindOtpPhone(server.DB, user.Mobile)

	// if err != nil {
	// 	err = errors.New("OTP_DOES_NOT_EXIST")
	// 	resultFail.Status = http.StatusUnprocessableEntity
	// 	resultFail.Msg = err.Error()
	// 	responses.JSON(w, http.StatusOK, resultFail)
	// 	return
	// }
	// if otpData.Code != user.Otp {
	// 	err = errors.New("OTP_WRONG")
	// 	resultFail.Status = http.StatusUnprocessableEntity
	// 	resultFail.Msg = err.Error()
	// 	responses.JSON(w, http.StatusOK, resultFail)
	// 	return
	// }

	_, err = user.SaveUser(server.DB)

	if err != nil {

		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	//w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	result := ResultLogin{}
	result.Status = http.StatusOK
	result.Data = "success"
	responses.JSON(w, http.StatusCreated, result)
}

func (server *Server) CreateUserRSA(w http.ResponseWriter, r *http.Request) {

	// privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	// if err != nil {
	// 	panic(err)
	// }
// Export the keys to pem string
//priv_pem := ExportRsaPrivateKeyAsPemStr(privateKey)
//pub_pem, _ := ExportRsaPublicKeyAsPemStr(&privateKey.PublicKey) 
// The public key is a part of the *rsa.PrivateKey struct
//publicKey := privateKey.PublicKey


	rsaModel := ParamRSA{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//formattedError := formaterror.FormatError(err.Error())
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	err = json.Unmarshal(body, &rsaModel)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

 //pub_pem:= "-----BEGIN RSA PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAo84jWTvgTX+msSNrYt1F\nOPmvrkJMh+1MOP1rM+bQD9JHSv21ivE4iJ6zUe7GIN5dCtcdWR1ZPrR7byUHtHS+\nBV4S9O/SIAJy+PzZu+sdGHDTQAbpEj+fcaunW0PsoHX+zf+yNRaKQhgiLx5QGwG1\nZve1tfjEkAxEJIPf14UGoyZx6Om3QIC6gxfU9lf1b0UXzllJVQQ8GN+Ts8u3KU0y\nj9JY8MwgkrPqIFL48hNL1krvwBQwTh4WZIoycYP5HaLtrKc2JHf8GROwbqvHiwOH\nBxAHQovCQvPCiD4YXwYenOGiKd+v8cO5zvedfhm5f6QdjxPWl+b/4pINS9hIZsQk\nxmtCPAzbI6JipV/ST0m4aqJxwMOA6nx70kvW0qmdRRvSOZt5bV9RHKQDca7ubPVH\nGbm2Wz1oSUaZQ3DE8zYgh+9SMcvnswp6iM6BqnwYA8BAXwp33TNda6QXfvYeK6cQ\nXxisYrF498QyfAyPY+fIlqEWBoCQNiIcR6wsNQ09MCLHozOtTOte+8K3ym96nWDj\nf88TWAQUniX882Ia5mKCfG1vtXfhwyDFda0fHwLkXaILbq1fluxNdpxN6Eb5qpyg\n513+ZbdLUOl8zckINeCFUdM/5knedRXkYwbtvimx3Do23y2M6S/yEpH2meQTPWeN\nb2l3ymqk4Z6MR62p07DSFKMCAwEAAQ==\n-----END RSA PUBLIC KEY-----\n"
 priv_pem:= "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"	
 //fmt.Printf("private:%s", priv_pem)
//	fmt.Printf("public:%s", pub_pem)
	priv_parsed, _ := ParseRsaPrivateKeyFromPemStr(priv_pem)
   // pub_parsed, _ := ParseRsaPublicKeyFromPemStr(pub_pem)	
	//fmt.Printf("publickey:%s", priv_parsed)
	// encryptedBytes, err := rsa.EncryptOAEP(
	// 	sha256.New(),
	// 	rand.Reader,
	// 	pub_parsed,
	// 	[]byte(`{
	// 	"nickname":"phrapsfu2323",
    //     "email": "hungndf3grvfnd@geeeemail.com",
    //     "mobile":"09496519143",
    //     "address":"hanoiphog",
    //     "password": "chaoemcogainamdhong",
    //     "token_devices":"AAAAJDDDHHJXNRdddXXDKDKDGMGSMSSO",
    //     "otp":"442545"
    //     }`),
	// 	nil)
	// if err != nil {
	// 	panic(err)
	// }
	//rs:= formatresult.ReturnGlobalArray(pub_pem +"[]"+priv_pem)
   // rs:= formatresult.ReturnGlobalArray(encryptedBytes)
	//responses.JSON(w, http.StatusUnprocessableEntity, encodeBase64(encryptedBytes))
	//fmt.Println("encrypted bytes: ", encryptedBytes)

	data, err := base64.StdEncoding.DecodeString(string(rsaModel.Data))
    if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
    }

	decryptedBytes, err := priv_parsed.Decrypt(nil, data, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}

	fmt.Println("decrypted message: ", string(decryptedBytes))
	//=========================================
	user := models.User{}
	err = json.Unmarshal(decryptedBytes, &user)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
		
	}
	user.Prepare()
	err = user.Validate("create")
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
		
	}
	// otpModel := models.Otp{}
	// otpData, err := otpModel.FindOtpPhone(server.DB, user.Mobile)

	// if err != nil {
	// 	err = errors.New("OTP_DOES_NOT_EXIST")
	// 	resultFail.Status = http.StatusUnprocessableEntity
	// 	resultFail.Msg = err.Error()
	// 	responses.JSON(w, http.StatusOK, resultFail)
	// 	return
	// }
	// if otpData.Code != user.Otp {
	// 	err = errors.New("OTP_WRONG")
	// 	resultFail.Status = http.StatusUnprocessableEntity
	// 	resultFail.Msg = err.Error()
	// 	responses.JSON(w, http.StatusOK, resultFail)
	// 	return
	// }

	_, err = user.SaveUser(server.DB)

	if err != nil {

		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	///w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	result := ResultLogin{}
	result.Status = http.StatusOK
	result.Data = "success"
	responses.JSON(w, http.StatusCreated, result)
}




func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint64(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userGotten)
}

func (server *Server) GetProfileUser(w http.ResponseWriter, r *http.Request) {

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		errResult := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, errResult)
		return
	}

	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint64(tokenID))
	//Get link for avartar
	if err != nil {
		responses.JSON(w, http.StatusOK, formaterror.ReturnErr(err))
		return
	}
	if userGotten.Avartar != "" {
		userGotten.Avartar = server.getLink(userGotten.Avartar, "avatar")

	}

	result := formatresult.ReturnUser(userGotten)

	responses.JSON(w, http.StatusOK, result)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	// tokenID, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	errResult := formaterror.ReturnErr(errors.New("Unauthorized"))
	// 	responses.JSON(w, http.StatusOK, errResult)
	// 	return
	// }

	// vars := mux.Vars(r)
	// uid, err := strconv.ParseUint(vars["id"], 10, 32)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	if tokenID == 0 {
		err := formaterror.ReturnErr(errors.New(http.StatusText(http.StatusUnauthorized)))
		responses.JSON(w, http.StatusOK, err)

		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, tokenID)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	if updatedUser.ID != 0 {
		sucess := formaterror.ReturnSuccess()
		responses.JSON(w, http.StatusOK, sucess)
	}

}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint64(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}
