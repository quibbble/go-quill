# go-quill

## TODO
- Add an endpoint to quibbble that is GET info gameKey={KEY}
  - Game key
  - Returns min and max players
  - Variants (if there are any)
  - Any other info needed when creating a game
    - In the case of quill - all the cards that currently exist in the game
  - These fields are all constants that should be returned from Builder.GetInfo() *bg.Info
- Create randomness
  - Write code that puts sorts all cards into categories by:
    - Cost
- Create cards that have synergy together
  - Add groups? Maybe this will help with building synergy 

# Ideas
- Add unstable
- Add cursed X trait
- Add dominion trait
- Add friendly fire trait (attacked by friendly unit heals it)
- Add invincible trait
- Add traps
- Add watcher
- Add bombard
- Add rumble
- Add big cards like imperial resouces, god hand, necro, VD, etc
