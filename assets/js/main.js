import {
  channelPage,
  homePage,
  pages,
  selectActiveFeed,
  userPage,
  postPage,
} from "./share.js";

// variables
//user information

// dark mode
const switchDl = document.querySelector("#switch-dl");
const darkSwitch = document.querySelector("#sidebar-option-darkmode");
// activity buttons
let actButtonContainer;
let actButtonsAll;
// activity feeds
let activityFeeds;
let activityFeedsContentAll;
// sidebar butons
const goHome = document.querySelector("#btn-go-home");
// right panel buttons
let rightPanelButtons;
// sidebar elements
const userProfileImage = document.querySelectorAll(".profile-pic");
const userProfileImageEmpty = document.querySelectorAll(".profile-pic--empty");
const sidebarOption = document.querySelector("#sidebar-options");
const sidebarOptionsList = document.querySelector(".sidebar-options-list");
// login/register buttons
// TODO overhaul the naming of these buttons
const registerForm = document.querySelector("#register-form");
const loginTitle = document.querySelector("#login-title");
const loginForm = document.querySelector("#login-form");
const loginFormButton = document.querySelector("#login");
const loginFormUser = document.querySelector("#login_username");
const loginFormPassword = document.querySelector("#login_password");
const logoutFormButton = document.querySelector("#logout");
const btnsLogin = document.querySelectorAll('[id^="btn_login-"]');
const btnsRegister = document.querySelectorAll('[id^="btn_register-"]');
const btnsForgot = document.querySelector("#btn_forgot");
// join channel
const joinChannelButton = document.querySelector("#join-channel-btn");
//input fields with fancy animations
const styledInputs = document.querySelectorAll(
  'textarea, input[type="text"], input[type="password"], input[type="email"]',
);
// modals
const modals = document.querySelectorAll(".modal");
// popovers
const popovers = document.querySelectorAll("[popover]");
// login/register forms
const formLogin = document.querySelector("#form-login");
const formRegister = document.querySelector("#form-register");
const formForgot = document.querySelector("#form-forgot");
const formEditUser = document.querySelector("#form-edit-user");
const formAccSettings = document.querySelector("#form-acc-settings");
const formViewStats = document.querySelector("#form-view-stats");
const formRemoveAcc = document.querySelector("#form-remove-acc");
let forgotVisible = false;
// login/register modal
const loginHeader = document.querySelector("#modal-header-logreg");
const loginModal = document.querySelector("#container-form-login");
// const editUserModal = document.querySelector('#container-form-edit-user');
const accSettingsModal = document.querySelector("#container-form-acc-settings");
const viewStatsModal = document.querySelector("#container-form-view-stats");
const removeAccModal = document.querySelector("#container-form-remove-acc");
// Get the buttons that open the modals
const openLoginModal = document.querySelector("#btn-open-login-modal");
const openLoginModalFallback = document.querySelector(
  "#btn-open-login-modal-fallback",
);
// const openEditUserModal = document.querySelector('#btn-open-edit-user-modal');
// const openAccSettingsModal = document.querySelector('#btn-open-acc-settings-modal');
// const openViewStatsModal = document.querySelector('#btn-open-view-stats-modal');
// const openRemoveAccModal = document.querySelector('#btn-open-remove-acc-modal');
// Get the <span> element that closes the modal
const closeLoginModal = loginModal.getElementsByClassName("close")[0];
// const closeEditUserModal = editUserModal.getElementsByClassName("close")[0];
// const closeAccSettingsModal = accSettingsModal.getElementsByClassName("close")[0];
// const closeViewStatsModal = viewStatsModal.getElementsByClassName("close")[0];
// const closeRemoveAccModal = removeAccModal.getElementsByClassName("close")[0];
// registration form
const regForm = document.querySelector("#form-register");
const regFormInputs = regForm.querySelectorAll("input");
const regFormSpans = regForm.querySelectorAll("span");
const regFormIcons = regForm.querySelectorAll(".validation-icon");
const regFormTooltips = regForm.querySelectorAll(".validation-tooltip");
const regPass = document.querySelector("#register_password");
const liValidNum = document.querySelector("#li-valid-num");
const liValidUpper = document.querySelector("#li-valid-upper");
const liValidLower = document.querySelector("#li-valid-lower");
const liValid8 = document.querySelector("#li-valid-8");
const regPassRpt = document.querySelector("#register_password-rpt");
const validList = regForm.querySelector("ul");
// drag and drop
// adapted from https://medium.com/@cwrworksite/drag-and-drop-file-upload-with-preview-using-javascript-cd85524e4a63
// user
const dropAreaUser = document.querySelector("#drop-zone--user-image");
let dropButtonUser;
let inputUser;
let uploadedFileUser;
if (dropAreaUser) {
  dropButtonUser = dropAreaUser.querySelector(".button");
  inputUser = dropAreaUser.querySelector("input");
  uploadedFileUser = document.querySelector("#uploaded-file--user-image");
}
// post
const dropAreaPost = document.querySelector("#drop-zone--post");
const dropButtonPost = dropAreaPost.querySelector(".button");
const inputPost = dropAreaPost.querySelector("input");
const uploadedFilePost = document.querySelector("#uploaded-file--post");
const dragText = document.querySelector(".dragText");
const dragButton = document.querySelector(".button");
// page select
const selectDropdown = document.getElementById("pagedrop");
// feeds
const feeds = document.querySelectorAll(".feeds-wrapper");
// ------
let file;
let filename;

