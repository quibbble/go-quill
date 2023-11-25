# Units

By default all units have:
- attack - the amount of damage dealt when attacking or retaliating.
- health - the hit points the unit can take before being destroyed.
- traits - special abilities unique to that unit.

## Creature

Creature units additionally have:
- cooldown - the number of turns before the unit may attack.
- movement - the number of spaces the unit may move a turn.
- codex - set of directions a unit may move and attack. 
  - An 8 bit value with each bit representing a direction: `up|down|left|right|upper-left|lower-right|lower-left|upper-right`.
  - Example: `11000000` - the unit may move up and down. 
- items - are held by units and provide additional traits.

## Structure

Structure units additonally:
- cannot move.
- cannot attack (but do retaliate from being attacked).
- cannot hold items.
