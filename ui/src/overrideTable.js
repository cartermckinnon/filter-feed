const { OverrideSpec } = require('./api/v1_pb.js');
const { getKeyForValue } = require('./util.js');

export class OverrideTable {
    table;
    footer;
    request;
    tableBody;
    requestChangeCallback;

    constructor(request) {
        this.request = request;

        let table = document.createElement('table');
        this.table = table;

        let thead = document.createElement('thead');
        thead.innerHTML = `
            <tr>
                <th>Target</th>
                <th>Value</th>
                <th></th>
            </tr>`;

        let tbody = document.createElement('tbody');
        this.tableBody = tbody;

        for (let spec of request.getOverridesList()) {
            let row = new OverrideRow(spec);
            row.setRemoveButtonCallback(() => this.callback_removeRow(row));
            tbody.appendChild(row.getElement());
        }

        let footer = new OverrideTableFooter();
        this.footer = footer;
        footer.setAddButtonCallback(() => this.callback_addRow());
        tbody.appendChild(footer.getElement());

        table.appendChild(thead);
        table.appendChild(tbody);
    }

    callback_addRow() {
        let row = this.footer.createRowFromInputs();
        if (row == null) {
            return;
        }
        this.request.addOverrides(row.getSpec());
        this.requestChangeCallback(this.request);
        row.setRemoveButtonCallback(() => this.callback_removeRow(row));
        this.footer.prependSibling(row.getElement());
    }

    removeSpecFromRequest(specToRemove) {
        let newList = [];
        let modified = false;
        for (let spec of this.request.getOverridesList()) {
            if (spec !== specToRemove) {
                newList.push(spec);
            } else {
                modified = true;
            }
        }
        if (modified) {
            this.request.setOverridesList(newList);
            this.requestChangeCallback(this.request);
        }
    }

    callback_removeRow(row) {
        let spec = row.getSpec();
        this.removeSpecFromRequest(spec);
        row.remove();
    }

    getElement() {
        return this.table;
    }

    setRequestChangeCallback(callback) {
        this.requestChangeCallback = callback;
    }
}

class OverrideRow {
    tableRow;
    removeButton;
    spec;

    constructor(spec) {
        this.spec = spec;

        let tr = document.createElement('tr');
        this.tableRow = tr;
        tr.classList.add('override-row');

        let targetCol = document.createElement('td');
        targetCol.innerText = getKeyForValue(OverrideSpec.OverrideTarget, spec.getTarget());
        tr.appendChild(targetCol);

        let valueCol = document.createElement('td');
        valueCol.innerText = spec.getValue();
        tr.appendChild(valueCol);

        let removeButtonCol = document.createElement('td');
        let removeButton = document.createElement('a');
        this.removeButton = removeButton;
        removeButton.classList.add('btn-floating');
        removeButton.classList.add('btn-small');
        removeButton.classList.add('red');
        removeButton.classList.add('icon-button');
        let removeIcon = document.createElement('i');
        removeIcon.classList.add('material-icons');
        removeIcon.innerText = "remove";
        removeButton.appendChild(removeIcon);
        removeButtonCol.appendChild(removeButton);
        tr.appendChild(removeButtonCol);
    }

    setRemoveButtonCallback(callback) {
        this.removeButton.onclick = (ev) => callback();
    }

    remove() {
        this.tableRow.parentElement.removeChild(this.tableRow);
    }

    getElement() {
        return this.tableRow;
    }

    getSpec() {
        return this.spec;
    }
}

class OverrideTableFooter {
    tableRow;
    addButton;
    inputs;

    constructor() {
        this.inputs = {};

        let tr = document.createElement('tr');
        this.tableRow = tr;
        tr.classList.add('table-footer-row');

        let targetSelectCol = this.createSelectColumn('Target', OverrideSpec.OverrideTarget);
        let valueInputCol = this.createInputColumn('Value', 'text');

        tr.appendChild(targetSelectCol);
        tr.appendChild(valueInputCol);

        let addButtonCol = document.createElement('td');
        let addButton = document.createElement('a');
        this.addButton = addButton;
        addButton.classList.add('btn-floating');
        addButton.classList.add('btn-small');
        addButton.classList.add('green');
        addButton.classList.add('icon-button');
        let addIcon = document.createElement('i');
        addIcon.classList.add('material-icons');
        addIcon.innerText = "add";
        addButton.appendChild(addIcon);
        addButtonCol.appendChild(addButton);
        tr.appendChild(addButtonCol);

    }

    getElement() {
        return this.tableRow;
    }

    setAddButtonCallback(callback) {
        this.addButton.onclick = callback;
    }

    prependSibling(element) {
        this.tableRow.parentElement.insertBefore(element, this.tableRow);
    }

    validateInputs() {
        return this.inputs["target"].value != "" &&
            this.inputs["value"].value != "";
    }

    createRowFromInputs() {
        if (!this.validateInputs()) {
            return null;
        }

        let target = this.inputs["target"].value;
        let value = this.inputs["value"].value;

        let spec = new OverrideSpec();
        spec.setTarget(OverrideSpec.OverrideTarget[target]);
        spec.setValue(value);

        return new OverrideRow(spec);
    }

    createSelectColumn(name, enumType) {
        let col = document.createElement('td');
        let select = document.createElement('select');
        let defaultOption = document.createElement('option');
        defaultOption.innerText = this.capitalizeFirstLetter(name);
        defaultOption.disabled = true;
        defaultOption.selected = true;
        defaultOption.value = "";
        select.appendChild(defaultOption);
        this.inputs[name.toLowerCase()] = select;
        for (let option in enumType) {
            let opt = document.createElement('option');
            opt.value = option;
            opt.innerText = this.capitalizeFirstLetter(option);
            select.appendChild(opt);
        }
        col.appendChild(select);
        return col;
    }

    createInputColumn(name, type) {
        let col = document.createElement('td');
        let input = document.createElement('input');
        input.type = type;
        input.placeholder = this.capitalizeFirstLetter(name);
        this.inputs[name.toLowerCase()] = input;
        col.appendChild(input);
        return col;
    }

    capitalizeFirstLetter(string) {
        return string.charAt(0).toUpperCase() + string.slice(1).toLowerCase();
    }
}