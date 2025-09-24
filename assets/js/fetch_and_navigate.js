import { setActivePage, newContentLoaded } from "./main.js";
import { pages } from "./share.js";
import { showMainNotification } from "./notifications.js";
import { selectActiveFeed } from "./navigation.js";

const expect =
  "background-color: rgb(108 207 93); color: #000000; font-weight: bold; padding: .1rem; border-radius: 1rem;";
const standard =
  "background-color: transparent; color: #ffffff; font-weight: normal;";
const warn =
  "background-color: #000000; color: #e3c144; font-weight: bold; border: 1px solid #e3c144; padding: 5px; border-radius: 5px;";

export function navigateToPage(dest, entity) {
  // dest is a string - "channel"
  // entity is the template element:work
  // console.info("%cdest:%o", expect, dest);
  // console.info("%centity:%o", expect, entity);

  let link;
  if (dest === "home") {
    link = entity.getAttribute(`data-dest`);
  } else {
    link = entity.getAttribute(`data-${dest}-id`);
  }
  // console.log("link: ", link);
  const page = dest + "Page";

  if (!link) {
    console.error(`${dest} ID is missing`);
    return Promise.reject(`${dest} ID is missing`);
  }

  setActivePage(dest);
  changePage(page);
  return fetchData(dest, link);
}

export function changePage(page) {
  // console.log("%cpage: ", expect, page);
  let pageId;
  if (typeof page != "string") {
    pageId = page.id;
  } else {
    const pageSlice = page.slice(0, -4);
    pageId = pageSlice + "-page";
  }

  pages.forEach((element) => {
    if (element.id === pageId) {
      element.classList.add("active-feed");
    } else if (pageId !== "home-page" && element.id === "home-page") {
      element.classList.remove("active-feed");
    } else {
      // console.log("%cpage: ",warn, page);
      // TODO need to modify home-page template to populate by injection
      // TODO when injected, the content can be cleared

      // clear the content of the page previously active
      element.innerHTML = "";

      // TODO with injected home-page, this can be removed
      // need to clear content of the calendar dropdown, or the calendar won't work on otehr pages after navigation

      element.classList.remove("active-feed");
    }
  });

  selectActiveFeed();
}

function renderContent(data, entity) {
  const content = data[`${entity}sHTML`];
  const target = document.getElementById(`${entity}-page`);
  if (target && content) {
    target.innerHTML = content;

    // Add error class if status is present
    if (data.status && data.status >= 400) {
      target.classList.add(`error-${data.status}`);
    } else {
      // Clean up previous error classes if re-rendering normal content
      target.className = target.className
        .split(" ")
        .filter((cls) => !cls.startsWith("error-"))
        .join(" ");
    }

    // Dispatch custom event with details
    const newContentLoaded = new CustomEvent("newContentLoaded", {
      detail: {
        entity,
        status: data.status || 200,
      },
    });
    document.dispatchEvent(newContentLoaded);
  } else {
    console.warn("Target element or content missing");
  }
}

export async function fetchData(entity, Id) {
  const stateObj = { entity, id: Id };
  const url = `/cdx/${entity}/${Id}`;
  history.pushState(stateObj, "", url);

  try {
    const response = await fetch(`/${entity}/${Id}`);
    const data = await response.json();

    // Always render what backend sent
    renderContent(data, entity);

    // If backend included a non-OK status, treat it as error
    if (!response.ok || (data.status && data.status >= 400)) {
      throw {
        error: new Error(`Backend error ${data.status || response.status}`),
        data,
      };
    }
  } catch (e) {
    if (e.data) {
      // Render the backend-provided error page
      renderContent(e.data, entity);
    } else {
      // Fallback if no data available
      const target = document.getElementById(`${entity}-page`);
      if (target) {
        target.innerHTML = `<div class="error">Something went wrong</div>`;
      }
    }
  }
}

export async function fetchHome() {
  const entity = "home";
  const stateObj = { entity: entity };
  const url = "/";
  history.pushState(stateObj, "", url);

  try {
    const response = await fetch(`/home`);
    const data = await response.json();
    renderContent(data, entity);
  } catch (e) {
    if (e.data) {
      // Render the backend-provided error page
      renderContent(e.data, entity);
    } else {
      // Fallback if no data available
      const target = document.getElementById(`home-page`);
      if (target) {
        target.innerHTML = `<div class="error">Something went wrong</div>`;
      }
    }
  }
}

function showLoader() {
  const backdrop = document.getElementById("loader-backdrop");
  backdrop?.classList.remove("hidden");
}

function hideLoader() {
  const backdrop = document.getElementById("loader-backdrop");
  if (!backdrop) return;

  backdrop.classList.add("hidden");
  setTimeout(() => backdrop.remove(), 500); // fade-out duration
}

document.addEventListener("DOMContentLoaded", async function () {
  showLoader(); // start the spinner

  const path = window.location.pathname;
  const [, dest, id] = path.match(/\/cdx\/(\w+)\/([^/]+)/) || [];
  const homeDiv = document.getElementById("home-page");

  try {
    if (!dest) {
      if (homeDiv && homeDiv.children.length > 0) {
        // 400 page already injected by server, do nothing
        return;
      } else {
        setActivePage("home");
        changePage("homePage");
        await fetchHome();
      }
    } else {
      const page = dest + "Page";
      setActivePage(dest);
      changePage(page);
      await fetchData(dest, id);
    }
  } catch (e) {
    console.error("Page load failed:", e);
  } finally {
    // Hide loader **after content is rendered**
    requestAnimationFrame(() => hideLoader());
  }
});
