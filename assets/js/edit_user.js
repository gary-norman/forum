document.addEventListener("DOMContentLoaded", function () {
  const userSettingBlocks = document.querySelectorAll('[id^="settings-user-"]');
  const nameInput = document.querySelector("#name");
  const nameContent = document.querySelector("#user-name-content");
  const bioInput = document.querySelector("#bio");
  const bioContent = document.querySelector("#user-bio-content");
  const dragDropImage = document.querySelector("#drop-zone--user-image");
  const inputs = [bioInput, nameInput];
  const adjustHeight = (element) => {
    element.style.height = "auto"; // Reset height to recalculate
    element.style.height = element.scrollHeight + "px"; // Set height to fit content
  };
  // const contents = [bioContent, nameContent]

  console.log("bioContent", bioContent);
  console.log("nameContent", nameContent);
  console.log("inputs", inputs);

  inputs.forEach((input) => {
    if (input) {
      input.addEventListener("input", () => adjustHeight(input));
      input.addEventListener("click", () => adjustHeight(input));
    }
  });

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

  userSettingBlocks.forEach((block) => {
    const editButton = block.querySelector('[id^="edit-user"]');
    const submitButton = block.querySelector('[id^="submit-user"]');
    const cancelButton = block.querySelector('[id^="cancel-user"]');
    editButton.addEventListener("click", function (e) {
      block.classList.add("editing");
      userSettingBlocks.forEach((otherBlock) => {
        if (block !== otherBlock) {
          otherBlock.classList.remove("editing");
        }
      });

      switch (editButton.id) {
        case "edit-user-name":
          nameContent.classList.add("editing");
          nameInput.classList.add("editing");
          bioContent.classList.remove("editing");
          bioInput.classList.remove("editing");
          dragDropImage.classList.remove("editing");
          nameInput.focus();
          break;
        case "edit-user-avatar":
          dragDropImage.classList.add("editing");
          nameContent.classList.remove("editing");
          nameInput.classList.remove("editing");
          bioContent.classList.remove("editing");
          bioInput.classList.remove("editing");
          break;
        case "edit-user-bio":
          bioContent.classList.add("editing");
          bioInput.classList.add("editing");
          nameInput.classList.remove("editing");
          nameContent.classList.remove("editing");
          dragDropImage.classList.remove("editing");
          bioInput.focus();
          break;
        default:
          console.log("edit button clicked but section not recognised");
      }
    });
    cancelButton.addEventListener("click", function (e) {
      block.classList.remove("editing");

      switch (cancelButton.id) {
        case "cancel-user-name":
          nameContent.classList.remove("editing");
          nameInput.classList.remove("editing");
          break;
        case "cancel-user-avatar":
          dragDropImage.classList.remove("editing");
          break;
        case "cancel-user-bio":
          bioContent.classList.remove("editing");
          bioInput.classList.remove("editing");
          break;
        default:
          console.log("cancel button clicked but section not recognised");
      }
    });
    submitButton.addEventListener("click", function (e) {
      block.classList.remove("editing");

      switch (cancelButton.id) {
        case "submit-user-name":
          nameContent.classList.remove("editing");
          nameInput.classList.remove("editing");
          break;
        case "submit-user-avatar":
          dragDropImage.classList.remove("editing");
          break;
        case "submit-user-bio":
          bioContent.classList.remove("editing");
          bioInput.classList.remove("editing");
          break;
        default:
          console.log("submit button clicked but section not recognised");
      }
    });
  });
});