// post channel selection dropdown
document.addEventListener("DOMContentLoaded", () => {
  const dropdownToggle = document.querySelector(".dropdown-toggle");
  const wrapperDropdown = document.querySelector(".wrapper-dropdown");

  dropdownToggle.addEventListener("click", () => {
    const isActive = dropdownToggle.classList.toggle("active");

    if (isActive) {
      dropdownToggle.style.background = "var(--color-hl-pink)";
      wrapperDropdown.style.display = "block";
    } else {
      dropdownToggle.style.background = "";
      wrapperDropdown.style.display = "none";
    }
  });
});

document.addEventListener("DOMContentLoaded", function () {
  feeds.forEach((feed) => {
    feed.classList.add("feeds-wrapper-loaded");
  });

  const joinedAndOwnedChannelContainer = document.querySelector(
    "#sidebar-channel-block",
  );
  let joinedAndOwnedChannels;

  if (joinedAndOwnedChannelContainer) {
    joinedAndOwnedChannels =
      joinedAndOwnedChannelContainer.querySelectorAll(".sidebar-channel");
    joinedAndOwnedChannels.forEach((channel) =>
      channel.addEventListener("click", (e) => {
        e.preventDefault();
        navigateToChannel(channel);
      }),
    );
  }

  toggleUserInteracted("add");
  actButtonContainer = document.querySelector("#activity-bar");
  if (!actButtonContainer) {
    console.error("Activity bar not found");
    return;
  }
  actButtonsAll = actButtonContainer.querySelectorAll("button");
  activityFeeds = document.querySelector("#user-activity-feeds");
  if (!activityFeeds) {
    console.error("User activity feeds not found");
    return;
  }
  activityFeedsContentAll = activityFeeds.querySelectorAll(
    '[id^="activity-feed-"]',
  );
  // Log the buttons for debugging purposes (optional)
  if (actButtonsAll) {
    actButtonsAll.forEach((button) =>
      button.addEventListener("click", (e) => {
        toggleFeed(
          document.getElementById("activity-" + e.target.id),
          document.getElementById("activity-feed-" + e.target.id),
          e.target,
        );
        // console.log('activity-' + e.target.id);
      }),
    );
    // right panel buttons
    rightPanelButtons = document.querySelectorAll(
      '[id^="right-panel-channel--"]',
    );
  }

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
  // console.log("activity buttons:", actButtonsAll)
});

