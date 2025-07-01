import {activePageElement} from "./main.js";

const angry =
    "background-color: #000000; color: #ff0000; font-weight: bold; border: 2px solid #ff0000; padding: 5px; border-radius: 5px;";
const expect =
    "background-color: #000000; color: #00ff00; font-weight: bold; border: 1px solid #00ff00; padding: 5px; border-radius: 5px;";
const warn =
    "background-color: #000000; color: #e3c144; font-weight: bold; border: 1px solid #e3c144; padding: 5px; border-radius: 5px;";

let filterButtons;
const currentDate = new Date();
const formattedDate = currentDate.toISOString().split("T")[0];
let allPagePosts = [];

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
    allPagePosts = Array.from(activePageElement.querySelectorAll(".card.link"));
}

// compareDates used to compare 2 dates, which then is used in the default js sort function
function compareDates(a, b) {
    const dateA = new Date(a.dataset.createdAt);
    const dateB = new Date(b.dataset.createdAt);
    return dateA - dateB;
}

// compareNumberValues used to compare 2 numbers, which then is used in the default js sort function
function compareNumberValues(sortBy) {
    return function(a,b){
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
    }
}

// compareActivity used to compare 2 dates per element, and chooses the most recent date for each
// This is then used in the default js sort function
function compareActivity(a, b) {
    const createdDateA = new Date(a.dataset.createdAt);
    const createdDateB = new Date(b.dataset.createdAt);
    const lastReactedDateA = new Date(a.dataset.lastReaction);
    const lastReactedDateB = new Date(b.dataset.lastReaction);

    // Pick the later date for a
    const latestDateA = (lastReactedDateA > createdDateA) ? lastReactedDateA : createdDateA;

    // Pick the later date for b
    const latestDateB = (lastReactedDateB > createdDateB) ? lastReactedDateB : createdDateB;

    // Return difference for sort
    return latestDateA - latestDateB;
}

// sortDescending sort by oldest date / most amount first
function sortDescending(func) {
    return allPagePosts.sort(func);
}

// sortAscending sort by newest date / least amount first
function sortAscending(func) {
    return allPagePosts.sort(func).reverse();
}

