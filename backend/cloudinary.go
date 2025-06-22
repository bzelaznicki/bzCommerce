package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (cfg *apiConfig) handleCloudinarySignUpload(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Folder string `json:"folder"`
	}

	req := request{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&req)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to decode parameters")
		return
	}
	folder := strings.TrimSpace(req.Folder)
	if folder == "" {
		respondWithError(w, http.StatusBadRequest, "Missing folder name")
		return
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	params := map[string]string{
		"timestamp": timestamp,
		"folder":    folder,
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var paramString string
	for _, k := range keys {
		paramString += fmt.Sprintf("%s=%s&", k, params[k])
	}
	paramString = paramString[:len(paramString)-1]

	stringToSign := paramString + os.Getenv("CLOUDINARY_API_SECRET")

	h := sha1.New()
	h.Write([]byte(stringToSign))
	signature := hex.EncodeToString(h.Sum(nil))

	resp := map[string]string{
		"timestamp":  timestamp,
		"signature":  signature,
		"api_key":    os.Getenv("CLOUDINARY_API_KEY"),
		"cloud_name": os.Getenv("CLOUDINARY_CLOUD_NAME"),
	}

	respondWithJSON(w, http.StatusOK, resp)
}
