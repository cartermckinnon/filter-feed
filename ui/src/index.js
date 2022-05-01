const { FetchFeedRequest } = require('./api/v1_pb.js');
const { Materialize } = require('./materialize.min.js');

var feedURLInput = null;

function checkFeedURL() {
  const value = feedURLInput.value;
  if (value == "") {
    resetFeedURLValidity();
  } else {
    let url;
    try {
      url = new URL(value);
    } catch (_) {
      url = null
    }
    if (url != null) {
      validFeedURL();
    } else {
      invalidFeedURL();
    }
  }
}

function resetFeedURLValidity() {
  feedURLInput.classList.remove("invalid");
  feedURLInput.classList.remove("valid");
  disableCreateButton();
}

function invalidFeedURL() {
  feedURLInput.classList.add("invalid");
  feedURLInput.classList.remove("valid");
  disableCreateButton();
}

function disableCreateButton() {
  const button = document.getElementById("create-url-button");
  if (button != null) {
    button.classList.add('disabled');
  }
}

function enableCreateButton() {
  const button = document.getElementById("create-url-button");
  if (button != null) {
    button.classList.remove('disabled');
  }
}

function validFeedURL() {
  feedURLInput.classList.add("valid");
  feedURLInput.classList.remove("invalid");
  enableCreateButton();
}

function createURL() {
  b64 = getRequestAsBase64();
  host = window.location.host;
  if (window.location.hostname !== "localhost") {
    host = "api." + host;
  }
  return location.protocol + "://" + host + "/v1/f/" + b64;
}

function getRequestAsBase64() {
  let url = feedURLInput.value;
  let request = new FetchFeedRequest();
  request.setFeedurl(url);
  bytes = request.serializeBinary();
  return toBase64(bytes);
}

function toBase64(dataArr) {
  return btoa(dataArr.reduce((data, val) => {
    return data + String.fromCharCode(val);
  }, ''));
}

function displayURL(url) {
  const elem = document.querySelector("#url-container");

  const card = document.createElement('div');
  card.classList.add("card");
  const cardContent = document.createElement('div');
  cardContent.classList.add('card-content');

  const copyBtn = document.createElement('a');
  copyBtn.id = "url-copy-icon";
  copyBtn.classList.add('btn-floating');
  copyBtn.classList.add('btn-small');
  copyBtn.classList.add('blue');
  copyBtn.classList.add('pulse');
  copyBtn.classList.add('icon-button');
  copyBtn.onclick = function () {
    copyBtn.classList.remove('pulse');
    copyBtn.classList.add('grey');
    copyBtn.classList.remove('blue');
    copyTextToClipboard(url);
    M.toast({
      html: 'URL copied to clipboard!'
    })
  };
  const copyIcon = document.createElement('i');
  copyIcon.classList.add('material-icons');
  copyIcon.innerText = "content_copy";
  copyBtn.appendChild(copyIcon);
  cardContent.appendChild(copyBtn);

  const urlText = document.createElement('pre');
  urlText.append(copyBtn);
  urlText.append(url);
  cardContent.appendChild(urlText);

  card.appendChild(cardContent);

  elem.appendChild(card);
}

function fallbackCopyTextToClipboard(text) {
  var textArea = document.createElement("textarea");
  textArea.value = text;

  // Avoid scrolling to bottom
  textArea.style.top = "0";
  textArea.style.left = "0";
  textArea.style.position = "fixed";

  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();

  try {
    var successful = document.execCommand('copy');
    var msg = successful ? 'successful' : 'unsuccessful';
    console.log('Fallback: Copying text command was ' + msg);
  } catch (err) {
    console.error('Fallback: Oops, unable to copy', err);
  }

  document.body.removeChild(textArea);
}

function copyTextToClipboard(text) {
  if (!navigator.clipboard) {
    fallbackCopyTextToClipboard(text);
    return;
  }
  navigator.clipboard.writeText(text).then(function () {
    console.log('Async: Copying to clipboard was successful!');
  }, function (err) {
    console.error('Async: Could not copy text: ', err);
  });
}

document.addEventListener('DOMContentLoaded', function () {
  document.querySelector('#create-url-button').onclick = function () {
    console.log("Create button clicked");
    let url = createURL();
    displayURL(url);
  };
  const elem = document.querySelector('#feed_url');
  elem.onblur = function () {
    checkFeedURL();
  };
  feedURLInput = elem;
});