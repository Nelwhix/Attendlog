// this function finds the distance between two points using their
// longitude and latitude
function distanceBtwTwoPoints(lat1, lat2, long1, long2) {
    long1 = long1 * Math.PI / 180;
    long2 = long2 * Math.PI / 180;
    lat1 = lat1 * Math.PI / 180;
    lat2 = lat2 * Math.PI / 180;

    const longitudeDiff = long2 - long1;
    const latitudeDiff = lat2 - lat1;
    const a = Math.pow(Math.sin(latitudeDiff / 2), 2)
        + Math.cos(lat1) * Math.cos(lat2)
        * Math.pow(Math.sin(longitudeDiff / 2),2);

    const c = 2 * Math.asin(Math.sqrt(a));
    const r = 6371;

    return(c * r);
}