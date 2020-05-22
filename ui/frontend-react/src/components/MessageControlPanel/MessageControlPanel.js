import React from 'react';
import './MessageControlPanel.css';
import downloadIcon from './i/download.svg';
import resendIcon from './i/resend.svg';
import backIcon from './i/back.svg';
import removeIcon from './i/cancel.svg';
import { withRouter } from 'react-router-dom';

class MessageControlPanel extends React.Component {
  state = {
    openSelected: false,
    selected: '',
  }

  handleOpenSelected = () => this.setState({
    openSelected: true,
  });

  handleCloseSelected = () => this.setState({
    openSelected: false,
  });

  handleItemSelect = e => {
    this.setState({ selected: e.target.value });
  };

  handleReleaseClick = e => {
    e.preventDefault();

    this.props.onReleaseMessage(this.state.selected);
  };

  renderSelectedForRelease = () => {
    // const { selected } = this.state;
    const { outgoingServers } = this.props;

    return (
      <select onChange={this.handleItemSelect}>
        <option value="">Choose one</option>
        {outgoingServers.map(outgoingServer => {
          return (
            <option key={outgoingServer} value={outgoingServer}>{outgoingServer}</option>
          );
        })}
      </select>
      // <div className="selected-mobile">
      //   <div
      //     className="selected"
      //     onClick={this.handleOpenSelectList}
      //   >
      //     <p className="p-normal btn-link">{selected || 'Choose one'}</p>
      //     <img
      //       src={downloadIcon}
      //       className="navigation-details__buttons-download-icon"
      //       alt="open"
      //     />
      //   </div>
      //   {openSelectList && this.renderSelectResendList()}
      // </div>
    )
  }

  render() {
    const { openSelected, selected } = this.state;
    const { onRemoveMessage, onDownloadMessage, history } = this.props;

    return (
      <div className="navigation-details">
        <div className="navigation-details-mobile">
          <div className="navigation-details__buttons">
            <button className="navigation-details__buttons-remove" onClick={onRemoveMessage}>
              <img src={removeIcon} className="navigation-details__buttons-remove-icon" alt="remove" />
              <span className="p-normal btn-link">Remove</span>
            </button>
            <button className="navigation-details__buttons-download" onClick={onDownloadMessage}>
              <img src={downloadIcon} className="navigation-details__buttons-download-icon" alt="download" />
              <span className="p-normal btn-link">Download</span>
            </button>
            <div className="navigation-details__buttons-separate" />
            {!openSelected && <div className="navigation-details__resend">
              <button className={`navigation-details__buttons-resend`} onClick={this.handleOpenSelected}>
                <img src={resendIcon} className="navigation-details__buttons-resend-icon" alt="resend"/>
                <span className="p-normal btn-link">Release</span>
              </button>
            </div>}
          </div>
          <div className="navigation-details__resend-mobile">
            {openSelected && this.renderSelectedForRelease()}
            {openSelected && (
              <div className="navigation-details__resend-select">
                <button onClick={this.handleReleaseClick}
                        className={`navigation-details__buttons-release ${selected && 'active'}`}>
                  <p className="p-normal">Release!</p>
                </button>
                <div className="navigation-details__resend-close" onClick={this.handleCloseSelected}>+</div>
              </div>
            )}
          </div>
        </div>
        <div className="navigation-details__buttons">
          <button onClick={() => history.push('/')} className="navigation-details__buttons-back p-normal btn-link">
            <img src={backIcon} className="navigation-details__buttons-back-icon" alt="back" />
            Back
          </button>
        </div>
      </div>
    );
  }
}

export default withRouter(MessageControlPanel);
