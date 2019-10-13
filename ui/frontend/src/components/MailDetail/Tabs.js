import React from 'react';
import Tab from './Tab';

class Tabs extends React.Component {
    render() {
        const { tabs } = this.props;

        return (
            <ul className="tabs">
                {tabs.map(tab => {
                    return (
                        tab.cond
                            ? <Tab key={tab.tab} tab={tab} onTabClick={this.props.onTabClick} />
                            : null
                    );
                })}
            </ul>
        );
    }
}

export default Tabs;
