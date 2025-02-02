document.addEventListener("DOMContentLoaded", function () {
    const textarea = document.querySelector("textarea");

    textarea.addEventListener("input", function () {
        this.style.height = "auto"; // Reset height to recalculate
        this.style.height = this.scrollHeight + "px"; // Set height to fit content
    });
});
