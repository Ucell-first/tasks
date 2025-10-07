package errors

import errs "errors"

var (
	// ErrConfigurationFileNotExists appears when configuration file
	// wasn't found using passed path.
	ErrConfigurationFileNotExists = errs.New("configuration file not found")

	// ErrConfigurationFilePathIsEmpty appears when no file path was passed
	// to configuration provider.
	ErrConfigurationFilePathIsEmpty = errs.New("configuration file path is empty")

	// ErrConfigurationIsNotSet appears when trying to start configuration
	// providers but configuration structure pointer wasn't set with RegisterServiceConfig.
	ErrConfigurationIsNotSet = errs.New("configuration structure pointer wasn't set")

	// ErrConfigurationIsNotValidatable appears when dynamically registered service
	// configuration structure cannot be asserted as configstruct.ValidatableServiceConfig
	// interface.
	ErrConfigurationIsNotValidatable = errs.New("configuration is not validatable, does not contain" +
		" 'Validate() error' function")

	// ErrConfigurationParsingFailed appears when configuration file
	// was loaded but failed to be parsed into structure.
	ErrConfigurationParsingFailed = errs.New("configuration file parsing failed")

	// ErrNoConfigurationRegistered appears when provider is trying to start
	// but unable to do so due to missing configuration (not registered).
	ErrNoConfigurationRegistered = errs.New("configuration wasn't registered")

	// ErrNotAConfigurationStruct appears when not a configuration struct
	// was passed to configuration provider's Parse() function.
	ErrNotAConfigurationStruct = errs.New("not configstruct.Struct was passed")

	// ErrNotAPtr appears when passed configuration structure isn't a pointer.
	ErrNotAPtr = errs.New("passed value isn't pointer")
)
