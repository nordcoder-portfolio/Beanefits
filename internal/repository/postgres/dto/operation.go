package dto

type OperationRecord struct {
	AccountID   int64
	OpType      OperationType
	OperationID string

	RequestJSON  JSON
	ResponseJSON *JSON
	HTTPStatus   *int

	CreatedAt Ts
}

type OperationPendingInsert struct {
	AccountID   int64
	OpType      OperationType
	OperationID string
	RequestJSON JSON
}

type OperationFinalize struct {
	AccountID   int64
	OpType      OperationType
	OperationID string

	HTTPStatus   int
	ResponseJSON JSON
}
