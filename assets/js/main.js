import {channelPage, homePage, pages, selectActiveFeed, userPage} from "./share.js";
// variables
//user information

// dark mode
const switchDl = document.querySelector('#switch-dl');
const darkSwitch = document.querySelector('#sidebar-option-darkmode');
// activity buttons
let actButtonContainer;
let actButtonsAll;
// activity feeds
let activityFeeds;
let activityFeedsContentAll;
// sidebar elements
const userProfileImage = document.querySelectorAll('.profile-pic');
const userProfileImageEmpty = document.querySelectorAll('.profile-pic--empty');
const sidebarOption = document.querySelector('#sidebar-options');
const sidebarOptionsList = document.querySelector('.sidebar-options-list');
// login/register buttons
// TODO overhaul the naming of these buttons
const loginTitle = document.querySelector('#login-title');
const loginForm = document.querySelector('#loginForm');
const loginFormButton = document.querySelector('#login');
const loginFormUser = document.querySelector('#login_username');
const loginFormPassword = document.querySelector('#login_password');
const logoutFormButton = document.querySelector('#logout');
const btnLogin = document.querySelectorAll('[id^="btn_login-"]');
const btnRegister = document.querySelectorAll('[id^="btn_register-"]');
const btnForgot = document.querySelector('#btn_forgot');
// join channel
const joinChannelButton = document.querySelector('#join-channel-btn')
//input fields with fancy animations
const styledInputs = document.querySelectorAll(
    'textarea, input[type="text"], input[type="password"], input[type="email"]'
)
// modals
const modals = document.querySelectorAll(".modal");
// popovers
const popovers = document.querySelectorAll("[popover]")
console.log("popovers: ", popovers)
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
const openLoginModalFallback = document.querySelector('#btn-open-login-modal-fallback');
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
// page select
const selectDropdown = document.getElementById('pagedrop');
// ------
let file;
let filename;
// document.addEventListener('DOMContentLoaded', () => {
//
//
// });
document.addEventListener('DOMContentLoaded', function () {
    toggleUserInteracted("add");
    actButtonContainer = document.querySelector('#activity-bar');
    if (!actButtonContainer) {
        console.error('Activity bar not found');
        return;
    }
    actButtonsAll = actButtonContainer.querySelectorAll('button');
    activityFeeds = document.querySelector('#user-activity-feeds');
    if (!activityFeeds) {
        console.error('User activity feeds not found');
        return;
    }
    activityFeedsContentAll = activityFeeds.querySelectorAll('[id^="activity-feed-"]');
    // Log the buttons for debugging purposes (optional)
    if (actButtonsAll) {
        actButtonsAll.forEach( button => button.addEventListener('click', (e) => {
            toggleFeed(document.getElementById("activity-" + e.target.id),document.getElementById("activity-feed-" + e.target.id),  e.target);
            console.log('activity-' + e.target.id);
        }) );
    }

    if (typeof getUserProfileImageFromAttribute === 'function') {
        getUserProfileImageFromAttribute();
    } else {
        console.error('getUserProfileImageFromAttribute is not defined');
    }
    if (typeof getInitialFromAttribute === 'function') {
        getInitialFromAttribute();
    } else {
        console.error('getInitialFromAttribute is not defined');
    }
    console.log("activity buttons:", actButtonsAll)
});

// SECTION ----- functions ------
// toggle user-interacted class to input fields to prevent label animation before they are selected
function toggleUserInteracted(action) {
    styledInputs.forEach(input => {
        if (action === "add") {
        input.addEventListener("focus", function () {
            this.closest('.input-wrapper').classList.add("user-interacted");
        });
        } if (action === "remove") {
            document.querySelectorAll('.input-wrapper').forEach(element => {
                element.classList.remove("user-interacted");
            })
        }
    });
}

// revert to original state if modals are closed
// TODO get this popover reset to work without constant click listeners

document.addEventListener("click", () => {
    popovers.forEach(popover => {
        if (popover.matches(":popover-open")) {
            console.log("Open popover: ", popover);
            popover.addEventListener('toggle', () => {toggleUserInteracted("remove")})
        }
    });
});



// revert to original state if modals are closed
modals.forEach(modal => {
    const observer = new MutationObserver(mutations => {
        mutations.forEach(mutation => {
            if (mutation.attributeName === "style") {
                const currentDisplay = getComputedStyle(modal).display;
                if (currentDisplay === "none") {
                    console.log("Modal closed:", modal);
                    toggleUserInteracted("remove")
                }
            }
        });
    });

    observer.observe(modal, { attributes: true, attributeFilter: ["style"] });
});


// organize z-indices
const makeZIndexes = (layers) =>
    layers.reduce((agg, layerName, index) => {
        const valueName = `z-index-${layerName}`;
        agg[valueName] = index * 100;

        return agg;
    }, {})

