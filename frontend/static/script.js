// htmx.logAll();
window.hideNotification = () => {
    const errorNotification = document.getElementById('error-notification');
    errorNotification.classList.remove('opacity-100');
    setTimeout(() => errorNotification.style.display = 'none', 2000)
}

htmx.onLoad(body => {
    const errorNotification = document.getElementById('error-notification');
    document.body.addEventListener('htmx:afterSettle', event => {
        const detail = event.detail;
        if (detail.xhr.status >= 400) {
            errorNotification.style.display = 'flex';
            errorNotification.classList.add('opacity-100');
            setTimeout(() => hideNotification(), 2000);
        }
    });
})
// document.body.addEventListener('htmx:load', () => {
//     htmx.logAll();
//
//     const generateButtons = document.getElementById('generate-buttons')?.querySelectorAll('button');
//     const questions = document.getElementById('questions');
//
//     const placeholder = document.getElementById('question-placeholder');
//     if (placeholder?.style != null) {
//         placeholder.style.display = 'none';
//     }
//
//     generateButtons?.forEach(generateButton => {
//         generateButton.addEventListener('htmx:beforeSend', () => {
//             generateButtons?.forEach(button => button.disabled = true);
//             placeholder.style.display = 'flex';
//             questions.insertBefore(placeholder, null);
//             placeholder.scrollIntoView({ behavior: 'smooth'});
//         });
//
//         generateButton.addEventListener('htmx:afterOnLoad', () => {
//             generateButtons?.forEach(button => button.disabled = false);
//             placeholder.style.display = 'none';
//             questions.lastElementChild?.scrollIntoView({ behavior: 'smooth'});
//         });
//     });
//
// });