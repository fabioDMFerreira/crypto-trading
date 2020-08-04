import React from 'react';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import {
  Link, RouteProps,
} from 'react-router-dom';

export default ({ location }:RouteProps) => (
  <Navbar bg="light" expand="lg">
    <Navbar.Brand href="/">Crypto Trading</Navbar.Brand>
    <Navbar.Toggle aria-controls="basic-navbar-nav" />
    <Navbar.Collapse id="basic-navbar-nav">
      <Nav className="mr-auto" activeKey={location && location.pathname}>
        <Nav.Item>
          <Nav.Link as={Link} to="/" href="/">
            Home
          </Nav.Link>
        </Nav.Item>
        <Nav.Item>
          <Nav.Link as={Link} to="/benchmarks" href="/benchmarks">
            Benchmarks
          </Nav.Link>
        </Nav.Item>
        <Nav.Item>
          <Nav.Link as={Link} to="/applications" href="/applications">
            Applications
          </Nav.Link>
        </Nav.Item>
      </Nav>
    </Navbar.Collapse>
  </Navbar>
);
