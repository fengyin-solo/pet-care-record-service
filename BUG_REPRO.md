# Temporary reminder error stops retry

## Bug
A typed temporary delivery error loses its classification, ending the retry and storing a failed state.

## Trigger
Return a temporary busy error on the first send and success on the second available send.

## Error
`temporary failure did not recover: gateway busy`
