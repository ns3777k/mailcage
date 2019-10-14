export function getMessages(start = 0) {
    return fetch(`/api/v1/messages?start=${start}`)
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

export function getOutgoingServers() {
    return fetch('/api/v1/outgoing-servers')
        .then(response => response.json());
}

export function release(server, messageId) {
    return fetch(`/api/v1/release?server=${server}&id=${messageId}`, { method: 'POST' })
        .then(response => response.json());
}
