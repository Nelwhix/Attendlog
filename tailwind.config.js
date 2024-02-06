/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./views/**/*.tmpl"],
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