// ------- channel rules -------
document.addEventListener("DOMContentLoaded", () => {
  const addButton = document.querySelector("#add-unsubmitted-rule");
  const submitButton = document.querySelector("#edit-channel-rules-btn");
  const addedRulesWrapper = document.querySelector("#rules-wrapper-added");
  const removedRulesWrapper = document.querySelector("#rules-wrapper-removed");
  const inputField = document.querySelector("#create-unsubmitted-rule");
  const hiddenInput = document.querySelector("#rules-hidden-input");
  const existingRulesContainer = document.querySelector(
    "#rules-wrapper-existing",
  );
  let existingRules = existingRulesContainer.querySelectorAll(
    '[id^="existing-channel-rule-"]',
  );

  existingRules.forEach((element) =>
    element.addEventListener("click", (e) => {
      removeExistingRule(e.target.id);
    }),
  );

  let rulesList = [];
  let addRuleCounter = 0;

  addButton.addEventListener("click", addRule);
  inputField.addEventListener("keydown", handleKeyPress);

  function addRule() {
    const ruleText = inputField.value.trim();

    if (ruleText) {
      const ruleId = `ruleItem-${addRuleCounter++}`;
      createRuleItem(ruleId, ruleText, "add");

      rulesList.push({ id: ruleId, text: ruleText });
      updateHiddenInput();

      inputField.value = "";
    }
  }
  function removeExistingRule(ruleId) {
    const item = document.getElementById(ruleId);
    const ruleText = item.innerText.trim();
    console.log("existing ruleText: ", ruleText);

    if (ruleText) {
      console.log("remove rule: ", ruleId);
      createRuleItem(ruleId, ruleText, "remove");
      rulesList.push({ id: ruleId, text: ruleText });
      updateHiddenInput();
    }
    console.log("removing ", item);
    item.remove();
  }

  function createRuleItem(ruleId, ruleText, process) {
    const ruleItem = document.createElement("li");
    ruleItem.classList.add("rule-item");
    ruleItem.classList.add("flex-space-between");
    ruleItem.id = ruleId;

    // Rule text span
    const ruleTextSpan = document.createElement("span");
    ruleTextSpan.textContent = ruleText;
    ruleTextSpan.classList.add("rule-text");
    if (process === "add") {
      console.log("process add: text = ", ruleText, "ID = ", ruleItem.id);
      ruleTextSpan.addEventListener("click", () => editRule(ruleId, ruleText));
    } else {
      console.log("process remove: text = ", ruleText, "ID = ", ruleItem.id);
      ruleItem.addEventListener("click", () => removeRule(ruleId));
    }

    // Delete button
    const deleteButton = document.createElement("button");
    deleteButton.classList.add("delete-rule-btn");
    deleteButton.classList.add("btn-channel");
    deleteButton.classList.add("btn-sm");
    deleteButton.classList.add("btn-icoonly");
    deleteButton.innerHTML = `<span class="btn-minus" role="contentinfo" aria-description="Remove Rule"></span>`;
    deleteButton.addEventListener("click", (event) => {
      event.stopPropagation(); // Prevent triggering edit on click
      removeRule(ruleId);
    });

    ruleItem.appendChild(ruleTextSpan);
    ruleItem.appendChild(deleteButton);
    console.log("process query: ", process);
    if (process === "remove") {
      removedRulesWrapper.appendChild(ruleItem);
    } else {
      addedRulesWrapper.appendChild(ruleItem);
    }
  }

  function editRule(ruleId, oldText) {
    const ruleItem = document.getElementById(ruleId);
    if (!ruleItem) return;

    const editInput = document.createElement("input");
    editInput.type = "text";
    editInput.value = oldText;
    editInput.classList.add("rule-edit-input");

    ruleItem.replaceChildren(editInput);
    editInput.focus();

    editInput.addEventListener("keydown", (event) => {
      if (event.key === "Enter") {
        saveEditedRule(ruleId, editInput.value);
      } else if (event.key === "Escape") {
        cancelEdit(ruleId, oldText);
      }
    });

    editInput.addEventListener("blur", () =>
      saveEditedRule(ruleId, editInput.value),
    );
  }

  function saveEditedRule(ruleId, newText) {
    if (!newText.trim())
      return cancelEdit(
        ruleId,
        rulesList.find((r) => r.id === ruleId)?.text || "",
      );

    const ruleIndex = rulesList.findIndex((rule) => rule.id === ruleId);
    if (ruleIndex !== -1) {
      rulesList[ruleIndex].text = newText;
      updateHiddenInput();
    }

    createRuleItem(ruleId, newText);
    document.querySelector(".rule-edit-input")?.remove();
  }

  function cancelEdit(ruleId, oldText) {
    createRuleItem(ruleId, oldText);
    document.querySelector(".rule-edit-input")?.remove();
  }

  function removeRule(ruleId) {
    console.log("removeRule: ", ruleId);
    rulesList = rulesList.filter((rule) => rule.id !== ruleId);
    const ruleItem = document.getElementById(ruleId);
    console.log("ruleItem ID: ", ruleItem.id);

    if (ruleItem) ruleItem.remove();

    if (ruleItem.id.startsWith("existing-channel-rule-")) {
      console.log("removing existing rule: ", ruleItem.innerText);
      const ruleTextSpan = document.createElement("span");
      ruleTextSpan.textContent = ruleItem.innerText;
      ruleTextSpan.id = ruleId;
      existingRulesContainer.appendChild(ruleTextSpan);

      existingRules = existingRulesContainer.querySelectorAll(
        '[id^="existing-channel-rule-"]',
      );
      existingRules.forEach((element) =>
        element.addEventListener("click", (e) => {
          removeExistingRule(e.target.id);
        }),
      );
    }

    updateHiddenInput();
  }

  function updateHiddenInput() {
    hiddenInput.value = JSON.stringify(rulesList);
  }

  function handleKeyPress(event) {
    if (event.key === "Enter") {
      if (event.ctrlKey || event.metaKey) {
        submitButton.click();
      } else {
        addButton.click();
        event.preventDefault();
      }
    }
  }
});

