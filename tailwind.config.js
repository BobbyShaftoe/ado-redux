/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./views/**/*.templ",
        "static/js/**/*.js",
    ],
    darkMode: 'class',
    theme: {

        extend: {
            colors: {
                'nav-blue': '#3187df',
                'nav-blue-light': '#7db3ea',
                'nav-blue-lightest': '#d5ebfa',
                'nav-blue-dark': '#1d4e81',
                'nav-blue-darker': '#14375a',
                'nav-blue-darkest': '#182d40',
                'button-blue': 'rgb(21 30 48 / var(--tw-bg-opacity))',
                'button-blue-mid': 'rgb(28 41 66 / var(--tw-bg-opacity))',
                'button-blue-dark': 'rgb(11 15 24 / var(--tw-bg-opacity))',
            }
        },
    },
    plugins: [],
}

