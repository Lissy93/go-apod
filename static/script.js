/**
 * Vanilla JS code for the homepage.
 * Fetches todays image and meta info from backend, and renders to UI
 */

/* API endpoint paths, using either current server or public instance */
const makeEndpointUrls = () => {
  const origin = window.location.origin;
  const hostname = origin && origin !== 'null' ? origin : 'https://go-apod.herokuapp.com';
  return {
    home: hostname,
    apod: `${hostname}/apod`,
    image: `${hostname}/image`,
  };
};

/* Fetch data from APOD API */
const makeRequest = () => {
  const apiUrl = makeEndpointUrls().apod;
  fetch(apiUrl)
  .then(response => response.json())
  .then(data => {
    updateDom(data);
  })
  .catch((error) => {
    showError(error);
  }).finally(() => {
    hideLoader();
  });
}

/* Hide loading spinner */
const hideLoader = () => {
  document.getElementById('loader').style.display = 'none';
};

/* Shows error message on UI */
const showError = (err) => {
  document.getElementById('error').style.display = 'block';
  document.getElementById('err-msg').innerText = err;
};

/* Converts timestamp into readable local date */
const formatDate = (date) => {
  if (!date) return '';
  return new Date().toLocaleDateString(
    "en-US",
    { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' },
  );
};

/* Using the response from APOD API, update the DOM to render results */
const updateDom = (apod) => {
  document.getElementsByClassName('apod-info')[0].style.display = 'block';
  const titleElem = document.getElementById('apod-title');
  const descriptionElem = document.getElementById('apod-explanation');
  const copyrightElem = document.getElementById('apod-copyright');
  const dateElem = document.getElementById('apod-date');
  const linkElem = document.getElementById('apod-hd-link');
  const iframeElem = document.getElementById('apod-dynamic-content');
  const imageElem = document.getElementById('apod-picture');

  titleElem.innerText = apod.title;
  descriptionElem.innerText = apod.explanation;
  copyrightElem.innerText = apod.copyright || '';
  dateElem.innerText = formatDate(apod.date);
  linkElem.innerText = 'View HD Image';
  linkElem.setAttribute('href', apod.hdurl || apod.url);

  if (apod.media_type !== 'image') {
    iframeElem.setAttribute('src', apod.url);
    iframeElem.style.display = 'block';
    imageElem.style.display = 'none';
    linkElem.innerText = 'View Dynamic Content';
  }

  document.getElementById('response').innerHTML = prettyPrint(apod);
}

/* Updates API docs with endpoint based on hostname */
const setApiEndPoints = () => {
  const { apod, image } = makeEndpointUrls();
  document.getElementById('get-apod').innerText = apod;
  document.getElementById('get-img').innerText = image;
}

/* Format API JSON response in nicely */
const prettyPrint = (json) => {
  if (typeof json != 'string') { json = JSON.stringify(json, undefined, 2); }
  json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
    let cls = 'number';
    if (/^"/.test(match)) {
      cls = (/:$/.test(match))? 'key' : 'string';
    } else if (/true|false/.test(match)) { cls = 'boolean'; }
    else if (/null/.test(match)) { cls = 'null'; }
    return '<span class="' + cls + '">' + match + '</span>';
  });
};

/* When page has loaded, make request then update the DOM  */
document.addEventListener('DOMContentLoaded', (e) => {
  makeRequest();
  setApiEndPoints();
});
