import { activePageElement } from "./main.js";
import {
  applyCustomTheme,
  showThemesClickable,
  pickTheme,
  themeList,
} from "./consoleThemes.js";

let filterButtons;
const currentDate = new Date();
const formattedDate = currentDate.toISOString().split("T")[0];
let allPageCards = [];

document.addEventListener("newContentLoaded", () => {
  if (activePageElement.id !== "post-page") {
    listenToClicksForPopovers();
    toggleFilters();
    setAllPosts();
    filterContent();
  }
});

document.addEventListener("DOMContentLoaded", () => {
  if (activePageElement.id !== "post-page") {
    listenToClicksForPopovers();
    toggleFilters();
    setAllPosts();
    filterContent();
  }
});

function setAllPosts() {
  if (activePageElement.id !== "user-page") {
    allPageCards = Array.from(activePageElement.querySelectorAll(".card.link"));
  } else {
    allPageCards = Array.from(
      activePageElement
        .querySelector(".collapsible-expanded")
        .querySelectorAll(".card.link"),
    );
  }
  // console.log("%csetAllPosts:", expect, allPageCards);
  // console.log("%cactivePageElement:", expect, activePageElement);
}

// compareDates used to compare 2 dates, which then is used in the default js sort function
function compareDates(a, b) {
  const dateA = new Date(a.dataset.createdAt);
  const dateB = new Date(b.dataset.createdAt);
  return dateA - dateB;
}

// compareNumberValues used to compare 2 numbers, which then is used in the default js sort function
function compareNumberValues(sortBy) {
  return function (a, b) {
    let valueA, valueB;

    switch (sortBy) {
      case "likes":
        valueA = Number(a.querySelector(".btn-likes").textContent);
        valueB = Number(b.querySelector(".btn-likes").textContent);
        break;
      case "dislikes":
        valueA = Number(a.querySelector(".btn-dislikes").textContent);
        valueB = Number(b.querySelector(".btn-dislikes").textContent);
        break;
      case "comments":
        valueA = Number(a.querySelector(".btn-reply").textContent);
        valueB = Number(b.querySelector(".btn-reply").textContent);
        break;
      default:
        valueA = 0;
        valueB = 0;
    }

    return valueB - valueA;
  };
}

// compareActivity used to compare 2 dates per element, and chooses the most recent date for each
// This is then used in the default js sort function
function compareActivity(a, b) {
  const createdDateA = new Date(a.dataset.createdAt);
  const createdDateB = new Date(b.dataset.createdAt);
  const lastReactedDateA = new Date(a.dataset.lastReaction);
  const lastReactedDateB = new Date(b.dataset.lastReaction);

  // Pick the later date for a
  const latestDateA =
    lastReactedDateA > createdDateA ? lastReactedDateA : createdDateA;

  // Pick the later date for b
  const latestDateB =
    lastReactedDateB > createdDateB ? lastReactedDateB : createdDateB;

  // Return difference for sort
  return latestDateA - latestDateB;
}

// sortDescending sort by oldest date / most amount first
function sortDescending(func) {
  return allPageCards.sort(func);
}

// sortAscending sort by newest date / least amount first
function sortAscending(func) {
  return allPageCards.sort(func).reverse();
}

