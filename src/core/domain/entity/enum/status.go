package enum

import "strconv"

type Status int

const (
	StatusRejected Status = -1
	StatusPending  Status = 0
	StatusAccepted Status = 1
)

func (this Status) ToString() string {
	return strconv.FormatInt(int64(this), 10)
}

func StatusFromInt(_status int) *Status {
	var status Status
	switch _status {
	case -1:
		status = StatusRejected
		return &status
	case 0:
		status = StatusPending
		return &status
	case 1:
		status = StatusAccepted
		return &status
	}
	return nil
}
