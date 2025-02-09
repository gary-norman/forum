document.addEventListener('DOMContentLoaded', function () {
    const userSettingBlocks = document.querySelectorAll('[id^="settings-user-"]');

    console.log(userSettingBlocks)


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