import {
  setActivePage,
  newContentLoaded,
  activePage,
  activePageElement,
} from "./main.js";
import { changePage, navigateToPage, fetchHome } from "./fetch_and_navigate.js";
import { toggleReplyForm } from "./comments.js";
import { pageData } from "./share.js";
import { showInlineNotification } from "./notifications.js";
import {
  applyCustomTheme,
  showThemesClickable,
  pickTheme,
  themeList,
} from "./consoleThemes.js";

// sidebar butons
let goHomeBtns;
export let scrollWindow;

document.addEventListener("newContentLoaded", () => {
  listenToChannelLinks();
});

document.addEventListener("DOMContentLoaded", () => {
  listenToChannelLinks();
  goHomeBtns = document.querySelectorAll(".btn-go-home");
  console.log("goHomeBtns:", goHomeBtns)
  goHomeBtns.forEach((button) =>
      button.addEventListener("click", () => {
        window.location.href = '/';
      }),
  );
});

// --- go home ---
export function goHome() {
  const stateObj = { entity: "home", id: "home" };
  history.pushState(stateObj, "", `/`);
  setActivePage("home");
  changePage(pageData["homePage"]);
  // navigateToPage("home", e.target)
}



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
        pageData["homePage"].innerHTML = "";
        target = activePageElement.querySelector(`#post-card-${postID}`);
        if (sidebarProfile) {
          toggleReplyForm(target);
        } else {
          return;
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
      return;
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

  // // Handle navigation for elements with data-dest or .link class

  if (dest || !isButton || (isButton && dest)) {
    const parent = e.target.closest(".link");
    if (parent) {
      const newDest = parent.getAttribute("data-dest");
      navigateToPage(newDest, parent);
    }
    return;
  }
  if (e.target.matches(".link")) {
    const newDest = e.target.getAttribute("data-dest");
    navigateToPage(newDest, e.target);
  }
}

export function selectActiveFeed() {
  switch (activePage) {
    case "home-page":
      console.custom.angryinfo("ON HOME PAGE");
      scrollWindow = document.querySelector(`#home-feed`);
      document.dispatchEvent(newContentLoaded);
      // console.log("scrollWindow:", scrollWindow)
      break;
    case "user-page":
      console.custom.angryinfo("ON USER PAGE");
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
      console.custom.angryinfo("ON CHANNEL PAGE");

      scrollWindow = document.querySelector(`#channel-feed`);
      // console.log(scrollWindow)
      break;
    case "post-page":
      console.custom.angryinfo("ON POST PAGE");

      // console.log(scrollWindow)
      break;
    default:
      console.custom.angryinfo("NO ACTIVE FEED... switching to home-page");
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
