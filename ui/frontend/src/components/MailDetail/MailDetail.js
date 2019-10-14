import React from 'react';
import Tabs from './Tabs';
import { getRecipients, getSender } from '../../utils/formatter';
import { isHtmlMessage, getHtmlMessage, formatMessagePlain } from '../../utils/helpers';
import { withRouter } from 'react-router-dom';
import { getOutgoingServers, release } from '../../api/mailcage';

const TAB_HTML = 1;
const TAB_PLAIN = 2;
const TAB_SOURCE = 3;
const TAB_MIME = 4;

class MailDetail extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            message: null,
            showAllHeaders: false,
            tab: false,
            showReleaseServersList: false,
            outgoingServers: [],
            outgoingServer: '',
        };
    }

    handleTabClick = tab => {
        this.setState((state, prev) => ({ ...state, tab }));
    };

    toggleRelease = e => {
        e.preventDefault();

        this.setState((state, prevState) => ({
            ...state,
            showReleaseServersList: !state.showReleaseServersList,
        }));
    };

    handleReleaseClick = e => {
        e.preventDefault();

        release(this.state.outgoingServer, this.state.message.ID);
    };

    handleSelectOutgoingServer = e => {
        const { value } = e.target;

        this.setState((state, prev) => ({
            ...state,
            outgoingServer: value,
        }));
    };

    handleDelete = e => {
        e.preventDefault();

        this.props.onDeleteMessage(this.state.message.ID)
            .then(() => {
                this.props.history.push(`/`);
            });
    };

    toggleHeaders = e => {
        this.setState((state, prevState) => ({
            ...state,
            showAllHeaders: !state.showAllHeaders,
        }));
    };

    goBack = e => {
        e.preventDefault();
        this.props.history.goBack();
    };

    componentDidMount() {
        const { id } = this.props.match.params;

        getOutgoingServers()
            .then(outgoingServers => {
                this.setState((state, prev) => ({
                    ...state,
                    outgoingServers,
                    outgoingServer: outgoingServers[0] || '',
                }));
            });

        this.props.onGetMessage(id)
            .then(message => {
                let currentTab = TAB_SOURCE;
                switch (true) {
                    case isHtmlMessage(message):
                        currentTab = TAB_HTML;
                        break;

                    case !isHtmlMessage(message):
                        currentTab = TAB_PLAIN;
                        break;
                }

                this.setState((state, prev) => ({
                    ...state,
                    message,
                    tab: currentTab
                }));
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
            { title: 'MIME', tab: TAB_MIME, current: TAB_MIME === this.state.tab, cond: !!this.state.message.MIME },
        ];

        return (
            <div>
                <table className="unstriped">
                    <tbody>
                    <tr>
                        <td colSpan={2}>
                            <div className="button-group">
                                <button onClick={this.goBack} type="button" className="button">Back</button>
                                <button onClick={this.toggleHeaders} type="button" className="button">
                                    {this.state.showAllHeaders ? 'Hide' : 'Show'} all headers
                                </button>
                                <button onClick={this.handleDelete} type="button" className="alert button">
                                    Remove
                                </button>
                                <button onClick={this.toggleRelease} type="button" className="warning button">
                                    {this.state.showReleaseServersList ? 'Close' : 'Release'}
                                </button>
                                {this.state.showReleaseServersList && <>
                                    <button onClick={this.handleReleaseClick} type="button" className="warning button">
                                        Release
                                    </button>
                                    <select onChange={this.handleSelectOutgoingServer}>
                                        {this.state.outgoingServers.map(outgoingServer => {
                                            return (
                                                <option key={outgoingServer} value={outgoingServer}>
                                                    {outgoingServer}
                                                </option>
                                            );
                                        })}
                                    </select>
                                    </>}
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
                    {this.state.tab === TAB_MIME &&
                        <div className="tabs-panel is-active">
                            {(this.state.message.MIME.Parts || []).map((part, index) => {
                                return (
                                    <div key={index}>
                                        <a href={`/api/v1/download-part?id=${this.state.message.ID}&part=${index}`}>
                                            Download {part.Headers['Content-Type'] || 'Unknown type'} ({part.Size}) bytes
                                        </a>
                                    </div>
                                );
                            })}
                        </div>}
                </div>
            </div>
        );
    }
}

export default withRouter(MailDetail);
