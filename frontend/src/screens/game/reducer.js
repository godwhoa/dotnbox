export const CONNECTION_STATE = {
  CONNECTING: 0,
  CONNECTED: 1,
  DISCONNECTED: 2,
};

export const GAME_STATE = {
  WAITING: 0,
  PLAYER_ONE_TURN: 1,
  PLAYER_TWO_TURN: 2,
  PAUSED: 3,
  GAME_OVER: 4,
};

export const PAYLOAD_TYPE = {
  GAMECONFIG: "GAMECONFIG",
  STATE: "STATE",
  ERROR: "ERROR",
};

export const initialState = {
  m: 0,
  n: 0,
  connection: CONNECTION_STATE.CONNECTING,
  player: 0,
  turn: 0,
  state: GAME_STATE.WAITING,
  grid: {},
  boxes: {},
  scores: {},
  error: null,
};

const CONNECTION_STATE_CHANGED = "CONNECTION_STATE_CHANGED";

export function reducer(state = initialState, action) {
  switch (action.type) {
    case PAYLOAD_TYPE.STATE:
    case PAYLOAD_TYPE.GAMECONFIG:
    case PAYLOAD_TYPE.ERROR:
      return {
        ...state,
        ...action.payload,
      };
    case CONNECTION_STATE_CHANGED:
      return {
        ...state,
        connection: action.payload,
      };
    default:
      return state;
  }
}
