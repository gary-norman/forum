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
console.log("%cformattedDate:", expect, formattedDate);

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
    let startDateInteracted = false;

    filterButtons.forEach(button => {
        const id = button.getAttribute('popovertarget');
        const popover = activePageElement.querySelector(`#${id}`);
        const clearButton = popover.querySelector(".clear-choices ");

        button.addEventListener("click", (e) => {
            const noneSelected= checkSelectedInputs(popover);

            if (!popover.matches(':popover-open') && button.contains(e.target)) {
                // console.log("%ctoggleFilters -- 1", expect);
                toggleFilterButtonState(button, noneSelected, "button");
                // console.log("%cpopover opened", expect, popover);
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

                    // console.log("%cbox:", expect, input);
                    // console.log("%cinput.checked:", expect, input.checked);

                    const noneSelected= checkSelectedInputs(popover);
                    // console.log("%cnoneSelected:", expect, noneSelected);

                    if (noneSelected) {
                        // console.log("%chiding the clearButton -- 3", angry);
                        toggleClearButton(popover, "hide");

                    } else {
                        // console.log("%cshowing the clearButton -- 3", angry);
                        toggleClearButton(popover, "show");
                    }

                    // console.log("test1")
                    // console.log("%ctoggleFilters -- 2", expect);
                    toggleFilterButtonState(button, noneSelected, "inside");
                });
            });
        }

        if (popover && popover.matches('[popover].card.date')) {
            const startInput = popover.querySelector('input[id^="creation-year-start"]');
            const endInput =  popover.querySelector('input[id^="creation-year-end"]');

            popover.querySelectorAll("input").forEach(function (input) {
                input.addEventListener("input", (e) => {

                    if (startInput.value !== "" && startDateInteracted === false) {
                        console.log("%ctest 1:", warn);

                        endInput.value = formattedDate;
                        console.log("%cend input:", expect, endInput.value);

                    }
                    const noneSelected= checkSelectedInputs(popover);
                    console.log("%cnoneSelected:", expect, noneSelected);
                    console.log("%cpopover -- 1:", expect, popover);
                    toggleFilterButtonState(button, noneSelected, "inside");

                    if (noneSelected) {
                        console.log("%chiding the clearButton -- 3", angry);
                        toggleClearButton(popover, "hide");

                    } else {
                        console.log("%cshowing the clearButton -- 3", angry);
                        toggleClearButton(popover, "show");
                    }
                });

                input.addEventListener("blur", () => {
                    console.log("%cstartDateInteracted -- 3: ", angry, startDateInteracted);

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
                    // console.log("test2")
                    // console.log("%ctoggleFilters -- 3", expect);

                    toggleFilterButtonState(button, noneSelected, "inside");
                    clearRadios(popover, button);
                });

                if (popover && popover.matches('[popover].card.date')) {
                    popover.querySelectorAll("input").forEach(function (input) {
                        input.value = "";

                        const noneSelected = checkSelectedInputs(popover);

                        startDateInteracted = false;
                        toggleClearButton(popover, "hide");
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

    // console.log("%cbutton:", warn, button);

    const buttonLabel = button.getAttribute("aria-label");
    // console.log("%caria label:", warn, buttonLabel);

    button.querySelector("span").innerText = buttonLabel.slice(0, -7);
}

function listenToClicksForPopovers(){
    const popovers = activePageElement.querySelectorAll('[popover].filter-popover');
    // console.log("%cpopovers:", expect, popovers);
    if (popovers) {
        document.addEventListener("click", (e) => {
            popovers.forEach((popover) => {
                // console.log("popover", popover);
                // console.log("___________________________________");
                // console.log("!popover.contains(e.target)", !popover.contains(e.target));
                // console.log("popover.matches(':popover-open')", popover.matches(':popover-open'));
                if (!popover.contains(e.target) && popover.matches(':popover-open')) {
                    const button = activePageElement.querySelector(`[popovertarget="${popover.id}"]`);
                    const noneSelected = checkSelectedInputs(popover);
                    // console.log("%cpopover:", warn, popover);
                    // console.log("%cnoneSelected", warn, noneSelected);

                    toggleFilterButtonState(button, noneSelected, "outside");
                } else {
                    // console.log("baabaabba")
                }
            });
        }, true);
    }
}

function toggleFilterButtonState(button, noneSelected, clickLocation) {
    const id = button.getAttribute('popovertarget');
    const popover = activePageElement.querySelector(`#${id}`);
    const cancelButton = popover.querySelector(".clear-choices ");

    // console.log("%cnoneSelected:", warn, noneSelected);

    if (clickLocation === "button" && noneSelected) {
        button.classList.add("active");
    }

    if (popover && !popover.matches(':popover-open')) {
        return;
    }

    // console.log("%cbutton:", warn, button);
    // console.log("%cnoneSelected:", warn, noneSelected);
    // console.log("%cclickLocation:", warn, clickLocation);

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
        // console.log("%csecond if statement:", angry);
        console.log("%ctest2", angry);

        button.querySelector("span").classList.add("btn-filters");
        button.querySelector("span").classList.remove("btn-active-filter");
        button.classList.remove("selected");
        button.classList.add("active");


        if (cancelButton && clickLocation === "outside") {
            // console.log("%csecond-B if statement:", angry);
            button.classList.remove("active");
        }

    } else {
        // console.log("%cstyling button", angry, button);
        // console.log("%ctest1", angry);
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
    // console.log("%ccheckboxes:", warn, checkboxes);

    // console.log("%cradios:", warn, radios);

    if (checkboxes.length > 0) {
        // console.log("%ccheckboxes:", expect, checkboxes);
        return Array.from(checkboxes).every(checkbox => !checkbox.checked);
    } else if (radios.length > 0) {
        // console.log("%cradios:", expect, radios);
        return Array.from(radios).every(radio => !radio.checked);
    } else if (dates.length > 0) {
        // console.log("%cdates:", expect, dates);
        return [...dates].every(date => !date.value);
    }
    return console.log("%cno checkboxes or radios", angry);
}


