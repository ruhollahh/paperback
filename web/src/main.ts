import '@carbon/web-components/es/components/modal/index.js';
import '@carbon/web-components/es/components/button/index.js';


document.getElementById('modal-example-button')?.addEventListener('click', () => {
    const modal = document.getElementById('modal-example')
    if (modal) {
        modal.open = true;
    }
});