import { fetchData, changePage } from "./fetch_and_navigate.js";
import { UpdateUI } from "./update_UI_elements.js";
import { fireCalendarListeners } from "./calendar.js";
import { selectActiveFeed, goHome } from "./navigation.js";
import { showMainNotification } from "./notifications.js";

const angry =
  "background-color: #000000; color: #ea4f92; font-weight: bold; border: 2px solid #ea4f92; padding: .2rem; border-radius: .8rem;";
const expect =
  "background-color: rgb(108 207 93); color: #000000; font-weight: bold; padding: .1rem; border-radius: 1rem;";

export let activePage, activePageElement;

// Create a new custom event
export const newContentLoaded = new CustomEvent("newContentLoaded");

window.addEventListener("popstate", (event) => {
  if (!event.state || typeof event.state !== "object") {
    // Handle missing or invalid state gracefully
    displayError("Navigation state is missing or invalid.");
    return;
  }
  const { entity, id } = event.state;
  const page = entity + "Page";
  console.info("%cPopped state:", expect, event.state);
  console.info("%centity:", expect, entity);
  console.info("%cid:", expect, id);
  if (entity === "home") {
    goHome();
  } else if (entity && id) {
    try {
      setActivePage(entity);
      changePage(page);
      fetchData(entity, id);
    } catch (err) {
      showMainNotification("Failed to fetch data.");
      // Optionally log error or handle further
    }
  } else {
    // Use a custom error class for HTTP-like errors
    class HttpError extends Error {
      constructor(message, status) {
        super(message);
        this.name = "HttpError";
        this.status = status;
      }
    }
    const error = new HttpError("Invalid URL: entity or id missing", 400);
    showMainNotification(error.message);
    // Optionally log error or handle further
  }
});

document.addEventListener("DOMContentLoaded", (event) => {
  // console.log("fired DOMContentLoaded");
  setActivePage("home");
  UpdateUI();
  selectActiveFeed();
});

document.addEventListener("newContentLoaded", (event) => {
  // console.log("%cfired", expect, "newContentLoaded");
  UpdateUI();
});

export function setActivePage(dest) {
  activePage = dest + "-page";
  activePageElement = document.querySelector(`#${dest}-page`);
  // console.info("%cactivePageElement:", expect, activePageElement);
}

// variables
//user information
// const homePageUserContainer = document.querySelector("#home-page-users");
// if (homePageUserContainer !== null) {
//   homePageUserContainer.addEventListener("click", (e) => {
//     console.log("%cclicked ", expect, e.target);
//     if (e.target.matches(".card")) {
//       navigateToPage("user", e.target);
//     }
//   });
// }
