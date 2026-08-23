# Health update retry duplicates history

## Bug
A successful second publish still returns the first busy error and creates another logical update history entry.

## Trigger
Update one pet to `recovering`, make the publisher fail temporarily once, then let its retry succeed.

## Error
`retry succeeded downstream but caller still received failure: health publisher busy`
