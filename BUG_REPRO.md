# Recovered card build publishes partial state

## Bug
A recovered section-render panic returns and caches the partially built pet card.

## Trigger
Build a card with one valid section followed by `panic-section`, then query the same key.

## Error
`failed build did not return a clean error: card=&{Key:card-6 Name:Nori Sections:[profile] Ready:false}`
