import React from 'react';

class Tab extends React.Component {
    handleTabClick = e => {
        e.preventDefault();

        if (this.props.onTabClick) {
            this.props.onTabClick(this.props.tab.tab);
        }
    };

    render() {
        const { tab } = this.props;
        const activeClass = tab.current ? 'is-active' : '';

        return (
            <li onClick={this.handleTabClick} className={`tabs-title ${activeClass}`}>
                {/*eslint-disable-next-line jsx-a11y/role-supports-aria-props*/}
                <a href={`#${tab.tab}`} aria-selected={tab.current}>{tab.title}</a>
            </li>
        );
    }
}

export default Tab;
