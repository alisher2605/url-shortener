package snowflake

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/snowflake"
)

func GenerateSnowflakeId() (string, error) {
	id, err := snowflakeId()
	if err != nil {
		return "", err
	}

	return base62Conversion(id), nil
}

func snowflakeId() (int64, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Errorf("GenerateShortURL failed create new node. err = %v", err)

		return 0, err
	}

	ID := node.Generate().Int64()

	return ID, nil
}

func base62Conversion(id int64) string {
	const chars string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var buf bytes.Buffer
	uID := uint64(id)
	if uID == 0 {
		buf.WriteByte(chars[0])
	} else {
		encode(uID, &buf, chars)
	}

	return buf.String()
}

func encode(id uint64, buf *bytes.Buffer, chars string) {
	l := uint64(len(chars))
	if id/l != 0 {
		encode(id/l, buf, chars)
	}

	buf.WriteByte(chars[id%l])
}
