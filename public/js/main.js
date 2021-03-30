function sidebarClick(p1,p2,p3,p4) {
    var rs = getComputedStyle(document.querySelector(':root'))
    document.getElementById(p1).style.backgroundColor = rs.getPropertyValue('--clr-primary-200');
    document.getElementById(p2).style.backgroundColor = rs.getPropertyValue('--clr-primary-600');
    document.getElementById(p3).style.display = 'none';
    document.getElementById(p4).style.display = 'block';
}

function sidebarChores() {

}

function sidebarGroups() {

}