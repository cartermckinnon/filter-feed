const { FetchFeedRequest } = require('./api/v1_pb.js');
const { Materialize } = require('./materialize.min.js');
const { FilterTable } = require('./filterTable.js');
const { OverrideTable } = require('./overrideTable.js');

export class URLCard {
    element;
    request;
    header;
    filterTable;
    OverrideTable;

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

        let titleRow = document.createElement('div');
        titleRow.classList.add('card-content');
        titleRow.classList.add('card-row');
        titleRow.style = 'margin-bottom: 0;';
        let title = document.createElement('h4');
        title.style = 'margin: 0';
        title.innerText = "Manage your feed";
        titleRow.appendChild(title);
        div.appendChild(titleRow);

        let header = new URLCardHeader(request.getFeedurl(), "-");
        this.header = header;
        header.setFilteredURL(createFilteredURL(request));
        div.appendChild(header.getElement());

        let filterTable = new FilterTable(request);
        this.filterTable = filterTable;
        filterTable.setRequestChangeCallback((request) => this.onRequestChange(request));
        div.appendChild(this.createSection('Filters', 'Decide which items appear in the feed.', filterTable.getElement()));

        let overrideTable = new OverrideTable(request);
        this.overrideTable = overrideTable;
        overrideTable.setRequestChangeCallback((request) => this.onRequestChange(request));
        div.appendChild(this.createSection('Overrides', 'Change how the feed appears in your reader.', overrideTable.getElement()));        

        // this is a ridiculous hack to get the <select>-s rendered properly by Materialize
        setTimeout(() => {
            M.FormSelect.init(document.querySelectorAll('select'));
        }, 100);
        // this is a ridiculous hack to get the tooltip-s rendered properly by Materialize
        setTimeout(() => {
            M.Tooltip.init(document.querySelectorAll('.tooltipped'));
        }, 100);        
    }

    onRequestChange(request) {
        this.request = request;
        this.header.setFilteredURL(createFilteredURL(request));
    }

    createSection(headerText, tooltipText, element) {
        let div = document.createElement('div');
        div.classList.add('card-content');
        div.classList.add('card-row');
        let header = document.createElement('h5');
        header.style = 'margin-top: 0';
        header.innerText = headerText;
        let tooltip = document.createElement('a');
        tooltip.classList.add('btn-flat');
        tooltip.classList.add('tooltipped');
        tooltip.setAttribute('data-tooltip', tooltipText);
        tooltip.setAttribute('data-position', 'top');
        let tooltipIcon = document.createElement('i');
        tooltipIcon.classList.add('material-icons');
        tooltipIcon.innerText = 'info';
        tooltip.appendChild(tooltipIcon);
        header.appendChild(tooltip);
        div.appendChild(header);
        div.appendChild(element);
        return div;
    }
}

const FETCH_FEED_PATH = "/v1/ff/";

function createFilteredURL(request) {
    if (request.getFiltersList().length == 0 && request.getOverridesList().length == 0) {
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

function base64StringToUint8Array(base64String) {
    let binaryString = atob(base64String);
    let len = binaryString.length;
    let bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes.buffer;
}

function parseRequestFromURL(urlString) {
    let url = new URL(urlString);
    if (url.hostname.endsWith(window.location.hostname)) {
        let b64 = url.pathname.substring(FETCH_FEED_PATH.length);
        let bytes = base64StringToUint8Array(b64);
        let request = FetchFeedRequest.deserializeBinary(bytes);
        return request;
    }
    return null;
}

class URLCardHeader {
    element;
    filteredURLSpan;

    constructor(sourceUrl, filteredURL) {
        let row = document.createElement('div');
        row.classList.add('row');
        row.style = 'margin-bottom: 0;';

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
        head.classList.add('card-row');
        head.appendChild(row);
        this.element = head;
    }

    setFilteredURL(url) {
        this.filteredURLSpan.innerText = url;
    }

    getElement() {
        return this.element;
    }
}