import { listenToReplies, toggleReplyForm } from "./comments.js";
import { listenToShare } from "./share.js";
import { selectActiveFeed } from "./navigation.js";
import { listenToLikeDislike, listenToNoUser } from "./reactions.js";
import { listenToDropdowns } from "./post.js";
import { saveColourScheme } from "./colour_scheme.js";
import { activePage, setActivePage } from "./main.js";
import { getRandomInt } from "./helper_functions.js";
import { listenToRules } from "./channel_rules.js";
import { toggleModals, togglePopovers } from "./popups.js";
import { agreeToJoin, showJoinPopoverRules } from "./join_channel.js";
import {
  displayCalendars,
  getCalendarVars,
  processDateRange,
} from "./calendar.js";

// activity buttons
let actButtonContainer, actButtonsAll;

// activity feeds
let activityFeeds, activityFeedsContentAll;

export function UpdateUI() {
  // console.log("updating UI");
  listenToRules();
  listenToReplies();
  listenToShare();
  listenToLikeDislike();
  listenToNoUser();
  listenToDropdowns();
  listenToPageSetup();
  saveColourScheme();
  updateProfileImages();
  toggleModals();
  resetInputStyle();
  togglePopovers();
  agreeToJoin();
  showJoinPopoverRules();
  // getCalendarVars();
  // displayCalendars();
  // processDateRange();
}

export function updateProfileImages() {
  if (typeof getUserProfileImageFromAttribute === "function") {
    getUserProfileImageFromAttribute();
  } else {
    console.error("getUserProfileImageFromAttribute is not defined");
  }
  if (typeof getInitialFromAttribute === "function") {
    getInitialFromAttribute();
  } else {
    console.error("getInitialFromAttribute is not defined");
  }
}

function resetInputStyle() {
  const modals = document.querySelectorAll(".modal");

  // revert to original state if modals are closed
  modals.forEach((modal) => {
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.attributeName === "style") {
          const currentDisplay = getComputedStyle(modal).display;
          if (currentDisplay === "none") {
            console.log("Modal closed:", modal);
            toggleUserInteracted("remove");
          }
        }
      });
    });

    observer.observe(modal, { attributes: true, attributeFilter: ["style"] });
  });
}

// toggleUserInteracted toggles user-interacted class on input fields to prevent label animation before they are selected
export function toggleUserInteracted(action) {
  //input fields with fancy animations
  const styledInputs = document.querySelectorAll(
    'textarea, input[type="text"], input[type="password"], input[type="email"]',
  );

  styledInputs.forEach((input) => {
    if (action === "add") {
      input.addEventListener("focus", function () {
        this.closest(".input-wrapper").classList.add("user-interacted");
      });
    }
    if (action === "remove") {
      document.querySelectorAll(".input-wrapper").forEach((element) => {
        element.classList.remove("user-interacted");
      });
    }
  });
}

function getUserProfileImageFromAttribute() {
  const userProfileImage = document.querySelectorAll(".profile-pic");

  for (let i = 0; i < userProfileImage.length; i++) {
    let attArr = ["user", "auth", "channel"];
    attArr[0] = userProfileImage[i].getAttribute("data-image-user");
    attArr[1] = userProfileImage[i].getAttribute("data-image-auth");
    attArr[2] = userProfileImage[i].getAttribute("data-image-channel");

    // console.table(attArr)
    if (attArr[0]) {
      // Ensure the `data-image-user` attribute has a value
      userProfileImage[i].style.background =
        `url('${attArr[0]}') no-repeat center`;
      userProfileImage[i].style.backgroundSize = "cover"; // Add `cover` for background sizing
    } else if (attArr[1]) {
      // Ensure the `data-image-auth` attribute has a value
      userProfileImage[i].style.background =
        `url('${attArr[1]}') no-repeat center`;
      userProfileImage[i].style.backgroundSize = "cover"; // Add `cover` for background sizing
    } else if (attArr[2]) {
      // Ensure the `data-image-channel` attribute has a value
      userProfileImage[i].style.background =
        `url('${attArr[2]}') no-repeat center`;
      userProfileImage[i].style.backgroundSize = "cover"; // Add `cover` for background sizing
    } else {
      console.warn(
        "No data-image- attribute value found for element:",
        userProfileImage[i],
      );
    }
  }
}

