// variables
const switchDl = document.getElementById('switch-dl');
const darkSwitch = document.getElementById('sidebar-option-darkmode');
// activity buttons
const actButtonContainer = document.querySelector('#activity-bar')
const actButtonsAll = actButtonContainer.querySelectorAll('button')
// activity feeds
const activityFeeds = document.querySelector('#activity-feeds')
const activityFeedsAll = activityFeeds.querySelectorAll('[id^="activity"]')
const activityFeedsContentAll = activityFeeds.querySelectorAll('[id^="activity-feed-"]')

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
    actButtonsAll.forEach( button => button.classList.remove('btn-active') );
    // activityFeedsAll.forEach(feed => {
    //     feed.classList.remove('collapsible-expanded');
    //     feed.classList.add('collapsible-collapsed');
    // });
    activityFeedsContentAll.forEach(feed => {
        feed.classList.remove('collapsible-expanded');
        feed.classList.add('collapsible-collapsed');
    });
    setTimeout(() => {
        // targetFeed.classList.remove('collapsible-collapsed');
        targetFeedContent.classList.remove('collapsible-collapsed');
        // targetFeed.classList.add('collapsible-expanded');
        targetFeedContent.classList.add('collapsible-expanded');
        targetButton.classList.toggle('btn-active'); }, 1000);
}

// event listeners
// switchDl.addEventListener('click', toggleColorScheme);
darkSwitch.addEventListener('click', toggleDarkMode);
actButtonsAll.forEach( button => button.addEventListener('click', (e) => {
    toggleFeed(document.getElementById("activity-" + e.target.id),document.getElementById("activity-feed-" + e.target.id),  e.target);
    console.log('activity-' + e.target.id);
}) );