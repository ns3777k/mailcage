import React from 'react';
import { getSender, getRecipients } from '../../utils/formatter';
import './MessageHeaders.css';

class MessageHeaders extends React.Component {
  state = {
    showAll: false,
  };

  toggleHeaders = () => {
    this.setState({ showAll: !this.state.showAll });
  };

  renderBriefHeaders() {
    const { message } = this.props;

    const headers = [
      { header: 'Subject', value: message.content.Headers['Subject'][0] },
      { header: 'From', value: getSender(message.content.Headers, message.from) },
      { header: 'To', value: getRecipients(message.content.Headers, message.to).join(', ') },
    ];

    const renderItem = header => {
      return (
        <div key={header.header} className="headers__line">
          <div className="headers__line-label">{header.header}</div>
          <div className="headers__line-text">{header.value}</div>
        </div>
      );
    };

    return (
      <div>
        {headers.map(renderItem)}
      </div>
    );
  }

  renderAllHeaders() {
    const { message } = this.props;
    const renderItem = headerName => {
      const headerValue = message.content.Headers[headerName].join(', ');
      return (
        <div key={headerName} className="headers__line">
          <div className="headers__line-label">{headerName}</div>
          <div className="headers__line-text">{headerValue}</div>
        </div>
      );
    };

    return (
      <div>
        {Object.keys(message.content.Headers).map(renderItem)}
      </div>
    );
  }

  render() {
    const toggleText = this.state.showAll ? 'Show less' : 'Show all';
    return (
      <div className="headers">
        {!this.state.showAll && this.renderBriefHeaders()}
        {this.state.showAll && this.renderAllHeaders()}
        <span onClick={this.toggleHeaders} className="headers__line-show-all">{toggleText}</span>
      </div>
    );
  }
}

export default MessageHeaders;
