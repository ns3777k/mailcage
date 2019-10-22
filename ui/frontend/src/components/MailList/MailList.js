import React from 'react';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';
import { Button } from 'primereact/button';
import { PER_PAGE } from '../../constants';
import { getRecipients, getSender } from '../../utils/formatter';
import parseISO from 'date-fns/parseISO';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import { withRouter } from 'react-router-dom';

class MailList extends React.Component {
    componentDidMount() {
        this.props.onGetMessages(this.props.start);
    }

    onPage = e => {
        this.props.onGetMessages(e.first);
    };

    handleRefresh = e => {
        this.props.onGetMessages(this.props.start);
    };

    handleSelect = e => {
        this.props.history.push(`/${e.data.id}`);
    };

    actionTemplate = e => {
        const { id } = e;
        return (
            <Button type="button"
                    icon="pi pi-trash"
                    onClick={e => {
                        e.stopPropagation();
                        this.props.onDeleteMessage(id);
                    }}
                    className="p-button-danger"/>
        );
    };

    render() {
        const { start, total } = this.props;
        const paginatorLeft = <Button onClick={this.handleRefresh} icon="pi pi-refresh"/>;
        const paginatorRight = <Button label={`Total messages: ${total}`} />;
        const messages = this.props.messages.map(message => ({
            id: message.ID,
            from: getSender(message),
            to: getRecipients(message).join(', '),
            subject: message.Content.Headers['Subject'][0],
            when: `${formatDistanceToNow(parseISO(message.CreatedAt))} ago`
        }));

        return (
            <>
                <DataTable lazy
                           paginator
                           responsive
                           onRowSelect={this.handleSelect}
                           selectionMode="single"
                           paginatorLeft={paginatorLeft}
                           paginatorRight={paginatorRight}
                           value={messages}
                           first={start}
                           rows={PER_PAGE}
                           onPage={this.onPage}
                           totalRecords={total}>
                    <Column field="from" header="From" />
                    <Column field="to" header="To" />
                    <Column field="when" header="When" />
                    <Column field="subject" header="Subject" />
                    <Column body={this.actionTemplate} style={{ textAlign:'center', width: '8em' }}/>
                </DataTable>
            </>
        );
    }
}

export default withRouter(MailList);
