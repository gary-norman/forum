// INFO was a DOMContentLoaded function
import { activePageElement } from "./main.js";
import {
  applyCustomTheme,
  showThemesClickable,
  pickTheme,
  themeList,
} from "./consoleThemes.js";

document.addEventListener("newContentLoaded", () => {
  if (
    activePageElement.id === "channel-page" ||
    activePageElement.id === "user-page"
  ) {
    listenToEditDetails();
  }
});

document.addEventListener("DOMContentLoaded", () => {
  if (
    activePageElement.id === "channel-page" ||
    activePageElement.id === "user-page"
  ) {
    listenToEditDetails();
  }
});

function listenToEditDetails() {
  const settingBlocks = activePageElement.querySelectorAll('[id*="settings-"]');
  const nameInput = activePageElement.querySelector('[id$="name-input"]');
  const nameContent = activePageElement.querySelector('[id$="name-content"]');
  const nameLabel = activePageElement.querySelector('[id$="name-label"]');
  const bioInput = activePageElement.querySelector('[id$="bio-input"]');
  const bioContent = activePageElement.querySelector('[id$="bio-content"]');
  const dragDropImage = activePageElement.querySelector(
    "#drop-zone--user-image",
  );
  const inputs = [bioInput, nameInput];
  const adjustHeight = (element) => {
    element.style.height = "auto"; // Reset height to recalculate
    element.style.height = element.scrollHeight + "px"; // Set height to fit content
  };
  // const contents = [bioContent, nameContent]

  // console.log("bioContent", bioContent);
  // console.log("nameContent", nameContent);
  // console.log("inputs", inputs);

  // inputs.forEach((input) => {
  //   if (input) {
  //     input.addEventListener("focus", () => {
  //       console.log("%cadjusting height - input", warn );
  //       adjustHeight(input)
  //     });
  //
  //     input.addEventListener("input", () => {
  //       console.log("%cadjusting height - input", warn );
  //       adjustHeight(input)
  //     });
  //
  //     input.addEventListener("click", () => {
  //       console.log("%cadjusting height - click", warn );
  //       adjustHeight(input)
  //     });
  //   }
  // });

  if (bioInput) {
    bioInput.addEventListener("focus", () => adjustHeight(bioInput));
  }

  // inputs.forEach((input) => {
  //   input.addEventListener("input", function () {
  //     this.style.height = "auto"; // Reset height to recalculate
  //     this.style.height = this.scrollHeight + "px"; // Set height to fit content
  //   });
  //
  //   input.addEventListener("click", function () {
  //     this.style.height = "auto"; // Reset height to recalculate
  //     this.style.height = this.scrollHeight + "px"; // Set height to fit content
  //   });
  //
  //   bioInput.addEventListener("focus", function () {
  //     this.style.height = "auto"; // Reset height to recalculate
  //     this.style.height = this.scrollHeight + "px"; // Set height to fit content
  //   });
  // });
  function removeEditingState() {
    console.custom.info("nameContent:", nameContent);
    console.custom.info("nameInput:", nameInput);
    console.custom.info("nameLabel:", nameLabel);

    nameContent.classList.remove("editing");
    nameInput.classList.remove("editing");
    nameLabel.classList.remove("editing");
    bioContent.classList.remove("editing");
    bioInput.classList.remove("editing");
    dragDropImage.classList.remove("editing");
  }

  console.log("settingBlocks", settingBlocks);

  settingBlocks.forEach((block) => {
    const editButton = block.querySelector('[id*="edit-"]');
    const submitButton = block.querySelector('[id*="submit-"]');
    const cancelButton = block.querySelector('[id*="cancel-"]');

    editButton.addEventListener("click", function (e) {
      block.classList.add("editing");
      settingBlocks.forEach((otherBlock) => {
        if (block !== otherBlock) {
          otherBlock.classList.remove("editing");
        }
      });

      switch (true) {
        case editButton.id.endsWith("edit-name"):
          removeEditingState();
          nameContent.classList.add("editing");
          nameInput.classList.add("editing");
          nameLabel.classList.add("editing");
          nameInput.focus();
          break;
        case editButton.id.endsWith("edit-avatar"):
          removeEditingState();
          dragDropImage.classList.add("editing");

          break;
        case editButton.id.endsWith("edit-bio"):
          removeEditingState();
          bioContent.classList.add("editing");
          bioInput.classList.add("editing");
          bioInput.focus();
          break;
        default:
          console.log("edit button clicked but section not recognised");
      }
    });
    cancelButton.addEventListener("click", function (e) {
      block.classList.remove("editing");

      switch (true) {
        case cancelButton.id.endsWith("cancel-name"):
          nameContent.classList.remove("editing");
          nameInput.classList.remove("editing");
          nameLabel.classList.remove("editing");
          break;
        case cancelButton.id.endsWith("cancel-avatar"):
          dragDropImage.classList.remove("editing");
          break;
        case cancelButton.id.endsWith("cancel-bio"):
          bioContent.classList.remove("editing");
          bioInput.classList.remove("editing");
          break;
        default:
          console.log("cancel button clicked but section not recognised");
      }
    });
    submitButton.addEventListener("click", function (e) {
      block.classList.remove("editing");

      switch (true) {
        case submitButton.id.endsWith("submit-name"):
          // nameContent.classList.remove("editing");
          // nameInput.classList.remove("editing");
          removeEditingState();
          break;
        case submitButton.id.endsWith("submit-avatar"):
          // dragDropImage.classList.remove("editing");
          removeEditingState();
          break;
        case submitButton.id.endsWith("submit-bio"):
          // bioContent.classList.remove("editing");
          // bioInput.classList.remove("editing");
          removeEditingState();
          break;
        default:
          console.log("submit button clicked but section not recognised");
      }
    });
  });
}