function filterContent() {
  const filtersRow = activePageElement.querySelector(".filters-row");
  const channelCheckboxes = filtersRow
    .querySelector(`[id$="dropdown-channel"]`)
    .querySelectorAll('input[type="checkbox"]:checked');
  const reactionCheckboxes = filtersRow
    .querySelector(`[id$="dropdown-reaction"]`)
    .querySelectorAll('input[type="checkbox"]:checked');
  // TODO commenting out type as not needed for base audit
  // const typeCheckboxes = filtersRow.querySelector(`[id$="dropdown-type"]`).querySelectorAll('input[type="checkbox"]:checked');
  const commentRadios = filtersRow
    .querySelector(`[id$="dropdown-comments"]`)
    .querySelectorAll('input[type="radio"]:checked');
  const sortRadios = filtersRow
    .querySelector(`[id$="dropdown-sort"]`)
    .querySelectorAll('input[type="radio"]:checked');
  const startDateInput = filtersRow.querySelector(
    `[id$="dropdown-date"] input[name="creation-year-start"]`,
  );
  const endDateInput = filtersRow.querySelector(
    `[id$="dropdown-date"] input[name="creation-year-end"]`,
  );

  // Build active filters
  const activeFilters = {
    channels: Array.from(channelCheckboxes).map((cb) => cb.value),
    reactions: Array.from(reactionCheckboxes).map((rb) => rb.value),
    // TODO commenting out type as not needed for base audit
    // types: Array.from(typeCheckboxes).map(cb => cb.value.slice(5)), //slice to remove "type-" before type of card from the template
    comments: commentRadios.length > 0 ? commentRadios[0].value : null,
    sort: sortRadios.length > 0 ? sortRadios[0].value : null,
    startDate: startDateInput?.value ? new Date(startDateInput.value) : null,
    endDate: endDateInput?.value ? new Date(endDateInput.value) : null,
  };

  // Sort
  // Default sort
  if (activeFilters.sort === null) {
    allPageCards = sortAscending(compareDates);
    reorderVisiblePosts();
  } else {
    // sort by date
    //sort newest first
    if (activeFilters.sort.includes("most-new")) {
      allPageCards = sortAscending(compareDates);
    }

    //sort oldest first
    if (activeFilters.sort.includes("most-old")) {
      allPageCards = sortDescending(compareDates);
    }

    //sort by recent activity
    //sort newest first
    if (activeFilters.sort.includes("most-new-activity")) {
      allPageCards = sortAscending(compareActivity);
    }

    //sort oldest first
    if (activeFilters.sort.includes("most-old-activity")) {
      allPageCards = sortDescending(compareActivity);
    }

    //sort by likes
    //sort most first
    if (activeFilters.sort.includes("most-likes")) {
      allPageCards = sortDescending(compareNumberValues("likes"));
    }

    //sort least first
    if (activeFilters.sort.includes("least-likes")) {
      allPageCards = sortAscending(compareNumberValues("likes"));
    }

    //sort by dislikes
    //sort most first
    if (activeFilters.sort.includes("most-dislikes")) {
      allPageCards = sortDescending(compareNumberValues("dislikes"));
    }

    //sort least first
    if (activeFilters.sort.includes("least-dislikes")) {
      allPageCards = sortAscending(compareNumberValues("dislikes"));
    }

    //sort by comments
    //sort most first
    if (activeFilters.sort.includes("most-comments")) {
      allPageCards = sortDescending(compareNumberValues("comments"));
    }

    //sort least first
    if (activeFilters.sort.includes("least-comments")) {
      allPageCards = sortAscending(compareNumberValues("comments"));
    }

    reorderVisiblePosts();
  }

  allPageCards.forEach((card) => {
    const cardChannel = card.dataset.channelId;
    // TODO commenting out type as not needed for base audit
    // const cardType = card.dataset.dest;
    const cardCommentsCount = Number(
      card.querySelector(".btn-reply").textContent,
    );
    const cardLikesCount = Number(card.querySelector(".btn-likes").textContent);
    const cardDislikesCount = Number(
      card.querySelector(".btn-dislikes").textContent,
    );
    const cardCreatedDate = new Date(card.dataset.createdAt);

    let visible = true;

    // Filter by channel
    if (
      activeFilters.channels.length > 0 &&
      !activeFilters.channels.includes(cardChannel)
    ) {
      visible = false;
    }

    // TODO commenting out type as not needed for base audit
    // // Filter by type
    // if (activeFilters.types.length > 0 && !activeFilters.types.includes(cardType)) {
    //     visible = false;
    // }

    // Filter by comments
    if (activeFilters.comments !== null) {
      let commentMatch = false;

      if (
        activeFilters.comments.includes("has-comments") &&
        cardCommentsCount > 0
      ) {
        commentMatch = true;
      }

      if (
        activeFilters.comments.includes("no-comments") &&
        cardCommentsCount <= 0
      ) {
        commentMatch = true;
      }

      if (commentMatch === false) {
        visible = false;
      }
    }

    // Filter by reaction
    if (activeFilters.reactions.length > 0) {
      let reactionMatch = false;

      if (activeFilters.reactions.includes("liked") && cardLikesCount > 0) {
        reactionMatch = true;
      }

      if (
        activeFilters.reactions.includes("disliked") &&
        cardDislikesCount > 0
      ) {
        reactionMatch = true;
      }

      if (
        activeFilters.reactions.includes("no-reaction") &&
        cardLikesCount === 0 &&
        cardDislikesCount === 0
      ) {
        reactionMatch = true;
      }

      if (reactionMatch === false) {
        visible = false;
      }
    }

    // Filter by date range
    if (activeFilters.startDate !== null && activeFilters.endDate !== null) {
      //set the time as the end of the day, because it was counting the date at 1:00:00
      // and wouldn't include dates within the day
      activeFilters.endDate.setHours(23, 59, 59, 999);

      if (cardCreatedDate <= activeFilters.startDate) {
        visible = false;
      }
      if (cardCreatedDate >= activeFilters.endDate) {
        visible = false;
      }
    }

    // Show or hide card, by hiding the container holding it
    card.parentElement.classList.toggle("hide", !visible);
  });
}

