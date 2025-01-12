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
let forgotVisible = false;
// login/register modal
const modal = document.querySelector('#form-login-container');
// Get the button that opens the modal
const openLoginModal = document.querySelector('#btn-open-login-modal');
// Get the <span> element that closes the modal
const closeLoginModal = document.getElementsByClassName("close")[0];
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
const dropArea = document.querySelector('#drop_zone');
const uploadedFile = document.querySelector('#uploadedFile')
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
// sendRequest('/protected', 'GET').then((response) => {
//     console.log('Use of protected route successful:', response);
//     })
//     .catch((error) => {
//         console.error('Use of protected route failed:', error);
// });

// ---- event listeners -----

// drag and drop
// dropButton.addEventListener('click', input.click.bind(input), false);
// when browse
// loginFormButton.addEventListener('click', () => {
//     sendRequest('/login', 'POST')
//         .then((response) => {
//             console.log('Login successful:', response);
//         })
//         .catch((error) => {
//             console.error('Login failed:', error);
//         });
// })


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
    displayFile();
});
function displayFile() {
    let fileType = file.type;
    // console.log(fileType);
    let validExtensions = ["image/*"];
    if (validExtensions.includes(fileType)) {
        let fileReader = new FileReader();
        fileReader.onload = () => {
            uploadedFile.innerHTML = `<div class="dragText">uploaded</div>
        <div class="uploadedFile">${file.name}</div>`;
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
// login register modal
openLoginModal.addEventListener('click', () => modal.style.display = 'block');
closeLoginModal.addEventListener('click', () => modal.style.display = 'none');
window.addEventListener('click', ({ target }) => {
    if (target === modal) {
        modal.style.display = 'none';
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