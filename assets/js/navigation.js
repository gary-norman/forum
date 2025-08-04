import { setActivePage, activePageElement, newContentLoaded } from "./main.js";
import { changePage, navigateToPage } from "./fetch_and_navigate.js";
import { toggleReplyForm } from "./comments.js";
import { data } from "./share.js";
import { UpdateUI } from "./update_UI_elements.js";

// sidebar butons
const goHome = document.querySelector("#btn-go-home");

// --- go home ---
goHome.addEventListener("click", (e) => {
  setActivePage("home");
  history.pushState({}, "", `/`);
  changePage(data["homePage"]);
  // navigateToPage("home", e.target)
});

// INFO was inside a DOMContentLoaded function
export function listenToChannelLinks() {
  const joinedAndOwnedChannelContainer = document.querySelector(
    "#sidebar-channel-block",
  );

  // console.log(joinedAndOwnedChannelContainer);
  let joinedAndOwnedChannels;

  if (joinedAndOwnedChannelContainer) {
    joinedAndOwnedChannels =
      joinedAndOwnedChannelContainer.querySelectorAll(".sidebar-channel");
    joinedAndOwnedChannels.forEach((channel) =>
      channel.addEventListener("click", (e) => {
        e.preventDefault();
        navigateToPage("channel", channel);
      }),
    );
  }
}

function navigateToEntity(e) {
  console.info("navigateToEntity called");
  const target = e.target;
  const isButton = target.nodeName.toLowerCase() === "button";
  const commentAction = target.getAttribute("data-action");
  const dest = target.getAttribute("data-dest");
  const postID = target.getAttribute("data-post-id");

  if (commentAction === "navigate--comment-post" && dest) {
    console.info("Navigating to comment post:", dest);
    navigateToPage(dest, target)
      .then(() => {
        toggleReplyForm(target);
      })
      .catch((err) => {
        console.error("Navigation failed:", err);
      });
    return;
  }

  if (commentAction === "comment-post") {
    console.info("Toggling reply form for post ID:", postID);
    toggleReplyForm(target);
    return;
  }

  // Handle navigation for elements with data-dest or .link class
  let navTarget = target.closest(".link");
  if (navTarget && navTarget.getAttribute("data-dest")) {
    navigateToPage(navTarget.getAttribute("data-dest"), navTarget);
  }
}

document.addEventListener("click", navigateToEntity, { capture: false });

// addGlobalEventListener is an event listener on a parent element that only runs your callback
// when the event happens on a specific type of child element (that matches the selector you give)
// good for event delegation and for new elements populated dynamically
export function addGlobalEventListener(type, selector, callback, parent) {
  parent.addEventListener(type, (e) => {
    if (e.target.matches(selector)) {
      callback(e);
    }
  });
}
