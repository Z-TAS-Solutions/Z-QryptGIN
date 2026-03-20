package database

// Add this interface at the top
type Validatable interface {
	Validate() error
}

// ... All your custom types (Email, UserCustomID, etc.) go here ...
