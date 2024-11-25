// variables
const switchDl = document.getElementById('switch-dl');
const darkSwitch = document.getElementById('sidebar-option-darkmode');
// activity buttons
const actButtonContainer = document.querySelector('#activity-bar')
const actButtonsAll = actButtonContainer.querySelectorAll('button')
// activity feeds
const activityFeeds = document.querySelector('#activity-feeds')
const filterButtonsContainer = activityFeeds.querySelectorAll('.button-row')
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
    const timeOut = 400;
    actButtonsAll.forEach( button => button.classList.remove('btn-active') );
    activityFeedsContentAll.forEach(feed => {
        feed.classList.remove('collapsible-expanded');
        feed.classList.add('collapsible-collapsed');
    });
    setTimeout(() => {
        // targetFeed.classList.remove('collapsible-collapsed');
        targetFeedContent.classList.remove('collapsible-collapsed');
        // targetFeed.classList.add('collapsible-expanded');
        targetFeedContent.classList.add('collapsible-expanded');
        targetButton.classList.toggle('btn-active'); }, timeOut);
    setTimeout(() => {
        filterButtonsContainer.forEach(feed => {
        // const currentButtonDisplay = feed.getAttribute('display');
        // const newButtonDisplay = currentButtonDisplay === 'none' ? 'flex' : 'none';
        // feed.setAttribute('display', newButtonDisplay);
        feed.classList.add('hide-feed');
        });
        targetFeed.querySelector('.button-row').classList.remove('hide-feed');}, timeOut)

}

// event listeners
// switchDl.addEventListener('click', toggleColorScheme);
darkSwitch.addEventListener('click', toggleDarkMode);
actButtonsAll.forEach( button => button.addEventListener('click', (e) => {
    toggleFeed(document.getElementById("activity-" + e.target.id),document.getElementById("activity-feed-" + e.target.id),  e.target);
    console.log('activity-' + e.target.id);
}) );