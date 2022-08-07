const { FetchFeedRequest, FilterType, FilterTarget, FilterEffect } = require('./api/v1_pb.js');
const { Materialize } = require('./materialize.min.js');
const { URLCard } = require('./urlCard.js');
const regexColorizer = require('regex-colorizer');

var feedURLInput = null;

const CREATE_URL_CARD_BUTTON = "create-url-card-button";

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
  const button = document.getElementById(CREATE_URL_CARD_BUTTON);
  if (button != null) {
    button.classList.add('disabled');
  }
}

function enableCreateButton() {
  const button = document.getElementById(CREATE_URL_CARD_BUTTON);
  if (button != null) {
    button.classList.remove('disabled');
  }
}

function validFeedURL() {
  feedURLInput.classList.add("valid");
  feedURLInput.classList.remove("invalid");
  enableCreateButton();
}


function displayURL(url) {
  const elem = document.querySelector("#url-container");

  const card = document.createElement('div');
  card.classList.add("card");
  const cardContent = document.createElement('div');
  cardContent.classList.add('card-content');

  const row = document.createElement('div');
  row.classList.add('row');
  const colA = document.createElement('div');
  colA.classList.add('col');
  colA.classList.add('s1');
  const colB = document.createElement('div');
  colB.classList.add('col');
  colB.classList.add('s11');

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
  colA.appendChild(copyBtn);

  const urlText = document.createElement('pre');
  urlText.append(url);
  colB.appendChild(urlText);

  row.appendChild(colA);
  row.appendChild(colB);
  cardContent.appendChild(row);
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
  document.querySelector('#' + CREATE_URL_CARD_BUTTON).onclick = function () {
    console.log("Create button clicked");
    let url = document.querySelector('#feed_url').value;
    let urlCard = new URLCard(url);
    container = document.querySelector("#url-container");
    container.innerHTML = "";
    container.appendChild(urlCard.element);
  };
  const elem = document.querySelector('#feed_url');
  elem.onblur = function () {
    checkFeedURL();
  };
  feedURLInput = elem;

  let elems = document.querySelectorAll('select');
  M.FormSelect.init(elems);

  M.Sidenav.init(document.querySelectorAll('.sidenav'));

  regexColorizer.addStyleSheet();
});