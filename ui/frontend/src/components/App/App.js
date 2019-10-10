import React from 'react';
import MailList from '../MailList/MailList';

class App extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            messages: [],
        };

        this.onMessage = this.onMessage.bind(this);

        this.ws = new WebSocket(`ws://${window.location.host}/api/v1/ws`);
        this.ws.addEventListener('message', this.onMessage);
    }

    onMessage(e) {
        const message = JSON.parse(e.data);
        switch (message.Type) {
            case 'stored':
                this.setState((state, prevState) => {
                    return {
                        ...state,
                        messages: [
                            message.Payload,
                            ...state.messages
                        ],
                    };
                });
                break;

            case 'delete_message':
                break;

            case 'delete_messages':
                break;

            default:
                console.error('unknown message: ', e.data);
        }
    }

    componentDidMount() {
        fetch('/api/v1/messages')
            .then(response => response.json())
            .then(response => {
                this.setState((state, prevState) => ({
                    ...state,
                    messages: response.Items,
                }));
            });
    }

    render() {
        return (
            <div className="App">
                <MailList messages={this.state.messages} />
            </div>
        );
    }
}

export default App;
