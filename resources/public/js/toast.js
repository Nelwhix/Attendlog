function showToast({type, message}) {
    const flashContainer = document.getElementById('flashContainer')
    const flashText = document.getElementById('flashText')
    const alertBox = document.getElementById('alertBox');

    flashText.innerHTML = message
    alertBox.classList.add(`alert-${type}`)
    flashContainer.style.display = 'block'

    setTimeout(() => {
        flashContainer.style.display = 'hidden'
    }, 5 * 1000) // 5 seconds
}