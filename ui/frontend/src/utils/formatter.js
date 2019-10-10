export function getSender(message) {
    if (Array.isArray(message.Content.Headers['From']) && message.Content.Headers['From'].length > 0) {
        return message.Content.Headers['From'][0];
    }

    return `${message.From.Mailbox}@${message.From.Domain}`;
}

export function getRecipients(message) {
    const inHeaders = Array.isArray(message.Content.Headers['To']) && message.Content.Headers['To'].length > 0;
    if (inHeaders) {
        return message.Content.Headers['To'];
    }

    return message.To.map(to => `${to.Mailbox}@${to.Domain}`);
}
