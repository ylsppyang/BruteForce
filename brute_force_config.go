package bruteforce

import "errors"

const (
	BF_PROTO_HTTP = iota
	BF_PROTO_RFB
	BF_PROTO_SMTP
	BF_PROTO_SSH
	BF_PROTO_END
)

var (
	Bf_setting           = make([]uint, BF_PROTO_END)
	ErrProtocolIdInvalid = errors.New("protocol Id is invalid")
)

func Config_bf_setting(proto, threshold uint) error {
	if proto >= BF_PROTO_END {
		return ErrProtocolIdInvalid
	}
	Bf_setting[proto] = threshold
	return nil
}
