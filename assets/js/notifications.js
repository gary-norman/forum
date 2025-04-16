// showMainNotification changes the element ID to provide feedback to user
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
// showNotification changes the element ID to provide feedback to user
export function showInlineNotification(elementID, messageOld, messageNew, success) {
    const notification = document.getElementById(elementID);
    notification.textContent = messageNew;
    notification.style.color = "var(--color-hl-green)";
    if (!success) {
        notification.style.color = "var(--color-error)";
        setTimeout(() => {
            notification.textContent = messageOld;
            notification.style.color = "var(--color-fg-1)";
        }, 2500); // Hide after 3 seconds
    }
}