let mode = localStorage.getItem("dark-mode");
const lightSwitches = document.querySelectorAll('.light-switch');
if (mode) {
    console.log('On load dark mode is set in local storage');
    document.documentElement.classList.add("dark");
    lightSwitches.forEach((lightSwitch) => {
        lightSwitch.checked = true;
    });
}
if (mode !== 'true') {
    console.log('On load light mode is set in local storage');
    document.documentElement.classList.remove("dark");
    document.documentElement.classList.add("light");
    lightSwitches.forEach((lightSwitch) => {
        lightSwitch.checked = false;
    });
}


if (lightSwitches.length > 0) {
    lightSwitches.forEach((lightSwitch, i) => {
        if (localStorage.getItem('dark-mode') === 'true') {
            console.log('Local storage value is set to dark');
            lightSwitch.checked = true;
        }
        lightSwitch.addEventListener('change', () => {
            console.log('Light switch clicked');
            const {checked} = lightSwitch;
            lightSwitches.forEach((el, n) => {
                if (n !== i) {
                    el.checked = checked;
                }
            });
            if (lightSwitch.checked) {
                console.log('Dark mode option is checked');
                document.documentElement.classList.add('dark');
                localStorage.setItem('dark-mode', true);
            } else {
                console.log('Dark mode option is not checked');
                document.documentElement.classList.remove('dark');
                localStorage.setItem('dark-mode', false);
            }
        });
    });
}
