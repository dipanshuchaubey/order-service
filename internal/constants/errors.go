package constants

// Error Codes
const (
	ErrorInvalidQueueUrl  = "INVALID_QUEUE_URL"
	ErrorInvalidRegion    = "INVALID_REGION"
	ErrorDecondingMessage = "DECODING_MESSAGE_ERROR"
	MySQLFetchError       = "MYSQL_FETCH_ERROR"
	RedisError            = "REDIS_ERROR"
	MarshalError          = "MARSHAL_ERROR"
)

// Error Messages
const (
	ErrorFetchingOrders = "Error fetching orders: %s"
	ErrorMarshalling    = "Error marshalling data:: %s"
	RedisSetError       = "Error setting value in redis: %s"
	RedisGetError       = "Error getting value from redis: %s"
)
