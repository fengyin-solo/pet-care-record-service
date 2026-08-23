# Stale health review overwrites retry

## Bug
A delayed callback from version one is accepted after version two succeeds, splitting detail and summary states.

## Trigger
Start version one, complete a version-two retry, then deliver the version-one callback last.

## Error
`late callback from the first attempt was accepted`
