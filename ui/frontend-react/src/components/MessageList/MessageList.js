import React from 'react';
import './MessageList.css';
import sendToIcon from './i/sendTo.svg';
import Swipeout from 'rc-swipeout';
import { withRouter } from 'react-router-dom';

class MessageList extends React.Component {

  handleRemoveMessage = (e, id) => {
    e.stopPropagation();

    this.props.onRemoveMessage(id);
  };

  renderTable() {
    const { messages } = this.props;
    const renderItem = item => {
      return (
        <tr key={item.id} onClick={() => this.props.history.push(`/message/${item.id}`)} className={`mail-table__edit`}>
          <td className="mail-table__bullet">
            <div className={`mail-table__bullet-${item.unread ? 'new' : 'old'}`}/>
          </td>
          <td className="mail-table__from">
            <p className={item.unread ? "p-title" : "p-normal"}>{item.from}</p>
            <p className="p-small">{item.from}</p>
          </td>
          <td className="mail-table__to">
            <p className={item.unread ? "p-title" : "p-normal"}>{item.to}</p>
            <p className="p-small">{item.to}</p>
          </td>
          <td className="mail-table__title">
            <p className={item.unread ? "p-title" : "p-normal"}>{item.subject}</p>
            <p className="p-small">{item.subject}</p>
          </td>
          <td className="mail-table__time">
            <p className="p-small">{item.when}</p>
          </td>
          <td>
            <button onClick={e => this.handleRemoveMessage(e, item.id)} className="mail-table__remove">+</button>
          </td>
        </tr>
      )
    };

    return (
      <table className="mail-table__table">
        <thead>
        <tr>
          <th/>
          <th>from</th>
          <th>to</th>
          <th>subject</th>
          <th>when</th>
          <th/>
        </tr>
        </thead>
        <tbody>
        {messages.map(renderItem)}
        </tbody>
      </table>
    );
  }

  renderMobileTable() {
    const { messages, onRemoveMessage } = this.props;
    const renderItem = item => {
      return (
        <Swipeout key={item.id} right={[
          {
            text: <button className="mail-table__remove_white">+</button>,
            onPress: () => onRemoveMessage(item.id),
            style: { backgroundColor: '#000AFF', color: 'white' }
          }
        ]}>
          <div className="mail-table__mobile-row" onClick={() => this.props.history.push(`/message/${item.id}`)}>
            <div className="mail-table__mobile-from">
              <div className={`mail-table__mobile-bullet-${item.unread ? 'new' : 'old'}`}/>
              <p className={`mail-table__mobile-title ${item.unread && 'mail-table__mobile-title_weight_bold'}`}>{item.to}</p>
              <p className="mail-table__mobile-small">{item.from}</p>
            </div>
            <div className="mail-table__mobile-to">
              <img src={sendToIcon} className="mail-table__mobile-to-icon" alt="to"/>
              <p className={`mail-table__mobile-title ${item.unread && 'mail-table__mobile-title_weight_bold'}`}>{item.to}</p>
              <p className="mail-table__mobile-small">{item.to}</p>
            </div>
            <div className="mail-table__mobile-subject">
              <p className={`mail-table__mobile-subject-title ${item.unread && 'mail-table__mobile-subject-title_weight_bold'}`}>{item.subject}</p>
              <p className="mail-table__mobile-small">{item.subject}</p>
            </div>
            <div className="mail-table__mobile-time">
              <p className="mail-table__mobile-small-time">{item.when}</p>
            </div>
          </div>
        </Swipeout>
      );
    };

    return (
      <div>
        {messages.map(renderItem)}
      </div>
    );
  }

  render() {
    return (
      <div className="mail-table">
        {this.renderTable()}
        {this.renderMobileTable()}
      </div>
    );
  }
}

export default withRouter(MessageList);
