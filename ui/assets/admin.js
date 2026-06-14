document.addEventListener('DOMContentLoaded', () => {
  const tbody = document.getElementById('feedback-table-body');
  const refreshBtn = document.getElementById('refresh-btn');

  const loadFeedbacks = async () => {
    const response = await fetch('/api/feedbacks', {
      credentials: 'include'
    });

    if (!response.ok) {
      tbody.innerHTML = `<tr><td colspan="5" class="text-danger p-4">Unable to load feedbacks.</td></tr>`;
      return;
    }

    const feedbacks = await response.json();
    tbody.innerHTML = feedbacks.map((feedback) => `
      <tr>
        <td>${escapeHtml(feedback.title)}</td>
        <td>${escapeHtml(feedback.message)}</td>
        <td><span class="badge text-bg-secondary">${escapeHtml(feedback.status)}</span></td>
        <td>${new Date(feedback.created_at).toLocaleString()}</td>
        <td>
          <select class="form-select form-select-sm status-select" data-id="${feedback.id}">
            <option value="reviewing" ${feedback.status === 'reviewing' ? 'selected' : ''}>reviewing</option>
            <option value="resolved" ${feedback.status === 'resolved' ? 'selected' : ''}>resolved</option>
          </select>
        </td>
      </tr>
    `).join('');

    tbody.querySelectorAll('.status-select').forEach((select) => {
      select.addEventListener('change', async () => {
        await updateStatus(select.dataset.id, select.value);
        await loadFeedbacks();
      });
    });
  };

  const updateStatus = async (id, status) => {
    await fetch(`/api/feedbacks/${id}/status`, {
      method: 'PATCH',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ status }),
    });
  };

  const escapeHtml = (value) => value
      .replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll('"', '&quot;')
      .replaceAll("'", '&#39;');

  refreshBtn.addEventListener('click', loadFeedbacks);
  loadFeedbacks();
});