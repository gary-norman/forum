
import {navigateToPage} from "./fetch_and_navigate.js";
import {UpdateUI} from "./update_UI_elements.js";
import {fireCalendarListeners} from "./calendar.js";

export let activePage, activePageElement;

// Create a new custom event
export const newContentLoaded = new CustomEvent("newContentLoaded");

document.addEventListener("DOMContentLoaded", (event) => {
  console.log("fired DOMContentLoaded");
  setActivePage("home")
  UpdateUI();
  fireCalendarListeners(activePageElement)
  fireCalendarListeners();
});

document.addEventListener("newContentLoaded", (event) => {
  console.log("fired newContentLoaded");
  UpdateUI();
});


export function setActivePage(dest) {
  activePage = dest + "-page";
  activePageElement = document.querySelector(`#${dest}-page`);
  console.warn("activePageElement:", activePageElement);
}

// variables
//user information
const homePageUserContainer = document.querySelector("#home-page-users");
if (homePageUserContainer !== null) {
  homePageUserContainer.addEventListener("click", (e) => {
    console.log("clicked ", e.target);
    if (e.target.matches(".card")) {
      navigateToPage("user", e.target);
    }
  });
}