function getRandomInt(max) {
    return Math.floor(Math.random() * max);
}

function getUserProfileImageFromAttribute() {
    for (let i = 0; i < userProfileImage.length; i++) {
        let attArr = ['user', 'auth', 'channel'];
        attArr[0] = userProfileImage[i].getAttribute('data-image-user');
        attArr[1] = userProfileImage[i].getAttribute('data-image-auth');
        attArr[2] = userProfileImage[i].getAttribute('data-image-channel');
        // console.table(attArr)
        if (attArr[0]) { // Ensure the `data-image-user` attribute has a value
            userProfileImage[i].style.background = `url('${attArr[0]}') no-repeat center`;
            userProfileImage[i].style.backgroundSize = 'cover'; // Add `cover` for background sizing
        } else if (attArr[1]) { // Ensure the `data-image-auth` attribute has a value
            userProfileImage[i].style.background = `url('${attArr[1]}') no-repeat center`;
            userProfileImage[i].style.backgroundSize = 'cover'; // Add `cover` for background sizing
        } else if (attArr[2]) { // Ensure the `data-image-channel` attribute has a value
            userProfileImage[i].style.background = `url('${attArr[2]}') no-repeat center`;
            userProfileImage[i].style.backgroundSize = 'cover'; // Add `cover` for background sizing
        }
        else {
            console.warn('No data-image- attribute value found for element:', userProfileImage[i]);
        }
    }
}

