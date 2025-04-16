import {
   pages,
  selectActiveFeed,
  data,
  listenToShare,
} from "./share.js";
import {changePage, navigateToPage} from "./fetch_and_navigate.js";
import {UpdateUI} from "./update_UI_elements.js";

export let activePage;

// Create a new custom event
export const newContentLoaded = new CustomEvent("newContentLoaded");

document.addEventListener("DOMContentLoaded", (event) => {
  // console.log("fired DOMContentLoaded");
  UpdateUI();
});

document.addEventListener("newContentLoaded", (event) => {
  // console.log("fired newContentLoaded");
  UpdateUI();
});

export function modifyActivePage(dest) {
  activePage = dest + "-page";
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
