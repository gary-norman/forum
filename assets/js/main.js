// variables
const switchDl = document.getElementById('switch-dl');
const darkSwitch = document.getElementById('sidebar-option-darkmode');
// activity buttons
let actButtonContainer;
let actButtonsAll;
// activity feeds
let activityFeeds;
let activityFeedsContentAll;

document.addEventListener('DOMContentLoaded', function () {
    actButtonContainer = document.querySelector('#activity-bar');
    actButtonsAll = actButtonContainer.querySelectorAll('button');
    activityFeeds = document.querySelector('#activity-feeds');
    activityFeedsContentAll = activityFeeds.querySelectorAll('[id^="activity-feed-"]');
    // console.log(actButtonsAll);


    // event listeners
// switchDl.addEventListener('click', toggleColorScheme);
    darkSwitch.addEventListener('click', toggleDarkMode);
    actButtonsAll.forEach( button => button.addEventListener('click', (e) => {
        toggleFeed(document.getElementById("activity-" + e.target.id),document.getElementById("activity-feed-" + e.target.id),  e.target);
        console.log('activity-' + e.target.id);
    }) );
});

// functions
function toggleColorScheme() {
    // Get the current color scheme
    const currentScheme = document.documentElement.getAttribute('color-scheme');
    // Toggle between light and dark
    const newScheme = currentScheme === 'light' ? 'dark' : 'light';
    // Set the new color scheme
    document.documentElement.setAttribute('color-scheme', newScheme);
}
function toggleDarkMode() {
    const checkbox = document.getElementById("darkmode-checkbox");
    checkbox.checked = !checkbox.checked;
    console.log("toggle dark mode")
}

function toggleFeed(targetFeed, targetFeedContent, targetButton) {
    const timeOut = 400;
    actButtonsAll.forEach( button => button.classList.remove('btn-active') );
    activityFeedsContentAll.forEach(feed => {
        feed.classList.remove('collapsible-expanded');
        feed.classList.add('collapsible-collapsed');
    });
    setTimeout(() => {
        targetFeedContent.classList.remove('collapsible-collapsed');
        targetFeedContent.classList.add('collapsible-expanded');
        targetButton.classList.toggle('btn-active'); }, timeOut);
    setTimeout(() => {
        targetFeed.querySelector('.button-row').forEach(feed => {
        feed.classList.add('hide-feed');
        });
        targetFeed.querySelector('.button-row').classList.remove('hide-feed');}, timeOut)

}
// drag and drop
// adapted from https://medium.com/@cwrworksite/drag-and-drop-file-upload-with-preview-using-javascript-cd85524e4a63
const dropArea = document.querySelector("#drop_zone");
const dragText = document.querySelector(".dragText");
const dragButton = document.querySelector(".button");

let button = dropArea.querySelector(".button");
let input = dropArea.querySelector("input");
let file;
let filename
button.onclick = () => {
    input.click();
};

// when browse
input.addEventListener("change", function () {
    file = this.files[0];
    dropArea.classList.add("active");
});

// when file is inside drag area
dropArea.addEventListener("dragover", (event) => {
    event.preventDefault();
    dropArea.classList.add("active");
    dragText.textContent = "Release to Upload";
    dragButton.style.display = "none";
    // console.log('File is inside the drag area');
});

// when file leaves the drag area
dropArea.addEventListener("dragleave", () => {
    dropArea.classList.remove("active");
    // console.log('File left the drag area');
    dragText.textContent = "Drag your file here";
});

// when file is dropped
dropArea.addEventListener("drop", (event) => {
    event.preventDefault();
    dropArea.classList.add("dropped");
    // console.log('File is dropped in drag area');
    file = event.dataTransfer.files[0]; // grab single file even if user selects multiple files
    // console.log(file);
});

