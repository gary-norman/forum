import { activePageElement } from "./main.js";

const angry =
  "background-color: #000000; color: #ff0000; font-weight: bold; border: 2px solid #ff0000; padding: 5px; border-radius: 5px;";
const expect =
  "background-color: #000000; color: #00ff00; font-weight: bold; border: 1px solid #00ff00; padding: 5px; border-radius: 5px;";
const warn =
  "background-color: #000000; color: #e3c144; font-weight: bold; border: 1px solid #e3c144; padding: 5px; border-radius: 5px;";

let filterButtons;

document.addEventListener("newContentLoaded", () => {
  toggleFilters();
});
document.addEventListener("DOMContentLoaded", () => {
  toggleFilters();
});

function toggleFilters() {
  let filterContainer = activePageElement.querySelector(".filters-row");
  if (len.filterContainer > 0) {
    filterButtons = filterContainer.querySelectorAll("button");

    // console.log("%cfilters:", expect, filterButtons);

    filterButtons.forEach((button) => {
      const id = button.getAttribute("popovertarget");
      const popover = activePageElement.querySelector(`#${id}`);

      // console.log("%cid:", expect, id);
      // console.log("%cbutton:", expect, button);
      // console.log("%cpopover:", warn, popover);

      popover.addEventListener("toggle", () => {
        if (popover.matches(":popover-open")) {
          button.classList.add("active");
          // console.log("%cpopover open", warn, popover);
        } else {
          button.classList.remove("active");
          // console.log("%cpopover closed", warn, popover);
        }
      });
    });
  }
}
