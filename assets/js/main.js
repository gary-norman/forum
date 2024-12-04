// variables
const switchDl = document.querySelector('#switch-dl');
const darkSwitch = document.querySelector('#sidebar-option-darkmode');
// activity buttons
const actButtonContainer = document.querySelector('#activity-bar')
const actButtonsAll = actButtonContainer.querySelectorAll('button')
// activity feeds
const activityFeeds = document.querySelector('#activity-feeds')
const activityFeedsContentAll = activityFeeds.querySelectorAll('[id^="activity-feed-"]')
// login/register butons
const btnLogin = document.querySelectorAll('[id^="btn_login"]');
const btnRegister = document.querySelectorAll('[id^="btn_register-"]');
const btnForgot = document.querySelector('#btn_forgot');
// login/register forms
const formLogin = document.querySelector('#form-login');
const formRegister = document.querySelector('#form-register');
const formForgot = document.querySelector('#form-forgot');
let forgotVisible = false;
// login/register modal
const modal = document.querySelector('#form-login-container');
// Get the button that opens the modal
const openLoginModal = document.querySelector('#btn-open-login-modal');
// Get the <span> element that closes the modal
const closeLoginModal = document.getElementsByClassName("close")[0];
// drag and drop
// adapted from https://medium.com/@cwrworksite/drag-and-drop-file-upload-with-preview-using-javascript-cd85524e4a63
const dropArea = document.querySelector('#drop_zone');
const dragText = document.querySelector('.dragText');
const dragButton = document.querySelector('.button');
let dropButton = dropArea.querySelector('.button');
let input = dropArea.querySelector('input');
let file;
let filename;

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
    const checkbox = document.querySelector('#darkmode-checkbox');
    checkbox.checked = !checkbox.checked;
    console.log('toggle dark mode')
}

function toggleFeed(targetFeed, targetFeedContent, targetButton) {
    actButtonsAll.forEach(button => button.classList.remove('btn-active'));
    activityFeedsContentAll.forEach(feed => feed.classList.replace('collapsible-expanded', 'collapsible-collapsed'));
    targetFeedContent.classList.replace('collapsible-collapsed', 'collapsible-expanded');
    targetButton.classList.toggle('btn-active');
    targetFeed.querySelector('.button-row').classList.toggle('hide-feed', false);
}

function logReg() {
    console.log('Toggling login and register forms');
    formLogin.classList.toggle('display-off');
    formRegister.classList.toggle('display-off');

    if (forgotVisible === true) {
        console.log('Form was visible, hiding it now');
        formForgot.classList.remove('display-off');
        forgotVisible = false;
    } else {
        console.log('Form was not visible');
    }
}
function forgot() {
    formLogin.classList.add('display-off');
    formRegister.classList.add('display-off');
    formForgot.classList.toggle('display-off');
    forgotVisible = true;
}

// ---- event listeners -----
// drag and drop
dropButton.addEventListener('click', input.click.bind(input), false);
// when browse
input.addEventListener('change', function () {
    file = this.files[0];
    dropArea.classList.add('active');
});
// when file is inside drag area
dropArea.addEventListener('dragover', (event) => {
    event.preventDefault();
    dropArea.classList.add('active');
    dragText.textContent = 'release to Upload';
    dragButton.style.display = 'none';
    // console.log('File is inside the drag area');
});
// when file leaves the drag area
dropArea.addEventListener('dragleave', () => {
    dropArea.classList.remove('active');
    // console.log('File left the drag area');
    dragText.textContent = 'drag your file here';
});
// when file is dropped
dropArea.addEventListener('drop', (event) => {
    event.preventDefault();
    dropArea.classList.add('dropped');
    // console.log('File is dropped in drag area');
    file = event.dataTransfer.files[0]; // grab single file even if user selects multiple files
    // console.log(file);
});
// switchDl.addEventListener('click', toggleColorScheme);
darkSwitch.addEventListener('click', toggleDarkMode);
actButtonsAll.forEach( button => button.addEventListener('click', (e) => {
    toggleFeed(document.getElementById("activity-" + e.target.id),document.getElementById("activity-feed-" + e.target.id),  e.target);
    console.log('activity-' + e.target.id);
}) );
// login register modal
openLoginModal.addEventListener('click', () => modal.style.display = 'block');
closeLoginModal.addEventListener('click', () => modal.style.display = 'none');
window.addEventListener('click', ({ target }) => {
    if (target === modal) {
        modal.style.display = 'none';
    }
});
btnLogin.forEach(button => button.addEventListener('click', () => logReg));
btnRegister.forEach(button => button.addEventListener('click', () => logReg));
btnForgot.addEventListener('click', forgot);