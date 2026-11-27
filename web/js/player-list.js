export class PlayerList {
  constructor(root) {
    this.players = new Map();
    this.available = false;
    this.list = root.querySelector("#player-list");
    this.count = root.querySelector("#player-count");
    this.empty = root.querySelector("#players-empty");
    this.emptyTitle = root.querySelector("#empty-title");
    this.emptyDescription = root.querySelector("#empty-description");
    this.status = root.querySelector("#players-status");
    this.render();
  }

  // Player shape: { id: string, nickname: string, online: boolean }.
  // A complete snapshot replaces previous data; only online players appear.
  replace(players) {
    this.players.clear();
    for (const player of players) this.store(player);
    this.available = true;
    this.render();
  }

  upsert(player) {
    this.store(player);
    this.render();
  }

  remove(playerId) {
    this.players.delete(playerId);
    this.render();
  }

  // Keep the count unknown until a complete snapshot has arrived.
  setUnavailable() {
    this.players.clear();
    this.available = false;
    this.render();
  }

  store(player) {
    if (!player || typeof player.id !== "string" || !player.id) return;
    if (player.online !== true || typeof player.nickname !== "string" || !player.nickname.trim()) {
      this.players.delete(player.id);
      return;
    }
    this.players.set(player.id, { id: player.id, nickname: player.nickname });
  }

  render() {
    const players = this.available ? [...this.players.values()] : [];
    players.sort((a, b) =>
      a.nickname.localeCompare(b.nickname, "tr") || a.id.localeCompare(b.id, "en")
    );
    const fragment = document.createDocumentFragment();
    for (const player of players) {
      const item = document.createElement("li");
      item.className = "player-item";
      const avatar = document.createElement("span");
      avatar.className = "player-avatar";
      avatar.setAttribute("aria-hidden", "true");
      avatar.textContent = Array.from(player.nickname.trim())[0].toLocaleUpperCase("tr-TR");
      const name = document.createElement("span");
      name.className = "player-name";
      // Nicknames are untrusted input: never interpolate them into HTML.
      name.textContent = player.nickname;
      item.append(avatar, name);
      fragment.append(item);
    }
    this.list.replaceChildren(fragment);
    this.list.hidden = players.length === 0;
    this.empty.hidden = players.length > 0;
    this.count.textContent = this.available ? String(players.length) : "—";
    this.count.parentElement.setAttribute("aria-label", this.available
      ? `${players.length} çevrimiçi oyuncu`
      : "Çevrimiçi oyuncu sayısı henüz bilinmiyor");
    this.emptyTitle.textContent = this.available ? "Henüz kimse katılmadı" : "Oyuncu listesi henüz kullanılamıyor";
    this.emptyDescription.textContent = this.available
      ? "Katılan oyuncular burada görünecek."
      : "Canlı bağlantı bekleniyor.";
    this.status.textContent = this.available
      ? `${players.length} çevrimiçi oyuncu.`
      : "Oyuncu listesi henüz kullanılamıyor.";
  }
}
