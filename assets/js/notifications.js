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
  }, timeout);
}
// showNotification changes the given element to provide feedback to user. notification is the element that accepts the message, originalText returns any pre-existing (eg instructional) text, messageSuccess is the text to show when the action is successful, success is a boolean indicating if the action was successful, and timeout is how long to show the notification for.
export function showInlineNotification(
  notification,
  originalText,
  messageSuccess,
  success,
  timeout = 2500,
) {
  // const notification = activePageElement.getElementById(elementID);
  notification.textContent = messageSuccess;
  notification.style.color = "var(--clr-accent--1)";
  if (!success) {
    notification.style.color = "var(--clr--error)";
    setTimeout(() => {
      notification.textContent = originalText;
      notification.style.color = "var(--clr-fg-1)";
    }, timeout);
  }
}
