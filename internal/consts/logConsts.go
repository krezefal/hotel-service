package consts

// env vars error messages
const (
	EnvVarNotSet = "env var is not set: %s"
)

// config error messages
const (
	CfgNotExists = "config file doesn't exist: %s"
	CfgErrorGet  = "error getting config file: %s"
	CfgIsInvalid = "invalid config: %s"

	CfgNoEnv          = "no environment specified"
	CfgInvalidEnv     = "invalid env: %s"
	CfgNoStorage      = "no storage path specified"
	CfgInvalidStorage = "invalid storage path: %s"
)

// request error messages
const (
	ReqDecodeFail = "failed to decode request"
)

// request info messages
const (
	ReqBodyDecoded = "request body decoded"
)

// storage error messages
const (
	StorageInternalError    = "internal storage error"
	StorageExistedOrderID   = "order with provided id already exists"
	StorageInvalidOrderID   = "order not found by provided id"
	StorageInvalidHotelID   = "hotel not found by provided id"
	StorageInvalidRoomID    = "room not found by provided id"
	StorageUnavailableDate  = "unavailable date"
	StorageNoRoomsForPeriod = "no available rooms for provided period"
)

// storage info messages
const (
	StorageDataFetched      = "Data was successfully fetched from the storage"
	StorageInventoryUpdated = "Hotels & rooms inventory was successfully updated"
	StorageOrderCreated     = "New order with id %s was successfully created"
	StorageOrderCancelled   = "Order with id %s was successfully cancelled by user request"
)
