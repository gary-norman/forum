import {activePageElement} from "./main.js";

const angry =
    "background-color: #000000; color: #ff0000; font-weight: bold; border: 2px solid #ff0000; padding: 5px; border-radius: 5px;";
const expect =
    "background-color: #000000; color: #00ff00; font-weight: bold; border: 1px solid #00ff00; padding: 5px; border-radius: 5px;";
const warn =
    "background-color: #000000; color: #e3c144; font-weight: bold; border: 1px solid #e3c144; padding: 5px; border-radius: 5px;";

let filterButtons;

document.addEventListener("newContentLoaded", () => {
    if (activePageElement.id !== "post-page") {
        listenToClicksForPopovers();
        toggleFilters();
    }
});

document.addEventListener("DOMContentLoaded", () => {
    if (activePageElement.id !== "post-page") {
        listenToClicksForPopovers();
        toggleFilters();
    }
});

function toggleFilters() {
    let filterContainer = activePageElement.querySelector(".filters-row");
    filterButtons = filterContainer.querySelectorAll("button.btn-action");

    filterButtons.forEach(button => {
        const id = button.getAttribute('popovertarget');
        const popover = activePageElement.querySelector(`#${id}`);
        // debugger;
        const cancelButton = popover.querySelector(".clear-checkboxes");

        button.addEventListener("click", () => {
            const noneSelected= updateCheckboxes(popover);

            if (!popover.matches(':popover-open')) {

                toggleFilterButtonState(button, noneSelected, "inside");
                console.log("%cpopover opened", expect, popover);
            }
        })
        if (popover && popover.matches('[popover].wrapper-filter-dropdown')) {
            popover.querySelectorAll('.dropdown-option').forEach(function (option) {
                option.addEventListener("click", (e) => {
                    const checkbox = option.querySelector("input[type='checkbox']");
                    checkbox.checked = !checkbox.checked;

                    const noneSelected= updateCheckboxes(popover);

                    if (noneSelected) {
                        toggleClearButton(popover, "hide");
                    } else {
                        toggleClearButton(popover, "show");
                    }
                    // console.log("test1")
                    toggleFilterButtonState(button, noneSelected, "inside");
                });
            });
        }

        if (cancelButton){
            cancelButton.addEventListener("click", () => {
                popover.querySelectorAll('.dropdown-option').forEach(function (option) {
                    const checkbox = option.querySelector("input[type='checkbox']");
                    checkbox.checked = false;

                    const noneSelected= updateCheckboxes(popover);

                    toggleClearButton(popover, "hide");
                    // console.log("test2")
                    toggleFilterButtonState(button, noneSelected, "inside");
                });
            });
        }
    })
}

function listenToClicksForPopovers(){
    const popovers = activePageElement.querySelectorAll('[popover].filter-popover');
    // console.log("%cpopovers:", expect, popovers);
    if (popovers) {
        document.addEventListener("click", (e) => {
            popovers.forEach((popover) => {
                if (e.target !== popover) {
                    const button = activePageElement.querySelector(`[popovertarget="${popover.id}"]`);
                    const noneSelected = updateCheckboxes(popover);
                    // console.log("test3 - id: ", popover.id);
                    toggleFilterButtonState(button, noneSelected, "outside");
                }
            });
        }, true);
    }
}

function toggleFilterButtonState(button, noneSelected, clickLocation) {
    const id = button.getAttribute('popovertarget');
    const popover = activePageElement.querySelector(`#${id}`);
    const cancelButton = popover.querySelector(".clear-checkboxes");


    // console.log("%cnoneSelected:", warn, noneSelected);
    // console.log("%cIs false strictly? ", warn, noneSelected === false);
    // console.log("%cclickLocation:", warn, clickLocation);

    if (clickLocation === "outside" && noneSelected) {
        // console.log("%cfirst if statement:", angry);
        button.querySelector("span").classList.add("btn-filters");
        button.querySelector("span").classList.remove("btn-active-filter");
        button.classList.remove("selected");
        button.classList.remove("active");
        return;
    } else {

    }

    if (noneSelected) {
        // console.log("%csecond if statement:", angry);

        button.querySelector("span").classList.add("btn-filters");
        button.querySelector("span").classList.remove("btn-active-filter");
        button.classList.remove("selected");
        button.classList.add("active");

        if (cancelButton) {
            // console.log("%ccancelButton:", expect, cancelButton);
            button.classList.remove("active");
        }

    } else {
        // console.log("%cthird if statement:", angry);
        button.querySelector("span").classList.remove("btn-filters");
        button.querySelector("span").classList.add("btn-active-filter");
        button.classList.add("selected");
    }
}

function toggleClearButton(popover, state) {
    const button = popover.querySelector(".clear-checkboxes");
    const scrollContainer = popover.querySelector(".content");
    if (state === "show") {
        scrollContainer.style.maxHeight = `calc(28rem - 0.8rem - 4.4rem)`;
        button.classList.remove("hide");
    } else {
        scrollContainer.style.maxHeight = `calc(28rem - 0.8rem)`;
        button.classList.add("hide");
    }
}

function updateCheckboxes(popover) {
    //update state of checkboxes
    const checkboxes = popover.querySelectorAll("input[type='checkbox']");

    return Array.from(checkboxes).every(checkbox => !checkbox.checked);
}


