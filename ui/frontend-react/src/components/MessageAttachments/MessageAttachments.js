import React from 'react';
import './MessageAttachments.css';
import downloadIcon from './i/download.svg';

const MAX_FILE_NAME_LENGTH = 15;

class MessageAttachments extends React.Component {
  renderFile(message, file, index) {
    let fileName = (file.Headers['Content-Type'] ||[])[0] || 'Unknown type';
    const contentDisposition = (file.Headers['Content-Disposition'] || [])[0] || '';
    const contentDispositionPrefix = 'attachment; filename=';

    if (contentDisposition.includes(contentDispositionPrefix)) {
      fileName = contentDisposition.substring(contentDispositionPrefix.length);
    }

    if (fileName.length > MAX_FILE_NAME_LENGTH) {
      fileName = fileName.substring(0, MAX_FILE_NAME_LENGTH) + '...';
    }

    return (
      <a href={`/api/v1/download-part?id=${message.id}&part=${index}`} className="file" key={index}>
        <img src={downloadIcon} className="file__download-icon" alt="download"/>
        {/*<div className="file__type">{`.${file.type.toUpperCase()}`}</div>*/}
        <div className={`file__name`}>{fileName}</div>
        <div className="file__size">{`(${file.Size} bytes)`}</div>
      </a>
    )
  }

  render() {
    const { message } = this.props;
    if (!message.MIME) {
      return null;
    }

    const parts = message.MIME.Parts || [];

    return (
      <div className="attached-files">
        <div className="attached-files__label">MIME</div>
        <div className="attached-files__box">
          {parts.map((file, index) => this.renderFile(message, file, index))}
        </div>
      </div>
    );
  }
}

export default MessageAttachments;
