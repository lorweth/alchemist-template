package pagination

const (
	defaultLimit  = 30
	defaultOffset = 0
)

type Input struct {
	Page      int  `json:"page"`
	Size      int  `json:"size"`
	WithTotal bool `json:"with_total"`
}

func (pagination Input) IsValid() error {
	if pagination.Page < 1 {
		return ErrPageNumberInvalid
	}

	if pagination.Size < 1 {
		return ErrPageSizeInvalid
	}

	return nil
}

func ToOffsetLimit(pagination Input) (int, int) {
	limit := pagination.Size
	if limit == 0 {
		limit = defaultLimit
	}

	offset := limit * (pagination.Page - 1)
	if offset < 0 {
		offset = defaultOffset
	}

	return offset, limit
}
