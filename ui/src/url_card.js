const { FetchFeedRequest, FilterTarget, FilterType, FilterEffect, FilterSpec } = require('./api/v1_pb.js');
const { Materialize } = require('./materialize.min.js');

export class URLCard {
    element;
    request;
    header;
    filterTable;

    constructor(url) {
        let request = parseRequestFromURL(url);
        if (request == null) {
            request = new FetchFeedRequest();
            request.setFeedurl(url);
        }
        this.request = request;

        let div = document.createElement('div');
        this.element = div;
        div.classList.add('card');

        let header = new URLCardHeader(url, "-");
        this.header = header;
        div.appendChild(header.getElement());


        div.appendChild(this.createFilterSection(request));

        // this is a ridiculous hack to get the <select>-s rendered properly by Materialize
        setTimeout(() => {
            M.FormSelect.init(document.querySelectorAll('select'));
        }, 250);
    }

    onRequestChange(request) {
        console.log("request changed");
        this.request = request;
        this.header.setFilteredURL(createFilteredURL(request));
    }

    createFilterSection(request) {
        let div = document.createElement('div');
        div.classList.add('card-content');
        div.classList.add('card-row');

        let filterHeader = document.createElement('h5');
        filterHeader.innerText = "Filters";
        div.appendChild(filterHeader);

        let filterTable = new FilterTable(request);
        this.filterTable = filterTable;
        filterTable.setRequestChangeCallback((request) => this.onRequestChange(request));
        div.appendChild(filterTable.getElement());

        return div;
    }
}

const FETCH_FEED_PATH = "/v1/ff/";

function createFilteredURL(request) {
    console.log(request);
    if (request.getFiltersList().length == 0) {
        return "-";
    }
    let bytes = request.serializeBinary();
    let b64 = toBase64(bytes);
    let host = window.location.host;
    if (window.location.hostname !== "localhost") {
        host = "api." + host;
    }
    return location.protocol + "//" + host + FETCH_FEED_PATH + b64;
}

function toBase64(dataArr) {
    return btoa(dataArr.reduce((data, val) => {
        return data + String.fromCharCode(val);
    }, ''));
}

function parseRequestFromURL(urlString) {
    let url = new URL(urlString);
    if (url.hostname.endsWith(window.location.hostname)) {
        let path = url.pathname.substring(FETCH_FEED_PATH.length);
        let bytes = atob(path);
        let request = FilterRequest.deserializeBinary(bytes);
        return request;
    }
    return null;
}

function getEnumValues(enumType) {
    return Object.keys(enumType).filter(key => !isNaN(Number(MotifIntervention[key])));
}

class URLCardHeader {
    element;
    filteredURLSpan;

    constructor(sourceUrl, filteredURL) {
        let row = document.createElement('div');
        row.classList.add('row');

        let sourceURLCol = document.createElement('div');
        sourceURLCol.classList.add('col');
        sourceURLCol.classList.add('s12');
        sourceURLCol.classList.add('m6');
        sourceURLCol.classList.add('l4');
        let sourceURLLabel = document.createElement('label');
        sourceURLLabel.innerText = "Source URL";
        sourceURLCol.appendChild(sourceURLLabel);
        let sourceURLSpan = document.createElement('span');
        sourceURLSpan.innerText = sourceUrl;
        sourceURLSpan.classList.add('code-text');
        sourceURLCol.appendChild(sourceURLSpan);
        row.appendChild(sourceURLCol);

        let filteredURLCol = document.createElement('div');
        filteredURLCol.classList.add('col');
        filteredURLCol.classList.add('s12');
        filteredURLCol.classList.add('m6');
        filteredURLCol.classList.add('l8');
        let filteredURLLabel = document.createElement('label');
        filteredURLLabel.innerText = "Filtered URL";
        filteredURLCol.appendChild(filteredURLLabel);
        let filteredURLSpan = document.createElement('span');
        this.filteredURLSpan = filteredURLSpan;
        filteredURLSpan.innerText = filteredURL;
        filteredURLSpan.classList.add('code-text');
        filteredURLCol.appendChild(filteredURLSpan);
        row.appendChild(filteredURLCol);

        let head = document.createElement('div');
        head.classList.add('card-content');
        head.appendChild(row);
        this.element = head;
    }

