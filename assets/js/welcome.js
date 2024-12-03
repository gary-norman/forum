const btnLogin = document.getElementById("btn_login");
const btnRegister = document.getElementById("btn_register");
const btnGuest = document.getElementById("btn_guest");
const door = document.getElementById("access-door");
const doorShadow = document.getElementById("access-door-shadow");
const shadowClip = document.getElementById("shadow-clip");
let login = false;
let register = false;
let guest = false;

btnLogin.addEventListener("click", () => {
  if (!login && !register) {
    console.log("!l !r a", login);
    door.classList.remove("grid-welcome-access-door-close");
    door.classList.add("grid-welcome-access-door-open");
    // shadowClip.classList.remove("grid-welcome-access-door-close");
    // shadowClip.classList.add("grid-welcome-access-door-open");
    doorShadow.classList.remove("grid-welcome-access-door-push-close");
    doorShadow.classList.add("grid-welcome-access-door-push-open");
    btnLogin.classList.toggle("button-3d-pressed");
    login = true;
    console.log("!l !r b", login);
  } else if (login && !register) {
    console.log("l !r a", login);
    door.classList.remove("grid-welcome-access-door-open");
    door.classList.add("grid-welcome-access-door-close");
    // shadowClip.classList.remove("grid-welcome-access-door-open");
    // shadowClip.classList.add("grid-welcome-access-door-close");
    doorShadow.classList.remove("grid-welcome-access-door-push-open");
    doorShadow.classList.add("grid-welcome-access-door-push-close");
    btnLogin.classList.toggle("button-3d-pressed");
    login = false;
    console.log("l !r b", login);
  } else if (!login && register) {
    console.log("!l r a", login);
    door.classList.remove("grid-welcome-access-door-open");
    door.classList.add("grid-welcome-access-door-close");
    // shadowClip.classList.remove("grid-welcome-access-door-open");
    // shadowClip.classList.add("grid-welcome-access-door-close");
    doorShadow.classList.remove("grid-welcome-access-door-push-open");
    doorShadow.classList.add("grid-welcome-access-door-push-close");
    setTimeout(() => { door.classList.remove("grid-welcome-access-door-close"); }, 500);
    doorShadow.classList.remove("grid-welcome-access-door-push-close");
    // shadowClip.classList.remove("grid-welcome-access-door-close");
    door.classList.add("grid-welcome-access-door-open");
    // shadowClip.classList.add("grid-welcome-access-door-open");
    doorShadow.classList.add("grid-welcome-access-door-push-open");
    btnLogin.classList.toggle("button-3d-pressed");
    btnRegister.classList.toggle("button-3d-pressed");
    login = true;
    register = false;
    console.log("!l r b", login);
  } else {
    console.log("else", login);
  }
});

btnRegister.addEventListener("click", () => {
  if (!login && !register) {
    console.log("!l !r a", login);
    door.classList.remove("grid-welcome-access-door-close");
    door.classList.add("grid-welcome-access-door-open");
    doorShadow.classList.remove("grid-welcome-access-door-push-close");
    doorShadow.classList.add("grid-welcome-access-door-push-open");
    btnRegister.classList.toggle("button-3d-pressed");
    register = true;
    console.log("!l !r b", login);
  } else if (!login && register) {
    console.log("!l r a", login);
    door.classList.remove("grid-welcome-access-door-open");
    door.classList.add("grid-welcome-access-door-close");
    doorShadow.classList.remove("grid-welcome-access-door-push-open");
    doorShadow.classList.add("grid-welcome-access-door-push-close");
    btnRegister.classList.toggle("button-3d-pressed");
    register = false;
    console.log("!l r b", login);
  } else if (login && !register) {
    console.log("l !r a", login);
    door.classList.remove("grid-welcome-access-door-open");
    doorShadow.classList.remove("grid-welcome-access-door-push-open");
    door.classList.add("grid-welcome-access-door-close");
    doorShadow.classList.add("grid-welcome-access-door-push-close");
    setTimeout(() => { door.classList.remove("grid-welcome-access-door-close"); }, 500);
    doorShadow.classList.remove("grid-welcome-access-door-push-close");
    door.classList.add("grid-welcome-access-door-open");
    doorShadow.classList.add("grid-welcome-access-door-push-open");
    btnLogin.classList.toggle("button-3d-pressed");
    btnRegister.classList.toggle("button-3d-pressed");
    register = true;
    login = false;
    console.log("l !r b", login);
  } else {

  }
});