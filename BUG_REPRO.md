# Cancelled refresh keeps retrying

## Bug
A cancelled pet refresh returns to its caller while background work continues retrying and delays shutdown.

## Trigger
Let the first refresh call start, cancel the request, release the client call, and wait for coordinator shutdown.

## Error
`refresh calls kept growing after cancellation: 3`
