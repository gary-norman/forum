// variables
const switchDl = document.querySelector('#switch-dl');
const darkSwitch = document.querySelector('#sidebar-option-darkmode');
// activity buttons
const actButtonContainer = document.querySelector('#activity-bar')
const actButtonsAll = actButtonContainer.querySelectorAll('button')
// activity feeds
const activityFeeds = document.querySelector('#activity-feeds')
const activityFeedsContentAll = activityFeeds.querySelectorAll('[id^="activity-feed-"]')
// login/register buttons
// sidebar elements
const sidebarOption = document.querySelector('#sidebar-options');
const sidebarOptionsList = document.querySelector('.sidebar-options-list');
// TODO overhaul the naming of these buttons
const loginTitle = document.querySelector('#login-title');
const loginFormButton = document.querySelector('#login');
const logoutFormButton = document.querySelector('#logout');
const btnLogin = document.querySelectorAll('[id^="btn_login-"]');
const btnRegister = document.querySelectorAll('[id^="btn_register-"]');
const btnForgot = document.querySelector('#btn_forgot');
// login/register forms
const formLogin = document.querySelector('#form-login');
const formRegister = document.querySelector('#form-register');
const formForgot = document.querySelector('#form-forgot');
const formEditUser = document.querySelector('#form-edit-user');
const formAccSettings = document.querySelector('#form-acc-settings');
const formViewStats = document.querySelector('#form-view-stats');
const formRemoveAcc = document.querySelector('#form-remove-acc');
let forgotVisible = false;
// login/register modal
const loginModal = document.querySelector('#container-form-login');
const editUserModal = document.querySelector('#container-form-edit-user');
const accSettingsModal = document.querySelector('#container-form-acc-settings');
const viewStatsModal = document.querySelector('#container-form-view-stats');
const removeAccModal = document.querySelector('#container-form-remove-acc');
// Get the buttons that open the modals
const openLoginModal = document.querySelector('#btn-open-login-modal');
const openEditUserModal = document.querySelector('#btn-open-edit-user-modal');
// const openAccSettingsModal = document.querySelector('#btn-open-acc-settings-modal');
// const openViewStatsModal = document.querySelector('#btn-open-view-stats-modal');
// const openRemoveAccModal = document.querySelector('#btn-open-remove-acc-modal');
// Get the <span> element that closes the modal
const closeLoginModal = loginModal.getElementsByClassName("close")[0];
const closeEditUserModal = editUserModal.getElementsByClassName("close")[0];
// const closeAccSettingsModal = accSettingsModal.getElementsByClassName("close")[0];
// const closeViewStatsModal = viewStatsModal.getElementsByClassName("close")[0];
// const closeRemoveAccModal = removeAccModal.getElementsByClassName("close")[0];
// registration form
const regForm = document.querySelector('#form-register');
const regFormInputs = regForm.querySelectorAll('input');
const regFormSpans = regForm.querySelectorAll('span');
const regFormIcons = regForm.querySelectorAll('.validation-icon');
const regFormTooltips = regForm.querySelectorAll('.validation-tooltip');
const regPass = document.querySelector('#register_password');
const liValidNum = document.querySelector('#li-valid-num');
const liValidUpper = document.querySelector('#li-valid-upper');
const liValidLower = document.querySelector('#li-valid-lower');
const liValid8 = document.querySelector('#li-valid-8');
const regPassRpt = document.querySelector('#register_password-rpt');
const validList = regForm.querySelector('ul');
// drag and drop
// adapted from https://medium.com/@cwrworksite/drag-and-drop-file-upload-with-preview-using-javascript-cd85524e4a63
// user
const dropAreaUser = document.querySelector('#drop-zone--user-image');
const dropButtonUser = dropAreaUser.querySelector('.button');
const inputUser = dropAreaUser.querySelector('input');
const uploadedFileUser = document.querySelector('#uploaded-file--user-image');
// post
const dropAreaPost = document.querySelector('#drop-zone--post');
const dropButtonPost = dropAreaPost.querySelector('.button');
const inputPost = dropAreaPost.querySelector('input');
const uploadedFilePost = document.querySelector('#uploaded-file--post');
const dragText = document.querySelector('.dragText');
const dragButton = document.querySelector('.button');
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
// toggle the various feeds on user-page
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
// toggle login and register forms
function logReg(target) {
    if (target === 'btn_login-2') {
        console.log('btn-log-2');
        formLogin.classList.remove('display-off');
        formForgot.classList.add('display-off');
        forgotVisible = false;
    } else if (target === 'btn_register-2') {
        console.log('btn-reg-2');
        formRegister.classList.remove('display-off');
        formForgot.classList.add('display-off');
        forgotVisible = false;
    } else {
        console.log('Toggling login and register forms');
        formLogin.classList.toggle('display-off');
        formRegister.classList.toggle('display-off');
    }
}
// toggle forgot password form
function forgot() {
    formLogin.classList.add('display-off');
    formRegister.classList.add('display-off');
    formForgot.classList.remove('display-off');
    forgotVisible = true;
}
// validate each requirement of the password
function validatePass() {
    if (regPass.value.match(/[0-9]/g)) {
        liValidNum.classList.add('li-valid');
    }
    if (regPass.value.match(/[A-Z]/g)) {
        liValidUpper.classList.add('li-valid');
    }
    if (regPass.value.match(/[a-z]/g)) {
        liValidLower.classList.add('li-valid');
    }
    if (regPass.value.length >= 8) {
        liValid8.classList.add('li-valid');
    }
    if (!regPass.value.match(/[0-9]/g)) {
        liValidNum.classList.remove('li-valid');
    }
    if (!regPass.value.match(/[A-Z]/g)) {
        liValidUpper.classList.remove('li-valid');
    }
    if (!regPass.value.match(/[a-z]/g)) {
        liValidLower.classList.remove('li-valid');
    }
    if (regPass.value.length < 8) {
        liValid8.classList.remove('li-valid');
    }
}
// confirm password validation
function confirmPass() {
    if (regPassRpt.value !== regPass.value || regPassRpt.value === "") {
        return regPassRpt.classList.add("pass-nomatch");
    }
    regPassRpt.classList.remove("pass-nomatch");
}

