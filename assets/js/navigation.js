import {setActivePage, newContentLoaded, activePage, activePageElement} from "./main.js";
import { changePage, navigateToPage } from "./fetch_and_navigate.js";
import { toggleReplyForm } from "./comments.js";
import { data } from "./share.js";

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
  // debugger;
  const hasDestination = e.target.getAttribute("data-dest");
  const isButton = e.target.nodeName.toLowerCase() === "button";
  const commentAction = e.target.getAttribute("data-action");
  const postID = e.target.getAttribute("data-post-id");

  if (hasDestination || !isButton || (isButton && hasDestination)) {
    if (commentAction === "navigate--comment-post") {
      const dest = e.target.getAttribute("data-dest");
      navigateToPage(dest, e.target).then(() => {
        toggleReplyForm();
      });
    } else if (commentAction === "comment-post") {
      toggleReplyForm();
    } else if (hasDestination || !isButton || (isButton && hasDestination)) {
      const parent = e.target.closest(".link");
      if (parent) {
        const dest = parent.getAttribute("data-dest");
        navigateToPage(dest, parent);
      }
    } else if (e.target.matches(".link")) {
      const dest = e.target.getAttribute("data-dest");
      navigateToPage(dest, e.target);
    }
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
