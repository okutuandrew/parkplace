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
     socket.send(JSON.stringify({ type: "connected", role: "driver"}));


  };

  socket.onclose = () => {
    console.log("🔴 WebSocket disconnected");

    // reconnect after 3 seconds
    setTimeout(connectWebSocket, 3000);
  };

  socket.onerror = (err) => {
    console.error("WebSocket error:", err);
  };

  socket.onmessage = (event) => {

    console.log("📨 Message from server:", event.data);

    // We'll process messages later.
  };
}

// 👇 Add this line
connectWebSocket();