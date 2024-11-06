const btnLogin = document.getElementById("btn_login");
const btnRegister = document.getElementById("btn_register");
const btnGuest = document.getElementById("btn_guest");
const door = document.getElementById("access-door");
let doorClosed = true;
let clickCount = 0;

btnLogin.addEventListener("click", () => {
  if (clickCount === 0) {
    door.classList.add("grid-welcome-access-door-open");
    doorClosed = false;
    clickCount = 1;
  } else if (doorClosed) {
    door.classList.remove("grid-welcome-access-door-close");
    door.classList.add("grid-welcome-access-door-open");
    doorClosed = false;
  } else {
    door.classList.remove("grid-welcome-access-door-open");
    door.classList.add("grid-welcome-access-door-close");
    doorClosed = true;
  }
});
