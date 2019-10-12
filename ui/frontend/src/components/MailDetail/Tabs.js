import React from 'react';
import Tab from './Tab';

class Tabs extends React.Component {
    render() {
        const { tabs } = this.props;

        return (
            <ul className="tabs">
                {tabs.map(tab => {
                    return (
                        <Tab key={tab.tab} tab={tab} onTabClick={this.props.onTabClick} />
                    );
                })}
            </ul>
        );
    }
}

export default Tabs;
