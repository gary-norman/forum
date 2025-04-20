import {listenToReplies} from "./comments.js";
import {listenToEditDetails} from "./edit_user.js";
import {listenToShare, selectActiveFeed} from "./share.js";
import {listenToLikeDislike} from "./reactions.js";
import {listenToDropdowns} from "./post.js";
import {saveColourScheme} from "./colour_scheme.js";
import {activePage} from "./main.js";
import {getRandomInt} from "./helper_functions.js";
import {listenToRules} from "./channel_rules.js";
import {toggleModals, togglePopovers} from "./popups.js";
import {listenToChannelLinks} from "./navigation.js";

// activity buttons
let actButtonContainer, actButtonsAll;

// activity feeds
let activityFeeds, activityFeedsContentAll;

export function UpdateUI() {
    // console.log("updating UI");
    listenToRules();
    listenToReplies();
    listenToEditDetails();
    listenToShare();
    listenToLikeDislike();
    listenToDropdowns();
    listenToPageSetup();
    saveColourScheme();
    updateProfileImages();
    toggleModals();
    resetInputStyle();
    togglePopovers();
    listenToChannelLinks();
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
            if (userProfileImage[i].parentNode.classList.contains("result-card")) {
                console.log("updating", userProfileImage[i].parentNode.children[1].textContent, "to ", `url('${attArr[2]}')`, "with parent: ", userProfileImage[i].parentNode);
            }

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
    const userProfileImage = document.querySelectorAll(".profile-pic");

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

// INFO was a DOMContentLoaded function
export function listenToPageSetup() {
    // feeds
    const feeds = document.querySelectorAll(".feeds-wrapper");

    feeds.forEach((feed) => {
        feed.classList.add("feeds-wrapper-loaded");
    });

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
                console.log('activity-' + e.target.id);
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