// SECTION ----- functions ------
function fetchAndNavigate(url, page, errorMessage) {
  fetch(url, { method: "GET" })
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return response.json();
    })
    .then((data) => {
      console.table(data);
      changePage(page);
    })
    .catch((error) => {
      console.error(errorMessage, error);
    });
}

function fetchChannelData(channelId) {
  console.log("Fetching channel data for ID:", channelId);
  fetch(`/channels/${channelId}`)
    .then((response) => response.json())
    .then((data) => {
      document.getElementById("channel-page-banner").innerHTML = data.postsHTML;
    })
    .catch((error) => console.error("Error fetching channel data:", error));
}

export function navigateToChannel(channel) {
  const link = channel.getAttribute("data-channel-id");
  if (!link) {
    console.error("Channel ID is missing");
    return;
  }

  // Change page view to channel page
  changePage(channelPage);

  // Now fetch and inject the updated channel data
  fetchChannelData(link);
}

export function navigateToPost(post) {
  const link = post.getAttribute("data-post-id");
  if (!link) {
    console.error("Post ID is missing");
    return;
  }
  fetchAndNavigate(`/posts/${link}`, postPage, "Error navigating to post:");
}

export function navigateToAuthor(author) {
  const link = author.getAttribute("data-author-id");
  if (!link) {
    console.error("Author ID is missing");
    return;
  }
  fetchAndNavigate(`/users/${link}`, userPage, "Error navigating to author:");
}

