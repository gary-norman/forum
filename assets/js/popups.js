// modals
// import { resetInputStyle, toggleUserInteracted } from "./update_UI_elements.js";
import { activePageElement } from "./main.js";

const angry =
  "background-color: #000000; color: #ff0000; font-weight: bold; border: 2px solid #ff0000; padding: 5px; border-radius: 5px;";
const expect =
  "background-color: #000000; color: #00ff00; font-weight: bold; border: 1px solid #00ff00; padding: 5px; border-radius: 5px;";
const warn =
  "background-color: #000000; color: #e3c144; font-weight: bold; border: 1px solid #e3c144; padding: 5px; border-radius: 5px;";
export function closePostForm() {
  console.info("closePostForm");
  const formPost = document.querySelector("#form-post");
  const cancelPost = document.querySelector("#cancel-post");
  if (formPost && cancelPost) {
    cancelPost.popoverTargetElement = formPost;
    cancelPost.popoverTargetAction = "hide";
  }
}

const formChannel = document.querySelector("#form-channel");
const cancelChannel = document.querySelector("#cancel-channel");
cancelChannel.popoverTargetElement = formChannel;
cancelChannel.popoverTargetAction = "hide";
const formChannelRules = document.querySelector("#form-edit-channel-rules");
const cancelChannelRules = document.querySelector("#cancel-channel-rules");
cancelChannelRules.popoverTargetElement = formChannelRules;
cancelChannelRules.popoverTargetAction = "hide";

document.addEventListener("DOMContentLoaded", () => {
  toggleModals();
  // popoverJumpToDefaultInput();
});

// export function togglePopoverUserInteracted() {
//   const popovers = document.querySelectorAll("[popover]:has(.user-label)");
//   popovers.forEach((popover) => {
//     // Ensure the listener is only added once
//     if (!popover._toggleListenerAdded) {
//       popover.addEventListener("toggle", () => {
//         if (popover.matches(":popover-open")) {
//           toggleUserInteracted("remove");
//         }
//       });
//       popover._toggleListenerAdded = true;
//     }
//   });
// }

// export function popoverJumpToDefaultInput() {
//   const popovers = document.querySelectorAll("[popover]:has(.default-input)");
//   // Use WeakSet to track popovers with listeners
//   if (!popoverJumpToDefaultInput._initialized) {
//     popoverJumpToDefaultInput._initialized = new WeakSet();
//   }
//   const initialized = popoverJumpToDefaultInput._initialized;
//
//   let i = 1;
//
//   popovers.forEach((popover) => {
//     console.log("%cpopover", expect, i, popover);
//     i++;
//     if (!initialized.has(popover)) {
//       popover.addEventListener("toggle", (e) => {
//         if (e.newState === "open") {
//           const defaultInput = popover.querySelector(".defaultInput");
//           if (defaultInput) defaultInput.focus();
//         }
//       });
//       initialized.add(popover);
//     }
//   });
// }

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
    openLoginModal.addEventListener("click", () => {
      loginModal.style.display = "block";

      // Delay attaching the outside click listener
      setTimeout(() => {
        document.addEventListener("click", handleClickOutsideModals);
      }, 0);
    });

    openLoginModalFallback.addEventListener(
      "click",
      () => (loginModal.style.display = "block"),
    );
  } else {
    console.warn("openLoginModal not found");
  }

  function handleClickOutsideModals(e) {
    console.log("%ce.target", warn, e.target);

    // if the target is not the modal itself, which is of class modal-content
    if (e.target === loginModal || e.target === closeLoginModal) {
      console.log("%cclicked outside modals", expect);

      loginModal.style.display = "none";
      resetInputStyle();
      document.removeEventListener("click", handleClickOutsideModals);
    }
  }
}
