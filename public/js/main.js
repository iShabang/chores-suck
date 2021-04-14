function sideClick(sel,show) {
    var rs = getComputedStyle(document.querySelector(':root'));
    var c1 = rs.getPropertyValue('--clr-primary-600');
    var c2 = rs.getPropertyValue('--clr-primary-200');
    //find each header. Change to default style
    headers = document.getElementsByClassName('s-head');
    for (h of headers) {
        h.style.backgroundColor = c1;
    }
    //find the selected header. Change to selected style.
    document.getElementById(sel).style.backgroundColor = c2;
    //find each content div and hide it
    divs = document.getElementsByClassName('v-content');
    for (d of divs) {
        d.style.display = 'none';
    }
    //show the selected content
    document.getElementById(show).style.display = 'block';
}

function hamburger(nav) {
    document.getElementById(nav).classList.toggle('flex')
}
