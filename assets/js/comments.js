// INFO was a DOMContentLoaded function
import {activePage} from "./main.js";

export function listenToReplies() {
    let postID;
    let commentID;
    let msg;
    let title;
    const textareas = document.querySelectorAll('[id^="comment-form-textarea-"]');
    const replyButtons = document.querySelectorAll('[id^="reply-button-"]')
    const allReplyForms = document.querySelectorAll('.form-reply');


    if (activePage === "post-page") {

        replyButtons.forEach(button => {
            button.addEventListener("click", function (e) {
                const card = button.closest('.card');
                const channelID = card.getAttribute('data-channel-id');
                const postID = card.getAttribute('data-post-id');
                const commentID = card.getAttribute('data-comment-id');
                const targetForm = card.querySelector('.form-reply');

                allReplyForms.forEach(form => {
                    if (form !== targetForm)
                        form.classList.remove('replying');
                });

                setTimeout(() => {
                    if (targetForm.classList.contains('replying')) {
                        targetForm.classList.remove('replying');
                        const textarea = targetForm.querySelector('[id^="comment-form-textarea-"]');
                        textarea.value = '';
                        textarea.style.height = "5rem";
                    } else {
                        targetForm.classList.add('replying');
                    }


                }, 200);

            })
        })

        textareas.forEach(textarea => {
            postID = textarea.closest('.card').getAttribute('data-post-id');
            commentID = textarea.closest('.card').getAttribute('data-comment-id');

            textarea.addEventListener("input", function () {
                this.style.height = "auto"; // Reset height to recalculate
                this.style.height = this.scrollHeight + "px"; // Set height to fit content
            });
        })

        // Get the form by its name
        let replyForms = document.querySelectorAll('form[name="replyForm"]'); // Returns a NodeList

        // console.log("reply forms:", replyForms)

        replyForms.forEach(form => {
            // Add submit event listener
            form.addEventListener("submit", function (event) {
                let content = form["content"].value.trim(); // Trim spaces

                if (content === "") {
                    alert("Reply cannot be empty!");
                    event.preventDefault(); // Stop form submission
                }
            });
        })
    }
}
