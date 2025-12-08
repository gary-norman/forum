let socket;



function connectWebSocket() {
    socket = new WebSocket("ws://localhost:8888");

    socket.onopen = function() {
        console.log("WebSocket connected");
    };

    socket.onclose = () => {
        console.log("Disconnected from WebSocket server");
    };

    socket.onerror = function(error) {
        console.error("WebSocket error:", error);
    };

    socket.onmessage = (e) => {
        console.log("Message from server: " + e.data);
        const outputMessage = document.getElementById("outputMessage");
        outputMessage.innerHTML += `<p>${e.data}</p>`;
    };


}




document.addEventListener("DOMContentLoaded", () => {
    // Initial connection
    connectWebSocket();



    const sendMessageButton = document.getElementById("sendMessage");
    sendMessageButton.addEventListener("click", sendMessage);
});

export function sendMessage() {
    const messageInput = document.getElementById("messageInput");
    const message = messageInput.value;

    console.log("sent!")

    // Check if WebSocket is open before sending
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(message);
        messageInput.value = "";
    } else {
        console.error("WebSocket is not open. State:", socket.readyState);
        alert("Connection is closed. Reconnecting...");
        // Optionally reconnect
        connectWebSocket();
    }
}
