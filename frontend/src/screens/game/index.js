import Dot from "./components/Dot";
import Line from "./components/Line";
import Marker from "./components/Marker";
import Row from "./components/Row";
import { useEffect, useReducer, useRef } from "react";
import { lineKey, order, interleave, range, playerColor } from "./helpers";
import { reducer, initialState, CONNECTION_STATE, GAME_STATE } from "./reducer";
import styled from "styled-components";
import { useParams } from "react-router-dom";

function DotsAndLines({ N, y, handlePlace, lineColor }) {
  let dots = [...range(N + 1)].map((x) => <Dot key={`dot-${x}-${y}`} />);
  let lines = [...range(N)].map((x) => {
    const coordinates = { from: { x, y }, to: { x: x + 1, y } };
    return (
      <Line
        onClick={() => handlePlace(coordinates)}
        key={lineKey(coordinates)}
        color={lineColor(coordinates)}
      />
    );
  });
  return <Row>{interleave(dots, lines)}</Row>;
}

function LinesAndMarkers({
  N,
  y,
  handlePlace,
  lineColor,
  markerText,
  markerColor,
}) {
  let lines = [...range(N + 1)].map((x) => {
    const coordinates = { from: { x, y }, to: { x, y: y + 1 } };
    return (
      <Line
        onClick={() => handlePlace(coordinates)}
        vertical={true}
        key={lineKey(coordinates)}
        color={lineColor(coordinates)}
      />
    );
  });
  let markers = [...range(N)].map((x) => {
    return (
      <Marker
        key={`marker-${x}-${y}`}
        text={markerText(x, y)}
        color={markerColor(x, y)}
      />
    );
  });
  return <Row>{interleave(lines, markers)}</Row>;
}

const FlexContainer = styled.div`
  height: 100%;
  padding: 0;
  margin: 0;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const FlexRow = styled.div`
  width: auto;
`;

function Winner(score, player) {
  if (score[1] === score[2]) {
    return "Tie!";
  }
  if (score[1] > score[2] && player === 1) {
    return `You win!`;
  }
  if (score[2] > score[1] && player === 2) {
    return `You win!`;
  }
  return `You lose!`;
}

function GameHeader({ connection, turn, player, state, scores }) {
  if (connection === CONNECTION_STATE.CONNECTING) {
    return <h1>Connecting...</h1>;
  }
  switch (state) {
    case GAME_STATE.WAITING:
      return <h1>Waiting...</h1>;
    case GAME_STATE.PLAYER_ONE_TURN:
    case GAME_STATE.PLAYER_TWO_TURN:
      return <h1>{turn === player ? "Your turn" : "Their turn"}</h1>;
    case GAME_STATE.GAME_OVER:
      return <h1>{Winner(scores, player)}</h1>;
    default:
      return <></>;
  }
}

export default function Game() {
  const { roomID } = useParams();
  const [state, dispatch] = useReducer(reducer, initialState);
  const { n, m, connection } = state;
  const [N, M] = [n, m];
  const ws = useRef(null);

  useEffect(() => {
    ws.current = new WebSocket(
      `ws://${window.location.hostname}:8080/room/${roomID}`
    );
    ws.current.onmessage = (event) => {
      const action = JSON.parse(event.data);
      dispatch(action);
    };

    ws.current.onopen = () => {
      dispatch({
        type: "CONNECTION_STATE_CHANGED",
        payload: CONNECTION_STATE.CONNECTED,
      });
    };

    ws.current.onclose = (event) => {
      dispatch({
        type: "CONNECTION_STATE_CHANGED",
        payload: CONNECTION_STATE.DISCONNECTED,
      });
    };

    return () => {
      ws.current.close();
    };
  }, [roomID]);

  const handlePlace = (coordinates) => {
    if (connection === CONNECTION_STATE.CONNECTED) {
      ws.current.send(
        JSON.stringify({ type: "PLACE", payload: order(coordinates) })
      );
    }
  };

  const handlePlayAgain = () => {
    if (state.state === GAME_STATE.GAME_OVER) {
      ws.current.send(JSON.stringify({ type: "REMATCH" }));
    }
  };

  const lineColor = (coordinates) => {
    return playerColor(state.grid[lineKey(coordinates)] || 0);
  };

  let dotsnlines = [...range(M + 1)].map((y) => (
    <DotsAndLines
      key={`dl-${y}`}
      N={N}
      y={y}
      handlePlace={handlePlace}
      lineColor={lineColor}
    />
  ));
  let linesnmarkers = [...range(M)].map((y) => (
    <LinesAndMarkers
      key={`lm-${y}`}
      N={N}
      y={y}
      handlePlace={handlePlace}
      lineColor={lineColor}
      markerText={(x, y) => {
        const owner = state.boxes[`${x}-${y}`] || 0;
        return owner === 0 ? "" : `P${owner}`;
      }}
      markerColor={(x, y) => {
        const owner = state.boxes[`${x}-${y}`] || 0;
        return playerColor(owner);
      }}
    />
  ));
  return (
    <>
      <FlexContainer>
        <FlexRow>
          <GameHeader {...state} />
          <div>{interleave(dotsnlines, linesnmarkers)}</div>
          {state.state === GAME_STATE.GAME_OVER && (
            <button onClick={handlePlayAgain}>Play Again</button>
          )}
        </FlexRow>
      </FlexContainer>
    </>
  );
}
