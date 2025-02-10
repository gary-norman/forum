
document.addEventListener('DOMContentLoaded', function () {
    const userSettingBlocks = document.querySelectorAll('[id^="settings-user-"]');
    const bioInput = document.querySelector('#bio');

    console.log(bioInput);
    console.log("bioInput tag:", bioInput.tagName);
    console.log("is bioInput disabled?", bioInput.disabled);
    console.log("is bioInput readonly?", bioInput.readOnly);


    // console.log("found: ", textarea)
    bioInput.addEventListener("input", function () {
        this.style.height = "auto"; // Reset height to recalculate
        this.style.height = this.scrollHeight + "px"; // Set height to fit content
    });


    userSettingBlocks.forEach(block => {
        const editButton = block.querySelector('[id^="edit-user"]')
        const submitButton = block.querySelector('[id^="submit-user"]')
        const cancelButton = block.querySelector('[id^="cancel-user"]')
        editButton.addEventListener("click", function (e) {
            console.log("edit button clicked")
            block.classList.add("editing")




        })
        cancelButton.addEventListener("click", function (e) {
            console.log("cancel button clicked")
            block.classList.remove("editing")
        })
        submitButton.addEventListener("click", function (e) {
            console.log("submit button clicked")
            block.classList.remove("editing")
        })
    })
})