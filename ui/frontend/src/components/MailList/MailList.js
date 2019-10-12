import React from 'react';
import MailItem from './MailItem';
import { Link } from 'react-router-dom';

const PER_PAGE = 50;

class MailList extends React.Component {
    componentDidMount() {
        this.props.onGetMessages(this.props.start);
    }

    componentDidUpdate(prevProps) {
        if (prevProps.start !== this.props.start) {
            this.props.onGetMessages(this.props.start);
        }
    }

    render() {
        const { start, total } = this.props;
        const endRange = start + PER_PAGE > total ? total : start + PER_PAGE;
        const hasPrev = start > 0;
        const hasNext = start + PER_PAGE < total;

        return (
            <>
                <nav aria-label="Pagination">
                    <ul className="pagination">
                        <li className={`${hasPrev ? '' : 'disabled'}`}>
                            {hasPrev
                                ? <Link to={`?start=${start - PER_PAGE}`}>Prev</Link>
                                : <span>Prev</span>}
                        </li>
                        <li>{start + 1} - {endRange} of {total}</li>
                        <li className={`${hasNext ? '' : 'disabled'}`}>
                            {hasNext
                                ? <Link to={`?start=${start + PER_PAGE}`}>Next</Link>
                                : <span>Next</span>}
                        </li>
                    </ul>
                </nav>
                <table className="table-expand hover">
                    <thead>
                    <tr className="table-expand-row">
                        <th>From</th>
                        <th>To</th>
                        <th>When</th>
                        <th>Subject</th>
                        <th>&nbsp;</th>
                    </tr>
                    </thead>
                    <tbody>
                    {this.props.messages.map(message => {
                        return (
                            <MailItem onDeleteMessage={this.props.onDeleteMessage}
                                      key={message.ID}
                                      message={message} />
                        );
                    })}
                    </tbody>
                </table>
            </>
        );
    }
}

export default MailList;