function getInitialFromAttribute() {
  const userProfileImageEmpty = document.querySelectorAll(
    ".profile-pic--empty",
  );
  const colorsArr = [
    ["var(--clr-accent--blue)", "var(--clr-light--1)"],
    ["var(--clr-accent--green)", "var(--clr-dark--1)"],
    ["var(--clr-accent--orange)", "var(--clr-light--1)"],
    ["var(--clr-accent--pink)", "var(--clr-light--1)"],
    ["var(--clr-accent--yellow)", "var(--clr-dark--1)"],
    ["var(--clr-accent--red)", "var(--clr-light--1)"],
  ];

  function getConsistentThemeIndex(name) {
    let hash = 0;
    for (let i = 0; i < name.length; i++) {
      hash = name.charCodeAt(i) + ((hash << 5) - hash);
    }
    return Math.abs(hash) % colorsArr.length;
  }

  for (const el of userProfileImageEmpty) {
    const nameUser = el.getAttribute("data-name-user");
    const nameSidebar = el.getAttribute("data-name-user-sidebar");
    const nameChannel = el.getAttribute("data-name-channel");

    const name = nameUser || nameSidebar || nameChannel;

    if (name) {
      const theme = getConsistentThemeIndex(name);
      const [bg, fg] = colorsArr[theme];
      el.style.background = bg;
      el.style.color = fg;

      if (nameUser) {
        el.style.fontSize = "2rem";
      } else if (nameSidebar) {
        el.style.fontSize = "5rem";
      }

      el.setAttribute("data-initial", Array.from(name)[0]);
    } else {
      console.warn("No data-name- attribute value found for element:", el);
    }
  }
}

// INFO was a DOMContentLoaded function
export function listenToPageSetup() {
  // feeds
  const feeds = document.querySelectorAll(".feeds-wrapper");

  // feeds.forEach((feed) => {
  //     feed.classList.add("feeds-wrapper-loaded");
  // });

  toggleUserInteracted("add");
  if (activePage === "user-page") {
    actButtonContainer = document.querySelector("#activity-bar");
    actButtonsAll = actButtonContainer.querySelectorAll("button");
    activityFeeds = document.querySelector("#user-activity-feeds");
    activityFeedsContentAll = activityFeeds.querySelectorAll(
      '[id^="activity-feed-"]',
    );
    // Toggle the feed based on the chosen activity button
    actButtonsAll.forEach((button) =>
      button.addEventListener("click", (e) => {
        toggleFeed(
          document.getElementById("activity-" + e.target.id),
          document.getElementById("activity-feed-" + e.target.id),
          e.target,
        );
        console.log("activity-" + e.target.id);
      }),
    );
  }
}

// toggle the various feeds on user-page
function toggleFeed(targetFeed, targetFeedContent, targetButton) {
  const timeOut = 400;
  // const allFeedsExceptTarget = Array.from(activityFeedsContentAll).filter(feed => feed.id !== targetFeedContent.id);

  actButtonsAll.forEach((button) => button.classList.remove("btn-active"));
  activityFeedsContentAll.forEach((feed) => {
    feed.classList.remove("collapsible-expanded");
    feed.classList.add("collapsible-collapsed");
    setTimeout(() => {
      feed.parentElement.classList.add("collapsible-collapsed");
    }, timeOut);
  });
  setTimeout(() => {
    targetFeedContent.classList.remove("collapsible-collapsed");
    targetFeedContent.classList.add("collapsible-expanded");
    targetFeedContent.parentElement.classList.remove("collapsible-collapsed");
    selectActiveFeed();
    targetButton.classList.toggle("btn-active");
  }, timeOut);
  setTimeout(() => {
    Array.from(activityFeedsContentAll).forEach((feed) => {
      const filtersRow = feed.parentElement.querySelector(".filters-row");
      if (filtersRow) {
        filtersRow.classList.add("hide-feed");
      }
    });
    targetFeed.querySelector(".button-row").classList.remove("hide-feed");
  }, timeOut);
}
