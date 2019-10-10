import React from 'react';
import MailItem from './MailItem';

class MailList extends React.Component {
    componentDidMount() {
        if (this.props.onGetMessages) {
            this.props.onGetMessages();
        }
    }

    render() {
        return (
            <table className="table-expand hover">
                <thead>
                <tr className="table-expand-row">
                    <th>From</th>
                    <th>To</th>
                    <th>When</th>
                    <th>Subject</th>
                    <th>&nbsp;</th>
                </tr>
                </thead>
                <tbody>
                {this.props.messages.map(message => {
                    return (
                        <MailItem onDeleteMessage={this.props.onDeleteMessage}
                                  key={message.ID}
                                  message={message} />
                    );
                })}
                </tbody>
            </table>
        );
    }
}

export default MailList;
