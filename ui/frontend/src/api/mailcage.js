export function getMessages() {
    return fetch('/api/v1/messages')
        .then(response => response.json());
}

export function deleteMessage(id) {
    return fetch(`/api/v1/message?id=${id}`, { method: 'DELETE' });
}

export function createWebSocket() {
    return new WebSocket(`ws://${window.location.host}/api/v1/ws`);
}

export function getMessage(id) {
    return fetch(`/api/v1/message?id=${id}`)
        .then(response => response.json());
}
