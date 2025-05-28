import { navigateToPage } from "./fetch_and_navigate.js";
import { UpdateUI } from "./update_UI_elements.js";
import { fireCalendarListeners } from "./calendar.js";

const angry =
  "background-color: #000000; color: #ff0000; font-weight: bold; border: 2px solid #ff0000; padding: 5px; border-radius: 5px;";
const expect =
  "background-color: #000000; color: #00ff00; font-weight: bold; border: 1px solid #00ff00; padding: 5px; border-radius: 5px;";

export let activePage, activePageElement;

// Create a new custom event
export const newContentLoaded = new CustomEvent("newContentLoaded");

document.addEventListener("DOMContentLoaded", (event) => {
  console.log("fired DOMContentLoaded");
  setActivePage("home");
  UpdateUI();
  // fireCalendarListeners(activePageElement)
  // fireCalendarListeners();
});

document.addEventListener("newContentLoaded", (event) => {
  console.log("%cfired", expect, "newContentLoaded");
  UpdateUI();
});

export function setActivePage(dest) {
  activePage = dest + "-page";
  activePageElement = document.querySelector(`#${dest}-page`);
  console.info("%cactivePageElement:", expect, activePageElement);
}

// variables
//user information
const homePageUserContainer = document.querySelector("#home-page-users");
if (homePageUserContainer !== null) {
  homePageUserContainer.addEventListener("click", (e) => {
    console.log("%cclicked ", expect, e.target);
    if (e.target.matches(".card")) {
      navigateToPage("user", e.target);
    }
  });
}
