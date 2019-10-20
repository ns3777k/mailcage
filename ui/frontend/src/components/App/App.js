import React from 'react';
import 'primereact/resources/themes/nova-light/theme.css';
import 'primereact/resources/primereact.min.css';
import 'primeflex/primeflex.css';
import 'primeicons/primeicons.css';
import { Switch, Route } from 'react-router-dom';
import MailList from '../MailList/MailList';
import MailDetail from '../MailDetail/MailDetail';
import * as api from '../../api/mailcage';
import { withRouter } from 'react-router-dom';
import { getSender } from '../../utils/formatter';

class App extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            messages: [],
            total: 0,
            count: 0,
            start: 0,
        };

        this.ws = api.createWebSocket();
        this.ws.addEventListener('message', this.onMessage);

        if (typeof Notification !== 'undefined') {
            Notification.requestPermission();
        }
    }

    createNotification(message) {
        if (typeof Notification === 'undefined') {
            return;
        }

        const title = `Mail from ${getSender(message)}`;
        const options = {
            body: message.Content.Headers['Subject'][0],
            tag: 'MailCage'
        };

        const notification = new Notification(title, options);
        notification.addEventListener('click', e => {
            this.props.history.push(`/${message.ID}`);
            notification.close();
        });
    }

    onMessage = e => {
        const message = JSON.parse(e.data);
        switch (message.Type) {
            case 'stored':
                this.handleGetMessages(this.state.start);
                this.createNotification(message.Payload);
                break;

            case 'deleted_message':
                const id = message.Payload;
                this.setState((state, prevState) => ({
                    ...state,
                    messages: state.messages.filter(msg => msg.ID !== id),
                    total: state.total - 1,
                }));
                break;

            case 'deleted_messages':
                this.setState((state, prevState) => ({
                    ...state, messages: [],
                    total: 0,
                }));
                break;

            default:
                console.error('unknown message: ', e.data);
        }
    };

    handleDeleteMessage = id => {
        return api.deleteMessage(id);
    };

    handleGetMessages = start => {
        api.getMessages(start)
            .then(response => {
                this.setState((state, prevState) => ({
                    ...state,
                    messages: response.Items,
                    total: response.Total,
                    count: response.Count,
                    start: response.Start,
                }));
            });
    };

    handleGetMessage = id => {
        return api.getMessage(id);
    };

    render() {
        return (
            <div className="p-grid">
                <div className="p-col-8 p-offset-2">
                    <Switch>
                        <Route path="/:id">
                            <MailDetail onGetMessage={this.handleGetMessage}
                                        onGetOutgoingServers={api.getOutgoingServers}
                                        onDeleteMessage={this.handleDeleteMessage} />
                        </Route>
                        <Route path="/">
                            <MailList
                                onDeleteMessage={this.handleDeleteMessage}
                                onGetMessages={this.handleGetMessages}
                                start={this.state.start}
                                total={this.state.total}
                                messages={this.state.messages} />
                        </Route>
                    </Switch>
                </div>
            </div>
        );
    }
}

export default withRouter(App);
