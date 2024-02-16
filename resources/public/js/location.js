function showPosition(position) {
    latitudeInput.value = position.coords.latitude.toFixed(6)
    longitudeInput.value = position.coords.longitude.toFixed(6)
}

const showError = (error) => {
    switch(error.code) {
        case error.PERMISSION_DENIED:
            showToast({
                'type': 'error',
                'message': "User denied the request for Geolocation."
            })
            break;
        case error.POSITION_UNAVAILABLE:
            showToast({
                'type': 'error',
                'message': "Location information is unavailable."
            })
            break;
        case error.TIMEOUT:
            showToast({
                'type': 'error',
                'message': "The request to get user location timed out."
            })
            break;
        default:
            showToast({
                'type': 'error',
                'message': "An unknown error occurred."
            })
            break;
    }
}