// modals
import {toggleUserInteracted} from "./update_UI_elements.js";

export function togglePopovers() {
    const popovers = document.querySelectorAll("[popover]");

    document.addEventListener("click", () => {
        popovers.forEach((popover) => {
            if (popover.matches(":popover-open")) {
                // console.log("Open popover: ", popover);
                popover.addEventListener("toggle", () => {
                    toggleUserInteracted("remove");
                });
            }
        });
    });
}

export function toggleModals() {
    const loginModal = document.querySelector("#container-form-login");

    // Get the <span> element that closes the modal
    const closeLoginModal = loginModal.getElementsByClassName("close")[0];

    // Get the buttons that open the modals
    const openLoginModal = document.querySelector("#btn-open-login-modal");
    const openLoginModalFallback = document.querySelector(
        "#btn-open-login-modal-fallback",
    );

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

    window.addEventListener("click", ({ target }) => {
        switch (target) {
            case loginModal:
                loginModal.style.display = "none";
                break;
        }
    });
}