import { showMainNotification } from "./notifications.js";
import { getCSRFToken } from "./authentication.js";

export function joinChannel() {
  const joinChannelButton = document.querySelector("#join-channel-btn");
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

export function agreeToJoin() {
  const checkbox = document.getElementById("rules-agree-checkbox");
  const button = document.getElementById("join-channel-btn");

  const updateButtonClasses = () => {
    if (checkbox.checked) {
      button.classList.add("btn-primary", "btn-action-primary");
    } else {
      button.classList.remove("btn-primary", "btn-action-primary");
    }
  };

  // Run once on load in case checkbox is pre-checked
  updateButtonClasses();

  // Listen for checkbox changes
  checkbox.addEventListener("change", updateButtonClasses);
}

