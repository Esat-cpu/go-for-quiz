import { PlayerList } from "./player-list.js";

// No requests are made until the backend provides a real room transport.
// Both screens expose the same presentation API for the future adapter.
window.quizRoom = Object.freeze({
  players: new PlayerList(document.querySelector(".players-card")),
});

// Host controls stay disabled until authenticated host actions exist.
