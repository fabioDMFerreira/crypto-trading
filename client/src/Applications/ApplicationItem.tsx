import React from 'react';
import Card from 'react-bootstrap/Card';
import CardDeck from 'react-bootstrap/CardDeck';
import Col from 'react-bootstrap/Col';
import ListGroup from 'react-bootstrap/ListGroup';
import Row from 'react-bootstrap/Row';

import AssetsTable from '../components/AssetsTable';
import JsonDisplayer from '../components/JsonDisplayer';
import {
  Account, Application, ApplicationExecutionState, Asset, LogEvent,
} from '../types';
import ApplicationChart from './ApplicationChart';
import LogEventsTable from './LogEventsTable';


interface Props {
  application: Application,
  account: Account,
  lastApplicationState: ApplicationExecutionState
  assets: Asset[],
  logEvents: LogEvent[],
}

export default ({
  application, account, lastApplicationState, assets, logEvents,
}: Props) => (
  <div>
    <CardDeck>

      <Card>
        <Card.Body>
          <Card.Subtitle>Current Amount</Card.Subtitle>
          <Card.Title>
            {account.amount}
            {' '}
            €
          </Card.Title>
        </Card.Body>
      </Card>
      <Card>
        <Card.Body>
          <Card.Subtitle>Asset</Card.Subtitle>
          <Card.Title>{application.asset}</Card.Title>
        </Card.Body>
      </Card>
      <Card>
        <Card.Body>
          <Card.Subtitle>Current Price</Card.Subtitle>
          <Card.Title>
            {lastApplicationState.state.currentPrice}
            {' '}
            €
          </Card.Title>
        </Card.Body>
      </Card>
    </CardDeck>
    <Row className="mt-4">
      <Col xs={3}>
        <ListGroup>
          <ListGroup.Item>
            Average:
            {' '}
            {lastApplicationState.state.average}
          </ListGroup.Item>
          <ListGroup.Item>
            Current Change:
            {' '}
            {lastApplicationState.state.currentChange}
          </ListGroup.Item>
          <ListGroup.Item>
            Standard Deviation:
            {' '}
            {lastApplicationState.state.standardDeviation}
          </ListGroup.Item>
          <ListGroup.Item variant="success">
            Higher Bollinger Band:
            {' '}
            {lastApplicationState.state.higherBollingerBand}
          </ListGroup.Item>
          <ListGroup.Item variant="danger">
            Lower Bollinger Band:
            {' '}
            {lastApplicationState.state.lowerBollingerBand}
          </ListGroup.Item>
        </ListGroup>
      </Col>
      <Col xs={9}>
        <JsonDisplayer json={application} />
      </Col>
    </Row>
    <ApplicationChart asset={application.asset} appID={application._id} accountID={application.accountID} />
    {
        logEvents
        && (
        <div className="mt-4">
          <LogEventsTable
            logEvents={logEvents}
          />
        </div>
        )
      }
    <div className="mt-4">
      <AssetsTable assets={assets} />
    </div>
  </div>
);
