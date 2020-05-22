import React from 'react';
import './MessageContainer.css';
import {withRouter} from 'react-router-dom';
import {getMessage, deleteMessage, getOutgoingServers, markAsRead, releaseMessage, getDownloadMessageLink} from '../../api/mailcage';
import MessageHeaders from '../MessageHeaders/MessageHeaders';
import MessageControlPanel from '../MessageControlPanel/MessageControlPanel';
import MessageAttachments from '../MessageAttachments/MessageAttachments';
import MessagePreview from '../MessagePreview/MessagePreview';

class MessageContainer extends React.Component {
  state = {
    outgoingServers: [],
    message: {
      id: '',
      unread: true,
      content: {},
      to: [],
      from: {},
      MIME: {},
    },
    error: '',
    isFetching: true,
  };

  async componentDidMount() {
    const {id} = this.props.match.params;
    let newState = {
      message: {
        id: '',
        unread: true,
        content: {},
        to: [],
        from: {},
        MIME: {},
      },
      error: '',
      isFetching: false,
    };

    try {
      const outgoingServers = await getOutgoingServers();
      this.setState({ outgoingServers });
    } catch (e) {
      //
    }

    try {
      const response = await getMessage(id);
      newState = Object.assign({}, this.state, newState, {
        message: {
          id: response.ID,
          unread: response.Unread,
          content: response.Content,
          from: response.From,
          to: response.To,
          MIME: response.MIME,
        },
        error: '',
      });

      markAsRead(id);

    } catch (e) {
      newState.error = `Error while fetching the message: ${e.statusText}`;
    }

    this.setState(newState);
  }

  handleReleaseMessage = async server => {
    await releaseMessage(this.state.message.id, server);
  };

  handleRemoveMessage = async () => {
    await deleteMessage(this.state.message.id);
    this.props.history.push('/');
  };

  handleDownloadMessage = async () => {
    window.location.href = getDownloadMessageLink(this.state.message.id);
  };

  render() {
    const {isFetching, error} = this.state;

    return (
      <div className="details">
        {error && <h3>{error}</h3>}
        {isFetching && <h2>Loading</h2>}
        {!isFetching && error.length === 0 && <>
          <MessageControlPanel onDownloadMessage={this.handleDownloadMessage}
                               onRemoveMessage={this.handleRemoveMessage}
                               onReleaseMessage={this.handleReleaseMessage}
                               outgoingServers={this.state.outgoingServers}/>
          <MessageHeaders message={this.state.message}/>
          <MessageAttachments message={this.state.message}/>
          <MessagePreview message={this.state.message}/>
        </>}
      </div>
    );
  }
}

export default withRouter(MessageContainer);
