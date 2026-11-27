import { PlayerList } from "./player-list.js";

window.quizRoom = Object.freeze({
  players: new PlayerList(document.querySelector(".players-card")),
});

// Host controls stay disabled until authenticated host actions exist.

const ws = new WebSocket(`ws://${location.host}/ws`);

ws.onopen = () => console.log("Connected");
ws.onclose = e => console.log("Closed:", e.code);
ws.onerror = e => console.log("Error:", e);
ws.onmessage = (event) => {
  let message;
  try {
    message = JSON.parse(event.data);
  } catch {
    return;
  }
  if (message?.type !== "lobby_players" || !Array.isArray(message.players)) return;

  const players = message.players.map(player => ({
    id: player?.ID,
    nickname: player?.Nickname,
    online: true,
  }));
  window.quizRoom.players.replace(players);
};
