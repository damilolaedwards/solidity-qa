Notiflix.Loading.dots();

document.addEventListener('DOMContentLoaded', () => Notiflix.Loading.remove());

// Handle htmx error responses
document.querySelector('main').addEventListener('htmx:responseError', (evt) => {
  console.error(evt);
  Notiflix.Notify.failure(
    'An error occurred. Please check your logs for the error message.'
  );
});

// Handle prompt form loading states
document
  .getElementById('promptForm')
  .addEventListener('htmx:beforeRequest', (_evt) =>
    Notiflix.Loading.hourglass()
  );
document
  .getElementById('promptForm')
  .addEventListener('htmx:afterRequest', (_evt) => Notiflix.Loading.remove());

const askLLMAboutFunction = (e) => {
  const messageBox = document.getElementById('messageBox');
  if (!messageBox) return;

  messageBox.value = e.currentTarget.getAttribute('data-action');
  triggerFormSubmit();
};

const askLLMAboutSuggestion = (e) => {
  const messageBox = document.getElementById('messageBox');
  if (!messageBox) return;

  messageBox.value = e.currentTarget.querySelector('p').innerText;
};

const copyToClipboard = (e) => {
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
};
