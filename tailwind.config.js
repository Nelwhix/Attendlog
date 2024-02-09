/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./templates/**/*.tmpl"],
    theme: {
        extend: {
            fontFamily: {
                'inter': ['Inter', 'sans-serif']
            }
        },
    },
    daisyui: {
        themes: ["light"],
    },
    plugins: [
        require('daisyui'),
    ],
}