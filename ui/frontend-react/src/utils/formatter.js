export function getSender(headers, from) {
    if (Array.isArray(headers['From']) && headers['From'].length > 0) {
        return headers['From'][0];
    }

    return `${from.Mailbox}@${from.Domain}`;
}

export function getRecipients(headers, to) {
    const inHeaders = Array.isArray(headers['To']) && headers['To'].length > 0;
    if (inHeaders) {
        return headers['To'];
    }

    return to.map(to => `${to.Mailbox}@${to.Domain}`);
}
