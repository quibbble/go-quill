# Choose

## Types

| **Name**         | **Description**                                                             | **Args**                                  |
|------------------|-----------------------------------------------------------------------------|-------------------------------------------|
| `Adjacent`       | Retrieve all `Types` adjacent to `ChooseUnitOrTile`.                        | `Types`, `ChooseUnitOrTile`               |
| `Codex`          | Retrieve all `Types` matching `Codex` to `ChooseUnitOrTile`.                | `Types`, `Codex`, `ChooseUnitOrTile`      |
| `Composite`      | Apply `SetFunction` to all `ChooseChain`.                                   | `SetFunction`, `ChooseChain`              |
| `Connected`      | Retrieve all `Types` connected to `ChooseUnit` using `ConnectionType`.      | `Types`, `ConnectionType`, `ChooseUnit`   |
| `CurrentPlayer`  | Retrieve the player who has the active turn.                                |                                           |
| `HookEventTile`  | Retrieve the tile affected by event found in the `HookEvent` context.       |                                           |
| `HookEventUnit`  | Retrieve the unit affected by event found in the `HookEvent` context.       |                                           |
| `ItemHolder   `  | Retrieve the unit that holds the `ChooseItem`.                              | `ChooseItem`                              |
| `OpposingPlayer` | Retrieve the player who does not have the active turn.                      |                                           |
| `OwnedTiles`     | Retrieve the set of tiles owned by `ChoosePlayer`.                          | `ChoosePlayer`                            |
| `OwnedUnits`     | Retrieve the set of units owned by `ChoosePlayer`.                          | `ChoosePlayer`                            |
| `Owner`          | Retrieve the owner of `ChooseCard`.                                         | `ChooseCard`                              |
| `Random`         | Retrieve one random choice from `Choose`. Return zero if `Choose` is empty. | `Choose`                                  |
| `Self`           | Retrieve self.                                                              |                                           |
| `Ranged`         | Retrieve all `Types` within `Range` of `ChooseUnitOrTile`.                  | `Types`, `Range`, `ChooseUnitOrTile`      |
| `Target`         | Retrieve the target found in targets list at index `Index`.                 | `Index`                                   |
| `Tiles`          | Retrieve a set of tiles that are optionally `Empty`.                        | `Empty`                                   |
| `TraitEventTile` | Retrieve the tile affected by event found in the `TraitEvent` context.      |                                           |
| `TraitEventUnit` | Retrieve the unit affected by event found in the `TraitEvent` context.      |                                           |
| `Units`          | Retrieve a set of units on the board that have a type in `Types`.           | `Types`                                   |
| `UUID`           | Retrieve the given `UUID`.                                                  | `UUID`                                    |

# Argsx

| **Name**           | **Requirements**                                                                                 |
|--------------------|--------------------------------------------------------------------------------------------------|
| `ChooseChain`      | List of [Choose](./choose.md).                                                                   |
| `Choose{X}`        | [Choose](./choose.md).                                                                           |
| `Codex`            | An eight character string containing only 0 or 1 i.e. `11001111`.                                |
| `ConnectionType`   | `Adjacent` or `Codex`                                                                            |
| `Empty`            | `true` or `false`.                                                                               |
| `Index`            | An integer.                                                                                      |
| `Range`            | An integer.                                                                                      |
| `SetFunction`      | `Union`, `Intersect`, `Diff`.                                                                    |
| `Types`            | A list containing one or more of the following: `Tile`, `Unit`, `Creature`, `Structure`, `Base`. |
| `UUID`             | A UUID.                                                                                          |
