import React from 'react';
import Tabs from './Tabs';
import { getRecipients, getSender } from '../../utils/formatter';
import { isHtmlMessage, getHtmlMessage, formatMessagePlain } from '../../utils/helpers';
import { withRouter } from 'react-router-dom';

const TAB_HTML = 1;
const TAB_PLAIN = 2;
const TAB_SOURCE = 3;

class MailDetail extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            message: null,
            showAllHeaders: false,
            tab: TAB_HTML,
        };
    }

    handleTabClick = tab => {
        this.setState((state, prev) => ({ ...state, tab }));
    };

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
        if (this.state.message === null) {
            return null;
        }

        const tabs = [
            { title: 'HTML', tab: TAB_HTML, current: TAB_HTML === this.state.tab, cond: isHtmlMessage(this.state.message) },
            { title: 'Plain', tab: TAB_PLAIN, current: TAB_PLAIN === this.state.tab, cond: !isHtmlMessage(this.state.message) },
            { title: 'Source', tab: TAB_SOURCE, current: TAB_SOURCE === this.state.tab, cond: true },
        ];

        return (
            <div>
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
                </table>

                <Tabs tabs={tabs} onTabClick={this.handleTabClick} />

                <div className="tabs-content">
                    {this.state.tab === TAB_HTML &&
                        <div className="tabs-panel is-active">
                            <iframe seamless srcDoc={`${getHtmlMessage(this.state.message)}`}
                                    frameBorder="0" style={{ width: '100%' }}/>
                        </div>}
                    {this.state.tab === TAB_PLAIN &&
                        <div className="tabs-panel is-active">
                            <p>{formatMessagePlain(this.state.message)}</p>
                        </div>}
                    {this.state.tab === TAB_SOURCE &&
                        <div className="tabs-panel is-active">
                            <pre>
                                {Object.keys(this.state.message.Content.Headers).map(header => {
                                    const value = this.state.message.Content.Headers[header];
                                    return (
                                        <div key={header}>{header}: {value}</div>
                                    );
                                })}
                                <p>{this.state.message.Content.Body}</p>
                            </pre>
                        </div>}
                </div>
            </div>
        );
    }
}

export default withRouter(MailDetail);
