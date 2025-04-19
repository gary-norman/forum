import {changePage, navigateToPage} from "./fetch_and_navigate.js";
import {data} from "./share.js";

// sidebar butons
const goHome = document.querySelector("#btn-go-home");

// --- go home ---
goHome.addEventListener("click", () => {
    changePage(data["homePage"]);
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

function navigateEntity(e) {
    if (!(e.target.getAttribute("data-dest"))) {
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

document.addEventListener("click", navigateEntity, {capture: false});

