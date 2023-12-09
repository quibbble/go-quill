# Condition

## Types

| **Name**     | **Description**                                    | **Args**                       |
|---------------|----------------------------------------------------|--------------------------------|
| `Contains`    | Pass when `Choice` is in `Choices`.                | `Choices`, `Choice`            |
| `Fail`        | Condition always fails.                            |                                |
| `ManaAbove`   | Pass when `ChoosePlayer` mana is above `Amount`.   | `ChoosePlayer`, `Amount`       |
| `ManaBelow`   | Pass when `ChoosePlayer` mana is below `Amount`.   | `ChoosePlayer`, `Amount`       |
| `Match`       | Pass when `ChooseA` matches `ChooseB`.             | `ChooseA`, `ChooseB`           |
| `StatAbove`   | Pass when `ChooseCard`'s `Stat` is above `Amount`. | `ChooseCard`, `Stat`, `Amount` |
| `StatBelow`   | Pass when `ChooseCard`'s `Stat` is below `Amount`. | `ChooseCard`, `Stat`, `Amount` |
| `UnitMissing` | Pass when `ChooseUnit` is not on the board.        | `ChooseUnit`                   |

## Args

| **Name**           | **Requirements**                                             |
|--------------------|--------------------------------------------------------------|
| `Amount`           | An integer.                                                  |
| `Choice`           | `Choose`.                                                    |
| `Choices`          | List of `Choose`.                                            |
| `Choose{X}`        | `Choose`.                                                    |
| `Stat`             | `Cost`, `Attack`, `Health`, `Movement`, `Cooldown`, `Range`. |
