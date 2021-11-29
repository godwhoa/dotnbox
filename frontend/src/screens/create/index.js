import { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function CreateRoom() {
  let navigate = useNavigate();
  const [roomName, setRoomName] = useState("");
  const [n, setN] = useState(4);
  const [m, setM] = useState(4);
  const createRoom = () => {
    fetch(`http://${window.location.hostname}:8080/room/${roomName}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ n, m }),
    }).then((res) => navigate(`/room/${roomName}`));
  };

  return (
    <div>
      <input
        type="text"
        value={roomName}
        placeholder="Room name"
        onChange={(e) => setRoomName(e.target.value)}
      />
      <input
        type="number"
        value={m}
        placeholder="M"
        onChange={(e) => setM(parseInt(e.target.value))}
        min="1"
        max="10"
        step="1"
      />
      <input
        type="number"
        value={n}
        placeholder="N"
        onChange={(e) => setN(parseInt(e.target.value))}
        min="1"
        max="10"
        step="1"
      />
      <button onClick={createRoom}>Create</button>
    </div>
  );
}
