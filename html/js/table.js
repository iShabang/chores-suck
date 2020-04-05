function generateTableHeader(table, data) {
  let thead = table.createTHead();
  let row = thead.insertRow();
  for (let key of data) {
    let th = document.createElement("th");
    let text = document.createTextNode(key);
    th.appendChild(text);
    row.appendChild(th);
  }
}

function generateTableBody(table, data) {
  for (let section of data) {
    let row = table.insertRow();
    for (let key in section) {
      let cell = row.insertCell();
      let text = document.createTextNode(section[key]);
      cell.appendChild(text);
    }
  }
}
