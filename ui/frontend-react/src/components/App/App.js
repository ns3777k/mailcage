import React from 'react';
import './App.css';
import { BrowserRouter, Redirect, Switch, Route } from 'react-router-dom';
import MessageListContainer from '../MessageListContainer/MessageListContainer';
import MessageContainer from "../MessageContainer/MessageContainer";

class App extends React.Component {
  render() {
    return (
      <div className="app">
        <div className="line"/>
        <BrowserRouter forceRefresh={false}>
          <Switch>
            <Route path="/" exact><MessageListContainer />
            </Route>
            <Route path="/message/:id">
              <MessageContainer />
            </Route>
            <Route path="*">
              <Redirect to="/" />
            </Route>
          </Switch>
        </BrowserRouter>
      </div>
    );
  }
}

export default App;