// retrieve the csrf_token cookie and explicitly set the X-CSRF-Token header in requests
function getCSRFToken() {

    const match = document.cookie
    .split('; ')
    .find((row) => row.startsWith('csrf_token'))
    return match ? match.substring('csrf_token='.length) : null;
}

// sendRequest('/protected', 'GET').then((response) => {
//     console.log('Use of protected route successful:', response);
//     })
//     .catch((error) => {
//         console.error('Use of protected route failed:', error);
// });

// ---- event listeners -----

sidebarOption.addEventListener('click', function (event) {
    sidebarOptionsList.classList.toggle('sidebar-options-reveal')
    sidebarOptionsList.classList.toggle('ul-forwards')
    sidebarOptionsList.classList.toggle('ul-reverse')
});

loginFormButton.addEventListener('submit', function (event) {
    event.preventDefault(); // Prevent the default form submission

    const form = event.target;
    const formData = new FormData(form); // Collect form data
    const csrfToken = getCSRFToken();
    console.log("csrfToken: ", csrfToken)

    fetch('/login', {
        method: 'POST',
        headers: {
            'x-csrf-token': csrfToken
        },
        body: formData // Send the form data
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.text(); // Or response.text() if the response isn't JSON
        })
        .then(data => {
            console.log('Success:', data);
        })
        .catch(error => {
            console.error('Error:', error);
        });
});

logoutFormButton.addEventListener('click', function (event) {
    event.preventDefault();

    const csrfToken = getCSRFToken();
    console.log("csrfToken: ", csrfToken)

    fetch('/logout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'x-csrf-token': csrfToken
        },
        body: JSON.stringify({})
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.text();
        })
        .then(data => {
            console.log('Success:', data);
        })
        .catch(error => {
            console.error('Error:', error);
        });
})

