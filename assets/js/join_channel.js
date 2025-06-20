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

  if (checkbox === null || button === null) {
    console.warn("Checkbox or button not found. Skipping agreeToJoin setup.");
    return;
  }

  const updateButtonClasses = () => {
    if (checkbox.checked) {
      button.classList.add("btn-primary", "btn-action-primary");
      button.classList.remove("btn-filled-nohover");
    } else {
      button.classList.remove("btn-primary", "btn-action-primary");
      button.classList.add("btn-filled-nohover");
    }
  };

  // Run once on load in case checkbox is pre-checked
  updateButtonClasses();

  // Listen for checkbox changes
  checkbox.addEventListener("change", updateButtonClasses);
}

export function showJoinPopoverRules() {
  const showRules = document.querySelector('[id^="showRules"]');
  const rulesContainer = document.querySelector('[id^="rulesContainer"]');

  if (showRules === null || rulesContainer === null) {
    console.warn(
      "Show rules button or rules container not found. Skipping showJoinPopoverRules setup.",
    );
    return;
  }

  showRules.addEventListener("click", function (event) {
    event.preventDefault();
    rulesContainer.classList.toggle("hidden");
  });
}
