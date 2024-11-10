package pokemon

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
)

func GenerateEtag(data any) string {
	jsonBytes, _ := json.Marshal(data)
    hash := sha1.Sum(jsonBytes)
    return fmt.Sprintf(`"%x"`, hash[:4])
}