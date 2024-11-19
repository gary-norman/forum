// variables
const switchDl = document.getElementById('switch-dl');
const darkSwitch = document.getElementById('sidebar-option-darkmode');
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

// event listeners
// switchDl.addEventListener('click', toggleColorScheme);
darkSwitch.addEventListener('click', toggleDarkMode);