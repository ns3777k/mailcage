import React from 'react';
import { getRecipients, getSender } from '../../utils/formatter';
import parseISO from 'date-fns/parseISO';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import { withRouter } from 'react-router-dom';

class MailItem extends React.Component {
    handleDelete = e => {
        e.stopPropagation();

        if (this.props.onDeleteMessage) {
            this.props.onDeleteMessage(this.props.message.ID);
        }
    };

    handleClick = e => {
        e.preventDefault();

        this.props.history.push(`/${this.props.message.ID}`);
    };

    render() {
        const { message } = this.props;

        return (
            <tr className="table-expand-row pointer" onClick={this.handleClick}>
                <td>{getSender(message)}</td>
                <td>{getRecipients(message).join(', ')}</td>
                <td>{formatDistanceToNow(parseISO(message.CreatedAt))} ago</td>
                <td>{message.Content.Headers['Subject'][0]}</td>
                <td>
                    <button onClick={this.handleDelete} type="button" className="alert button">
                        Remove
                    </button>
                </td>
            </tr>
        );
    }
}

export default withRouter(MailItem);
