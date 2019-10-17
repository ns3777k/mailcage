import React from 'react';
import MailItem from './MailItem';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';
import { Button } from 'primereact/button';
import { Link } from 'react-router-dom';
import { PER_PAGE } from '../../constants';
import {getRecipients, getSender} from "../../utils/formatter";
import parseISO from 'date-fns/parseISO';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';

class MailList extends React.Component {
    componentDidMount() {
        this.props.onGetMessages(this.props.start);
    }

    // componentDidUpdate(prevProps) {
    //     if (prevProps.start !== this.props.start) {
    //         this.props.onGetMessages(this.props.start);
    //     }
    // }

    onPage = e => {
        console.log(e);

        this.props.onGetMessages(e.first);
    };

    render() {
        const { total } = this.props;
        // const endRange = start + PER_PAGE > total ? total : start + PER_PAGE;
        // const hasPrev = start > 0;
        // const hasNext = start + PER_PAGE < total;

        // const paginatorLeft = <Button icon="pi pi-refresh"/>;
        // const paginatorRight = <Button icon="pi pi-cloud-upload"/>;
        const messages = this.props.messages.map(message => ({
            from: getSender(message),
            to: getRecipients(message).join(', '),
            subject: message.Content.Headers['Subject'][0],
            when: `${formatDistanceToNow(parseISO(message.CreatedAt))} ago`
        }));

        return (
            <>
                <DataTable value={messages} paginator lazy rows={PER_PAGE} onPage={this.onPage} totalRecords={total}>
                    <Column field="from" header="From" />
                    <Column field="to" header="To" />
                    <Column field="when" header="When" />
                    <Column field="subject" header="Subject" />
                </DataTable>
                {/*<nav aria-label="Pagination">*/}
                {/*    <ul className="pagination">*/}
                {/*        <li className={`${hasPrev ? '' : 'disabled'}`}>*/}
                {/*            {hasPrev*/}
                {/*                ? <Link to={`?start=${start - PER_PAGE}`}>Prev</Link>*/}
                {/*                : <span>Prev</span>}*/}
                {/*        </li>*/}
                {/*        <li>{start + 1} - {endRange} of {total}</li>*/}
                {/*        <li className={`${hasNext ? '' : 'disabled'}`}>*/}
                {/*            {hasNext*/}
                {/*                ? <Link to={`?start=${start + PER_PAGE}`}>Next</Link>*/}
                {/*                : <span>Next</span>}*/}
                {/*        </li>*/}
                {/*    </ul>*/}
                {/*</nav>*/}
                {/*<table className="table-expand hover">*/}
                {/*    <thead>*/}
                {/*    <tr className="table-expand-row">*/}
                {/*        <th>From</th>*/}
                {/*        <th>To</th>*/}
                {/*        <th>When</th>*/}
                {/*        <th>Subject</th>*/}
                {/*        <th>&nbsp;</th>*/}
                {/*    </tr>*/}
                {/*    </thead>*/}
                {/*    <tbody>*/}
                {/*    {this.props.messages.map(message => {*/}
                {/*        return (*/}
                {/*            <MailItem onDeleteMessage={this.props.onDeleteMessage}*/}
                {/*                      key={message.ID}*/}
                {/*                      message={message} />*/}
                {/*        );*/}
                {/*    })}*/}
                {/*    </tbody>*/}
                {/*</table>*/}
            </>
        );
    }
}

export default MailList;
