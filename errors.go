package roscatalog

import "errors"

// Error sentinels that callers can inspect with errors.Is.
var (
	// ErrInvalidResource indicates an inconsistent resource definition.
	ErrInvalidResource = errors.New("roscatalog: invalid resource definition")
	// ErrDuplicateResource indicates a duplicate resource ID in a catalog.
	ErrDuplicateResource = errors.New("roscatalog: duplicate resource")
	// ErrResourceNotFound is returned by Catalog.Find when an ID does not exist.
	ErrResourceNotFound = errors.New("roscatalog: resource not found")
	// ErrViewNotFound indicates a view name that does not exist on the resource.
	ErrViewNotFound = errors.New("roscatalog: view not found")
	// ErrDeliveryUnsupported indicates a Plan delivery mode not valid for the view.
	ErrDeliveryUnsupported = errors.New("roscatalog: delivery unsupported")
	// ErrUnknownFlag indicates an undeclared print flag.
	ErrUnknownFlag = errors.New("roscatalog: unknown flag")
	// ErrInvalidFlag indicates an invalid or mutually exclusive print flag combination.
	ErrInvalidFlag = errors.New("roscatalog: invalid flag")
	// ErrUnknownColumn indicates a column outside the view field list.
	ErrUnknownColumn = errors.New("roscatalog: unknown column")
	// ErrMissingArgument indicates a required argument is missing or empty.
	ErrMissingArgument = errors.New("roscatalog: missing required argument")
	// ErrSelectorRequired indicates an action requires a selector that was not provided.
	ErrSelectorRequired = errors.New("roscatalog: selector required")
	// ErrInvalidInterval indicates a duration outside the allowed interval or poll bounds.
	ErrInvalidInterval = errors.New("roscatalog: invalid interval")
	// ErrInvalidQuery indicates a malformed or unsupported query filter.
	ErrInvalidQuery = errors.New("roscatalog: invalid query")
)
