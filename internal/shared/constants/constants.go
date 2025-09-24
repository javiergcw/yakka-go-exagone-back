package constants

// Server constants
const (
	DefaultPort        = "8080"
	DefaultEnvironment = "development"
	DefaultLogLevel    = "info"
)

// Database constants
const (
	DefaultDBPort  = 5432
	DefaultSSLMode = "disable"
	DefaultDBHost  = "localhost"
	DefaultDBUser  = "postgres"
	DefaultDBName  = "yakka_db"
)

// API constants
const (
	APIVersion = "v1"
	APIPrefix  = "/api/" + APIVersion
)

// HTTP status messages
const (
	StatusHealthy = "healthy"
	StatusRunning = "running"
)

// Application info
const (
	AppName    = "yakka-backend"
	AppVersion = "1.2.0"
)

// Database table names
const (
	UsersTable = "users"
)

// HTTP headers
const (
	ContentTypeJSON = "application/json"
	Authorization   = "Authorization"
)

// Error messages
const (
	ErrInvalidJSON   = "Invalid JSON format"
	ErrInvalidUserID = "Invalid user ID"
	ErrUserNotFound  = "User not found"
	ErrUserExists    = "User already exists"
	ErrValidation    = "Validation error"
	ErrDatabase      = "Database error"
	ErrInternal      = "Internal server error"
)

// Success messages
const (
	MsgUserCreated    = "User created successfully"
	MsgUserUpdated    = "User updated successfully"
	MsgUserDeleted    = "User deleted successfully"
	MsgUsersRetrieved = "Users retrieved successfully"
	MsgUserFound      = "User found"
)
