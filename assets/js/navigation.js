import {
  setActivePage,
  newContentLoaded,
  activePage,
  activePageElement,
} from "./main.js";
import { changePage, navigateToPage } from "./fetch_and_navigate.js";
import { toggleReplyForm } from "./comments.js";
import { data } from "./share.js";
import { showInlineNotification } from "./notifications.js";
// sidebar butons
const goHome = document.querySelector("#btn-go-home");
export let scrollWindow;

document.addEventListener("newContentLoaded", () => {
  listenToChannelLinks();
});

document.addEventListener("DOMContentLoaded", () => {
  listenToChannelLinks();
});

// --- go home ---
goHome.addEventListener("click", (e) => {
  setActivePage("home");
  history.pushState({}, "", `/`);
  changePage(data["homePage"]);
  // navigateToPage("home", e.target)
});

// INFO was inside a DOMContentLoaded function
function listenToChannelLinks() {
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
  let target = e.target;
  const isButton = target.nodeName.toLowerCase() === "button";
  const commentAction = target.getAttribute("data-action");
  const dest = target.getAttribute("data-dest");
  const postID = target.getAttribute("data-post-id");

  if (commentAction === "navigate--comment-post" && dest) {
    const sidebarProfile = document.querySelector(".sidebarProfile");
    console.info("Navigating to comment post:", dest);
    navigateToPage(dest, target)
      .then(() => {
        target = activePageElement.querySelector(`#post-card-${postID}`);
        if (sidebarProfile) {
          toggleReplyForm(target);
        } else {
          const notifier = activePageElement.querySelector(
            `#post-card-info-${postID}`,
          );
          console.info("notifier:", notifier);
          showInlineNotification(
            notifier,
            "",
            "You must be logged in to reply to a post.",
            false,
          );
        }
      })
      .catch((err) => {
        console.error("Navigation failed:", err);
      });
    return;
  }

  if (commentAction === "comment-post") {
    const sidebarProfile = document.querySelector(".sidebarProfile");
    console.info("Toggling reply form for post ID:", postID);
    if (sidebarProfile) {
      toggleReplyForm(target);
    } else {
      const notifier = activePageElement.querySelector(
        `#post-card-info-${postID}`,
      );
      console.info("notifier:", notifier);
      showInlineNotification(
        notifier,
        "",
        "You must be logged in to reply to a post.",
        false,
      );
    }
    return;
  }

  // Handle navigation for elements with data-dest or .link class
  let navTarget = target.closest(".link");
  if (navTarget && navTarget.getAttribute("data-dest")) {
    navigateToPage(navTarget.getAttribute("data-dest"), navTarget);
  }
}

export function selectActiveFeed() {
  const angry =
    "background-color: #000000; color: #ea4f92; font-weight: bold; border: 2px solid #ea4f92; padding: .2rem; border-radius: .8rem;";
  switch (activePage) {
    case "home-page":
      console.log("%cON HOME PAGE BITCH", angry);
      scrollWindow = document.querySelector(`#home-feed`);
      // console.log("scrollWindow:", scrollWindow)
      break;
    case "user-page":
      console.log("%cON USER PAGE BITCH", angry);
      const userFeeds = Array.from(
        activePageElement.querySelectorAll('[id$="-feed"]'),
        // activePageElement.querySelector('[id="user-activity-feeds"]').querySelectorAll('[id$="-feed"]'),
      );
      const activeFeed = userFeeds.find((feed) =>
        feed.classList.contains("collapsible-expanded"),
      );

      scrollWindow = activeFeed;
      // console.log(scrollWindow)
      break;
    case "channel-page":
      console.log("%cON CHANNEL PAGE BITCH", angry);

      scrollWindow = document.querySelector(`#channel-feed`);
      // console.log(scrollWindow)
      break;
    case "post-page":
      console.log("%cON POST PAGE BITCH", angry);

      // console.log(scrollWindow)
      break;
    default:
      console.log("%cNO ACTIVE FEED BITCH... switching to home-page", angry);
      setActivePage("home");
      scrollWindow = document.querySelector(`#home-feed`);
      break;
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