// toggleUserInteracted toggles user-interacted class on input fields to prevent label animation before they are selected
function toggleUserInteracted(action) {
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

// revert to original state if modals are closed
// TODO get this popover reset to work without constant click listeners

document.addEventListener("click", () => {
  popovers.forEach((popover) => {
    if (popover.matches(":popover-open")) {
      console.log("Open popover: ", popover);
      popover.addEventListener("toggle", () => {
        toggleUserInteracted("remove");
      });
    }
  });
});

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

// organize z-indices
const makeZIndexes = (layers) =>
  layers.reduce((agg, layerName, index) => {
    const valueName = `z-index-${layerName}`;
    agg[valueName] = index * 100;

    return agg;
  }, {});

function getRandomInt(max) {
  return Math.floor(Math.random() * max);
}

function getUserProfileImageFromAttribute() {
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
  // console.log('getInitialFromAttribute running...')
  // console.log('Elements found: ', userProfileImageEmpty.length)
  const colorsArr = [
    ["var(--color-hl-blue)", "var(--color-light-1)"],
    ["var(--color-hl-green)", "var(--color-dark-1)"],
    ["var(--color-hl-orange)", "var(--color-light-1)"],
    ["var(--color-hl-pink)", "var(--color-light-1)"],
    ["var(--color-hl-yellow)", "var(--color-dark-1)"],
    ["var(--color-hl-red)", "var(--color-light-1)"],
  ];
  let userTheme = getRandomInt(6);
  for (let i = 0; i < userProfileImageEmpty.length; i++) {
    let attArr = ["user", "sidebar-user", "channel"];
    attArr[0] = userProfileImageEmpty[i].getAttribute("data-name-user");
    attArr[1] = userProfileImageEmpty[i].getAttribute("data-name-user-sidebar");
    attArr[2] = userProfileImageEmpty[i].getAttribute("data-name-channel");
    if (attArr[0]) {
      userProfileImageEmpty[i].style.background = colorsArr[userTheme][0];
      userProfileImageEmpty[i].style.color = colorsArr[userTheme][1];
      userProfileImageEmpty[i].style.fontSize = "2rem";
      userProfileImageEmpty[i].setAttribute(
        "data-initial",
        Array.from(`${attArr[0]}`)[0],
      );
      // console.log('attArr[0]');
    } else if (attArr[1]) {
      userProfileImageEmpty[i].style.background = colorsArr[userTheme][0];
      userProfileImageEmpty[i].style.color = colorsArr[userTheme][1];
      userProfileImageEmpty[i].style.fontSize = "5rem";
      userProfileImageEmpty[i].setAttribute(
        "data-initial",
        Array.from(`${attArr[1]}`)[0],
      );
      // console.log('attArr[0]');
    } else if (attArr[2]) {
      let theme = getRandomInt(6);
      userProfileImageEmpty[i].style.background = colorsArr[theme][0];
      userProfileImageEmpty[i].style.color = colorsArr[theme][1];
      userProfileImageEmpty[i].setAttribute(
        "data-initial",
        Array.from(`${attArr[2]}`)[0],
      );
      // console.log('attArr[1]');
    } else {
      console.warn(
        "No data-name- attribute value found for element:",
        userProfileImage[i],
      );
    }
  }
}

function toggleColorScheme() {
  // Get the current color scheme
  const currentScheme = document.documentElement.getAttribute("color-scheme");
  // Toggle between light and dark
  const newScheme = currentScheme === "light" ? "dark" : "light";
  // Set the new color scheme
  document.documentElement.setAttribute("color-scheme", newScheme);
  localStorage.setItem("color-scheme", newScheme);
}

// Apply persisted color scheme on page load
document.addEventListener("DOMContentLoaded", () => {
  const savedScheme = localStorage.getItem("color-scheme");
  if (savedScheme) {
    document.documentElement.setAttribute("color-scheme", savedScheme);
  }
});

function toggleDarkMode() {
  const checkbox = document.querySelector("#darkmode-checkbox");
  checkbox.checked = !checkbox.checked;
  console.log("toggle dark mode");
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
// toggle login and register forms
function logReg(target) {
  if (target === "btn_login-1" || target === "btn_login-2") {
    loginHeader.innerText = "sign in to codex";
  }
  if (target === "btn_register-1" || target === "btn_register-2") {
    loginHeader.innerText = "register for codex";
  }
  if (target === "btn_login-2") {
    console.log("btn-log-2");
    formLogin.classList.remove("display-off");
    formForgot.classList.add("display-off");
    forgotVisible = false;
  } else if (target === "btn_register-2") {
    console.log("btn-reg-2");
    formRegister.classList.remove("display-off");
    formForgot.classList.add("display-off");
    forgotVisible = false;
  } else {
    console.log("Toggling login and register forms");
    formLogin.classList.toggle("display-off");
    formRegister.classList.toggle("display-off");
  }
}
// toggle forgot password form
function forgot() {
  formLogin.classList.add("display-off");
  formRegister.classList.add("display-off");
  formForgot.classList.remove("display-off");
  loginHeader.innerText = "reset password";
  forgotVisible = true;
}
// validate each requirement of the password
function validatePass() {
  if (regPass.value.match(/[0-9]/g)) {
    liValidNum.classList.add("li-valid");
  }
  if (regPass.value.match(/[A-Z]/g)) {
    liValidUpper.classList.add("li-valid");
  }
  if (regPass.value.match(/[a-z]/g)) {
    liValidLower.classList.add("li-valid");
  }
  if (regPass.value.length >= 8) {
    liValid8.classList.add("li-valid");
  }
  if (!regPass.value.match(/[0-9]/g)) {
    liValidNum.classList.remove("li-valid");
  }
  if (!regPass.value.match(/[A-Z]/g)) {
    liValidUpper.classList.remove("li-valid");
  }
  if (!regPass.value.match(/[a-z]/g)) {
    liValidLower.classList.remove("li-valid");
  }
  if (regPass.value.length < 8) {
    liValid8.classList.remove("li-valid");
  }
}
// confirm password validation
function confirmPass() {
  if (regPassRpt.value !== regPass.value || regPassRpt.value === "") {
    return regPassRpt.classList.add("pass-nomatch");
  }
  regPassRpt.classList.remove("pass-nomatch");
}
// showMainNotification changes the element ID to provide feedback to user
function showMainNotification(message, timeout = 2500) {
  const notification = document.getElementById("notification-main");
  const notificationContent = document.getElementById(
    "notification-main-content",
  );
  notificationContent.textContent = message;
  notification.style.display = "flex";
  setTimeout(() => {
    notification.style.display = "none";
  }, timeout); // Hide after 3 seconds
}
// showNotification changes the element ID to provide feedback to user
function showNotification(elementID, messageOld, messageNew, success) {
  const notification = document.getElementById(elementID);
  notification.textContent = messageNew;
  notification.style.color = "var(--color-hl-green)";
  if (!success) {
    notification.style.color = "var(--color-error)";
    setTimeout(() => {
      notification.textContent = messageOld;
      notification.style.color = "var(--color-fg-1)";
    }, 2500); // Hide after 3 seconds
  }
}
// retrieve the csrf_token cookie and explicitly set the X-CSRF-Token header in requests
function getCSRFToken() {
  const match = document.cookie
    .split("; ")
    .find((row) => row.startsWith("csrf_token"));
  if (!match) {
    console.warn("CSRF token not found in cookies");
    return null;
  }
  return match.substring("csrf_token=".length);
}

// sendRequest('/protected', 'GET').then((response) => {
//     console.log('Use of protected route successful:', response);
//     })
//     .catch((error) => {
//         console.error('Use of protected route failed:', error);
// });
// switch the <p> elements in right panel to <textarea> for editing
// and change the edit button to submit

function changePage(page) {
  let pageId = page;
  // console.log("type", typeof page);
  if (typeof page != "string") {
    pageId = page.id;
  }
  pages.forEach((element) => {
    if (element.id === pageId) {
      element.classList.add("active-feed");
      console.log("set", element.id, "to active-feed");
    } else {
      element.classList.remove("active-feed");
    }
  });
}
// SECTION ---- event listeners -----

// --- go home ---
goHome.addEventListener("click", () => {
  changePage(homePage);
});

// --- select page ---
selectDropdown.addEventListener("change", () => {
  const selectedValue = selectDropdown.value;
  console.log("selectedValue: ", selectedValue);
  changePage(selectedValue);
});

// --- sidebar options dropdown ---
sidebarOption.addEventListener("click", function (event) {
  sidebarOptionsList.classList.toggle("sidebar-options-reveal");
  sidebarOptionsList.classList.toggle("ul-forwards");
  sidebarOptionsList.classList.toggle("ul-reverse");
});

// --- register ---

if (registerForm) {
  registerForm.addEventListener("submit", function (event) {
    event.preventDefault(); // Prevent the default form submission
    const form = event.target;
    const formData = new FormData(form); // Collect form data
    fetch("/register", {
      method: "POST",
      // headers: {
      //   "Content-Type": "application/json",
      // },
      body: formData,
      // body: JSON.stringify({
      //   register_user: document.getElementById("register_user").value,
      //   register_email: document.getElementById("register_email").value,
      //   register_password: document.getElementById("register_password").value,
      // }),
      cache: "no-store",
    })
      .then((response) => {
        if (response.ok) {
          return response.json(); // Parse JSON response
        } else {
          throw new Error("Registration failed.");
        }
      })
      .then((data) => {
        if (data.message === "registration failed!") {
          showNotification("login-title", "", data.message, false);
        } else {
          showNotification("login-title", "", data.message, true);
          setTimeout(() => {
            window.location.href = "/";
          }, 2000);
        }
      })
      .catch((error) => {
        console.error("Registration failed:", error);
      });
  });
}

// --- login ---
if (loginForm) {
  loginForm.addEventListener("submit", function (event) {
    event.preventDefault(); // Prevent the default form submission
    // const form = event.target;
    // const formData = new FormData(form); // Collect form data
    const csrfToken = getCSRFToken();
    console.log("csrfToken: ", csrfToken);

    fetch("/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "x-csrf-token": csrfToken,
      },
      // body: formData, // Send the form data
      body: JSON.stringify({
        username: document.getElementById("loginFormUser").value,
        password: document.getElementById("loginFormPassword").value,
      }),
      cache: "no-store",
    })
      .then((response) => {
        if (response.ok) {
          console.table(response.json);
          return response.json();
        } else {
          console.log("(error) response", response.json);
          throw new Error("Login failed.");
        }
      })
      .then((data) => {
        if (data.message === "incorrect password") {
          showNotification("login-title", "", data.message, false);
        } else if (data.message === "user not found") {
          showNotification("login-title", "", data.message, false);
        } else if (data.message === "failed to create cookies") {
          showNotification("login-title", "", data.message, false);
        } else {
          showNotification("login-title", "", data.message, true);
          setTimeout(() => {
            window.location.href = "/";
          }, 2000);
        }
      })
      .catch((error) => {
        console.error("Login failed:", error);
        // showMainNotification("An error occurred during login.");
      });
  });
}

