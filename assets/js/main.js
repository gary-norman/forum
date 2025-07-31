import {fetchData} from "./fetch_and_navigate.js";
import { UpdateUI } from "./update_UI_elements.js";
import { fireCalendarListeners } from "./calendar.js";

const angry =
  "background-color: #000000; color: #ea4f92; font-weight: bold; border: 2px solid #ea4f92; padding: .2rem; border-radius: .8rem;";
const expect =
  "background-color: rgb(108 207 93); color: #000000; font-weight: bold; padding: .1rem; border-radius: 1rem;";

export let activePage, activePageElement;

// Create a new custom event
export const newContentLoaded = new CustomEvent("newContentLoaded");

window.addEventListener("popstate", () => {
  const pathParts = window.location.pathname.split("/").filter(Boolean);
  const [entity, id] = pathParts;
  if (entity && id) {
    fetchData(entity.slice(0, -1), id); // Convert 'posts' to 'post'
  }
});

document.addEventListener("DOMContentLoaded", (event) => {
  // console.log("fired DOMContentLoaded");
  setActivePage("home");
  UpdateUI();
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
