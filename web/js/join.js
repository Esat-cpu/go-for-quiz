"use strict";

const form = document.querySelector("#join-form");
const nicknameInput = document.querySelector("#nickname");
const joinButton = document.querySelector("#join-button");
const errorMessage = document.querySelector("#join-error");
const statusMessage = document.querySelector("#join-status");

function showError(message) {
  errorMessage.textContent = message;
  errorMessage.hidden = false;
}

nicknameInput.addEventListener("input", () => {
  nicknameInput.removeAttribute("aria-invalid");
  errorMessage.hidden = true;
});

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  if (joinButton.disabled) return;

  errorMessage.hidden = true;
  const nickname = nicknameInput.value.trim();
  if (!nickname) {
    nicknameInput.setAttribute("aria-invalid", "true");
    showError("Lütfen boşluklardan oluşmayan bir takma ad yaz.");
    nicknameInput.focus();
    return;
  }

  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 12000);
  let joined = false;
  joinButton.disabled = true;
  joinButton.textContent = "Katılınıyor…";
  form.setAttribute("aria-busy", "true");
  statusMessage.textContent = "Yarışmaya katılınıyor.";

  try {
    const response = await fetch("/join", {
      method: "POST",
      credentials: "same-origin",
      headers: { "Content-Type": "application/json", "Accept": "application/json, text/plain" },
      body: JSON.stringify({ nickname }),
      signal: controller.signal,
    });

    if (!response.ok) {
      // Existing Go handlers return plain-text errors through http.Error.
      const serverMessage = (await response.text()).trim();
      showError(serverMessage || "Yarışmaya katılamadın. Lütfen tekrar dene.");
      statusMessage.textContent = "";
      return;
    }

    joined = true;
    statusMessage.textContent = "Katıldın. Bekleme odası açılıyor.";
    window.location.assign("/lobby");
  } catch (error) {
    showError(error.name === "AbortError"
      ? "Sunucu zamanında yanıt vermedi. Lütfen tekrar dene."
      : "Sunucuya ulaşılamadı. Ağ bağlantını kontrol edip tekrar dene.");
    statusMessage.textContent = "";
  } finally {
    clearTimeout(timeout);
    if (!joined) {
      joinButton.disabled = false;
      joinButton.textContent = "Yarışmaya katıl";
      form.removeAttribute("aria-busy");
    }
  }
});