function reorderVisiblePosts() {
  const feedContainer = activePageElement.querySelector(`[id$="feed"]`);
  // console.log("%cfeedContainer:", expect, feedContainer);
  // console.log("%cactivePageElement:", expect, activePageElement);
  // console.log("%callPageCards:", expect, allPageCards);

  // Remove existing posts from DOM
  allPageCards.forEach((post) => {
    // console.log("%cpost.parentElement: ", warn, post.parentElement)
    if (post.parentElement) {
      feedContainer.removeChild(post.parentElement);
    }
  });

  // Append posts back in sorted order
  allPageCards.forEach((post) => {
    // console.log("appending")
    feedContainer.appendChild(post.parentElement);
  });
}

function toggleFilters() {
  let filterContainer = activePageElement.querySelector(".filters-row");
  filterButtons = filterContainer.querySelectorAll("button.btn-action");
  let startDateInteracted = false;

  filterButtons.forEach((button) => {
    const id = button.getAttribute("popovertarget");
    const popover = activePageElement.querySelector(`#${id}`);
    const clearButton = popover.querySelector(".clear-choices ");

    button.addEventListener("click", (e) => {
      const noneSelected = checkSelectedInputs(popover);

      if (!popover.matches(":popover-open") && button.contains(e.target)) {
        toggleFilterButtonState(button, noneSelected, "button");
      }
    });

    if (popover && popover.matches("[popover].wrapper-filter-dropdown")) {
      popover.querySelectorAll(".dropdown-option").forEach(function (option) {
        // store handler on the element itself to reuse later
        if (!option._handler) {
          option._handler = optionHandlerFactory(popover, button);
          option.addEventListener("click", option._handler);
        }
      });
    }

    if (popover && popover.matches("[popover].card.date")) {
      const startInput = popover.querySelector(
        'input[id^="creation-year-start"]',
      );
      const endInput = popover.querySelector('input[id^="creation-year-end"]');

      popover.querySelectorAll("input").forEach(function (input) {
        input.addEventListener("input", (e) => {
          if (startInput.value !== "" && startDateInteracted === false) {
            endInput.value = formattedDate;
          }
          const noneSelected = checkSelectedInputs(popover);

          filterContent();

          toggleFilterButtonState(button, noneSelected, "inside");

          if (noneSelected) {
            toggleClearButton(popover, "hide");
          } else {
            toggleClearButton(popover, "show");
          }
        });

        input.addEventListener("blur", () => {
          startDateInteracted = true;
        });
      });
    }

    if (clearButton) {
      clearButton.addEventListener("click", () => {
        let noneSelected;
        popover.querySelectorAll(".dropdown-option").forEach(function (option) {
          const checkbox = option.querySelector("input[type='checkbox']");
          const radio = option.querySelector("input[type='radio']");
          let input;
          if (radio) {
            input = radio;
          } else if (checkbox) {
            input = checkbox;
          }
          input.checked = false;

          noneSelected = checkSelectedInputs(popover);
        });

        if (popover && popover.matches("[popover].card.date")) {
          popover.querySelectorAll("input").forEach(function (input) {
            input.value = "";

            noneSelected = checkSelectedInputs(popover);

            startDateInteracted = false;
          });
        }

        toggleClearButton(popover, "hide");
        toggleFilterButtonState(button, noneSelected, "inside");
        filterContent();
        clearRadios(popover, button);
      });
    }
  });
}

function optionHandlerFactory(popover, button) {
  return function (e) {
    console.custom.info("REATTACHED:", e.currentTarget);
    handleOptionClick(e, popover, button);
  };
}

