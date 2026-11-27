import { PlayerList } from "./player-list.js";

// No requests are made until the backend provides a real room transport.
// Both screens expose the same presentation API for the future adapter.
window.quizRoom = Object.freeze({
  players: new PlayerList(document.querySelector(".players-card")),
});

// Host controls stay disabled until authenticated host actions exist.

const ws = new WebSocket(`ws://${location.host}/ws`);

ws.onopen = () => console.log("Connected");
ws.onclose = e => console.log("Closed:", e.code);
ws.onerror = e => console.log("Error:", e);
