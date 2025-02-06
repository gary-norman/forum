document.addEventListener("DOMContentLoaded", function () {
    let postID;
    let commentID;
    let msg;
    let title;

    const textareas = document.querySelectorAll('[id^="comment-form-textarea-"]');
    console.log("found areas: ", textareas)

    textareas.forEach(textarea  => {
        postID = textarea.closest('.card').getAttribute('data-post-id');
        commentID = textarea.closest('.card').getAttribute('data-comment-id');

        console.log("found: ", textarea)
        textarea.addEventListener("input", function () {
            this.style.height = "auto"; // Reset height to recalculate
            this.style.height = this.scrollHeight + "px"; // Set height to fit content
        });
    })


});
