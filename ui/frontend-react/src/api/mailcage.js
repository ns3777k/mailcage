/**
 * Fetches all messages from start.
 *
 * @param {number} [start=0]
 * @returns {Promise<Promise<any>>}
 */
export function getMessages(start = 0) {
  return fetch(`/api/v1/messages?start=${start}`)
    .then(response => response.json())
    .then(response => ({
      Items: response.Items || [],
      Total: response.Total,
      Count: response.Count,
      Start: response.Start,
    }));
}

export function deleteAllMessages() {
  return fetch('/api/v1/messages', { method: 'DELETE' });
}

/**
 * @param {string} id
 */
export function deleteMessage(id) {
  return fetch(`/api/v1/message?id=${id}`, { method: 'DELETE' });
}

/**
 * @param {string} id
 */
export function getMessage(id) {
  return fetch(`/api/v1/message?id=${id}`)
    .then(response => {
      if (!response.ok) {
        throw response;
      }

      return response.json()
    });
}

export function getDownloadMessageLink(id) {
  return `/api/v1/download?id=${id}`;
}

export function getOutgoingServers() {
  return fetch('/api/v1/outgoing-servers')
    .then(response => response.json());
}

export function releaseMessage(id, server) {
  return fetch(`/api/v1/release?id=${id}&server=${server}`, { method: 'POST' });
}

export function markAsRead(id) {
  return fetch(`/api/v1/read?id=${id}`, { method: 'POST' });
}
