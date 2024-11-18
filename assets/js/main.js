// variables
const switchDl = document.getElementById('switch-dl');

// functions
function toggleColorScheme() {
    // Get the current color scheme
    const currentScheme = document.documentElement.getAttribute('color-scheme');

    // Toggle between light and dark
    const newScheme = currentScheme === 'light' ? 'dark' : 'light';

    // Set the new color scheme
    document.documentElement.setAttribute('color-scheme', newScheme);
}

// event listeners
switchDl.addEventListener('click', toggleColorScheme);