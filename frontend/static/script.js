htmx.logAll();

const generateButtons = document.getElementById('generate-buttons')?.querySelectorAll('button');
const questions = document.getElementById('questions');

const placeholder = document.getElementById('question-placeholder');
if (placeholder?.style != null) {
    placeholder.style.display = 'none';
}

generateButtons?.forEach(generateButton => {
    generateButton.addEventListener('htmx:beforeSend', () => {
        generateButtons?.forEach(button => button.disabled = true);
        placeholder.style.display = 'flex';
        questions.insertBefore(placeholder, null);
        placeholder.scrollIntoView({ behavior: 'smooth'});
    });

    generateButton.addEventListener('htmx:afterOnLoad', () => {
        generateButtons?.forEach(button => button.disabled = false);
        placeholder.style.display = 'none';
        questions.lastElementChild?.scrollIntoView({ behavior: 'smooth'});
    });
});

setTimeout(() => {
    const errorNotification = document.getElementById('error-notification');
    console.log('errorNotification', errorNotification);

    document.body.addEventListener('htmx:responseError', function(event) {
        const detail = event.detail;
        if (detail.xhr.status >= 400) {
            errorNotification.style.display = 'fixed'
            errorNotification.textContent = detail.xhr.toString()
            setTimeout(() => {
                errorNotification.style.display = 'none'
                errorNotification.textContent = ''
            },2000)
        }
    });
}, 1000)