    setFilteredURL(url) {
        console.log(url);
        this.filteredURLSpan.innerText = url;
    }

    getElement() {
        return this.element;
    }
}

class FilterTable {
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
                <th>Effect</th>
                <th>Type</th>
                <th>Expression</th>
                <th></th>
            </tr>`;

        let tbody = document.createElement('tbody');
        this.tableBody = tbody;

        let footer = new FilterTableFooter();
        this.footer = footer;
        footer.setAddButtonCallback(() => this.callback_addFilterRow());
        tbody.appendChild(footer.getElement());

        table.appendChild(thead);
        table.appendChild(tbody);
    }

    callback_addFilterRow() {
        let filterRow = this.footer.createFilterRowFromInputs();
        if (filterRow == null) {
            return;
        }
        this.request.addFilters(filterRow.getSpec());
        this.requestChangeCallback(this.request);
        filterRow.setRemoveButtonCallback(() => this.callback_removeFilterRow(filterRow));
        this.footer.prependSibling(filterRow.getElement());
    }

    removeSpecFromRequest(spec) {
        let newList = [];
        let modified = false;
        for (let filter of this.request.getFiltersList()) {
            if (filter !== spec) {
                newList.push(filter);
            } else {
                modified = true;
            }
        }
        if (modified) {
            this.request.setFiltersList(newList);
            this.requestChangeCallback(this.request);
        }
    }

    callback_removeFilterRow(filterRow) {
        let spec = filterRow.getSpec();
        this.removeSpecFromRequest(spec);
        filterRow.remove();
    }



    getElement() {
        return this.table;
    }

    setRequestChangeCallback(callback) {
        this.requestChangeCallback = callback;
    }
}

class FilterRow {
    tableRow;
    removeButton;
    spec;

    constructor(target, effect, type, expression, spec) {
        this.spec = spec;

        let tr = document.createElement('tr');
        this.tableRow = tr;
        tr.classList.add('filter-row');

        let targetCol = document.createElement('td');
        targetCol.innerText = target;
        tr.appendChild(targetCol);

        let effectCol = document.createElement('td');
        effectCol.innerText = effect;
        tr.appendChild(effectCol);

        let typeCol = document.createElement('td');
        typeCol.innerText = type;
        tr.appendChild(typeCol);

        let expressionCol = document.createElement('td');
        expressionCol.innerText = expression;
        if (type == "REGEX") {
            expressionCol.classList.add('regex');
        }
        tr.appendChild(expressionCol);

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

class FilterTableFooter {
    tableRow;
    addButton;
    inputs;

    constructor() {
        this.inputs = {};

        let tr = document.createElement('tr');
        this.tableRow = tr;
        tr.classList.add('table-footer-row');

        let targetSelectCol = this.createSelectColumn('Target', FilterTarget);
        let effectSelectCol = this.createSelectColumn('Effect', FilterEffect);
        let typeSelectCol = this.createSelectColumn('Type', FilterType);
        let expressionInput = this.createInputColumn('Expression', 'text');

        tr.appendChild(targetSelectCol);
        tr.appendChild(effectSelectCol);
        tr.appendChild(typeSelectCol);
        tr.appendChild(expressionInput);

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

    createFilterTableFooter() {
        addButton.onclick = () => this.addFilter(this.filterInputs);
    }

    prependSibling(element) {
        this.tableRow.parentElement.insertBefore(element, this.tableRow);
    }

    validateInputs() {
        return this.inputs["target"].value != "" &&
            this.inputs["effect"].value != "" &&
            this.inputs["type"].value != "" &&
            this.inputs["expression"].value != "";
    }

    createFilterRowFromInputs() {
        if (!this.validateInputs()) {
            return null;
        }

        let target = this.inputs["target"].value;
        let effect = this.inputs["effect"].value;
        let type = this.inputs["type"].value;
        let expression = this.inputs["expression"].value;

        let spec = new FilterSpec();
        spec.setTarget(FilterTarget[target]);
        spec.setEffect(FilterEffect[effect]);
        spec.setType(FilterType[type]);
        spec.setExpression(expression);

        return new FilterRow(
            target,
            effect,
            type,
            expression,
            spec);
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