export const PlayerOneColor = "#BAE89B";
export const PlayerTwoColor = "#B9BAC7";
export const PlayerNoneColor = "#f6eded";

export const playerColor = (player) => {
  switch (player) {
    case 1:
      return PlayerOneColor;
    case 2:
      return PlayerTwoColor;
    default:
      return PlayerNoneColor;
  }
};

export function lineKey(line) {
  const {
    from: { x, y },
    to: { x: x2, y: y2 },
  } = order(line);
  return `from-${x}-${y}-to-${x2}-${y2}`;
}

export function order({ from, to }) {
  if (from.x > to.x) {
    return { from: to, to: from };
  }
  return { from, to };
}

export function interleave(a, b) {
  return a.flatMap((item, index) => [item, b[index]]);
}

export function* range(N) {
  for (let i = 0; i < N; i++) {
    yield i;
  }
}
