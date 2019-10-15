// Multiple decoding functions copied from https://github.com/mailhog/MailHog-UI

import {
    unescapeFromBase64,
    convertUnicodeCodePointsToString,
    decodeQuotedPrintableHelper,
    convertBytesToUnicodeCodePoints
} from 'strutil';

// "=E3=81=82=E3=81=84" => [ 0xE3, 0x81, 0x82, 0xE3, 0x81, 0x84 ]
function decodeQuotedPrintableWithoutRFC2047(str) {
    return decodeQuotedPrintableHelper(str, "=");
}

function unescapeFromQuotedPrintableWithoutRFC2047(str, encoding) {
    const decodedBytes = decodeQuotedPrintableWithoutRFC2047(str);
    const unicodeBytes = convertBytesToUnicodeCodePoints(decodedBytes, encoding);
    return convertUnicodeCodePointsToString(unicodeBytes);
}

function hasMatchingHeader(message, header, value) {
    header = header.toLowerCase();
    const headersKeys = Object.keys(message.Content.Headers);

    return headersKeys.some(hk => {
        if (header !== message.Content.Headers[hk][0].toLowerCase()) {
            return false;
        }

        return message.Content.Headers[hk][0].match(value);
    });
}

function findMatchingMIME(message, mime) {
    if (!message.MIME || !message.MIME.Parts) {
        return null;
    }

    for (let i = 0; i < message.MIME.Parts.length; i++) {
        const part = message.MIME.Parts[i];
        const contentType = part.Headers['Content-Type'] || [];
        if (contentType.length === 0) {
            continue;
        }

        if (contentType[0].match(`${mime};?.*`)) {
            return part;
        }

        if (contentType[0].match(/multipart\/.*/)) {
            return findMatchingMIME(part, mime);
        }
    }

    return null;
}

function tryDecodeContent(message) {
    const charset = 'UTF-8';
    let content = message.Content.Body;
    const contentTransferEncoding = message.Content.Headers['Content-Transfer-Encoding'][0];

    if (contentTransferEncoding) {
        switch (contentTransferEncoding.toLowerCase()) {
            case 'quoted-printable':
                content = content.replace(/=[\r\n]+/gm,"");
                content = unescapeFromQuotedPrintableWithoutRFC2047(content, charset);
                break;
            case 'base64':
                // remove line endings to give original base64-encoded string
                content = content.replace(/\r?\n|\r/gm,"");
                content = unescapeFromBase64(content, charset);
                break;
            default:
        }
    }

    return content;
}

function tryDecode(l) {
    if (l.Headers && l.Headers['Content-Type'] && l.Headers['Content-Transfer-Encoding']) {
        return tryDecodeContent({ Content: l });
    }

    return l.Body;
}

function escapeHtml(html) {
    const entityMap = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#39;'
    };

    return html.replace(/[&<>"']/g, s => entityMap[s]);
}

function getPlainMessage(message) {
    if (message.Content.Headers && message.Content.Headers['Content-Type'] && message.Content.Headers['Content-Type'][0].match('text/plain')) {
        return tryDecode(message.Content);
    }

    const m = findMatchingMIME(message, 'text/plain');
    if (m !== null) {
        return tryDecode(m);
    }

    return message.Content.Body;
}

export function formatMessagePlain(message) {
    const body = getPlainMessage(message);
    const escaped = escapeHtml(body);
    return escaped.replace(/(https?:\/\/)([-[\]A-Za-z0-9._~:/?#@!$()*+,;=%]|&amp;|&#39;)+/g, '<a href="$&" target="_blank">$&</a>');
}

export function isHtmlMessage(message) {
    if (hasMatchingHeader(message, 'content-type', 'text/html')) {
        return true;
    }

    return findMatchingMIME(message, 'text/html') !== null;
}

export function getHtmlMessage(message) {
    if (hasMatchingHeader(message, 'content-type', 'text/html')) {
        return tryDecode(message.Content);
    }

    const m = findMatchingMIME(message, 'text/html');
    if (m !== null) {
        return tryDecode(m);
    }

    return 'Cannot render HTML';
}