// --- logout ---
if (logoutFormButton) {
  logoutFormButton.addEventListener("click", function (event) {
    event.preventDefault();

    fetch("/logout", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({}),
      cache: "no-store",
    })
      .then((response) => {
        if (response.ok) {
          return response.json(); // Parse JSON response
        } else {
          throw new Error("Logout failed.");
        }
      })
      .then((data) => {
        // Show notification based on the response from the server
        showMainNotification(data.message);
        // Optional: Redirect the user after showing the notification
        setTimeout(() => {
          window.location.href = "/"; // Replace with your desired location
        }, 3500);
      })
      .catch((error) => {
        console.error("Error:", error);
        showMainNotification("An error occurred during logout.");
      });
  });
}
// ---- join channel ----
joinChannelButton.addEventListener("submit", function (event) {
  event.preventDefault();
  // const form = event.target;
  // const formData = new FormData(form); // Collect form data
  const csrfToken = getCSRFToken();
  console.log("csrfToken: ", csrfToken);

  fetch("/channels/join", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "x-csrf-token": csrfToken,
    },
    body: JSON.stringify({
      channelId: document.getElementById("join-channel-id").value,
      agree: document.getElementById("rules-agree-checkbox").value,
    }),
    cache: "no-store",
  })
    .then((response) => {
      if (response.ok) {
        return response.json(); // Parse JSON response
      } else {
        throw new Error("Join channel failed.");
      }
    })
    .then((data) => {
      // Show notification based on the response from the server
      showMainNotification(data.message);
      // Optional: Redirect the user after showing the notification
      setTimeout(() => {
        window.location.href = "/"; // Replace with your desired location
      }, 3500);
    })
    .catch((error) => {
      console.error("Error:", error);
      showMainNotification("An error occurred while joining channel.");
    });
});

