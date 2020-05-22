import React from 'react';
import ControlPanel from './../ControlPanel/ControlPanel';
import MessageList from './../MessageList/MessageList';
import {getMessages, deleteAllMessages, deleteMessage} from '../../api/mailcage';
import {humanizeDateDistance} from '../../utils/date';
import {getRecipients, getSender} from '../../utils/formatter';

class MessageListContainer extends React.Component {
  state = {
    messages: [],
    total: 0,
    count: 0,
    start: 0,
    isFetching: true,
    error: '',
  };

  async queryMessageList(start = 0) {
    let newState = {
      messages: [],
      total: 0,
      count: 0,
      start: 0,
      isFetching: false,
      error: '',
    };

    try {
      const response = await getMessages(start);

      newState = Object.assign({}, this.state, newState, {
        messages: (response.Items || []).map(message => ({
          id: message.ID,
          from: getSender(message.Content.Headers, message.From),
          to: getRecipients(message.Content.Headers, message.To).join(', '),
          subject: message.Content.Headers['Subject'][0],
          when: humanizeDateDistance(message.CreatedAt),
          unread: message.Unread,
        })),
        total: response.Total,
        count: response.Count,
        start: response.Start,
      });
    } catch (e) {
        newState.error = e;
    }

    this.setState(newState);
  }

  componentDidMount() {
    this.queryMessageList(this.state.start);
  }

  handleRemoveAll = async () => {
    await deleteAllMessages();
    this.queryMessageList(this.state.start);
  };

  handleRefresh = () => {
    this.queryMessageList();
  };

  handleRemoveMessage = async id => {
    await deleteMessage(id);
    this.queryMessageList(this.state.start);
  };

  render() {
    const {isFetching, error, total} = this.state;

    return (
      <div>
        {error && <h3>{error}</h3>}
        {isFetching && <h2>Loading</h2>}
        {!isFetching &&
        <>
          <ControlPanel
            onRemoveAll={this.handleRemoveAll}
            onRefresh={this.handleRefresh}
            total={total}
          />
          <MessageList onRemoveMessage={this.handleRemoveMessage} messages={this.state.messages}/>
        </>}
      </div>
    );
  }
}

export default MessageListContainer;
