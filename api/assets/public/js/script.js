Notiflix.Loading.hourglass();

document.addEventListener('DOMContentLoaded', () => Notiflix.Loading.remove());

function askLLMAboutFunction(e) {
  const messageBox = document.getElementById('messageBox');
  if (!messageBox) return;

  messageBox.value = e.currentTarget.getAttribute('data-action');
  triggerFormSubmit();
}

function askLLMAboutSuggestion(e) {
  const messageBox = document.getElementById('messageBox');
  if (!messageBox) return;

  messageBox.value = e.currentTarget.querySelector('p').innerText;
}

const submitBtn = document.getElementById('submitBtn');
const promptForm = document.getElementById('promptForm');

function triggerFormSubmit() {
  const form = document.getElementById('promptForm');
  const event = new Event('submit', {
    bubbles: true,
    cancelable: true,
  });

  form.dispatchEvent(event);
}

function copyToClipboard(e) {
  const messageId = e.getAttribute('data-message-id');

  const content = document
    .getElementById(messageId)
    .getAttribute('data-message');

  navigator.clipboard.writeText(content).then(
    function () {
      Notiflix.Notify.success('Content copied to clipboard!');
    },
    function (err) {
      console.error('Could not copy text: ', err);
    }
  );
}

const cancelFormRequest = () => {
  htmx.trigger('#promptForm', 'htmx:abort');
};

document
  .querySelector('main')
  .addEventListener('htmx:responseError', function (evt) {
    console.error(evt);
    Notiflix.Notify.failure(
      'An error occurred. Please check your logs for the error message.'
    );
  });
