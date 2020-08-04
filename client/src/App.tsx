import 'bootstrap/dist/css/bootstrap.min.css';
import 'react-day-picker/lib/style.css';
import 'rc-time-picker/assets/index.css';
import './App.css';

import React from 'react';
import {
  BrowserRouter as Router,
  Route,
  Switch,
} from 'react-router-dom';

import Applications from './Applications';
import Benchmark from './Benchmark';
import Navbar from './components/Navbar';

function App() {
  return (
    <Router>
      <Navbar />
      <div className="App">
        <Switch>
          <Route exact path="/" />
          <Route path="/benchmarks">
            <Benchmark />
          </Route>
          <Route path="/applications">
            <Applications />
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

export default App;
