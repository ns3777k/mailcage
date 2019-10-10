import React from 'react';
import { Switch, Route } from 'react-router-dom';
import MailList from '../MailList/MailList';
import MailDetail from '../MailDetail/MailDetail';
import * as api from '../../api/mailcage';

class App extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            messages: [],
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
                }));
                break;

            case 'deleted_message':
                const id = message.Payload;
                this.setState((state, prevState) => ({
                    ...state,
                    messages: state.messages.filter(msg => msg.ID !== id)
                }));
                break;

            case 'deleted_messages':
                this.setState((state, prevState) => ({
                    ...state, messages: []
                }));
                break;

            default:
                console.error('unknown message: ', e.data);
        }
    };

    handleDeleteMessage = id => {
        return api.deleteMessage(id);
    };

    handleGetMessages = () => {
        api.getMessages()
            .then(response => {
                this.setState((state, prevState) => ({
                    ...state,
                    messages: response.Items,
                }));
            });
    };

    handleGetMessage = id => {
        return api.getMessage(id);
    };

    render() {
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
                                <MailList onDeleteMessage={this.handleDeleteMessage}
                                          onGetMessages={this.handleGetMessages}
                                          messages={this.state.messages} />
                            </Route>
                        </Switch>
                    </div>
                </div>
            </div>
        );
    }
}

export default App;
