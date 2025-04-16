import {showMainNotification} from "./notifications.js";
import {getCSRFToken} from "./authentication.js";

export function joinChannel() {

// join channel
    const joinChannelButton = document.querySelector("#join-channel-btn");

// ---- join channel ----
    joinChannelButton.addEventListener("submit", function (event) {
        event.preventDefault();
        // const form = event.target;
        // const formData = new FormData(form); // Collect form data
        const csrfToken = getCSRFToken();
        console.log("csrfToken: ", csrfToken);

        fetch("/channels/join", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "x-csrf-token": csrfToken,
            },
            body: JSON.stringify({
                channelId: document.getElementById("join-channel-id").value,
                agree: document.getElementById("rules-agree-checkbox").value,
            }),
            cache: "no-store",
        })
            .then((response) => {
                if (response.ok) {
                    return response.json(); // Parse JSON response
                } else {
                    throw new Error("Join channel failed.");
                }
            })
            .then((data) => {
                // Show notification based on the response from the server
                showMainNotification(data.message);
                // Optional: Redirect the user after showing the notification
                setTimeout(() => {
                    window.location.href = "/"; // Replace with your desired location
                }, 3500);
            })
            .catch((error) => {
                console.error("Error:", error);
                showMainNotification("An error occurred while joining channel.");
            });
    });
}