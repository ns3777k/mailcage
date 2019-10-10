import React from 'react';
import parseISO from 'date-fns/parseISO';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';

function getSender(message) {
    if (Array.isArray(message.Content.Headers['From']) && message.Content.Headers['From'].length > 0) {
        return message.Content.Headers['From'][0];
    }

    return `${message.From.Mailbox}@${message.From.Domain}`;
}

function getRecipients(message) {
    const inHeaders = Array.isArray(message.Content.Headers['To']) && message.Content.Headers['To'].length > 0;
    if (inHeaders) {
        return message.Content.Headers['To'];
    }

    return message.To.map(to => {
        return `${to.Mailbox}@${to.Domain}`;
    });
}

class MailList extends React.Component {
    render() {
        return (
            <table className="table-expand">
                <thead>
                <tr className="table-expand-row">
                    <th width="300">From</th>
                    <th width="300">To</th>
                    <th width="400">When</th>
                    <th>Subject</th>
                    <th>&nbsp;</th>
                </tr>
                </thead>
                <tbody>
                {this.props.messages.map(message => {
                    return (
                        <tr className="table-expand-row" key={message.ID}>
                            <td>{getSender(message)}</td>
                            <td>{getRecipients(message).join(', ')}</td>
                            <td>{formatDistanceToNow(parseISO(message.CreatedAt))} ago</td>
                            <td>{message.Content.Headers["Subject"][0]}</td>
                            <td>
                                <div>Remove</div>
                            </td>
                        </tr>
                    );
                })}
                </tbody>
            </table>
        );
    }
}

export default MailList;
