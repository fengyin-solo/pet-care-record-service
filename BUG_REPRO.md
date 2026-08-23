# Request scope identity crosses pet lookups

## Bug
An asynchronous audit reads a pooled request scope after another lookup has overwritten its owner, pet, and label.

## Trigger
Pause lookup A's audit, return its scope to the pool, reuse it for lookup B, then release both audit writes.

## Error
`first lookup inherited the next request identity: [{OwnerID:owner-B PetID:pet-B Label:alert} ...]`
