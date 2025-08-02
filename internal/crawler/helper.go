package crawler

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func InjectCredentials(header *http.Header) error {
	uuid, err := uuid.NewRandom()

	if err != nil {
		return err
	}

	epoch := time.Now().UnixMilli()
	epochStr := strconv.FormatInt(epoch, 10)
	secret := epochStr + ":Bus2019M@p_"
	proof := md5.Sum([]byte(secret))

	header.Add("device-id", uuid.String())
	header.Add("epoch", epochStr)
	header.Add("proof", hex.EncodeToString(proof[:]))
	header.Add("client-version", "ios|22")

	return nil
}
