# Care batch ownership leak

## Bug
Saved care labels share mutable storage with both the caller input and later query results.

## Trigger
Save two labels, reuse the input slice, fetch the batch, edit that response, and fetch the same pet again.

## Error
`saved labels changed after input reuse: saved=[caller-reused-buffer soft-food] fetched=[caller-reused-buffer soft-food]`
