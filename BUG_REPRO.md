# Follow-up delivery slot is not released

## Bug
The first item in a batch keeps the only delivery slot, so the second item cannot finish.

## Trigger
Send two follow-up reminders in one call while the delivery slot capacity is one.

## Error
`second follow-up never completed after the first delivery`
