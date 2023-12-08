# Choose

## Types

| **Type**         | **Description**                                                             | **Args**                                  |
|------------------|-----------------------------------------------------------------------------|-------------------------------------------|
| `Adjacent`       | Retrieve all `Types` adjacent to `ChooseUnitOrTile`.                        | `Types` - `ChooseUnitOrTile`              |
| `Codex`          | Retrieve all `Types` matching `Codex` to `ChooseUnitOrTile`.                | `Types` - `Codex` - `ChooseUnitOrTile`    |
| `Composite`      | Apply `SetFunction` to all `Choices`.                                       | `SetFunction` - `Choices`                 |
| `Connected`      | Retrieve all `Types` connected to `ChooseUnit` using `ConnectionType`.      | `Types` - `ConnectionType` - `ChooseUnit` |
| `CurrentPlayer`  | Retrieve the player who has the active turn.                                |                                           |
| `EventUnit`      | Retrieve the unit affected by event found in the `HookEvent` context.       |                                           |
| `OpposingPlayer` | Retrieve the player who does not have the active turn.                      |                                           |
| `OwnedTiles`     | Retrieve the set of tiles owned by `ChoosePlayer`.                          | `ChoosePlayer`                            |
| `OwnedUnits`     | Retrieve the set of units owned by `ChoosePlayer`.                          | `ChoosePlayer`                            |
| `Random`         | Retrieve one random choice from `Choose`. Return zero if `Choose` is empty. | `Choose`                                  |
| `SelfOwner`      | Retrieve the owner of self.                                                 |                                           |
| `Self`           | Retrieve self.                                                              |                                           |
| `TargetOwner`    | Retrieve the owner of target found in targets list at index `Index`.        | `Index`                                   |
| `Target`         | Retrieve the target found in targets list at index `Index`.                 | `Index`                                   |
| `Tiles`          | Retrieve a set of tiles that are optionally `Empty`.                        | `Empty`                                   |
| `Units`          | Retrieve a set of units on the board that have a type in `Types`.           | `Types`                                   |
| `UUID`           | Retrieve the given `UUID`.                                                  | `UUID`                                    |

# Args

| **Name**           | **Data Structure**         | **Additional Requirements**                            |
|--------------------|----------------------------|--------------------------------------------------------|
| `Types`            | `[]string`                 | `Tile`, `Unit`, `Creature`, `Structure`, and/or `Base` |
| `Choose`           | `map[string]interface{}`   | N/A                                                    |
| `ChooseUnit`       | `map[string]interface{}`   | N/A                                                    |
| `ChoosePlayer`     | `map[string]interface{}`   | N/A                                                    |
| `ChooseUnitOrTile` | `map[string]interface{}`   | N/A                                                    |
| `Choices`          | `[]map[string]interface{}` | N/A                                                    |
| `Codex`            | `string`                   | 8 char containing only 0 or 1 i.e. `11001111`          |
| `ConnectionType`   | `string`                   | `Adjacent` or `Codex`                                  |
| `Empty`            | `bool`                     | N/A                                                    |
| `SetFunction`      | `string`                   | `Union` or `Intersect`                                 |
| `Index`            | `int`                      | N/A                                                    |
| `UUID`             | `uuid.UUID`                | N/A                                                    |