if (inputUser) {
  inputUser.addEventListener("change", function () {
    file = this.files[0];
    dropAreaUser.classList.add("active");
  });

  // when file is inside drag area
  dropAreaUser.addEventListener("dragover", (event) => {
    event.preventDefault();
    dropAreaUser.classList.add("active");
    dragText.textContent = "release to Upload";
    dragButton.style.display = "none";
    // console.log('File is inside the drag area');
  });
  // when file leaves the drag area
  dropAreaUser.addEventListener("dragleave", () => {
    dropAreaUser.classList.remove("active");
    // console.log('File left the drag area');
    dragText.textContent = "drag your file here";
  });
  // when file is dropped
  dropAreaUser.addEventListener("drop", (event) => {
    event.preventDefault();
    dropAreaUser.classList.add("dropped");
    // console.log('File is dropped in drag area');
    file = event.dataTransfer.files[0]; // grab single file even if user selects multiple files
    // console.log(file);
    displayFile(uploadedFileUser, dropAreaUser);
  });
}
// SECTION ----- drag and drop ----
// get user image from manual click

// get post image from manual click
inputPost.addEventListener("change", function () {
  file = this.files[0];
  dropAreaUser.classList.add("active");
});

function displayFile(uploadedFile, dropArea) {
  let fileType = file.type;
  // console.log(fileType);
  let validExtensions = ["image/*"];
  if (validExtensions.includes(fileType)) {
    let fileReader = new FileReader();
    fileReader.onload = () => {
      uploadedFile.innerHTML = `<div class="dragText">uploaded</div>
        <div class="uploaded-file">${file.name}</div>`;
      dropArea.classList.add("dropped");
    };
    fileReader.readAsDataURL(file);
  } else {
    alert("This is not an Image File");
    dropArea.classList.remove("active");
    dragText.textContent = "Drag and drop your file, or";
    dragButton.style.display = "unset";
  }
}

