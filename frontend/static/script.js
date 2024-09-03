// htmx.logAll();

window.placeholder = () => document.getElementById("question-placeholder");
window.generateButtons = () =>
  document.getElementById("generate-buttons")?.querySelectorAll("button");
window.questions = () => document.getElementById("questions");

window.hideNotification = () => {
  const errorNotification = document.getElementById("error-notification");
  errorNotification.classList.remove("opacity-100");
  setTimeout(() => (errorNotification.style.display = "none"), 2000);
};

htmx.onLoad((body) => {
  const errorNotification = document.getElementById("error-notification");
  document.body.addEventListener("htmx:afterSettle", (event) => {
    const detail = event.detail;
    if (detail.xhr.status >= 400) {
      errorNotification.style.display = "flex";
      errorNotification.classList.add("opacity-100");
      setTimeout(() => hideNotification(), 2000);
    }
  });

  document.body.addEventListener("htmx:beforeSend", (event) => {
    const requestPath = event.detail.pathInfo.requestPath;
    if (requestPath.match(/\/generate\?type=/)) {
      questions().appendChild(placeholder());
      placeholder().style.display = "flex";
      placeholder().scrollIntoView({ behavior: "smooth" });
      generateButtons()?.forEach((x) => (x.disabled = true));
    }
  });

  document.body.addEventListener("htmx:afterOnLoad", (event) => {
    const requestPath = event.detail.pathInfo.requestPath;
    if (requestPath.match(/\/generate\?type=/)) {
      placeholder().style.display = "none";
      generateButtons()?.forEach((x) => (x.disabled = false));
    }
  });
});