function handleOptionClick(e, popover, button) {
  const option = e.currentTarget;

  let input;
  const checkbox = option.querySelector("input[type='checkbox']");
  const radio = option.querySelector("input[type='radio']");
  let optionText = option.querySelector("label").innerText;

  if (checkbox) {
    input = checkbox;
  } else if (radio) {
    input = radio;
    if (optionText === "has no comments") {
      optionText = optionText.slice(4);
    }
    clearRadios(popover, button);

    option.classList.add("selected");
    button.querySelector("span").innerText = optionText;
  }

  input.click();

  const noneSelected = checkSelectedInputs(popover);

  filterContent();

  if (noneSelected) {
    toggleClearButton(popover, "hide");
  } else {
    toggleClearButton(popover, "show");
  }
  toggleFilterButtonState(button, noneSelected, "inside");
}

function clearRadios(popover, button) {
  const allOptions = popover.querySelectorAll(".dropdown-option");

  allOptions.forEach((o) => {
    o.classList.remove("selected");
  });

  const buttonLabel = button.getAttribute("aria-label");
  button.querySelector("span").innerText = buttonLabel.slice(0, -7);
}

function listenToClicksForPopovers() {
  const popovers = activePageElement.querySelectorAll(
    "[popover].filter-popover",
  );
  let filterContainer = activePageElement.querySelector(".filters-row");
  filterButtons = filterContainer.querySelectorAll("button.btn-action");

  if (popovers) {
    filterButtons.forEach((button) => {
      button.addEventListener("click", (e) => {
        document.addEventListener("click", handleClicksOutsidePopovers, true);
      });
    });
  }

  function handleClicksOutsidePopovers(e) {
    popovers.forEach((popover) => {
      if (!popover.contains(e.target) && popover.matches(":popover-open")) {
        const button = activePageElement.querySelector(
          `[popovertarget="${popover.id}"]`,
        );
        const noneSelected = checkSelectedInputs(popover);

        toggleFilterButtonState(button, noneSelected, "outside");
        document.removeEventListener("click", handleClicksOutsidePopovers);
      }
    });
  }
}

function toggleFilterButtonState(button, noneSelected, clickLocation) {
  const id = button.getAttribute("popovertarget");
  const popover = activePageElement.querySelector(`#${id}`);
  const cancelButton = popover.querySelector(".clear-choices ");

  if (clickLocation === "button" && noneSelected) {
    button.classList.add("active");
  }

  if (popover && !popover.matches(":popover-open")) {
    return;
  }

  if (clickLocation === "outside" && noneSelected) {
    button.querySelector("span").classList.add("btn-filters");
    button.querySelector("span").classList.remove("btn-active-filter");
    button.classList.remove("selected");
    button.classList.remove("active");
    popover.hidePopover();
    return;
  }

  if (noneSelected) {
    button.querySelector("span").classList.add("btn-filters");
    button.querySelector("span").classList.remove("btn-active-filter");
    button.classList.remove("selected");
    button.classList.add("active");

    if (cancelButton && clickLocation === "outside") {
      button.classList.remove("active");
    }
  } else {
    button.querySelector("span").classList.remove("btn-filters");
    button.querySelector("span").classList.add("btn-active-filter");
    button.classList.add("selected");
    if (clickLocation === "outside") {
      popover.hidePopover();
    }
  }
}

function toggleClearButton(popover, state) {
  const button = popover.querySelector(".clear-choices ");
  const scrollContainer = popover.querySelector(".container-filter");
  if (button) {
    if (state === "show") {
      scrollContainer.style.maxHeight = `calc(28rem - 0.8rem - 4.4rem)`;
      button.classList.remove("hide");
    } else {
      scrollContainer.style.maxHeight = `calc(28rem - 0.8rem)`;
      button.classList.add("hide");
    }
  }
}

function checkSelectedInputs(popover) {
  //update state of checkboxes
  const checkboxes = popover.querySelectorAll("input[type='checkbox']");
  const radios = popover.querySelectorAll("input[type='radio']");
  const dates = popover.querySelectorAll("input[type='date']");

  if (checkboxes.length > 0) {
    return Array.from(checkboxes).every((checkbox) => !checkbox.checked);
  } else if (radios.length > 0) {
    return Array.from(radios).every((radio) => !radio.checked);
  } else if (dates.length > 0) {
    return [...dates].every((date) => !date.value);
  }
  return console.custom.info("no checkboxes, radios or dates");
}
