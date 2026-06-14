document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('feedback-form');
  const toastEl = document.getElementById('feedback-toast');
  const toast = new bootstrap.Toast(toastEl, { delay: 3500 });
  const toastTitle = document.getElementById('toast-title');
  const toastBody = document.getElementById('toast-body');

  const showToast = (title, body, isError = false) => {
    toastTitle.textContent = title;
    toastBody.textContent = body;
    toastEl.classList.toggle('text-bg-danger', isError);
    toastEl.classList.toggle('text-bg-success', !isError);
    toast.show();
  };

  form.addEventListener('submit', async (event) => {
    event.preventDefault();
    const title = document.getElementById('title').value.trim();
    const message = document.getElementById('message').value.trim();

    if (!title || !message) {
      showToast('Validation error', 'Title and message are required.', true);
      return;
    }

    const response = await fetch('/api/feedbacks', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, message }),
    });

    if (response.ok) {
      form.reset();
      showToast('Success', 'Your feedback has been submitted.');
      return;
    }

    const payload = await response.json().catch(() => ({}));
    showToast('Error', payload.error || 'Unable to submit feedback.', true);
  });
});
