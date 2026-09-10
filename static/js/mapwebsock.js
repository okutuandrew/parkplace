// =======================
// WEBSOCKET
// =======================
let socket = null;

console.log("✅ mapwebsock.js loaded");

function connectWebSocket() {

  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const wsURL = `${protocol}//${window.location.host}/ws`;

  socket = new WebSocket(wsURL);

  socket.onopen = () => {
    console.log("🟢 WebSocket connected");

    socket.send(JSON.stringify({
      type: "connected",
      role: "driver"
    }));
  };

  socket.onclose = () => {
    console.log("🔴 WebSocket disconnected");

    setTimeout(connectWebSocket, 3000);
  };

  socket.onerror = (err) => {
    console.error("WebSocket error:", err);
  };

  socket.onmessage = (event) => {
    console.log("📨 Message from server:", event.data);
  };
}

connectWebSocket();


// =======================
// PLATE AUTOCOMPLETE
// =======================

const plateInput = document.getElementById("maskednumberplate");
const plateSuggestions = document.getElementById("plateSuggestions");

plateInput.addEventListener("input", async function () {

  const query = plateInput.value.trim();

  if (query.length < 2) {
    plateSuggestions.innerHTML = "";
    return;
  }

  const response = await fetch(
    `/api/plates?q=${encodeURIComponent(query)}`
  );

  const plates = await response.json();

  plateSuggestions.innerHTML = "";

  plates.forEach(function (plate) {

    const option = document.createElement("div");

    option.textContent = plate;

    option.addEventListener("click", function () {
      plateInput.value = plate;
      plateSuggestions.innerHTML = "";
    });

    plateSuggestions.appendChild(option);
  });
});