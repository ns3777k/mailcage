import React from 'react';
import './ControlPanel.css';
import emailIcon from './i/email.svg';
import refreshIcon from './i/refresh.svg';
import removeIcon from './i/cancel.svg';

class ControlPanel extends React.Component {
  render() {
    const {
      onRemoveAll,
      onRefresh,
      total
    } = this.props;

    return (
      <div className="navigation-table">
        <div className="navigation-table__total">
          <img src={emailIcon} className="navigation-table__total-icon" alt="total messages"/>
          <p className="p-normal">Total messages: {total}</p>
        </div>
        <div className="navigation-table__buttons">
          <button onClick={onRefresh} className="navigation-table__buttons-refresh">
            <img src={refreshIcon} className="navigation-table__buttons-refresh-icon" alt="refresh"/>
            <p className="p-normal btn-link">Refresh</p>
          </button>
          <button className="navigation-table__buttons-remove" onClick={onRemoveAll}>
            <img src={removeIcon} className="navigation-table__buttons-remove-icon" alt="remove"/>
            <p className="p-normal btn-link">Remove all</p>
          </button>
        </div>
      </div>
    );
  }
}

export default ControlPanel;
