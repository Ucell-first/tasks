package errors

import errs "errors"

var (
	// ErrDatabaseConfigurationConnectionsCountInvalid appears when trying
	// to set configuration with invalid MaxIdleConnections or
	// MaxOpenedConnections counts.
	ErrDatabaseConfigurationConnectionsCountInvalid = errs.New("invalid connections count(s)")

	// ErrDatabaseConfigurationMasterDSNInvalid appears when trying to set
	// configuration with invalid connection URI (empty, invalid format, etc.)
	ErrDatabaseConfigurationMasterDSNInvalid = errs.New("invalid URI for database connection")

	// ErrDatabaseNoUsableConnections appears when trying to get database
	// connection but no usable connections was found (not initialized, not
	// ready, etc.)
	ErrDatabaseNoUsableConnections = errs.New("no usable database connections")

	// ErrDatabaseSlaveDSNsIsEmpty appears when we're trying to create read-only
	// connections, which are using slave_dsns parameter, but slave_dsns was
	// undefined.
	ErrDatabaseSlaveDSNsIsEmpty = errs.New("slave DSNs wasn't defined")
)
