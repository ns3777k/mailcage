import React from 'react';
import { TabView, TabPanel } from 'primereact/tabview';
import { Toolbar } from 'primereact/toolbar';
import { Button } from 'primereact/button';
import { Dropdown } from 'primereact/dropdown';
import { withRouter } from 'react-router-dom';
import { isHtmlMessage, getHtmlMessage, formatMessagePlain } from '../../utils/helpers';
import { getRecipients, getSender } from '../../utils/formatter';

const TAB_HTML = 0;
const TAB_PLAIN = 1;
const TAB_SOURCE = 2;

class MailDetail extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            message: null,
            activeTabIndex: TAB_SOURCE,
            outgoingServers: [],
            outgoingServer: '',
            releaseMode: false,
        };
    }

    componentDidMount() {
        const { id } = this.props.match.params;

        this.props.onGetOutgoingServers()
            .then(outgoingServers => {
                this.setState((state, prev) => ({
                    ...state,
                    outgoingServers: outgoingServers.map(s => ({ label: s, value: s })),
                    outgoingServer: outgoingServers[0] || '',
                }));
            });

        this.props.onGetMessage(id)
            .then(message => {
                const activeTabIndex = isHtmlMessage(message) ? TAB_HTML : TAB_PLAIN;
                this.setState((state, prev) => ({
                    ...state,
                    message,
                    activeTabIndex
                }));
            });
    }

    handleReleaseClick = e => {
        this.props.onReleaseMessage(this.state.outgoingServer, this.state.message.ID);
    };

    handleDeleteClick = e => {
        e.preventDefault();

        this.props.onDeleteMessage(this.state.message.ID)
            .then(() => {
                this.props.history.push(`/`);
            });
    };

    handleBackClick = e => {
        e.preventDefault();

        this.props.history.goBack();
    };

    handleTabChange = e => {
        this.setState({ activeTabIndex: e.index });
    };

    toggleOutgoingServers = e => {
        this.setState((state, prev) => ({
            ...state,
            releaseMode: !state.releaseMode,
        }));
    };

    handleSelectOutgoingServer = e => {
        this.setState((state, prev) => ({
            ...state,
            outgoingServer: e.value,
        }));
    };

    render() {
        if (this.state.message === null) {
            return null;
        }

        const columns = [
            { header: 'From', value: getSender(this.state.message) },
            { header: 'Subject', value: this.state.message.Content.Headers['Subject'][0] },
            { header: 'To', value: getRecipients(this.state.message).join(', ') },
        ];

        return (
            <div>
                <Toolbar>
                    <div className="p-toolbar-group-left">
                        <Button onClick={this.handleBackClick}
                                label="Back"
                                icon="pi pi-arrow-left"
                                style={{ marginRight: '.25em' }} />

                        <Button onClick={this.handleDeleteClick}
                                label="Remove"
                                icon="pi pi-trash"
                                className="p-button-danger" />
                    </div>
                    <div className="p-toolbar-group-right">
                        {this.state.releaseMode &&
                            <>
                                <Dropdown value={this.state.outgoingServer}
                                          options={this.state.outgoingServers}
                                          style={{ marginRight: '.25em' }}
                                          onChange={this.handleSelectOutgoingServer} />
                                <Button label="Release!"
                                        icon="pi pi-external-link"
                                        onClick={this.handleReleaseClick}
                                        style={{ marginRight: '.25em' }}
                                        className="p-button-success" />
                            </>}
                        <Button onClick={this.toggleOutgoingServers}
                                label={this.state.releaseMode ? 'Close' : 'Release'}
                                icon={`pi ${this.state.releaseMode ? 'pi-times' : 'pi-external-link'}`}
                                className="p-button-warning" />
                    </div>
                </Toolbar>
                <TabView>
                    <TabPanel header="Brief headers">
                        {columns.map(column => {
                            return (
                                <p key={column.header}><strong>{column.header}:&nbsp;</strong>{column.value}</p>
                            );
                        })}
                    </TabPanel>
                    <TabPanel header="All headers">
                        {Object.keys(this.state.message.Content.Headers).map(headerName => {
                            return (
                                <p key={headerName}>
                                    <strong>{headerName}:&nbsp;</strong>
                                    {this.state.message.Content.Headers[headerName].join(', ')}
                                </p>
                            );
                        })}
                    </TabPanel>
                </TabView>
                <TabView activeIndex={this.state.activeTabIndex} onTabChange={this.handleTabChange}>
                    <TabPanel disabled={!isHtmlMessage(this.state.message)} header="HTML">
                        <iframe seamless
                                srcDoc={`${getHtmlMessage(this.state.message)}`}
                                title="Message preview sandbox"
                                frameBorder="0"
                                style={{ width: '100%' }}/>
                    </TabPanel>

                    <TabPanel disabled={isHtmlMessage(this.state.message)} header="Plain">
                        {formatMessagePlain(this.state.message)}
                    </TabPanel>

                    <TabPanel header="Source">
                        <pre>
                            {Object.keys(this.state.message.Content.Headers).map(header => {
                                const value = this.state.message.Content.Headers[header];
                                return (
                                    <div key={header}>{header}: {value}</div>
                                );
                            })}
                            <p>{this.state.message.Content.Body}</p>
                        </pre>
                    </TabPanel>

                    <TabPanel disabled={!this.state.message.MIME} header="MIME">
                        {((this.state.message.MIME || {}).Parts || []).map((part, index) => {
                            return (
                                <div key={index}>
                                    <a href={`/api/v1/download-part?id=${this.state.message.ID}&part=${index}`}>
                                        Download {part.Headers['Content-Type'] || 'Unknown type'} ({part.Size}) bytes
                                    </a>
                                </div>
                            );
                        })}
                    </TabPanel>
                </TabView>
            </div>
        );
    }
}

export default withRouter(MailDetail);

// import { getRecipients, getSender } from '../../utils/formatter';
// import { getOutgoingServers, release } from '../../api/mailcage';
//
// class MailDetail extends React.Component {
//     constructor(props) {
//         super(props);
//
//         this.state = {
//             showAllHeaders: false,
//             showReleaseServersList: false,
//         };
//     }
//
//     handleReleaseClick = e => {
//         e.preventDefault();
//
//         release(this.state.outgoingServer, this.state.message.ID);
//     };
//
//     toggleHeaders = e => {
//         this.setState((state, prevState) => ({
//             ...state,
//             showAllHeaders: !state.showAllHeaders,
//         }));
//     };
//
//     render() {
//         return (
//             <div>
//                 <table className="unstriped">
//                     <tbody>
//                     {!this.state.showAllHeaders && <>
//                         <tr>
//                             <td>From</td>
//                             <td>{getSender(this.state.message)}</td>
//                         </tr>
//                         <tr>
//                             <td>Subject</td>
//                             <td>{this.state.message.Content.Headers['Subject'][0]}</td>
//                         </tr>
//                         <tr>
//                             <td>To</td>
//                             <td>{getRecipients(this.state.message).join(', ')}</td>
//                         </tr>
//                     </>}
//                     {this.state.showAllHeaders &&
//                         Object.keys(this.state.message.Content.Headers).map(headerName => {
//                             return (
//                                 <tr key={headerName}>
//                                     <td>{headerName}</td>
//                                     <td>
//                                         {this.state.message.Content.Headers[headerName].join(', ')}
//                                     </td>
//                                 </tr>
//                             );
//                         })
//                     }
//                     </tbody>
//                 </table>
//             </div>
//         );
//     }
// }
//
// export default withRouter(MailDetail);
