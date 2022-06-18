

const makeRequest = () => {
  const apiUrl = 'https://go-apod.herokuapp.com/apod';

  fetch(apiUrl)
  .then(response => response.json())
  .then(data => {
    console.log(data);
    updateDom(data)
  })
  .catch((error) => {
    console.error('Error:', error);
  });
}

const updateDom = (apod) => {
  const titleElem = document.getElementById('apod-title');
  const descriptionElem = document.getElementById('apod-explanation');
  const copyrightElem = document.getElementById('apod-copyright');
  const dateElem = document.getElementById('apod-date');
  const linkElem = document.getElementById('apod-hd-link');

  titleElem.innerText = apod.title;
  descriptionElem.innerText = apod.explanation;
  copyrightElem.innerText = apod.copyright;
  dateElem.innerText = apod.date;
  linkElem.innerText = 'View HD Image';
  linkElem.setAttribute('href', apod.hdurl);
}


// console.log('Hello');

makeRequest();
