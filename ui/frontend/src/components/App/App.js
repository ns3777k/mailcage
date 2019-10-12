import React from 'react';
import { Switch, Route } from 'react-router-dom';
import MailList from '../MailList/MailList';
import MailDetail from '../MailDetail/MailDetail';
import * as api from '../../api/mailcage';
import { withRouter } from 'react-router-dom';
import { parse } from 'query-string';

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
    }

    onMessage = e => {
        const message = JSON.parse(e.data);
        switch (message.Type) {
            case 'stored':
                this.setState((state, prevState) => ({
                    ...state,
                    messages: [ message.Payload, ...state.messages ],
                    total: state.total + 1,
                }));
                break;

            case 'deleted_message':
                const id = message.Payload;
                this.setState((state, prevState) => ({
                    ...state,
                    messages: state.messages.filter(msg => msg.ID !== id),
                    total: state.total + 1,
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

    handleGetMessages = (start) => {
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
        const query = parse(this.props.location.search);
        const currentStart = Number(query.start || 0) || 0;
        return (
            <div className="grid-container fluid">
                <div className="grid-x grid-margin-x">
                    <div className="cell small-offset-1 small-10">
                        <Switch>
                            <Route path="/:id">
                                <MailDetail onGetMessage={this.handleGetMessage}
                                            onDeleteMessage={this.handleDeleteMessage} />
                            </Route>
                            <Route path="/">
                                <MailList
                                    onDeleteMessage={this.handleDeleteMessage}
                                    onGetMessages={this.handleGetMessages}
                                    start={currentStart}
                                    total={this.state.total}
                                    messages={this.state.messages} />
                            </Route>
                        </Switch>
                    </div>
                </div>
            </div>
        );
    }
}

export default withRouter(App);