function filterContent() {
    const filtersRow = activePageElement.querySelector(".filters-row");
    const channelCheckboxes = filtersRow.querySelector(`[id$="dropdown-channel"]`).querySelectorAll('input[type="checkbox"]:checked');
    const reactionCheckboxes = filtersRow.querySelector(`[id$="dropdown-reaction"]`).querySelectorAll('input[type="checkbox"]:checked');
    const typeCheckboxes = filtersRow.querySelector(`[id$="dropdown-type"]`).querySelectorAll('input[type="checkbox"]:checked');
    const commentRadios = filtersRow.querySelector(`[id$="dropdown-comments"]`).querySelectorAll('input[type="radio"]:checked');
    const sortRadios = filtersRow.querySelector(`[id$="dropdown-sort"]`).querySelectorAll('input[type="radio"]:checked');
    const startDateInput = filtersRow.querySelector(`[id$="dropdown-date"] input[name="creation-year-start"]`);
    const endDateInput = filtersRow.querySelector(`[id$="dropdown-date"] input[name="creation-year-end"]`);

    // Build active filters
    const activeFilters = {
        channels: Array.from(channelCheckboxes).map(cb => cb.value),
        reactions: Array.from(reactionCheckboxes).map(rb => rb.value),
        types: Array.from(typeCheckboxes).map(cb => cb.value),
        comments: commentRadios.length > 0 ? commentRadios[0].value : null,
        sort: sortRadios.length > 0 ? sortRadios[0].value : null,
        startDate: startDateInput?.value ? new Date(startDateInput.value) : null,
        endDate: endDateInput?.value ? new Date(endDateInput.value) : null,
    };
    // Sort
    // Default sort
    if (activeFilters.sort === null) {
        allPagePosts = sortNewestDateFirst(compareDates);
        console.log("%cDEFAULT SORTING... allPagePosts:", warn, allPagePosts);
        reorderVisiblePosts();
    } else {

        // sort by date
        //sort newest first
        if (activeFilters.sort.includes("most-new")) {
            allPagePosts = sortNewestDateFirst(compareDates);
            // console.log("%cSORTING NEWEST FIRST... allPagePosts:", warn, allPagePosts);
        }


        //sort oldest first
        if (activeFilters.sort.includes("most-old")) {
            allPagePosts = sortOldestDateFirst(compareDates);
            // console.log("%cSORTING OLDEST FIRST... allPagePosts:", warn, allPagePosts);
        }

        //sort by recent activity
        //sort newest first
        if (activeFilters.sort.includes("most-new-activity")) {
            allPagePosts = sortNewestDateFirst(compareActivity);
            console.log("%cSORTING NEWEST ACTIVITY FIRST... allPagePosts:", warn, allPagePosts);
        }

        //sort oldest first
        if (activeFilters.sort.includes("most-old-activity")) {
            allPagePosts = sortOldestDateFirst(compareActivity);
            console.log("%cSORTING OLDEST ACTIVITY FIRST... allPagePosts:", warn, allPagePosts);
        }

        //sort oldest first

        //sort by likes
        //sort most first

        //sort least first

        //sort by dislikes
        //sort most first

        //sort least first

        //sort by comments
        //sort most first

        //sort least first

        reorderVisiblePosts();
    }




    allPagePosts.forEach(post => {
        const postChannel = post.dataset.channelId;
        const postComments = Number(post.querySelector(".btn-reply").textContent);
        const postLikes = Number(post.querySelector(".btn-likes").textContent);
        const postDislikes = Number(post.querySelector(".btn-dislikes").textContent);
        const postDate = new Date(post.dataset.createdAt);

        let visible = true;



        // Filter by channel
        if (activeFilters.channels.length > 0 && !activeFilters.channels.includes(postChannel)) {
            visible = false;
        }

        // Filter by comments
        if (activeFilters.comments !== null) {
            let commentMatch = false;

            if (activeFilters.comments.includes("has-comments") && postComments > 0) {
                commentMatch = true;
            }

            if (activeFilters.comments.includes("no-comments") && postComments <= 0) {
                commentMatch = true;
            }

            if (commentMatch === false) {
                visible = false;
            }
        }

        // Filter by reaction
        if (activeFilters.reactions.length > 0) {
            let reactionMatch = false;

            if (activeFilters.reactions.includes("liked") && postLikes > 0) {
                reactionMatch = true;
            }

            if (activeFilters.reactions.includes("disliked") && postDislikes > 0) {
                reactionMatch = true;
            }

            if (
                activeFilters.reactions.includes("no-reaction") &&
                postLikes === 0 &&
                postDislikes === 0
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

            if (postDate <= activeFilters.startDate) {
                visible = false;
            }
            if (postDate >= activeFilters.endDate) {
                visible = false;
            }
        }


        // Show or hide
        post.parentElement.classList.toggle('hide', !visible);
    });
}



function reorderVisiblePosts() {
    const feedContainer = activePageElement.querySelector(`[id$="feed"]`);

    // Remove existing posts from DOM
    allPagePosts.forEach(post => {
        console.log("removing")
        // console.log("%cREMOVING", angry, post.parentElement);
        feedContainer.removeChild(post.parentElement);
    });

    // Append posts back in sorted order
    allPagePosts.forEach(post => {
        console.log("appending")

        // console.log("%cAPPENDING", angry, post.parentElement);
        feedContainer.appendChild(post.parentElement)
    });
}



function toggleFilters() {
    let filterContainer = activePageElement.querySelector(".filters-row");
    filterButtons = filterContainer.querySelectorAll("button.btn-action");
    let startDateInteracted = false;

    filterButtons.forEach(button => {
        const id = button.getAttribute('popovertarget');
        const popover = activePageElement.querySelector(`#${id}`);
        const clearButton = popover.querySelector(".clear-choices ");

        button.addEventListener("click", (e) => {
            const noneSelected= checkSelectedInputs(popover);

            if (!popover.matches(':popover-open') && button.contains(e.target)) {
                toggleFilterButtonState(button, noneSelected, "button");
            }
        })

        if (popover && popover.matches('[popover].wrapper-filter-dropdown')) {
            popover.querySelectorAll(".dropdown-option").forEach(function (option) {
                option.addEventListener("click", (e) => {
                    let input;
                    const checkbox = option.querySelector("input[type='checkbox']");
                    const radio = option.querySelector("input[type='radio']");
                    let optionText = option.querySelector("label").innerText;

                    if (checkbox) {
                        input = checkbox;
                    } else if (radio) {
                        input = radio;
                        if (optionText === "has no comments" ) {
                            optionText = optionText.slice(4);
                        }
                        clearRadios(popover, button);

                        option.classList.add("selected");
                        button.querySelector("span").innerText = optionText;
                    }

                    input.click();

                    const noneSelected= checkSelectedInputs(popover);

                    filterContent();

                    if (noneSelected) {
                        toggleClearButton(popover, "hide");

                    } else {
                        toggleClearButton(popover, "show");
                    }
                    toggleFilterButtonState(button, noneSelected, "inside");
                });
            });
        }

        if (popover && popover.matches('[popover].card.date')) {
            const startInput = popover.querySelector('input[id^="creation-year-start"]');
            const endInput =  popover.querySelector('input[id^="creation-year-end"]');
            const pickers =

            popover.querySelectorAll("input").forEach(function (input) {
                input.addEventListener("input", (e) => {

                    if (startInput.value !== "" && startDateInteracted === false) {
                        endInput.value = formattedDate;
                    }
                    const noneSelected= checkSelectedInputs(popover);

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
            })

        }


        if (clearButton){
            clearButton.addEventListener("click", () => {
                popover.querySelectorAll('.dropdown-option').forEach(function (option) {
                    const checkbox = option.querySelector("input[type='checkbox']");
                    const radio = option.querySelector("input[type='radio']");
                    let input;
                    if (radio) {
                        input = radio;
                    } else if (checkbox) {
                        input = checkbox;
                    }
                    input.checked = false;

                    const noneSelected = checkSelectedInputs(popover);

                    toggleClearButton(popover, "hide");
                    toggleFilterButtonState(button, noneSelected, "inside");
                    console.log("test5");
                    filterContent();
                    clearRadios(popover, button);
                });

                if (popover && popover.matches('[popover].card.date')) {
                    popover.querySelectorAll("input").forEach(function (input) {
                        input.value = "";

                        const noneSelected = checkSelectedInputs(popover);

                        startDateInteracted = false;
                        toggleClearButton(popover, "hide");
                        filterContent();
                        toggleFilterButtonState(button, noneSelected, "inside");

                    });
                }
            });
        }
    })
}


function clearRadios(popover, button) {
    const allOptions = popover.querySelectorAll(".dropdown-option");

    allOptions.forEach((o) => {
        o.classList.remove("selected");
    })

    const buttonLabel = button.getAttribute("aria-label");
    button.querySelector("span").innerText = buttonLabel.slice(0, -7);
}

function listenToClicksForPopovers(){
    const popovers = activePageElement.querySelectorAll('[popover].filter-popover');

    if (popovers) {
        document.addEventListener("click", (e) => {
            popovers.forEach((popover) => {
                if (!popover.contains(e.target) && popover.matches(':popover-open')) {
                    const button = activePageElement.querySelector(`[popovertarget="${popover.id}"]`);
                    const noneSelected = checkSelectedInputs(popover);

                    toggleFilterButtonState(button, noneSelected, "outside");
                }
            });
        }, true);
    }
}

function toggleFilterButtonState(button, noneSelected, clickLocation) {
    const id = button.getAttribute('popovertarget');
    const popover = activePageElement.querySelector(`#${id}`);
    const cancelButton = popover.querySelector(".clear-choices ");

    if (clickLocation === "button" && noneSelected) {
        button.classList.add("active");
    }

    if (popover && !popover.matches(':popover-open')) {
        return;
    }

    if (clickLocation === "outside" && noneSelected) {
        // console.log("%cfirst if statement:", angry);
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
        return Array.from(checkboxes).every(checkbox => !checkbox.checked);
    } else if (radios.length > 0) {
        return Array.from(radios).every(radio => !radio.checked);
    } else if (dates.length > 0) {
        return [...dates].every(date => !date.value);
    }
    return console.log("%cno checkboxes or radios", angry);
}