inputUser.addEventListener('change', function () {
    file = this.files[0];
    dropAreaUser.classList.add('active');
});
// when file is inside drag area
dropAreaUser.addEventListener('dragover', (event) => {
    event.preventDefault();
    dropAreaUser.classList.add('active');
    dragText.textContent = 'release to Upload';
    dragButton.style.display = 'none';
    // console.log('File is inside the drag area');
});
// when file leaves the drag area
dropAreaUser.addEventListener('dragleave', () => {
    dropAreaUser.classList.remove('active');
    // console.log('File left the drag area');
    dragText.textContent = 'drag your file here';
});
// when file is dropped
dropAreaUser.addEventListener('drop', (event) => {
    event.preventDefault();
    dropAreaUser.classList.add('dropped');
    // console.log('File is dropped in drag area');
    file = event.dataTransfer.files[0]; // grab single file even if user selects multiple files
    // console.log(file);
    displayFile(uploadedFileUser, dropAreaUser);
});
function displayFile(uploadedFile, dropArea) {
    let fileType = file.type;
    // console.log(fileType);
    let validExtensions = ["image/*"];
    if (validExtensions.includes(fileType)) {
        let fileReader = new FileReader();
        fileReader.onload = () => {
            uploadedFile.innerHTML = `<div class="dragText">uploaded</div>
        <div class="uploaded-file">${file.name}</div>`;
            dropArea.classList.add("dropped");
        };
        fileReader.readAsDataURL(file);
    } else {
        alert("This is not an Image File");
        dropArea.classList.remove("active");
        dragText.textContent = "Drag and drop your file, or";
        dragButton.style.display = "unset";
    }
}
// switchDl.addEventListener('click', toggleColorScheme);
darkSwitch.addEventListener('click', toggleDarkMode);
actButtonsAll.forEach( button => button.addEventListener('click', (e) => {
    toggleFeed(document.getElementById("activity-" + e.target.id),document.getElementById("activity-feed-" + e.target.id),  e.target);
    console.log('activity-' + e.target.id);
}) );
// open modals
// TODO refactor the open and close modals
openLoginModal.addEventListener('click', () => loginModal.style.display = 'block');
openEditUserModal.addEventListener('click', () => editUserModal.style.display = 'block');
// openAccSettingsModal.addEventListener('click', () => accSettingsModal.style.display = 'block');
// openViewStatsModal.addEventListener('click', () => viewStatsModal.style.display = 'block');
// openRemoveAccModal.addEventListener('click', () => removeAccModal.style.display = 'block');
// close modals
closeLoginModal.addEventListener('click', () => loginModal.style.display = 'none');
closeEditUserModal.addEventListener('click', () => editUserModal.style.display = 'none');
// closeAccSettingsModal.addEventListener('click', () => accSettingsModal.style.display = 'none');
// closeViewStatsModal.addEventListener('click', () => viewStatsModal.style.display = 'none');
// closeRemoveAccModal.addEventListener('click', () => removeAccModal.style.display = 'none');
window.addEventListener('click', ({ target }) => {
    switch (target) {
        case loginModal:
            loginModal.style.display = 'none';
            break;
        case editUserModal:
            editUserModal.style.display = 'none';
            break;
        case accSettingsModal:
            accSettingsModal.style.display = 'none';
            break;
        case viewStatsModal:
            viewStatsModal.style.display = 'none';
            break;
        case removeAccModal:
            removeAccModal.style.display = 'none';
            break;
    }
});
// login / register / forgot
btnLogin.forEach(button =>
    button.addEventListener('click', (e) => logReg(e.target.id))
);
btnRegister.forEach(button =>
    button.addEventListener('click', (e) => logReg(e.target.id))
);
btnForgot.addEventListener('click', forgot);
// TODO get these working
// validate password
regPass.addEventListener('input', validatePass);
// check passwords match
regPassRpt.addEventListener('focusout', confirmPass);
// reverse the order of the password validation list delay
regPass.addEventListener('focusin', () => {
    setTimeout(() => {
        validList.classList.remove('ul-forwards');
        validList.classList.add('ul-reverse');
    }, 1000)
});
regPass.addEventListener('focusout', () => {
    setTimeout(() => {
        validList.classList.remove('ul-reverse');
        validList.classList.add('ul-forwards');
    }, 1000);
});