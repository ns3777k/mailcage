import React from 'react';
import { TabView, TabPanel } from 'primereact/tabview';
import { Toolbar } from 'primereact/toolbar';
import { Button } from 'primereact/button';
import { Dropdown } from 'primereact/dropdown';
import { withRouter } from 'react-router-dom';
import { isHtmlMessage, getHtmlMessage, formatMessagePlain } from '../../utils/helpers';

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

    handleDelete = e => {
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

        return (
            <div>
                <Toolbar>
                    <div className="p-toolbar-group-left">
                        <Button onClick={this.handleBackClick}
                                label="Back"
                                icon="pi pi-arrow-left"
                                style={{ marginRight: '.25em' }} />

                        <Button onClick={this.handleDelete}
                                label="Remove"
                                icon="pi pi-trash"
                                className="p-button-danger" />
                    </div>
                    <div className="p-toolbar-group-right">
                        <Dropdown value={this.state.outgoingServer} options={this.state.outgoingServers} onChange={this.handleSelectOutgoingServer} />
                        <Button label="Save" icon="pi pi-check" className="p-button-warning" />
                    </div>
                </Toolbar>
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
//     toggleRelease = e => {
//         e.preventDefault();
//
//         this.setState((state, prevState) => ({
//             ...state,
//             showReleaseServersList: !state.showReleaseServersList,
//         }));
//     };
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
//     goBack = e => {
//         e.preventDefault();
//         this.props.history.goBack();
//     };
//
//     componentDidMount() {
//         const { id } = this.props.match.params;
//
//         getOutgoingServers()
//             .then(outgoingServers => {
//                 this.setState((state, prev) => ({
//                     ...state,
//                     outgoingServers,
//                     outgoingServer: outgoingServers[0] || '',
//                 }));
//             });
//     }
//
//     render() {
//         return (
//             <div>
//                 <table className="unstriped">
//                     <tbody>
//                     <tr>
//                         <td colSpan={2}>
//                             <div className="button-group">
//                                 <button onClick={this.goBack} type="button" className="button">Back</button>
//                                 <button onClick={this.toggleHeaders} type="button" className="button">
//                                     {this.state.showAllHeaders ? 'Hide' : 'Show'} all headers
//                                 </button>
//                                 <button onClick={this.handleDelete} type="button" className="alert button">
//                                     Remove
//                                 </button>
//                                 <button onClick={this.toggleRelease} type="button" className="warning button">
//                                     {this.state.showReleaseServersList ? 'Close' : 'Release'}
//                                 </button>
//                                 {this.state.showReleaseServersList && <>
//                                     <button onClick={this.handleReleaseClick} type="button" className="warning button">
//                                         Release
//                                     </button>
//                                     <select onChange={this.handleSelectOutgoingServer}>
//                                         {this.state.outgoingServers.map(outgoingServer => {
//                                             return (
//                                                 <option key={outgoingServer} value={outgoingServer}>
//                                                     {outgoingServer}
//                                                 </option>
//                                             );
//                                         })}
//                                     </select>
//                                     </>}
//                             </div>
//                         </td>
//                     </tr>
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
