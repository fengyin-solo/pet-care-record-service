# Cancelled health sync is recorded

## Bug
A health sync continues after request cancellation and records a stale successful summary.

## Trigger
Start the downstream load, cancel its request context, then allow the load call to return.

## Error
`cancel did not reach source: err=<nil> sawCancel=false`
