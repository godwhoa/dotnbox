import React from "react";
import "./App.css";
import Game from "./screens/game";
import CreateRoom from "./screens/create";
import styled from "styled-components";
import { Routes, Route } from "react-router-dom";

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
`;

export default function App() {
  return (
    <Routes>
      <Route
        path="/"
        element={
          <Container>
            <CreateRoom />
          </Container>
        }
      />
      <Route
        path="/room/:roomID"
        element={
          <Container>
            <Game />
          </Container>
        }
      />
    </Routes>
  );
}
