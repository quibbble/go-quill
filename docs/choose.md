# Choose
<table>
<thead>
  <tr>
    <th>Type</th>
    <th>Description</th>
    <th>Args</th>
  </tr>
</thead>
<tbody>
  <tr>
    <td>`Adjacent`</td>
    <td>Retrieve all `Types` adjacent to `ChooseUnitOrTile`</td>
    <td>`Types []string`, `ChooseUnitOrTile Choose`</td>
  </tr>
  <tr>
    <td>`Codex`</td>
    <td>Retrieve all `Types` matching `Codex` to `ChooseUnitOrTile`</td>
    <td>`Types []string`, `Codex string`, `ChooseUnitOrTile Choose`</td>
  </tr>
  <tr>
    <td>`Composite`</td>
    <td>Apply `SetFunction` to all `Choices`</td>
    <td>`SetFunction string`, `Choices []Choose`</td>
  </tr>
  <tr>
    <td>`Connected`</td>
    <td>Retrieve all `Types` connected to `ChooseUnit` using `ConnectionType`</td>
    <td>`Types []string`, `ConnectionType string`, `ChooseUnit Choose`</td>
  </tr>
</tbody>
</table>
