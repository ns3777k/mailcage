import React from 'react';
import { getRecipients, getSender } from '../../utils/formatter';
import { withRouter } from 'react-router-dom';

class MailDetail extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            message: null,
            showAllHeaders: false,
        };
    }

    handleDelete = e => {
        e.preventDefault();

        if (this.props.onDeleteMessage) {
            this.props.onDeleteMessage(this.state.message.ID)
                .then(() => {
                    this.props.history.push(`/`);
                });
        }
    };

    toggleHeaders = e => {
        this.setState((state, prevState) => ({
            ...state,
            showAllHeaders: !state.showAllHeaders,
        }));
    };

    componentDidMount() {
        const { id } = this.props.match.params;

        this.props.onGetMessage(id)
            .then(message => {
                this.setState({ message });
            });
    }

    render() {
        return (
            <div>
                {this.state.message &&
                    <table className="unstriped">
                        <tbody>
                        <tr>
                            <td colSpan={2}>
                                <div className="button-group">
                                    <button onClick={this.handleDelete} type="button" className="alert button">
                                        Remove
                                    </button>
                                    <button onClick={this.toggleHeaders} type="button" className="button">
                                        {this.state.showAllHeaders ? 'Hide' : 'Show'} all headers
                                    </button>
                                </div>
                            </td>
                        </tr>
                        {!this.state.showAllHeaders && <>
                            <tr>
                                <td>From</td>
                                <td>{getSender(this.state.message)}</td>
                            </tr>
                            <tr>
                                <td>Subject</td>
                                <td>{this.state.message.Content.Headers['Subject'][0]}</td>
                            </tr>
                            <tr>
                                <td>To</td>
                                <td>{getRecipients(this.state.message).join(', ')}</td>
                            </tr>
                        </>}

                        {this.state.showAllHeaders &&
                            Object.keys(this.state.message.Content.Headers).map(headerName => {
                                return (
                                    <tr key={headerName}>
                                        <td>{headerName}</td>
                                        <td>
                                            {this.state.message.Content.Headers[headerName].join(', ')}
                                        </td>
                                    </tr>
                                );
                            })
                        }
                        </tbody>
                    </table>}
            </div>
        );
    }
}

export default withRouter(MailDetail);
