import React from 'react';
import './MessagePreview.css';
import {Tab, Tabs, TabList, TabPanel} from 'react-tabs';
import {getHtmlMessage, isHtmlMessage, formatMessagePlain} from '../../utils/helpers';

class MessagePreview extends React.Component {
  render() {
    const {message} = this.props;
    const isHTML = isHtmlMessage(message);

    return (
      <Tabs selectedTabPanelClassName="letter" selectedTabClassName="tabs__label_active">
        <TabList className="tabs">
          {isHTML && <Tab className="tabs__label">HTML</Tab>}
          {!isHTML && <Tab className="tabs__label">Plain</Tab>}
          <Tab className="tabs__label">Source</Tab>
        </TabList>
        {isHTML && <TabPanel>
          <iframe seamless
                  srcDoc={`${getHtmlMessage(message)}`}
                  title="Message preview sandbox"
                  frameBorder="0"
                  style={{width: '100%'}}/>
        </TabPanel>}
        {!isHTML && <TabPanel>
          {formatMessagePlain(message)}
        </TabPanel>}
        <TabPanel>
          <pre>
            {Object.keys(message.content.Headers).map(header => {
              return (
                <div key={header}>{header}: {message.content.Headers[header]}</div>
              );
            })}
            <p>{message.content.Body}</p>
            </pre>
        </TabPanel>
      </Tabs>
    );
  }
}

export default MessagePreview;
