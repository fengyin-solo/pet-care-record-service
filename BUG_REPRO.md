# Optional pet draft validation panics

## Bug
Disabling the optional validator still enters its method through a non-nil interface and leaves the draft stored.

## Trigger
Import a valid pet draft with optional validation disabled, then read the same draft key.

## Error
`pet draft validation failed: runtime error: invalid memory address or nil pointer dereference`