function getInitialFromAttribute() {
    console.log('getInitialFromAttribute running...')
    console.log('Elements found: ', userProfileImageEmpty.length)
    const colorsArr = [
        ['var(--color-hl-blue)', 'var(--color-light-1)'],
        ['var(--color-hl-green)', 'var(--color-dark-1)'],
        ['var(--color-hl-orange)', 'var(--color-light-1)'],
        ['var(--color-hl-pink)', 'var(--color-light-1)'],
        ['var(--color-hl-yellow)', 'var(--color-dark-1)'],
        ['var(--color-hl-red)',  'var(--color-light-1)']]
    let userTheme = getRandomInt(6)
    for (let i = 0; i < userProfileImageEmpty.length; i++) {
        let attArr = ['user', 'sidebar-user', 'channel'];
        attArr[0] = userProfileImageEmpty[i].getAttribute('data-name-user');
        attArr[1] = userProfileImageEmpty[i].getAttribute('data-name-user-sidebar');
        attArr[2] = userProfileImageEmpty[i].getAttribute('data-name-channel');
    if (attArr[0]) {
        userProfileImageEmpty[i].style.background = colorsArr[userTheme][0];
        userProfileImageEmpty[i].style.color = colorsArr[userTheme][1];
        userProfileImageEmpty[i].style.fontSize = '2rem';
        userProfileImageEmpty[i].setAttribute('data-initial', Array.from(`${attArr[0]}`)[0]);
        console.log('attArr[0]');
    } else if (attArr[1]) {
        userProfileImageEmpty[i].style.background = colorsArr[userTheme][0];
        userProfileImageEmpty[i].style.color = colorsArr[userTheme][1];
        userProfileImageEmpty[i].style.fontSize = '5rem';
        userProfileImageEmpty[i].setAttribute('data-initial', Array.from(`${attArr[1]}`)[0]);
        console.log('attArr[0]');
    } else if (attArr[2]) {
        let theme = getRandomInt(6)
        userProfileImageEmpty[i].style.background = colorsArr[theme][0];
        userProfileImageEmpty[i].style.color = colorsArr[theme][1];
        userProfileImageEmpty[i].setAttribute('data-initial', Array.from(`${attArr[2]}`)[0]);
        console.log('attArr[1]');
    }
    else {
        console.warn('No data-name- attribute value found for element:', userProfileImage[i]);
    }
  }
}

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
    // const allFeedsExceptTarget = Array.from(activityFeedsContentAll).filter(feed => feed.id !== targetFeedContent.id);

    actButtonsAll.forEach( button => button.classList.remove('btn-active') );
    activityFeedsContentAll.forEach(feed => {
        feed.classList.remove('collapsible-expanded');
        feed.classList.add('collapsible-collapsed');
    });
    setTimeout(() => {
        targetFeedContent.classList.remove('collapsible-collapsed');
        targetFeedContent.classList.add('collapsible-expanded');
        selectActiveFeed();
        targetButton.classList.toggle('btn-active'); }, timeOut);
    setTimeout(() => {
        Array.from(activityFeedsContentAll).forEach(feed => {
            const filtersRow = feed.parentElement.querySelector('.filters-row');
            if (filtersRow) {
                filtersRow.classList.add('hide-feed');
            }
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

// showMainNotification changes the element ID to provide feedback to user
function showMainNotification(message) {
    const notification = document.getElementById('notification-main');
    const notificationContent = document.getElementById('notification-main-content');
    notificationContent.textContent = message;
    notification.style.display = 'flex';
    setTimeout(() => {
        notification.style.display = 'none';
    }, 3000); // Hide after 3 seconds
}
function showNotification(elementID, message, success) {
    const notification = document.getElementById(elementID);
    notification.textContent = message;
    notification.style.color = "var(--color-hl-green)";
    if (!success) {
        setTimeout(() => {
            notification.textContent = "sign in to codex";
            notification.style.color = "var(--color-fg-1)";
        }, 3000); // Hide after 3 seconds
    }
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

// SECTION ---- event listeners -----

// --- select page ---
selectDropdown.addEventListener('change', () => {
    const selectedValue = selectDropdown.value;
    console.log("selectedValue: ", selectedValue);
    pages.forEach((element) => {
        console.log("elementID: ", element.id, "selectedPage: ", selectedValue);
        if (element.id === selectedValue) {
            element.classList.add('active-feed');
        } else {
            element.classList.remove('active-feed');
        }
    });
});
// --- sidebar options dropdown ---
sidebarOption.addEventListener('click', function (event) {
    sidebarOptionsList.classList.toggle('sidebar-options-reveal')
    sidebarOptionsList.classList.toggle('ul-forwards')
    sidebarOptionsList.classList.toggle('ul-reverse')
});

// --- login ---
if (loginForm) {
    loginForm.addEventListener('submit', function (event) {
        event.preventDefault(); // Prevent the default form submission
        // const form = event.target;
        // const formData = new FormData(form); // Collect form data
        const csrfToken = getCSRFToken();
        console.log("csrfToken: ", csrfToken)

        fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'x-csrf-token': csrfToken
            },
            // body: formData, // Send the form data
            body: JSON.stringify({
                username: document.getElementById('loginFormUser').value,
                password: document.getElementById('loginFormPassword').value,
            }),
            cache: "no-store"
        })
            .then(response => {
                if (response.ok) {
                    return response.json();
                } else {
                    throw new Error('Login failed.');
                }
            })
            .then(data => {
                if (data.message === "incorrect password") {
                    showNotification('login-title',data.message, false);
                } else {
                    showNotification('login-title', data.message, true);
                    setTimeout(() => {
                        window.location.href = '/';
                    }, 2000);
                }
            })
            .catch(error => {
                console.error('Error:', error);
                showMainNotification('An error occurred during login.');
            });
    });
}

if (logoutFormButton) {
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
            body: JSON.stringify({}),
            cache: 'no-store'
        })
            .then(response => {
                if (response.ok) {
                    return response.json(); // Parse JSON response
                } else {
                    throw new Error('Logout failed.');
                }
            })
            .then(data => {
                // Show notification based on the response from the server
                showMainNotification(data.message);
                // Optional: Redirect the user after showing the notification
                setTimeout(() => {
                    window.location.href = '/'; // Replace with your desired location
                }, 3500);
            })
            .catch(error => {
                console.error('Error:', error);
                showMainNotification('An error occurred during logout.');
            });
    })
}
joinChannelButton.addEventListener('submit', function (event) {
    event.preventDefault();
    // const form = event.target;
    // const formData = new FormData(form); // Collect form data
    const csrfToken = getCSRFToken();

    console.log("csrfToken: ", csrfToken)

    fetch('/channels/join', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'x-csrf-token': csrfToken
        },
        body: JSON.stringify( {
            channelId: document.getElementById('join-channel-id').value,
            agree: document.getElementById('rules-agree-checkbox').value,
        }),
        cache: 'no-store'
    })
        .then(response => {
            if (response.ok) {
                return response.json(); // Parse JSON response
            } else {
                throw new Error('Join channel failed.');
            }
        })
        .then(data => {
            // Show notification based on the response from the server
            showMainNotification(data.message);
            // Optional: Redirect the user after showing the notification
            setTimeout(() => {
                window.location.href = '/'; // Replace with your desired location
            }, 3500);
        })
        .catch(error => {
            console.error('Error:', error);
            showMainNotification('An error occurred while joining channel.');
        });
})

// SECTION ----- drag and drop ----
// get user image from manual click
inputUser.addEventListener('change', function () {
    file = this.files[0];
    dropAreaUser.classList.add('active');
});
// get post image from manual click
inputPost.addEventListener('change', function () {
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


// open modals
// TODO refactor the open and close modals
if (openLoginModal) {
    openLoginModal.addEventListener('click', () => loginModal.style.display = 'block');
}
openLoginModalFallback.addEventListener('click', () => loginModal.style.display = 'block');
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