// switchDl.addEventListener('click', toggleColorScheme);
darkSwitch.addEventListener("click", toggleDarkMode);

// open modals
// TODO refactor the open and close modals
if (openLoginModal) {
  openLoginModal.addEventListener(
    "click",
    () => (loginModal.style.display = "block"),
  );
}
openLoginModalFallback.addEventListener(
  "click",
  () => (loginModal.style.display = "block"),
);
// openEditUserModal.addEventListener('click', () => editUserModal.style.display = 'block');
// openAccSettingsModal.addEventListener('click', () => accSettingsModal.style.display = 'block');
// openViewStatsModal.addEventListener('click', () => viewStatsModal.style.display = 'block');
// openRemoveAccModal.addEventListener('click', () => removeAccModal.style.display = 'block');
// close modals
closeLoginModal.addEventListener(
  "click",
  () => (loginModal.style.display = "none"),
);
// closeEditUserModal.addEventListener('click', () => editUserModal.style.display = 'none');
// closeAccSettingsModal.addEventListener('click', () => accSettingsModal.style.display = 'none');
// closeViewStatsModal.addEventListener('click', () => viewStatsModal.style.display = 'none');
// closeRemoveAccModal.addEventListener('click', () => removeAccModal.style.display = 'none');
window.addEventListener("click", ({ target }) => {
  switch (target) {
    case loginModal:
      loginModal.style.display = "none";
      break;
    case accSettingsModal:
      accSettingsModal.style.display = "none";
      break;
    case viewStatsModal:
      viewStatsModal.style.display = "none";
      break;
    case removeAccModal:
      removeAccModal.style.display = "none";
      break;
  }
});
// TODO create the functionality in the function
// right panel buttons
if (rightPanelButtons) {
  rightPanelButtons.forEach((button) =>
    button.addEventListener("click", (e) => rightPanelEdit(e.target.id)),
  );
}
// login / register / forgot
btnsLogin.forEach((button) =>
  button.addEventListener("click", (e) => logReg(e.target.id)),
);
btnsRegister.forEach((button) =>
  button.addEventListener("click", (e) => logReg(e.target.id)),
);
btnsForgot.addEventListener("click", forgot);
// TODO get these working
// validate password
regPass.addEventListener("input", validatePass);
// check passwords match
regPassRpt.addEventListener("focusout", confirmPass);
// reverse the order of the password validation list delay
regPass.addEventListener("focusin", () => {
  setTimeout(() => {
    validList.classList.remove("ul-forwards");
    validList.classList.add("ul-reverse");
  }, 1000);
});
regPass.addEventListener("focusout", () => {
  setTimeout(() => {
    validList.classList.remove("ul-reverse");
    validList.classList.add("ul-forwards");
  }, 1000);
});
