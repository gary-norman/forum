import {changePage, navigateToPage} from "./fetch_and_navigate.js";
import {data} from "./share.js";
import {setActivePage} from "./main.js";

// sidebar butons
const goHome = document.querySelector("#btn-go-home");

// --- go home ---
goHome.addEventListener("click", (e) => {
    setActivePage("home");
    history.pushState({}, "", `/`);
    changePage(data["homePage"]);
    // navigateToPage("home", e.target)
});

// INFO was inside a DOMContentLoaded function
export function listenToChannelLinks() {
    const joinedAndOwnedChannelContainer = document.querySelector(
        "#sidebar-channel-block",
    );

    // console.log(joinedAndOwnedChannelContainer);
    let joinedAndOwnedChannels;

    if (joinedAndOwnedChannelContainer) {
        joinedAndOwnedChannels =
            joinedAndOwnedChannelContainer.querySelectorAll(".sidebar-channel");
        joinedAndOwnedChannels.forEach((channel) =>
            channel.addEventListener("click", (e) => {
                e.preventDefault();
                navigateToPage("channel", channel);
            }),
        );
    }
}

function navigateToEntity(e) {
    const hasDestination = e.target.getAttribute("data-dest");
    const isButton = e.target.nodeName.toLowerCase() === "button";

    if ((hasDestination || !isButton) || (isButton && hasDestination) ) {
        const parent = e.target.closest(".link");
        if (parent) {
            const dest = parent.getAttribute("data-dest");
            navigateToPage(dest, parent);
        }
        return;
    }
    if (e.target.matches(".link")) {
        const dest = e.target.getAttribute("data-dest");
        navigateToPage(dest, e.target);
    }
}

document.addEventListener("click", navigateToEntity, {capture: false});

// addGlobalEventListener is an event listener on a parent element that only runs your callback
// when the event happens on a specific type of child element (that matches the selector you give)
// good for event delegation and for new elements populated dynamically
export function addGlobalEventListener(type, selector, callback, parent) {
    parent.addEventListener(type, e => {
        if (e.target.matches(selector)) {
            callback(e)
        }
    })
}
