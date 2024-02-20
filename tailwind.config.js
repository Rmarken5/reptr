/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./service/internal/web/components/**/*.templ"
    ],
    theme: {
        extend: {},
        container: {
            center: true,
            padding: {
                DEFAULT: '1rem',
                sm: '2rem',
                lg: '4rem',
                xl: '5rem',
                '2xl': '5rem',
            },
        },
    },
    plugins: [
        require('@tailwindcss/forms'),
        require('@tailwindcss/typography'),
    ],
}

