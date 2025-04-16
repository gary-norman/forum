export function toggleColorScheme() {
    // dark mode
    const darkSwitch = document.querySelector("#sidebar-option-darkmode");

    // Get the current color scheme
    const currentScheme = document.documentElement.getAttribute("color-scheme");
    // Toggle between light and dark
    const newScheme = currentScheme === "light" ? "dark" : "light";
    // Set the new color scheme
    document.documentElement.setAttribute("color-scheme", newScheme);
    localStorage.setItem("color-scheme", newScheme);
}

// Apply persisted color scheme on page load
// INFO was a DOMContentLoaded function
export function saveColourScheme() {
    const savedScheme = localStorage.getItem("color-scheme");
    if (savedScheme) {
        document.documentElement.setAttribute("color-scheme", savedScheme);
    }
}


export function toggleDarkMode() {
    const checkbox = document.querySelector("#darkmode-checkbox");
    checkbox.checked = !checkbox.checked;
    console.log("toggle dark mode");
}