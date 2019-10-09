import React from 'react';
import MailList from '../MailList/MailList';

class App extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            messages: [],
        };
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
