import { activePageElement } from "./main.js";
// showMainNotification changes the main notification modal to provide feedback to user
export function showMainNotification(message, timeout = 2500) {
  const notification = document.getElementById("notification-main");
  const notificationContent = document.getElementById(
    "notification-main-content",
  );
  notificationContent.textContent = message;
  notification.style.display = "flex";
  setTimeout(() => {
    notification.style.display = "none";
  }, timeout); // Hide after 3 seconds
}
// showNotification changes the given element to provide feedback to user
export function showInlineNotification(
  notification,
  messageFail,
  messageSuccess,
  success,
) {
  // const notification = activePageElement.getElementById(elementID);
  notification.textContent = messageSuccess;
  notification.style.color = "var(--clr-accent--1)";
  if (!success) {
    notification.style.color = "var(--clr--error)";
    setTimeout(() => {
      notification.textContent = messageFail;
      notification.style.color = "var(--clr-fg-1)";
    }, 2500); // Hide after 3 seconds
  }
}
